package app

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

type ErrorCode string

const (
	ErrorCodeInvalidRoute           ErrorCode = "INVALID_ROUTE"
	ErrorCodeUnsupportedMessageKind ErrorCode = "UNSUPPORTED_MESSAGE_KIND"
	ErrorCodeRouteAlreadyRegistered ErrorCode = "ROUTE_ALREADY_REGISTERED"
	ErrorCodeRouteNotFound          ErrorCode = "ROUTE_NOT_FOUND"
	ErrorCodeNilHandler             ErrorCode = "NIL_HANDLER"
	ErrorCodeSessionInvalid         ErrorCode = "SESSION_INVALID"
)

type ApplicationError struct {
	Code    ErrorCode
	Message string
	Route   RouteKey
}

func (e *ApplicationError) Error() string {
	if e == nil {
		return ""
	}

	route := RenderRouteKey(e.Route)
	if route == "" {
		route = "<invalid>"
	}

	return fmt.Sprintf("%s: %s: %s", e.Code, route, e.Message)
}

type ApplicationEvent struct {
	Route       RouteKey
	PayloadType string
	Payload     any
}

type ApplicationResult struct {
	RequestID   string
	Route       RouteKey
	Target      Target
	Session     Session
	Identity    RequestIdentity
	PayloadType string
	Payload     any
	Events      []ApplicationEvent
	Error       *ApplicationError
}

type Handler interface {
	HandleRoute(context.Context, RouteRequest) (ApplicationResult, error)
}

type HandlerFunc func(context.Context, RouteRequest) (ApplicationResult, error)

func (f HandlerFunc) HandleRoute(ctx context.Context, request RouteRequest) (ApplicationResult, error) {
	if f == nil {
		return resultForRequest(request), &ApplicationError{
			Code:    ErrorCodeNilHandler,
			Message: "handler function is nil",
			Route:   request.Route,
		}
	}
	return f(ctx, request)
}

type Dispatcher struct {
	routes map[RouteKey]Handler
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		routes: make(map[RouteKey]Handler),
	}
}

func (d *Dispatcher) Register(route RouteKey, handler Handler) error {
	normalizedRoute, err := normalizeDispatchRoute(route)
	if err != nil {
		return err
	}
	if isNilHandler(handler) {
		return &ApplicationError{
			Code:    ErrorCodeNilHandler,
			Message: "handler is nil",
			Route:   normalizedRoute,
		}
	}

	if d.routes == nil {
		d.routes = make(map[RouteKey]Handler)
	}
	if _, exists := d.routes[normalizedRoute]; exists {
		return &ApplicationError{
			Code:    ErrorCodeRouteAlreadyRegistered,
			Message: "route is already registered",
			Route:   normalizedRoute,
		}
	}

	d.routes[normalizedRoute] = handler
	return nil
}

func (d *Dispatcher) Dispatch(ctx context.Context, request RouteRequest) (ApplicationResult, error) {
	if request.Identity.Status == "" {
		request.Identity = MetadataOnlyIdentityFromSession(request.Session)
	}
	normalizedRoute, err := normalizeDispatchRoute(request.Route)
	if err != nil {
		result := resultForRequest(request)
		if appErr, ok := err.(*ApplicationError); ok {
			result.Error = appErr
		}
		return result, err
	}

	request.Route = normalizedRoute
	result := resultForRequest(request)

	handler, exists := d.routes[normalizedRoute]
	if !exists {
		appErr := &ApplicationError{
			Code:    ErrorCodeRouteNotFound,
			Message: "route is not registered",
			Route:   normalizedRoute,
		}
		result.Error = appErr
		return result, appErr
	}

	handlerResult, err := handler.HandleRoute(ctx, request)
	handlerResult.RequestID = request.RequestID
	handlerResult.Route = normalizedRoute
	handlerResult.Target = request.Target
	handlerResult.Session = request.Session
	handlerResult.Identity = request.Identity
	if err != nil {
		if appErr, ok := err.(*ApplicationError); ok {
			handlerResult.Error = appErr
		}
		return handlerResult, err
	}

	return handlerResult, nil
}

func resultForRequest(request RouteRequest) ApplicationResult {
	return ApplicationResult{
		RequestID: request.RequestID,
		Route:     request.Route,
		Target:    request.Target,
		Session:   request.Session,
		Identity:  request.Identity,
	}
}

func normalizeDispatchRoute(route RouteKey) (RouteKey, error) {
	normalizedRoute := RouteKey{
		Kind:   MessageKind(strings.TrimSpace(string(route.Kind))),
		Module: strings.TrimSpace(route.Module),
		Name:   strings.TrimSpace(route.Name),
	}

	if normalizedRoute.Kind == "" || normalizedRoute.Module == "" || normalizedRoute.Name == "" {
		return normalizedRoute, &ApplicationError{
			Code:    ErrorCodeInvalidRoute,
			Message: "route kind, module, and name are required",
			Route:   normalizedRoute,
		}
	}

	if normalizedRoute.Kind != MessageKindCommand && normalizedRoute.Kind != MessageKindQuery {
		return normalizedRoute, &ApplicationError{
			Code:    ErrorCodeUnsupportedMessageKind,
			Message: "only command and query routes are dispatchable",
			Route:   normalizedRoute,
		}
	}

	return normalizedRoute, nil
}

func isNilHandler(handler Handler) bool {
	if handler == nil {
		return true
	}

	value := reflect.ValueOf(handler)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}
