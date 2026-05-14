package player

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	ModuleName = "player"

	CommandCreatePlayerAccount = "CreatePlayerAccount"
	QueryGetPlayerAccount      = "GetPlayerAccount"

	EventPlayerAccountCreated = "PlayerAccountCreated"
)

type AccountState string

const (
	AccountStateActive   AccountState = "active"
	AccountStateDisabled AccountState = "disabled"
	AccountStateDeleted  AccountState = "deleted"
)

func (s AccountState) IsValid() bool {
	switch s {
	case AccountStateActive, AccountStateDisabled, AccountStateDeleted:
		return true
	default:
		return false
	}
}

type Account struct {
	PlayerID     string
	DisplayName  string
	AccountState AccountState
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DisabledAt   *time.Time
	DeletedAt    *time.Time
}

type CreatePlayerAccountMutation struct {
	EventID      string
	OccurredAt   time.Time
	PlayerID     string
	DisplayName  string
	AccountState AccountState
	RequestedBy  string
}

type Repository interface {
	CreatePlayerAccount(context.Context, CreatePlayerAccountMutation) (Account, error)
	GetPlayerAccount(context.Context, string) (Account, error)
}

func NormalizeCreatePlayerAccountMutation(mutation CreatePlayerAccountMutation) (CreatePlayerAccountMutation, error) {
	var err error
	mutation.EventID, err = normalizeRequired("event_id", mutation.EventID)
	if err != nil {
		return CreatePlayerAccountMutation{}, err
	}
	mutation.PlayerID, err = normalizeRequired("player_id", mutation.PlayerID)
	if err != nil {
		return CreatePlayerAccountMutation{}, err
	}
	mutation.DisplayName, err = normalizeRequired("display_name", mutation.DisplayName)
	if err != nil {
		return CreatePlayerAccountMutation{}, err
	}
	mutation.RequestedBy, err = normalizeRequired("requested_by", mutation.RequestedBy)
	if err != nil {
		return CreatePlayerAccountMutation{}, err
	}
	if mutation.OccurredAt.IsZero() {
		return CreatePlayerAccountMutation{}, errors.New("player: occurred_at is required")
	}
	mutation.OccurredAt = mutation.OccurredAt.UTC()
	if mutation.AccountState == "" {
		mutation.AccountState = AccountStateActive
	}
	if mutation.AccountState != AccountStateActive {
		return CreatePlayerAccountMutation{}, errors.New("player: created accounts must start active")
	}
	return mutation, nil
}

func NormalizePlayerID(playerID string) (string, error) {
	return normalizeRequired("player_id", playerID)
}

func normalizeRequired(name string, value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("player: %s is required", name)
	}
	return value, nil
}
