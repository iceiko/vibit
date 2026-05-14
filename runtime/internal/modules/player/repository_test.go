package player

import (
	"context"
	"testing"
	"time"
)

func TestRepositoryInterfaceIsStorageNeutral(t *testing.T) {
	var _ Repository = recordingRepository{}
}

func TestNormalizeCreatePlayerAccountMutationTrimsAndDefaultsActiveState(t *testing.T) {
	occurredAt := time.Date(2026, 5, 14, 11, 12, 13, 0, time.FixedZone("test", 8*60*60))

	mutation, err := NormalizeCreatePlayerAccountMutation(CreatePlayerAccountMutation{
		EventID:     " event-1 ",
		OccurredAt:  occurredAt,
		PlayerID:    " player-1 ",
		DisplayName: " Player One ",
		RequestedBy: " maintainer ",
	})
	if err != nil {
		t.Fatalf("NormalizeCreatePlayerAccountMutation() error = %v, want nil", err)
	}

	if mutation.EventID != "event-1" || mutation.PlayerID != "player-1" || mutation.DisplayName != "Player One" || mutation.RequestedBy != "maintainer" {
		t.Fatalf("normalized mutation = %#v, want trimmed required fields", mutation)
	}
	if mutation.AccountState != AccountStateActive {
		t.Fatalf("AccountState = %q, want %q", mutation.AccountState, AccountStateActive)
	}
	if mutation.OccurredAt.Location() != time.UTC {
		t.Fatalf("OccurredAt location = %s, want UTC", mutation.OccurredAt.Location())
	}
}

func TestNormalizeCreatePlayerAccountMutationRejectsNonActiveInitialState(t *testing.T) {
	_, err := NormalizeCreatePlayerAccountMutation(CreatePlayerAccountMutation{
		EventID:      "event-1",
		OccurredAt:   time.Date(2026, 5, 14, 1, 2, 3, 0, time.UTC),
		PlayerID:     "player-1",
		DisplayName:  "Player One",
		AccountState: AccountStateDisabled,
		RequestedBy:  "maintainer",
	})
	if err == nil {
		t.Fatal("NormalizeCreatePlayerAccountMutation() error = nil, want non-active initial state rejection")
	}
}

func TestNormalizeCreatePlayerAccountMutationRequiresDurableEventMetadata(t *testing.T) {
	valid := CreatePlayerAccountMutation{
		EventID:     "event-1",
		OccurredAt:  time.Date(2026, 5, 14, 1, 2, 3, 0, time.UTC),
		PlayerID:    "player-1",
		DisplayName: "Player One",
		RequestedBy: "maintainer",
	}

	tests := []struct {
		name   string
		mutate func(*CreatePlayerAccountMutation)
	}{
		{name: "event_id", mutate: func(m *CreatePlayerAccountMutation) { m.EventID = " " }},
		{name: "occurred_at", mutate: func(m *CreatePlayerAccountMutation) { m.OccurredAt = time.Time{} }},
		{name: "player_id", mutate: func(m *CreatePlayerAccountMutation) { m.PlayerID = " " }},
		{name: "display_name", mutate: func(m *CreatePlayerAccountMutation) { m.DisplayName = " " }},
		{name: "requested_by", mutate: func(m *CreatePlayerAccountMutation) { m.RequestedBy = " " }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mutation := valid
			tt.mutate(&mutation)
			if _, err := NormalizeCreatePlayerAccountMutation(mutation); err == nil {
				t.Fatal("NormalizeCreatePlayerAccountMutation() error = nil, want required field rejection")
			}
		})
	}
}

func TestAccountStateIsValid(t *testing.T) {
	for _, state := range []AccountState{AccountStateActive, AccountStateDisabled, AccountStateDeleted} {
		if !state.IsValid() {
			t.Fatalf("%q IsValid() = false, want true", state)
		}
	}
	if AccountState("suspended").IsValid() {
		t.Fatal(`AccountState("suspended").IsValid() = true, want false`)
	}
}

type recordingRepository struct{}

func (recordingRepository) CreatePlayerAccount(context.Context, CreatePlayerAccountMutation) (Account, error) {
	return Account{}, nil
}

func (recordingRepository) GetPlayerAccount(context.Context, string) (Account, error) {
	return Account{}, nil
}
