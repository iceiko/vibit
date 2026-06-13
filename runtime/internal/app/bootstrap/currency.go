package bootstrap

import (
	"context"
	"errors"
	"strings"

	"github.com/iceiko/vibit/runtime/internal/app"
	appcurrency "github.com/iceiko/vibit/runtime/internal/app/currency"
)

type CurrencyService interface {
	EnsurePlayerWallet(context.Context, appcurrency.EnsurePlayerWalletRequest) (appcurrency.CurrencyWalletResult, error)
	GetOwnWallet(context.Context, appcurrency.GetOwnWalletRequest) (appcurrency.CurrencyWalletResult, error)
	ListOwnWalletBalances(context.Context, appcurrency.ListOwnWalletBalancesRequest) (appcurrency.CurrencyWalletResult, error)
	GrantCurrency(context.Context, appcurrency.GrantCurrencyRequest) (appcurrency.CurrencyWalletResult, error)
	SpendCurrency(context.Context, appcurrency.SpendCurrencyRequest) (appcurrency.CurrencyWalletResult, error)
	ListOwnWalletTransactions(context.Context, appcurrency.ListOwnWalletTransactionsRequest) (appcurrency.CurrencyWalletResult, error)
}

type CurrencyGrantPolicy interface {
	GrantPolicyForRoute(context.Context, app.RouteRequest, appcurrency.GrantCurrencyRequest) (CurrencyGrantPolicyResult, error)
}

type CurrencyGrantPolicyResult struct {
	Allowed       bool
	SystemActorID string
}

type StaticCurrencyGrantPolicy struct {
	Allowed       bool
	SystemActorID string
}

func (p StaticCurrencyGrantPolicy) GrantPolicyForRoute(context.Context, app.RouteRequest, appcurrency.GrantCurrencyRequest) (CurrencyGrantPolicyResult, error) {
	return CurrencyGrantPolicyResult{
		Allowed:       p.Allowed,
		SystemActorID: strings.TrimSpace(p.SystemActorID),
	}, nil
}

type CurrencyRouteHandlers struct {
	Service     CurrencyService
	GrantPolicy CurrencyGrantPolicy
}

func (h CurrencyRouteHandlers) RegisterRoutes(dispatcher *app.Dispatcher) error {
	if dispatcher == nil {
		return errors.New("currency bootstrap: dispatcher is nil")
	}
	if err := dispatcher.Register(appcurrency.EnsurePlayerWalletRoute(), app.HandlerFunc(h.HandleEnsurePlayerWalletRoute)); err != nil {
		return err
	}
	if err := dispatcher.Register(appcurrency.GetOwnWalletRoute(), app.HandlerFunc(h.HandleGetOwnWalletRoute)); err != nil {
		return err
	}
	if err := dispatcher.Register(appcurrency.ListOwnWalletBalancesRoute(), app.HandlerFunc(h.HandleListOwnWalletBalancesRoute)); err != nil {
		return err
	}
	if err := dispatcher.Register(appcurrency.GrantCurrencyRoute(), app.HandlerFunc(h.HandleGrantCurrencyRoute)); err != nil {
		return err
	}
	if err := dispatcher.Register(appcurrency.SpendCurrencyRoute(), app.HandlerFunc(h.HandleSpendCurrencyRoute)); err != nil {
		return err
	}
	return dispatcher.Register(appcurrency.ListOwnWalletTransactionsRoute(), app.HandlerFunc(h.HandleListOwnWalletTransactionsRoute))
}

func (h CurrencyRouteHandlers) HandleEnsurePlayerWalletRoute(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	service := h.Service
	if service == nil {
		return currencyApplicationErrorResult(request, appcurrency.PublicErrorCurrencyWalletUnavailable)
	}
	payload, ok := currencyPayload[appcurrency.EnsurePlayerWalletRequest](request.Payload)
	if !ok {
		return currencyApplicationErrorResult(request, appcurrency.PublicErrorCurrencyWalletInvalidRequest)
	}
	payload.Identity = request.Identity
	result, err := service.EnsurePlayerWallet(ctx, payload)
	return currencyResultForRequest(request, result, err)
}

func (h CurrencyRouteHandlers) HandleGetOwnWalletRoute(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	service := h.Service
	if service == nil {
		return currencyApplicationErrorResult(request, appcurrency.PublicErrorCurrencyWalletUnavailable)
	}
	payload, ok := currencyPayload[appcurrency.GetOwnWalletRequest](request.Payload)
	if !ok {
		return currencyApplicationErrorResult(request, appcurrency.PublicErrorCurrencyWalletInvalidRequest)
	}
	payload.Identity = request.Identity
	result, err := service.GetOwnWallet(ctx, payload)
	return currencyResultForRequest(request, result, err)
}

func (h CurrencyRouteHandlers) HandleListOwnWalletBalancesRoute(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	service := h.Service
	if service == nil {
		return currencyApplicationErrorResult(request, appcurrency.PublicErrorCurrencyWalletUnavailable)
	}
	payload, ok := currencyPayload[appcurrency.ListOwnWalletBalancesRequest](request.Payload)
	if !ok {
		return currencyApplicationErrorResult(request, appcurrency.PublicErrorCurrencyWalletInvalidRequest)
	}
	payload.Identity = request.Identity
	result, err := service.ListOwnWalletBalances(ctx, payload)
	return currencyResultForRequest(request, result, err)
}

func (h CurrencyRouteHandlers) HandleGrantCurrencyRoute(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	service := h.Service
	if service == nil {
		return currencyApplicationErrorResult(request, appcurrency.PublicErrorCurrencyWalletUnavailable)
	}
	payload, ok := currencyPayload[appcurrency.GrantCurrencyRequest](request.Payload)
	if !ok {
		return currencyApplicationErrorResult(request, appcurrency.PublicErrorCurrencyWalletInvalidRequest)
	}
	grantPolicy := h.GrantPolicy
	if grantPolicy == nil {
		return currencyApplicationErrorResult(request, appcurrency.PublicErrorCurrencyWalletInvalidRequest)
	}
	policy, err := grantPolicy.GrantPolicyForRoute(ctx, request, payload)
	if err != nil || !policy.Allowed || strings.TrimSpace(policy.SystemActorID) == "" {
		return currencyApplicationErrorResult(request, appcurrency.PublicErrorCurrencyWalletInvalidRequest)
	}
	payload.Identity = request.Identity
	payload.SystemActorID = strings.TrimSpace(policy.SystemActorID)
	result, err := service.GrantCurrency(ctx, payload)
	return currencyResultForRequest(request, result, err)
}

func (h CurrencyRouteHandlers) HandleSpendCurrencyRoute(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	service := h.Service
	if service == nil {
		return currencyApplicationErrorResult(request, appcurrency.PublicErrorCurrencyWalletUnavailable)
	}
	payload, ok := currencyPayload[appcurrency.SpendCurrencyRequest](request.Payload)
	if !ok {
		return currencyApplicationErrorResult(request, appcurrency.PublicErrorCurrencyWalletInvalidRequest)
	}
	payload.Identity = request.Identity
	result, err := service.SpendCurrency(ctx, payload)
	return currencyResultForRequest(request, result, err)
}

func (h CurrencyRouteHandlers) HandleListOwnWalletTransactionsRoute(ctx context.Context, request app.RouteRequest) (app.ApplicationResult, error) {
	service := h.Service
	if service == nil {
		return currencyApplicationErrorResult(request, appcurrency.PublicErrorCurrencyWalletUnavailable)
	}
	payload, ok := currencyPayload[appcurrency.ListOwnWalletTransactionsRequest](request.Payload)
	if !ok {
		return currencyApplicationErrorResult(request, appcurrency.PublicErrorCurrencyWalletInvalidRequest)
	}
	payload.Identity = request.Identity
	result, err := service.ListOwnWalletTransactions(ctx, payload)
	return currencyResultForRequest(request, result, err)
}

func currencyPayload[T any](payload any) (T, bool) {
	value, ok := payload.(T)
	if !ok {
		if pointerValue, pointerOK := payload.(*T); pointerOK && pointerValue != nil {
			value = *pointerValue
			ok = true
		}
	}
	return value, ok
}

func currencyResultForRequest(request app.RouteRequest, currencyResult appcurrency.CurrencyWalletResult, err error) (app.ApplicationResult, error) {
	result := resultForRequest(request)
	if err != nil {
		appErr := currencyApplicationError(request.Route, currencyPublicErrorCode(currencyResult, err))
		result.Error = appErr
		return result, appErr
	}
	result.PayloadType = "currency.CurrencyWalletResult"
	result.Payload = currencyResult
	return result, nil
}

func currencyPublicErrorCode(result appcurrency.CurrencyWalletResult, err error) appcurrency.PublicErrorCode {
	if result.PublicErrorCode != "" {
		return result.PublicErrorCode
	}
	var serviceErr *appcurrency.ServiceError
	if errors.As(err, &serviceErr) && serviceErr.PublicCode != "" {
		return serviceErr.PublicCode
	}
	return appcurrency.PublicErrorCurrencyWalletUnavailable
}

func currencyApplicationErrorResult(request app.RouteRequest, publicCode appcurrency.PublicErrorCode) (app.ApplicationResult, error) {
	result := resultForRequest(request)
	appErr := currencyApplicationError(request.Route, publicCode)
	result.Error = appErr
	return result, appErr
}

func currencyApplicationError(route app.RouteKey, publicCode appcurrency.PublicErrorCode) *app.ApplicationError {
	code := app.ErrorCode(publicCode)
	if code == "" {
		code = app.ErrorCode(appcurrency.PublicErrorCurrencyWalletUnavailable)
	}
	return &app.ApplicationError{
		Code:    code,
		Message: currencyErrorMessage(publicCode),
		Route:   route,
	}
}

func currencyErrorMessage(code appcurrency.PublicErrorCode) string {
	switch code {
	case appcurrency.PublicErrorCurrencyWalletInvalidRequest:
		return "currency wallet request is invalid"
	case appcurrency.PublicErrorCurrencyWalletUnauthenticated:
		return "currency wallet request is unauthenticated"
	case appcurrency.PublicErrorCurrencyWalletNotFound:
		return "currency wallet was not found"
	case appcurrency.PublicErrorCurrencyWalletAlreadyExists:
		return "currency wallet already exists"
	case appcurrency.PublicErrorCurrencyWalletNotActive:
		return "currency wallet is not active"
	case appcurrency.PublicErrorCurrencyWalletInsufficientBalance:
		return "currency wallet balance is insufficient"
	case appcurrency.PublicErrorCurrencyWalletDuplicateTransaction:
		return "currency wallet transaction already exists"
	case appcurrency.PublicErrorCurrencyWalletVersionMismatch:
		return "currency wallet version mismatch"
	default:
		return "currency wallet service is unavailable"
	}
}
