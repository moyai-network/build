package worlds

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/mcdb"
	"golang.org/x/exp/maps"
)

func Manager() *manager {
	return m
}

var m *manager

type manager struct {
	path string
	log  *slog.Logger
	w    *world.World

	worldsMu sync.Mutex
	worlds   map[string]*world.World
}

func NewManager(w *world.World, path string, log *slog.Logger) error {
	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("error creating world directory %s: %s", path, err)
	}

	dir, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("error loading world directory %s: %s", path, err)
	}

	m = &manager{
		path: path,
		log:  log,
		w:    w,

		worlds: map[string]*world.World{},
	}

	for _, d := range dir {
		if !d.IsDir() {
			continue
		}
		_, err := m.CreateWorld(d.Name())
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *manager) DefaultWorld() *world.World {
	return m.w
}

func (m *manager) World(name string) (*world.World, bool) {
	name = strings.ToLower(name)

	m.worldsMu.Lock()
	w, ok := m.worlds[name]
	m.worldsMu.Unlock()

	return w, ok
}

func (m *manager) Worlds() []*world.World {
	m.worldsMu.Lock()
	defer m.worldsMu.Unlock()

	return maps.Values(m.worlds)
}

func (m *manager) CreateWorld(name string) (*world.World, error) {
	name = strings.ToLower(name)

	prov, err := mcdb.Open(filepath.Join(m.path, name))
	if err != nil {
		return nil, fmt.Errorf("error loading world %s: %s", name, err)
	}
	prov.Settings().Name = name

	w := world.Config{
		Log:      m.log,
		Provider: prov,
		Entities: entity.DefaultRegistry,
	}.New()

	w.StopWeatherCycle()
	w.SetDefaultGameMode(world.GameModeCreative)
	w.SetTime(6000)
	w.StopTime()
	w.SetTickRange(0)
	w.StopThundering()
	w.StopRaining()

	m.worldsMu.Lock()
	m.worlds[name] = w
	m.worldsMu.Unlock()
	return w, nil
}

func PlayerWorld(p *player.Player) *world.World {
	if p == nil || p.Tx() == nil {
		return nil
	}
	return p.Tx().World()
}

func InDefaultWorld(p *player.Player) bool {
	man := Manager()
	return man != nil && PlayerWorld(p) == man.DefaultWorld()
}

func MovePlayer(p *player.Player, dst *world.World) bool {
	if p == nil || p.Tx() == nil || dst == nil {
		return false
	}
	if p.Tx().World() == dst {
		p.Teleport(dst.Spawn().Vec3Middle())
		return true
	}

	handle := p.Tx().RemoveEntity(p)
	if handle == nil {
		handle = p.H()
	}

	<-dst.Exec(func(tx *world.Tx) {
		e := tx.AddEntity(handle)
		moved, ok := e.(*player.Player)
		if !ok {
			return
		}
		moved.Teleport(tx.World().Spawn().Vec3Middle())
	})
	return true
}

func (m *manager) removeEntities(src *world.World) []*world.EntityHandle {
	var handles []*world.EntityHandle
	if src == nil {
		return handles
	}

	<-src.Exec(func(tx *world.Tx) {
		for e := range tx.Entities() {
			p, ok := e.(*player.Player)
			if ok {
				h := tx.RemoveEntity(p)
				if h == nil {
					h = p.H()
				}
				handles = append(handles, h)
				continue
			}
			_ = e.Close()
		}
	})
	return handles
}

func (m *manager) restorePlayers(handles []*world.EntityHandle) {
	if len(handles) == 0 || m.w == nil {
		return
	}

	<-m.w.Exec(func(tx *world.Tx) {
		spawn := tx.World().Spawn().Vec3Middle()
		for _, h := range handles {
			e := tx.AddEntity(h)
			p, ok := e.(*player.Player)
			if !ok {
				continue
			}
			p.Teleport(spawn)
		}
	})
}

func (m *manager) DeleteWorld(name string) error {
	name = strings.ToLower(name)

	m.worldsMu.Lock()
	w, ok := m.worlds[name]
	if ok {
		delete(m.worlds, name)
	}
	m.worldsMu.Unlock()

	if ok {
		handles := m.removeEntities(w)
		m.restorePlayers(handles)
		_ = w.Close()
	}

	err := os.RemoveAll(filepath.Join(m.path, name))
	if err != nil {
		return fmt.Errorf("error deleting world %s: %s", name, err)
	}
	return nil
}

func (m *manager) Close() {
	m.worldsMu.Lock()
	list := maps.Values(m.worlds)
	m.worlds = map[string]*world.World{}
	m.worldsMu.Unlock()

	for _, w := range list {
		handles := m.removeEntities(w)
		m.restorePlayers(handles)
		_ = w.Close()
	}
}
