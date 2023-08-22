package moyai

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"net"
	"strings"
)

type Allower struct {
	wl []string
}

func NewAllower(whitelist []string) *Allower {
	a := &Allower{wl: whitelist}
	return a
}

func (a Allower) Allow(_ net.Addr, d login.IdentityData, _ login.ClientData) (string, bool) {
	for _, u := range a.wl {
		if strings.EqualFold(u, d.DisplayName) {
			return "", true
		}
	}
	return text.Colourf("<red>You must be whitelisted in order to join the server.</red>"), false
}
