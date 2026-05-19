package protobuf

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/iceiko/vibit/runtime/internal/app"
	authenticationv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/authentication/v1"
	protocolv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/protocol/v1"
	"google.golang.org/protobuf/proto"
)

type FrameRequest struct {
	ConnectionID    string
	ConnectionEpoch uint64
	RemoteAddr      string
	Payload         []byte
}

type ApplicationDispatcher interface {
	Dispatch(context.Context, app.RouteRequest) (app.ApplicationResult, error)
}

type RouteProtector interface {
	ProtectRoute(context.Context, app.RouteProtectionRequest) (app.RouteProtectionResult, error)
}

type ConnectionBinder interface {
	BindConnection(context.Context, app.ConnectionBindingRequest) (app.ConnectionBindingResult, error)
}

type FrameHandler struct {
	Dispatcher       ApplicationDispatcher
	RouteProtector   RouteProtector
	ConnectionBinder ConnectionBinder
}

func (h FrameHandler) HandleFrame(ctx context.Context, frame FrameRequest) ([][]byte, error) {
	if h.Dispatcher == nil {
		return nil, fmt.Errorf("protobuf frame handler dispatcher is nil")
	}

	var envelope protocolv1.Envelope
	if err := proto.Unmarshal(frame.Payload, &envelope); err != nil {
		return nil, err
	}

	request, err := RouteRequestMetadataFromEnvelope(&envelope)
	if err != nil {
		return nil, err
	}
	request = requestWithFrameMetadata(request, frame)
	if request.Route == app.BindConnectionRoute() {
		return h.handleConnectionBinding(ctx, request)
	}
	request, routeProtectionResult, err := h.routeRequestForDispatch(ctx, request)
	if err != nil {
		var applicationError *app.ApplicationError
		if !errors.As(err, &applicationError) {
			return nil, err
		}
		if routeProtectionResult.Error == nil {
			routeProtectionResult.Error = applicationError
		}
		return marshalApplicationResult(routeProtectionResult)
	}

	result, err := h.Dispatcher.Dispatch(ctx, request)
	if err != nil {
		var applicationError *app.ApplicationError
		if !errors.As(err, &applicationError) {
			return nil, err
		}
		if result.Error == nil {
			result.Error = applicationError
		}
	}

	return marshalApplicationResult(result)
}

func (h FrameHandler) routeRequestForDispatch(ctx context.Context, request app.RouteRequest) (app.RouteRequest, app.ApplicationResult, error) {
	if request.Route == app.LogoutAccessTokenRoute() && PayloadTypeIsAuthenticatedRequest(request.PayloadType) {
		err := authenticationMalformedError(request.Route)
		return app.RouteRequest{}, applicationErrorResultForRequest(request, err), err
	}
	if PayloadTypeIsAuthenticatedRequest(request.PayloadType) {
		wrapperPayload, err := DecodePayload(request.PayloadType, request.PayloadBytes)
		if err != nil {
			return app.RouteRequest{}, applicationErrorResultForRequest(request, authenticationMalformedError(request.Route)), authenticationMalformedError(request.Route)
		}
		wrapper, ok := wrapperPayload.(*authenticationv1.AuthenticatedRequest)
		if !ok || wrapper == nil {
			return app.RouteRequest{}, applicationErrorResultForRequest(request, authenticationMalformedError(request.Route)), authenticationMalformedError(request.Route)
		}
		innerRequest, err := routeRequestWithAuthenticatedPayload(request, wrapper)
		if err != nil {
			return app.RouteRequest{}, applicationErrorResultForRequest(request, authenticationMalformedError(request.Route)), authenticationMalformedError(request.Route)
		}
		protectedRequest, result, err := h.protectRoute(ctx, innerRequest, wrapper.GetAccessToken(), true)
		if err != nil {
			return app.RouteRequest{}, result, err
		}
		return routeRequestWithDecodedDomainPayload(protectedRequest)
	}

	protectedRequest, result, err := h.protectRoute(ctx, request, "", false)
	if err != nil {
		return app.RouteRequest{}, result, err
	}
	return routeRequestWithDecodedDomainPayload(protectedRequest)
}

func (h FrameHandler) protectRoute(ctx context.Context, request app.RouteRequest, accessToken string, proofCarrierPresent bool) (app.RouteRequest, app.ApplicationResult, error) {
	if h.RouteProtector == nil {
		return request, app.ApplicationResult{}, nil
	}
	protection, err := h.RouteProtector.ProtectRoute(ctx, app.RouteProtectionRequest{
		Request:             request,
		AccessToken:         accessToken,
		ProofCarrierPresent: proofCarrierPresent,
	})
	if err != nil {
		var applicationError *app.ApplicationError
		if errors.As(err, &applicationError) {
			return app.RouteRequest{}, applicationErrorResultForRequest(request, applicationError), err
		}
		return app.RouteRequest{}, app.ApplicationResult{}, err
	}
	if protection.Allowed {
		request.Identity = protection.Identity
	}
	return request, app.ApplicationResult{}, nil
}

func routeRequestWithDecodedDomainPayload(request app.RouteRequest) (app.RouteRequest, app.ApplicationResult, error) {
	payload, err := DecodePayload(request.PayloadType, request.PayloadBytes)
	if err != nil {
		return app.RouteRequest{}, app.ApplicationResult{}, err
	}
	request.Payload = payload
	request, err = RouteRequestWithDomainPayload(request)
	return request, app.ApplicationResult{}, err
}

func routeRequestWithAuthenticatedPayload(request app.RouteRequest, wrapper *authenticationv1.AuthenticatedRequest) (app.RouteRequest, error) {
	if wrapper == nil {
		return app.RouteRequest{}, &DecodeEnvelopeError{Message: "authenticated request payload is nil"}
	}
	innerPayloadType := strings.TrimSpace(wrapper.GetInnerPayloadType())
	if innerPayloadType == "" {
		return app.RouteRequest{}, &DecodeEnvelopeError{Message: "authenticated request inner payload_type is empty"}
	}
	innerPayload := append([]byte(nil), wrapper.GetInnerPayload()...)
	if len(innerPayload) == 0 {
		return app.RouteRequest{}, &DecodeEnvelopeError{Message: "authenticated request inner payload is empty"}
	}
	request.PayloadType = innerPayloadType
	request.PayloadBytes = innerPayload
	request.Payload = nil
	return request, nil
}

func PayloadTypeIsAuthenticatedRequest(payloadType string) bool {
	return strings.TrimSpace(payloadType) == "vibit.authentication.v1.AuthenticatedRequest"
}

func authenticationMalformedError(route app.RouteKey) *app.ApplicationError {
	return &app.ApplicationError{
		Code:    app.ErrorCodeAuthenticationTokenMalformed,
		Message: "authenticated request payload is malformed",
		Route:   route,
	}
}

func applicationErrorResultForRequest(request app.RouteRequest, applicationError *app.ApplicationError) app.ApplicationResult {
	result := app.ApplicationResult{
		RequestID: request.RequestID,
		Route:     request.Route,
		Target:    request.Target,
		Session:   request.Session,
		Identity:  request.Identity,
		Error:     applicationError,
	}
	return result
}

func marshalApplicationResult(result app.ApplicationResult) ([][]byte, error) {
	responseEnvelope, err := BuildEnvelopeFromApplicationResult(result)
	if err != nil {
		return nil, err
	}
	responsePayload, err := proto.Marshal(responseEnvelope)
	if err != nil {
		return nil, err
	}

	return [][]byte{responsePayload}, nil
}

func requestWithFrameMetadata(request app.RouteRequest, frame FrameRequest) app.RouteRequest {
	if connectionID := strings.TrimSpace(frame.ConnectionID); connectionID != "" {
		request.Session.ConnectionID = connectionID
	}
	if frame.ConnectionEpoch != 0 {
		request.Session.ConnectionEpoch = frame.ConnectionEpoch
	}
	if request.Identity.Status == "" || request.Identity.Status == app.IdentityValidationMetadataOnly {
		request.Identity = app.MetadataOnlyIdentityFromSession(request.Session)
	}
	return request
}
