package currency

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestRepositoryInterfaceIsStorageNeutral(t *testing.T) {
	var _ Repository = recordingRepository{}
}

func TestClosedVocabularies(t *testing.T) {
	if !CurrencyWalletOwnerKindPlayer.IsValid() {
		t.Fatalf("%q IsValid() = false, want true", CurrencyWalletOwnerKindPlayer)
	}
	if CurrencyWalletOwnerKind("guild").IsValid() {
		t.Fatal(`CurrencyWalletOwnerKind("guild").IsValid() = true, want false`)
	}

	for _, state := range []CurrencyWalletLifecycleState{
		CurrencyWalletLifecycleActive,
		CurrencyWalletLifecycleSuspended,
		CurrencyWalletLifecycleClosed,
	} {
		if !state.IsValid() {
			t.Fatalf("%q IsValid() = false, want true", state)
		}
	}
	if CurrencyWalletLifecycleState("archived").IsValid() {
		t.Fatal(`CurrencyWalletLifecycleState("archived").IsValid() = true, want false`)
	}

	for _, kind := range []CurrencyWalletTransactionKind{
		CurrencyWalletTransactionGrant,
		CurrencyWalletTransactionSpend,
	} {
		if !kind.IsValid() {
			t.Fatalf("%q IsValid() = false, want true", kind)
		}
	}
	if CurrencyWalletTransactionKind("transfer").IsValid() {
		t.Fatal(`CurrencyWalletTransactionKind("transfer").IsValid() = true, want false`)
	}

	for _, actorKind := range []CurrencyWalletActorKind{
		CurrencyWalletActorSystem,
		CurrencyWalletActorPlayer,
		CurrencyWalletActorOperation,
	} {
		if !actorKind.IsValid() {
			t.Fatalf("%q IsValid() = false, want true", actorKind)
		}
	}
	if CurrencyWalletActorKind("session").IsValid() {
		t.Fatal(`CurrencyWalletActorKind("session").IsValid() = true, want false`)
	}
}

func TestNormalizeCurrencyWalletRecordTrimsAndNormalizesTimes(t *testing.T) {
	now := time.Date(2026, 5, 26, 9, 0, 0, 0, time.FixedZone("test", 8*60*60))
	suspendedAt := now.Add(time.Minute)

	record, err := NormalizeCurrencyWalletRecord(CurrencyWallet{
		WalletID:       " wallet-1 ",
		Owner:          CurrencyWalletOwner{Kind: CurrencyWalletOwnerKind(" player "), ID: " player-1 "},
		LifecycleState: CurrencyWalletLifecycleSuspended,
		WalletVersion:  CurrencyWalletVersion(4),
		CreatedAt:      now,
		UpdatedAt:      now.Add(2 * time.Minute),
		StateChangedAt: now.Add(time.Minute),
		SuspendedAt:    &suspendedAt,
	})
	if err != nil {
		t.Fatalf("NormalizeCurrencyWalletRecord() error = %v, want nil", err)
	}

	if record.WalletID != CurrencyWalletID("wallet-1") ||
		record.Owner.Kind != CurrencyWalletOwnerKindPlayer ||
		record.Owner.ID != "player-1" ||
		record.LifecycleState != CurrencyWalletLifecycleSuspended ||
		record.WalletVersion != CurrencyWalletVersion(4) {
		t.Fatalf("record = %#v, want trimmed suspended wallet", record)
	}
	if record.CreatedAt.Location() != time.UTC ||
		record.UpdatedAt.Location() != time.UTC ||
		record.StateChangedAt.Location() != time.UTC ||
		record.SuspendedAt == nil ||
		record.SuspendedAt.Location() != time.UTC {
		t.Fatalf("record times = %#v, want UTC-normalized times", record)
	}
}

func TestNormalizeCurrencyWalletRecordRejectsInvalidShape(t *testing.T) {
	valid := validCurrencyWalletRecord()
	tests := []struct {
		name   string
		mutate func(*CurrencyWallet)
	}{
		{name: "wallet_id", mutate: func(r *CurrencyWallet) { r.WalletID = " " }},
		{name: "owner_kind", mutate: func(r *CurrencyWallet) { r.Owner.Kind = "guild" }},
		{name: "owner_id", mutate: func(r *CurrencyWallet) { r.Owner.ID = " " }},
		{name: "lifecycle_state", mutate: func(r *CurrencyWallet) { r.LifecycleState = "archived" }},
		{name: "wallet_version", mutate: func(r *CurrencyWallet) { r.WalletVersion = 0 }},
		{name: "created_at", mutate: func(r *CurrencyWallet) { r.CreatedAt = time.Time{} }},
		{name: "updated_at", mutate: func(r *CurrencyWallet) { r.UpdatedAt = time.Time{} }},
		{name: "state_changed_at", mutate: func(r *CurrencyWallet) { r.StateChangedAt = time.Time{} }},
		{name: "updated_before_created", mutate: func(r *CurrencyWallet) { r.UpdatedAt = r.CreatedAt.Add(-time.Second) }},
		{name: "state_changed_before_created", mutate: func(r *CurrencyWallet) { r.StateChangedAt = r.CreatedAt.Add(-time.Second) }},
		{name: "suspended_at_without_suspended_state", mutate: func(r *CurrencyWallet) { at := r.CreatedAt; r.SuspendedAt = &at }},
		{name: "closed_at_without_closed_state", mutate: func(r *CurrencyWallet) { at := r.CreatedAt; r.ClosedAt = &at }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := valid
			tt.mutate(&record)
			if _, err := NormalizeCurrencyWalletRecord(record); err == nil {
				t.Fatal("NormalizeCurrencyWalletRecord() error = nil, want rejection")
			}
		})
	}
}

func TestNormalizeCurrencyWalletBalanceRecord(t *testing.T) {
	now := time.Date(2026, 5, 26, 10, 0, 0, 0, time.FixedZone("test", 8*60*60))

	record, err := NormalizeCurrencyWalletBalanceRecord(CurrencyWalletBalance{
		WalletID:       " wallet-1 ",
		CurrencyCode:   " Gems ",
		BalanceAmount:  CurrencyAmount(50),
		BalanceVersion: CurrencyBalanceVersion(2),
		CreatedAt:      now,
		UpdatedAt:      now.Add(time.Minute),
	})
	if err != nil {
		t.Fatalf("NormalizeCurrencyWalletBalanceRecord() error = %v, want nil", err)
	}
	if record.WalletID != CurrencyWalletID("wallet-1") ||
		record.CurrencyCode != CurrencyCode("Gems") ||
		record.BalanceAmount != CurrencyAmount(50) ||
		record.BalanceVersion != CurrencyBalanceVersion(2) ||
		record.CreatedAt.Location() != time.UTC {
		t.Fatalf("balance = %#v, want trimmed UTC record", record)
	}

	valid := validCurrencyWalletBalanceRecord()
	tests := []struct {
		name   string
		mutate func(*CurrencyWalletBalance)
	}{
		{name: "wallet_id", mutate: func(r *CurrencyWalletBalance) { r.WalletID = " " }},
		{name: "currency_code", mutate: func(r *CurrencyWalletBalance) { r.CurrencyCode = " " }},
		{name: "negative_balance", mutate: func(r *CurrencyWalletBalance) { r.BalanceAmount = -1 }},
		{name: "balance_version", mutate: func(r *CurrencyWalletBalance) { r.BalanceVersion = 0 }},
		{name: "updated_before_created", mutate: func(r *CurrencyWalletBalance) { r.UpdatedAt = r.CreatedAt.Add(-time.Second) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := valid
			tt.mutate(&record)
			if _, err := NormalizeCurrencyWalletBalanceRecord(record); err == nil {
				t.Fatal("NormalizeCurrencyWalletBalanceRecord() error = nil, want rejection")
			}
		})
	}
}

func TestNormalizeCurrencyWalletTransactionRecordTrimsIdempotencyAndActor(t *testing.T) {
	now := time.Date(2026, 5, 26, 11, 0, 0, 0, time.FixedZone("test", 8*60*60))
	metadata := []byte(`{"source":"daily_reward"}`)

	record, err := NormalizeCurrencyWalletTransactionRecord(CurrencyWalletTransaction{
		TransactionID:     " tx-1 ",
		WalletID:          " wallet-1 ",
		CurrencyCode:      " Gems ",
		TransactionKind:   CurrencyWalletTransactionGrant,
		AmountDelta:       CurrencyAmount(25),
		BalanceAfter:      CurrencyAmount(75),
		IdempotencyKey:    " grant-123 ",
		IdempotencyScope:  " daily-login ",
		Actor:             CurrencyWalletActor{Kind: CurrencyWalletActorKind(" system "), ID: " economy "},
		ReasonCode:        " daily_reward ",
		ExternalReference: " reward-run-1 ",
		MetadataJSON:      metadata,
		CreatedAt:         now,
	})
	if err != nil {
		t.Fatalf("NormalizeCurrencyWalletTransactionRecord() error = %v, want nil", err)
	}
	if record.TransactionID != CurrencyWalletTransactionID("tx-1") ||
		record.WalletID != CurrencyWalletID("wallet-1") ||
		record.CurrencyCode != CurrencyCode("Gems") ||
		record.IdempotencyKey != CurrencyWalletIdempotencyKey("grant-123") ||
		record.IdempotencyScope != CurrencyWalletIdempotencyScope("daily-login") ||
		record.Actor.Kind != CurrencyWalletActorSystem ||
		record.Actor.ID != "economy" ||
		record.ReasonCode != "daily_reward" ||
		record.ExternalReference != "reward-run-1" ||
		record.CreatedAt.Location() != time.UTC {
		t.Fatalf("transaction = %#v, want trimmed normalized transaction", record)
	}
	metadata[0] = '['
	if string(record.MetadataJSON) != `{"source":"daily_reward"}` {
		t.Fatalf("transaction metadata aliases caller bytes, got %q", string(record.MetadataJSON))
	}
}

func TestNormalizeCurrencyWalletTransactionRecordRejectsInvalidShape(t *testing.T) {
	valid := validCurrencyWalletTransactionRecord()
	tests := []struct {
		name   string
		mutate func(*CurrencyWalletTransaction)
	}{
		{name: "transaction_id", mutate: func(r *CurrencyWalletTransaction) { r.TransactionID = " " }},
		{name: "wallet_id", mutate: func(r *CurrencyWalletTransaction) { r.WalletID = " " }},
		{name: "currency_code", mutate: func(r *CurrencyWalletTransaction) { r.CurrencyCode = " " }},
		{name: "transaction_kind", mutate: func(r *CurrencyWalletTransaction) { r.TransactionKind = "refund" }},
		{name: "zero_delta", mutate: func(r *CurrencyWalletTransaction) { r.AmountDelta = 0 }},
		{name: "spend_positive_delta", mutate: func(r *CurrencyWalletTransaction) {
			r.TransactionKind = CurrencyWalletTransactionSpend
			r.AmountDelta = 10
		}},
		{name: "grant_negative_delta", mutate: func(r *CurrencyWalletTransaction) {
			r.TransactionKind = CurrencyWalletTransactionGrant
			r.AmountDelta = -10
		}},
		{name: "negative_balance_after", mutate: func(r *CurrencyWalletTransaction) { r.BalanceAfter = -1 }},
		{name: "idempotency_key", mutate: func(r *CurrencyWalletTransaction) { r.IdempotencyKey = " " }},
		{name: "idempotency_scope", mutate: func(r *CurrencyWalletTransaction) { r.IdempotencyScope = " " }},
		{name: "actor_kind", mutate: func(r *CurrencyWalletTransaction) { r.Actor.Kind = "session" }},
		{name: "actor_id_for_player", mutate: func(r *CurrencyWalletTransaction) { r.Actor.Kind = CurrencyWalletActorPlayer; r.Actor.ID = " " }},
		{name: "metadata_not_object", mutate: func(r *CurrencyWalletTransaction) { r.MetadataJSON = []byte(`[]`) }},
		{name: "created_at", mutate: func(r *CurrencyWalletTransaction) { r.CreatedAt = time.Time{} }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := valid
			tt.mutate(&record)
			if _, err := NormalizeCurrencyWalletTransactionRecord(record); err == nil {
				t.Fatal("NormalizeCurrencyWalletTransactionRecord() error = nil, want rejection")
			}
		})
	}
}

func TestNormalizeRepositoryInputs(t *testing.T) {
	create, err := NormalizeCreateCurrencyWalletInput(CreateCurrencyWalletInput{
		WalletID:       " wallet-1 ",
		Owner:          CurrencyWalletOwner{Kind: " player ", ID: " player-1 "},
		InitialState:   "",
		InitialVersion: 0,
		RequestedBy:    " currency_service ",
	})
	if err != nil {
		t.Fatalf("NormalizeCreateCurrencyWalletInput() error = %v, want nil", err)
	}
	if create.WalletID != CurrencyWalletID("wallet-1") ||
		create.Owner.Kind != CurrencyWalletOwnerKindPlayer ||
		create.Owner.ID != "player-1" ||
		create.InitialState != CurrencyWalletLifecycleActive ||
		create.InitialVersion != InitialCurrencyWalletVersion ||
		create.RequestedBy != "currency_service" {
		t.Fatalf("create = %#v, want default active create input", create)
	}

	get, err := NormalizeGetCurrencyWalletInput(GetCurrencyWalletInput{WalletID: " wallet-1 "})
	if err != nil {
		t.Fatalf("NormalizeGetCurrencyWalletInput() error = %v, want nil", err)
	}
	if get.WalletID != CurrencyWalletID("wallet-1") {
		t.Fatalf("get = %#v, want trimmed wallet id", get)
	}

	owner, err := NormalizeGetCurrencyWalletForOwnerInput(GetCurrencyWalletForOwnerInput{
		Owner: CurrencyWalletOwner{Kind: " player ", ID: " player-1 "},
	})
	if err != nil {
		t.Fatalf("NormalizeGetCurrencyWalletForOwnerInput() error = %v, want nil", err)
	}
	if owner.Owner.ID != "player-1" {
		t.Fatalf("owner = %#v, want trimmed owner", owner)
	}

	listBalances, err := NormalizeListCurrencyWalletBalancesInput(ListCurrencyWalletBalancesInput{
		WalletID:          " wallet-1 ",
		Limit:             0,
		AfterCurrencyCode: " Gems ",
	})
	if err != nil {
		t.Fatalf("NormalizeListCurrencyWalletBalancesInput() error = %v, want nil", err)
	}
	if listBalances.WalletID != CurrencyWalletID("wallet-1") ||
		listBalances.Limit != DefaultListCurrencyWalletBalancesLimit ||
		listBalances.AfterCurrencyCode != CurrencyCode("Gems") {
		t.Fatalf("list balances = %#v, want defaults and trimmed cursor", listBalances)
	}

	listTransactions, err := NormalizeListCurrencyWalletTransactionsInput(ListCurrencyWalletTransactionsInput{
		WalletID:             " wallet-1 ",
		CurrencyCode:         " Gems ",
		Limit:                0,
		AfterTransactionID:   " tx-1 ",
		AfterTransactionTime: validTime().Add(time.Minute),
	})
	if err != nil {
		t.Fatalf("NormalizeListCurrencyWalletTransactionsInput() error = %v, want nil", err)
	}
	if listTransactions.WalletID != CurrencyWalletID("wallet-1") ||
		listTransactions.CurrencyCode != CurrencyCode("Gems") ||
		listTransactions.Limit != DefaultListCurrencyWalletTransactionsLimit ||
		listTransactions.AfterTransactionID != CurrencyWalletTransactionID("tx-1") ||
		listTransactions.AfterTransactionTime.Location() != time.UTC {
		t.Fatalf("list transactions = %#v, want defaults and trimmed filters", listTransactions)
	}
}

func TestNormalizeGrantSpendInputsCopyMetadataAndValidateIdempotency(t *testing.T) {
	metadata := []byte(`{"campaign":"launch"}`)
	grant, err := NormalizeRecordCurrencyGrantInput(RecordCurrencyGrantInput{
		WalletID:          " wallet-1 ",
		TransactionID:     " tx-1 ",
		CurrencyCode:      " Gems ",
		Amount:            CurrencyAmount(25),
		IdempotencyKey:    " grant-1 ",
		IdempotencyScope:  " reward ",
		Actor:             CurrencyWalletActor{Kind: " system ", ID: " economy "},
		ReasonCode:        " launch_reward ",
		ExternalReference: " reward-run-1 ",
		MetadataJSON:      metadata,
	})
	if err != nil {
		t.Fatalf("NormalizeRecordCurrencyGrantInput() error = %v, want nil", err)
	}
	if grant.WalletID != CurrencyWalletID("wallet-1") ||
		grant.TransactionID != CurrencyWalletTransactionID("tx-1") ||
		grant.CurrencyCode != CurrencyCode("Gems") ||
		grant.Amount != CurrencyAmount(25) ||
		grant.IdempotencyKey != CurrencyWalletIdempotencyKey("grant-1") ||
		grant.IdempotencyScope != CurrencyWalletIdempotencyScope("reward") ||
		grant.Actor.Kind != CurrencyWalletActorSystem ||
		grant.Actor.ID != "economy" {
		t.Fatalf("grant = %#v, want trimmed grant input", grant)
	}
	metadata[0] = '['
	if string(grant.MetadataJSON) != `{"campaign":"launch"}` {
		t.Fatalf("grant metadata aliases caller bytes, got %q", string(grant.MetadataJSON))
	}

	expected := CurrencyBalanceVersion(3)
	spend, err := NormalizeRecordCurrencySpendInput(RecordCurrencySpendInput{
		WalletID:               " wallet-1 ",
		TransactionID:          " tx-2 ",
		CurrencyCode:           " Gems ",
		Amount:                 CurrencyAmount(5),
		IdempotencyKey:         " spend-1 ",
		IdempotencyScope:       " purchase ",
		Actor:                  CurrencyWalletActor{Kind: " player ", ID: " player-1 "},
		ReasonCode:             " shop_purchase ",
		ExpectedBalanceVersion: &expected,
	})
	if err != nil {
		t.Fatalf("NormalizeRecordCurrencySpendInput() error = %v, want nil", err)
	}
	if spend.Amount != CurrencyAmount(5) ||
		spend.ExpectedBalanceVersion == nil ||
		*spend.ExpectedBalanceVersion != CurrencyBalanceVersion(3) ||
		spend.Actor.Kind != CurrencyWalletActorPlayer ||
		spend.Actor.ID != "player-1" {
		t.Fatalf("spend = %#v, want positive spend and copied expected version", spend)
	}
	expected = 8
	if *spend.ExpectedBalanceVersion != CurrencyBalanceVersion(3) {
		t.Fatalf("spend expected version aliases caller pointer, got %d", *spend.ExpectedBalanceVersion)
	}

	for name, err := range map[string]error{
		"grant_zero_amount":        mustErrNormalizeGrant(RecordCurrencyGrantInput{WalletID: "wallet-1", TransactionID: "tx-1", CurrencyCode: "Gems", Amount: 0, IdempotencyKey: "key", IdempotencyScope: "scope", Actor: CurrencyWalletActor{Kind: CurrencyWalletActorSystem}}),
		"spend_zero_amount":        mustErrNormalizeSpend(RecordCurrencySpendInput{WalletID: "wallet-1", TransactionID: "tx-1", CurrencyCode: "Gems", Amount: 0, IdempotencyKey: "key", IdempotencyScope: "scope", Actor: CurrencyWalletActor{Kind: CurrencyWalletActorPlayer, ID: "player-1"}}),
		"grant_missing_idem_key":   mustErrNormalizeGrant(RecordCurrencyGrantInput{WalletID: "wallet-1", TransactionID: "tx-1", CurrencyCode: "Gems", Amount: 1, IdempotencyScope: "scope", Actor: CurrencyWalletActor{Kind: CurrencyWalletActorSystem}}),
		"spend_bad_expected_ver":   mustErrNormalizeSpend(RecordCurrencySpendInput{WalletID: "wallet-1", TransactionID: "tx-1", CurrencyCode: "Gems", Amount: 1, IdempotencyKey: "key", IdempotencyScope: "scope", Actor: CurrencyWalletActor{Kind: CurrencyWalletActorPlayer, ID: "player-1"}, ExpectedBalanceVersion: balanceVersionPtr(0)}),
		"metadata_must_be_object":  mustErrNormalizeGrant(RecordCurrencyGrantInput{WalletID: "wallet-1", TransactionID: "tx-1", CurrencyCode: "Gems", Amount: 1, IdempotencyKey: "key", IdempotencyScope: "scope", Actor: CurrencyWalletActor{Kind: CurrencyWalletActorSystem}, MetadataJSON: []byte(`[]`)}),
		"actor_player_requires_id": mustErrNormalizeSpend(RecordCurrencySpendInput{WalletID: "wallet-1", TransactionID: "tx-1", CurrencyCode: "Gems", Amount: 1, IdempotencyKey: "key", IdempotencyScope: "scope", Actor: CurrencyWalletActor{Kind: CurrencyWalletActorPlayer}}),
	} {
		if err == nil {
			t.Fatalf("%s error = nil, want rejection", name)
		}
	}
}

func TestNormalizeListResultsCopyRecords(t *testing.T) {
	balances := []CurrencyWalletBalance{validCurrencyWalletBalanceRecord()}
	balanceResult, err := NormalizeListCurrencyWalletBalancesResult(ListCurrencyWalletBalancesResult{
		Balances:         balances,
		NextCurrencyCode: " Gold ",
	})
	if err != nil {
		t.Fatalf("NormalizeListCurrencyWalletBalancesResult() error = %v, want nil", err)
	}
	if len(balanceResult.Balances) != 1 || balanceResult.NextCurrencyCode != CurrencyCode("Gold") {
		t.Fatalf("balance result = %#v, want normalized records and cursor", balanceResult)
	}
	balances[0].WalletID = "mutated"
	if balanceResult.Balances[0].WalletID != CurrencyWalletID("wallet-1") {
		t.Fatalf("balance result aliases caller slice, got %#v", balanceResult.Balances[0])
	}

	transactions := []CurrencyWalletTransaction{validCurrencyWalletTransactionRecord()}
	transactionResult, err := NormalizeListCurrencyWalletTransactionsResult(ListCurrencyWalletTransactionsResult{
		Transactions:            transactions,
		NextTransactionID:       " tx-2 ",
		NextTransactionCreateAt: validTime().Add(time.Minute),
	})
	if err != nil {
		t.Fatalf("NormalizeListCurrencyWalletTransactionsResult() error = %v, want nil", err)
	}
	if len(transactionResult.Transactions) != 1 ||
		transactionResult.NextTransactionID != CurrencyWalletTransactionID("tx-2") ||
		transactionResult.NextTransactionCreateAt.Location() != time.UTC {
		t.Fatalf("transaction result = %#v, want normalized records and cursor", transactionResult)
	}
	transactions[0].TransactionID = "mutated"
	if transactionResult.Transactions[0].TransactionID != CurrencyWalletTransactionID("tx-1") {
		t.Fatalf("transaction result aliases caller slice, got %#v", transactionResult.Transactions[0])
	}
}

func TestConflictAndRepositoryErrorsAreTypedAndRedacted(t *testing.T) {
	for _, class := range []CurrencyWalletConflictClass{
		CurrencyWalletConflictWalletNotFound,
		CurrencyWalletConflictWalletAlreadyExists,
		CurrencyWalletConflictWalletOwnerMismatch,
		CurrencyWalletConflictWalletNotActive,
		CurrencyWalletConflictBalanceNotFound,
		CurrencyWalletConflictInvalidCurrencyCode,
		CurrencyWalletConflictInvalidCurrencyAmount,
		CurrencyWalletConflictInsufficientBalance,
		CurrencyWalletConflictDuplicateTransaction,
		CurrencyWalletConflictConflictingDuplicateTransaction,
		CurrencyWalletConflictInvalidIdempotencyKey,
		CurrencyWalletConflictInvalidTransactionKind,
		CurrencyWalletConflictVersionMismatch,
		CurrencyWalletConflictStaleWalletVersion,
		CurrencyWalletConflictStaleBalanceVersion,
		CurrencyWalletConflictStorageUnavailable,
	} {
		if !class.IsValid() {
			t.Fatalf("%q IsValid() = false, want true", class)
		}
	}
	if CurrencyWalletConflictClass("raw_wallet_visible").IsValid() {
		t.Fatal(`CurrencyWalletConflictClass("raw_wallet_visible").IsValid() = true, want false`)
	}

	conflict := CurrencyWalletConflict{
		Class:          CurrencyWalletConflictInsufficientBalance,
		Retryable:      false,
		RedactedReason: "insufficient_balance",
	}
	if strings.Contains(conflict.Error(), "wallet-1") || strings.Contains(conflict.Error(), "player-1") || strings.Contains(conflict.Error(), "grant-1") {
		t.Fatalf("conflict error leaks private wallet material: %q", conflict.Error())
	}
	if !errors.Is(conflict, CurrencyWalletConflict{Class: CurrencyWalletConflictInsufficientBalance}) {
		t.Fatalf("errors.Is(conflict, same class) = false, want true")
	}

	cause := errors.New("driver detail: wallet-1 player-1 grant-1")
	repoErr := &CurrencyWalletRepositoryError{
		Kind:           ErrCurrencyWalletConflict,
		Conflict:       conflict,
		Operation:      "record_grant",
		RedactedReason: "insufficient_balance",
		Err:            cause,
	}
	if !errors.Is(repoErr, ErrCurrencyWalletConflict) ||
		!errors.Is(repoErr, CurrencyWalletConflict{Class: CurrencyWalletConflictInsufficientBalance}) ||
		!errors.Is(repoErr, cause) {
		t.Fatalf("repository error wrapping did not expose typed causes")
	}
	if strings.Contains(repoErr.Error(), "wallet-1") ||
		strings.Contains(repoErr.Error(), "player-1") ||
		strings.Contains(repoErr.Error(), "grant-1") ||
		strings.Contains(repoErr.Error(), "driver detail") {
		t.Fatalf("repository error leaks private material: %q", repoErr.Error())
	}
}

func TestStorageNeutralTypesAvoidDatabaseProtocolAndSessionCoupling(t *testing.T) {
	for _, typ := range []reflect.Type{
		reflect.TypeOf(CurrencyWallet{}),
		reflect.TypeOf(CurrencyWalletBalance{}),
		reflect.TypeOf(CurrencyWalletTransaction{}),
		reflect.TypeOf(CreateCurrencyWalletInput{}),
		reflect.TypeOf(RecordCurrencyGrantInput{}),
		reflect.TypeOf(RecordCurrencySpendInput{}),
		reflect.TypeOf(ListCurrencyWalletTransactionsInput{}),
	} {
		requireStorageNeutralFields(t, typ)
	}
}

type recordingRepository struct{}

func (recordingRepository) CreateCurrencyWallet(context.Context, CreateCurrencyWalletInput) (CurrencyWallet, error) {
	return CurrencyWallet{}, nil
}

func (recordingRepository) GetCurrencyWallet(context.Context, GetCurrencyWalletInput) (CurrencyWallet, error) {
	return CurrencyWallet{}, nil
}

func (recordingRepository) GetCurrencyWalletForOwner(context.Context, GetCurrencyWalletForOwnerInput) (CurrencyWallet, error) {
	return CurrencyWallet{}, nil
}

func (recordingRepository) ListCurrencyWalletBalances(context.Context, ListCurrencyWalletBalancesInput) (ListCurrencyWalletBalancesResult, error) {
	return ListCurrencyWalletBalancesResult{}, nil
}

func (recordingRepository) RecordCurrencyGrant(context.Context, RecordCurrencyGrantInput) (CurrencyWalletTransaction, error) {
	return CurrencyWalletTransaction{}, nil
}

func (recordingRepository) RecordCurrencySpend(context.Context, RecordCurrencySpendInput) (CurrencyWalletTransaction, error) {
	return CurrencyWalletTransaction{}, nil
}

func (recordingRepository) ListCurrencyWalletTransactions(context.Context, ListCurrencyWalletTransactionsInput) (ListCurrencyWalletTransactionsResult, error) {
	return ListCurrencyWalletTransactionsResult{}, nil
}

func validCurrencyWalletRecord() CurrencyWallet {
	now := validTime()
	return CurrencyWallet{
		WalletID:       CurrencyWalletID("wallet-1"),
		Owner:          CurrencyWalletOwner{Kind: CurrencyWalletOwnerKindPlayer, ID: "player-1"},
		LifecycleState: CurrencyWalletLifecycleActive,
		WalletVersion:  CurrencyWalletVersion(1),
		CreatedAt:      now,
		UpdatedAt:      now,
		StateChangedAt: now,
	}
}

func validCurrencyWalletBalanceRecord() CurrencyWalletBalance {
	now := validTime()
	return CurrencyWalletBalance{
		WalletID:       CurrencyWalletID("wallet-1"),
		CurrencyCode:   CurrencyCode("Gems"),
		BalanceAmount:  CurrencyAmount(25),
		BalanceVersion: CurrencyBalanceVersion(1),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func validCurrencyWalletTransactionRecord() CurrencyWalletTransaction {
	return CurrencyWalletTransaction{
		TransactionID:    CurrencyWalletTransactionID("tx-1"),
		WalletID:         CurrencyWalletID("wallet-1"),
		CurrencyCode:     CurrencyCode("Gems"),
		TransactionKind:  CurrencyWalletTransactionGrant,
		AmountDelta:      CurrencyAmount(25),
		BalanceAfter:     CurrencyAmount(50),
		IdempotencyKey:   CurrencyWalletIdempotencyKey("idem-1"),
		IdempotencyScope: CurrencyWalletIdempotencyScope("reward"),
		Actor:            CurrencyWalletActor{Kind: CurrencyWalletActorSystem},
		ReasonCode:       "test",
		CreatedAt:        validTime(),
	}
}

func validTime() time.Time {
	return time.Date(2026, 5, 26, 9, 0, 0, 0, time.UTC)
}

func mustErrNormalizeGrant(input RecordCurrencyGrantInput) error {
	_, err := NormalizeRecordCurrencyGrantInput(input)
	return err
}

func mustErrNormalizeSpend(input RecordCurrencySpendInput) error {
	_, err := NormalizeRecordCurrencySpendInput(input)
	return err
}

func balanceVersionPtr(version CurrencyBalanceVersion) *CurrencyBalanceVersion {
	return &version
}

func requireStorageNeutralFields(t *testing.T, typ reflect.Type) {
	t.Helper()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldName := strings.ToLower(field.Name)
		fieldType := strings.ToLower(field.Type.String())
		for _, forbidden := range []string{
			"sql",
			"pgx",
			"postgres",
			"proto",
			"websocket",
			"session",
			"token",
			"cookie",
			"authorization",
			"remoteaddr",
			"nakama",
			"pitaya",
		} {
			if strings.Contains(fieldName, forbidden) || strings.Contains(fieldType, forbidden) {
				t.Fatalf("%s.%s has forbidden storage/protocol/session coupling via %q", typ.Name(), field.Name, forbidden)
			}
		}
	}
}
