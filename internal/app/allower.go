package app

import (
	"context"
	"log/slog"
	"net"

	"github.com/google/uuid"
	flywayports "github.com/moyai-network/flyway/storage/ports"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
)

type Allower struct {
	whitelistStorage flywayports.WhitelistStorage
	whitelistEnabled bool
	log              *slog.Logger
}

func NewAllower(whitelistStorage flywayports.WhitelistStorage, whitelistEnabled bool, log *slog.Logger) *Allower {
	return &Allower{
		whitelistStorage: whitelistStorage,
		whitelistEnabled: whitelistEnabled,
		log:              log,
	}
}

func (a Allower) Allow(_ net.Addr, d login.IdentityData, _ login.ClientData) (string, bool) {
	if !a.whitelistEnabled || a.whitelistStorage == nil {
		return "", true
	}

	playerID, err := uuid.Parse(d.Identity)
	if err != nil || playerID == uuid.Nil {
		return "we couldn't verify your account identity", false
	}

	allowed, err := a.whitelistStorage.Whitelisted(context.Background(), playerID)
	if err != nil {
		if a.log != nil {
			a.log.Error("whitelist check failed", "error", err, "identity", d.Identity, "display_name", d.DisplayName, "xuid", d.XUID)
		}
		return "we couldn't verify whitelist access", false
	}
	if !allowed {
		return "you are not whitelisted on this server", false
	}
	return "", true
}
