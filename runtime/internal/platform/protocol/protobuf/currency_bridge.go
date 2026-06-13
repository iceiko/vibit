package protobuf

import (
	"time"

	"github.com/iceiko/vibit/runtime/internal/app"
	appcurrency "github.com/iceiko/vibit/runtime/internal/app/currency"
	currencyv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/currency/v1"
	currencymodule "github.com/iceiko/vibit/runtime/internal/modules/currency"
	"google.golang.org/protobuf/proto"
)

func routeRequestWithCurrencyPayload(request app.RouteRequest) (app.RouteRequest, bool, error) {
	switch request.Route {
	case appcurrency.EnsurePlayerWalletRoute():
		payload, ok := request.Payload.(*currencyv1.EnsurePlayerWalletRequest)
		if !ok || payload == nil {
			return app.RouteRequest{}, true, payloadBridgeError(request.Route, "payload must be vibit.currency.v1.EnsurePlayerWalletRequest")
		}
		request.Payload = appcurrency.EnsurePlayerWalletRequest{}
		return request, true, nil

	case appcurrency.GetOwnWalletRoute():
		payload, ok := request.Payload.(*currencyv1.GetOwnWalletRequest)
		if !ok || payload == nil {
			return app.RouteRequest{}, true, payloadBridgeError(request.Route, "payload must be vibit.currency.v1.GetOwnWalletRequest")
		}
		request.Payload = appcurrency.GetOwnWalletRequest{}
		return request, true, nil

	case appcurrency.ListOwnWalletBalancesRoute():
		payload, ok := request.Payload.(*currencyv1.ListOwnWalletBalancesRequest)
		if !ok || payload == nil {
			return app.RouteRequest{}, true, payloadBridgeError(request.Route, "payload must be vibit.currency.v1.ListOwnWalletBalancesRequest")
		}
		request.Payload = appcurrency.ListOwnWalletBalancesRequest{
			Limit:             int(payload.GetLimit()),
			AfterCurrencyCode: currencymodule.CurrencyCode(payload.GetAfterCurrencyCode()),
		}
		return request, true, nil

	case appcurrency.GrantCurrencyRoute():
		payload, ok := request.Payload.(*currencyv1.GrantCurrencyRequest)
		if !ok || payload == nil {
			return app.RouteRequest{}, true, payloadBridgeError(request.Route, "payload must be vibit.currency.v1.GrantCurrencyRequest")
		}
		request.Payload = appcurrency.GrantCurrencyRequest{
			CurrencyCode:           currencymodule.CurrencyCode(payload.GetCurrencyCode()),
			Amount:                 currencymodule.CurrencyAmount(payload.GetAmount()),
			IdempotencyKey:         currencymodule.CurrencyWalletIdempotencyKey(payload.GetIdempotencyKey()),
			IdempotencyScope:       currencymodule.CurrencyWalletIdempotencyScope(payload.GetIdempotencyScope()),
			ReasonCode:             payload.GetReasonCode(),
			ExternalReference:      payload.GetExternalReference(),
			MetadataJSON:           []byte(payload.GetMetadataJson()),
			ExpectedWalletVersion:  currencyWalletVersionFromOptionalInt64(payload.ExpectedWalletVersion),
			ExpectedBalanceVersion: currencyBalanceVersionFromOptionalInt64(payload.ExpectedBalanceVersion),
		}
		return request, true, nil

	case appcurrency.SpendCurrencyRoute():
		payload, ok := request.Payload.(*currencyv1.SpendCurrencyRequest)
		if !ok || payload == nil {
			return app.RouteRequest{}, true, payloadBridgeError(request.Route, "payload must be vibit.currency.v1.SpendCurrencyRequest")
		}
		request.Payload = appcurrency.SpendCurrencyRequest{
			CurrencyCode:           currencymodule.CurrencyCode(payload.GetCurrencyCode()),
			Amount:                 currencymodule.CurrencyAmount(payload.GetAmount()),
			IdempotencyKey:         currencymodule.CurrencyWalletIdempotencyKey(payload.GetIdempotencyKey()),
			IdempotencyScope:       currencymodule.CurrencyWalletIdempotencyScope(payload.GetIdempotencyScope()),
			ReasonCode:             payload.GetReasonCode(),
			ExternalReference:      payload.GetExternalReference(),
			MetadataJSON:           []byte(payload.GetMetadataJson()),
			ExpectedWalletVersion:  currencyWalletVersionFromOptionalInt64(payload.ExpectedWalletVersion),
			ExpectedBalanceVersion: currencyBalanceVersionFromOptionalInt64(payload.ExpectedBalanceVersion),
		}
		return request, true, nil

	case appcurrency.ListOwnWalletTransactionsRoute():
		payload, ok := request.Payload.(*currencyv1.ListOwnWalletTransactionsRequest)
		if !ok || payload == nil {
			return app.RouteRequest{}, true, payloadBridgeError(request.Route, "payload must be vibit.currency.v1.ListOwnWalletTransactionsRequest")
		}
		afterTime, err := parseCurrencyOptionalTime(payload.GetAfterTransactionTime())
		if err != nil {
			return app.RouteRequest{}, true, payloadBridgeError(request.Route, "after_transaction_time must be RFC3339Nano UTC text")
		}
		request.Payload = appcurrency.ListOwnWalletTransactionsRequest{
			CurrencyCode:         currencymodule.CurrencyCode(payload.GetCurrencyCode()),
			Limit:                int(payload.GetLimit()),
			AfterTransactionID:   currencymodule.CurrencyWalletTransactionID(payload.GetAfterTransactionId()),
			AfterTransactionTime: afterTime,
		}
		return request, true, nil

	default:
		return request, false, nil
	}
}

func protoPayloadFromCurrencyRoute(route app.RouteKey, payload any) (proto.Message, bool, error) {
	switch route {
	case appcurrency.EnsurePlayerWalletRoute():
		result, ok := currencyResultPayload(payload)
		if !ok {
			return nil, true, payloadBridgeError(route, "payload must be currency.CurrencyWalletResult")
		}
		return &currencyv1.EnsurePlayerWalletResponse{
			Wallet: protoCurrencyWallet(result.Wallet),
			Status: string(result.Status),
		}, true, nil

	case appcurrency.GetOwnWalletRoute():
		result, ok := currencyResultPayload(payload)
		if !ok {
			return nil, true, payloadBridgeError(route, "payload must be currency.CurrencyWalletResult")
		}
		return &currencyv1.GetOwnWalletResponse{
			Wallet: protoCurrencyWallet(result.Wallet),
			Status: string(result.Status),
		}, true, nil

	case appcurrency.ListOwnWalletBalancesRoute():
		result, ok := currencyResultPayload(payload)
		if !ok {
			return nil, true, payloadBridgeError(route, "payload must be currency.CurrencyWalletResult")
		}
		balances := make([]*currencyv1.CurrencyWalletBalance, 0, len(result.Balances))
		for _, balance := range result.Balances {
			balances = append(balances, protoCurrencyWalletBalance(balance))
		}
		return &currencyv1.ListOwnWalletBalancesResponse{
			Balances:         balances,
			NextCurrencyCode: string(result.NextCurrencyCode),
			Status:           string(result.Status),
		}, true, nil

	case appcurrency.GrantCurrencyRoute():
		result, ok := currencyResultPayload(payload)
		if !ok {
			return nil, true, payloadBridgeError(route, "payload must be currency.CurrencyWalletResult")
		}
		return &currencyv1.GrantCurrencyResponse{
			Transaction: protoCurrencyWalletTransaction(result.Transaction),
			Status:      string(result.Status),
		}, true, nil

	case appcurrency.SpendCurrencyRoute():
		result, ok := currencyResultPayload(payload)
		if !ok {
			return nil, true, payloadBridgeError(route, "payload must be currency.CurrencyWalletResult")
		}
		return &currencyv1.SpendCurrencyResponse{
			Transaction: protoCurrencyWalletTransaction(result.Transaction),
			Status:      string(result.Status),
		}, true, nil

	case appcurrency.ListOwnWalletTransactionsRoute():
		result, ok := currencyResultPayload(payload)
		if !ok {
			return nil, true, payloadBridgeError(route, "payload must be currency.CurrencyWalletResult")
		}
		transactions := make([]*currencyv1.CurrencyWalletTransaction, 0, len(result.Transactions))
		for _, transaction := range result.Transactions {
			transactions = append(transactions, protoCurrencyWalletTransaction(transaction))
		}
		return &currencyv1.ListOwnWalletTransactionsResponse{
			Transactions:        transactions,
			NextTransactionId:   string(result.NextTransactionID),
			NextTransactionTime: formatCurrencyTime(result.NextTransactionCreateAt),
			Status:              string(result.Status),
		}, true, nil

	default:
		return nil, false, nil
	}
}

func currencyResultPayload(payload any) (appcurrency.CurrencyWalletResult, bool) {
	result, ok := payload.(appcurrency.CurrencyWalletResult)
	if !ok {
		if pointerResult, pointerOK := payload.(*appcurrency.CurrencyWalletResult); pointerOK && pointerResult != nil {
			result = *pointerResult
			ok = true
		}
	}
	return result, ok
}

func currencyWalletVersionFromOptionalInt64(value *int64) *currencymodule.CurrencyWalletVersion {
	if value == nil {
		return nil
	}
	version := currencymodule.CurrencyWalletVersion(*value)
	return &version
}

func currencyBalanceVersionFromOptionalInt64(value *int64) *currencymodule.CurrencyBalanceVersion {
	if value == nil {
		return nil
	}
	version := currencymodule.CurrencyBalanceVersion(*value)
	return &version
}

func parseCurrencyOptionalTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}, err
	}
	return parsed.UTC(), nil
}

func protoCurrencyWallet(wallet currencymodule.CurrencyWallet) *currencyv1.CurrencyWallet {
	if wallet.WalletID == "" {
		return nil
	}
	return &currencyv1.CurrencyWallet{
		WalletId:       string(wallet.WalletID),
		OwnerKind:      string(wallet.Owner.Kind),
		LifecycleState: string(wallet.LifecycleState),
		WalletVersion:  int64(wallet.WalletVersion),
		CreatedAt:      formatCurrencyTime(wallet.CreatedAt),
		UpdatedAt:      formatCurrencyTime(wallet.UpdatedAt),
		StateChangedAt: formatCurrencyTime(wallet.StateChangedAt),
	}
}

func protoCurrencyWalletBalance(balance currencymodule.CurrencyWalletBalance) *currencyv1.CurrencyWalletBalance {
	if balance.CurrencyCode == "" {
		return nil
	}
	return &currencyv1.CurrencyWalletBalance{
		CurrencyCode:   string(balance.CurrencyCode),
		BalanceAmount:  int64(balance.BalanceAmount),
		BalanceVersion: int64(balance.BalanceVersion),
		CreatedAt:      formatCurrencyTime(balance.CreatedAt),
		UpdatedAt:      formatCurrencyTime(balance.UpdatedAt),
	}
}

func protoCurrencyWalletTransaction(transaction currencymodule.CurrencyWalletTransaction) *currencyv1.CurrencyWalletTransaction {
	if transaction.TransactionID == "" {
		return nil
	}
	return &currencyv1.CurrencyWalletTransaction{
		TransactionId:     string(transaction.TransactionID),
		CurrencyCode:      string(transaction.CurrencyCode),
		TransactionKind:   string(transaction.TransactionKind),
		AmountDelta:       int64(transaction.AmountDelta),
		BalanceAfter:      int64(transaction.BalanceAfter),
		ActorKind:         string(transaction.Actor.Kind),
		ReasonCode:        transaction.ReasonCode,
		ExternalReference: transaction.ExternalReference,
		MetadataJson:      string(transaction.MetadataJSON),
		CreatedAt:         formatCurrencyTime(transaction.CreatedAt),
	}
}

func formatCurrencyTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339Nano)
}
