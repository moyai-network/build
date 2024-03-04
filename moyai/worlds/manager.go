package worlds

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/mcdb"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	"os"
	"strings"
	"sync"
)

func Manager() *manager {
	return m
}

var m *manager

type manager struct {
	path string
	log  *logrus.Logger
	w    *world.World

	worldsMu sync.Mutex
	worlds   map[string]*world.World
}

func NewManager(w *world.World, path string, log *logrus.Logger) error {
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

	prov, err := mcdb.Open(m.path + "/" + name)
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

func (m *manager) DeleteWorld(name string) error {
	name = strings.ToLower(name)

	m.worldsMu.Lock()
	w, ok := m.worlds[name]

	if ok {
		for _, e := range w.Entities() {
			if p, ok := e.(*player.Player); ok {
				m.w.AddEntity(p)
				p.Teleport(m.w.Spawn().Vec3Middle())
				continue
			}
			_ = e.Close()
		}
		delete(m.worlds, name)
		_ = w.Close()
	}
	m.worldsMu.Unlock()

	err := os.RemoveAll(m.path + "/" + name)
	if err != nil {
		return fmt.Errorf("error deleting world %s: %s", name, err)
	}
	return nil
}

func (m *manager) Close() {
	m.worldsMu.Lock()
	defer m.worldsMu.Unlock()
	for _, w := range m.worlds {
		for _, e := range w.Entities() {
			if p, ok := e.(*player.Player); ok {
				m.w.AddEntity(p)
				p.Teleport(m.w.Spawn().Vec3Middle())
				continue
			}
			_ = e.Close()
		}
		_ = w.Close()
	}
}
