package protobuf

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/iceiko/vibit/runtime/internal/app"
	protocolv1 "github.com/iceiko/vibit/runtime/internal/generated/proto/vibit/protocol/v1"
	"google.golang.org/protobuf/proto"
)

type FrameRequest struct {
	ConnectionID string
	RemoteAddr   string
	Payload      []byte
}

type ApplicationDispatcher interface {
	Dispatch(context.Context, app.RouteRequest) (app.ApplicationResult, error)
}

type FrameHandler struct {
	Dispatcher ApplicationDispatcher
}

func (h FrameHandler) HandleFrame(ctx context.Context, frame FrameRequest) ([][]byte, error) {
	if h.Dispatcher == nil {
		return nil, fmt.Errorf("protobuf frame handler dispatcher is nil")
	}

	var envelope protocolv1.Envelope
	if err := proto.Unmarshal(frame.Payload, &envelope); err != nil {
		return nil, err
	}

	request, err := RouteRequestFromEnvelopeForDispatch(&envelope)
	if err != nil {
		return nil, err
	}
	request = requestWithFrameMetadata(request, frame)

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
	if request.Session.ConnectionID == "" {
		request.Session.ConnectionID = strings.TrimSpace(frame.ConnectionID)
	}
	if request.Identity.Status == "" || request.Identity.Status == app.IdentityValidationMetadataOnly {
		request.Identity = app.MetadataOnlyIdentityFromSession(request.Session)
	}
	return request
}
