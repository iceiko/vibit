package app

import "testing"

func TestRenderRouteKey(t *testing.T) {
	route := RouteKey{
		Kind:   MessageKindCommand,
		Module: " inventory ",
		Name:   " GrantItem ",
	}

	if got := RenderRouteKey(route); got != "inventory.GrantItem" {
		t.Fatalf("RenderRouteKey() = %q, want %q", got, "inventory.GrantItem")
	}
}

func TestRenderRouteKeyRequiresModuleAndName(t *testing.T) {
	if got := RenderRouteKey(RouteKey{Module: "inventory"}); got != "" {
		t.Fatalf("RenderRouteKey() = %q, want empty string", got)
	}

	if got := RenderRouteKey(RouteKey{Name: "GrantItem"}); got != "" {
		t.Fatalf("RenderRouteKey() = %q, want empty string", got)
	}
}

func TestMetadataOnlyIdentityFromSession(t *testing.T) {
	session := Session{
		ConnectionID:    " connection-1 ",
		SessionID:       " session-1 ",
		PlayerID:        " player-1 ",
		ConnectionEpoch: 7,
	}

	identity := MetadataOnlyIdentityFromSession(session)

	if identity.Status != IdentityValidationMetadataOnly {
		t.Fatalf("Status = %q, want %q", identity.Status, IdentityValidationMetadataOnly)
	}
	if identity.ActorKind != ActorKindPlayer || identity.ActorID != "player-1" {
		t.Fatalf("Actor = %q/%q, want player/player-1", identity.ActorKind, identity.ActorID)
	}
	if identity.PlayerID != "player-1" || identity.SessionID != "session-1" || identity.ConnectionID != "connection-1" {
		t.Fatalf("identity metadata = %#v, want normalized session metadata", identity)
	}
	if identity.SessionValidated || identity.PlayerIDValidated {
		t.Fatalf("metadata-only identity must not be marked validated: %#v", identity)
	}
	if identity.ConnectionEpoch != 7 {
		t.Fatalf("ConnectionEpoch = %d, want 7", identity.ConnectionEpoch)
	}
}

func TestMetadataOnlyIdentityWithoutPlayer(t *testing.T) {
	identity := MetadataOnlyIdentityFromSession(Session{ConnectionID: "connection-1"})

	if identity.ActorKind != ActorKindUnknown || identity.ActorID != "" || identity.PlayerID != "" {
		t.Fatalf("identity actor = %#v, want unknown actor without player_id", identity)
	}
}

func TestValidatedPlayerIdentity(t *testing.T) {
	identity := ValidatedPlayerIdentity(" player-2 ", Session{
		ConnectionID: "connection-1",
		SessionID:    "session-1",
		PlayerID:     "player-1",
	})

	if identity.Status != IdentityValidationValidated {
		t.Fatalf("Status = %q, want %q", identity.Status, IdentityValidationValidated)
	}
	if identity.PlayerID != "player-2" || identity.ActorID != "player-2" || identity.ActorKind != ActorKindPlayer {
		t.Fatalf("identity actor = %#v, want validated player-2", identity)
	}
	if !identity.SessionValidated || !identity.PlayerIDValidated {
		t.Fatalf("validated identity flags = %#v, want session and player validated", identity)
	}
}
