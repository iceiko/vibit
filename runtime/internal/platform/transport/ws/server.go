package ws

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/coder/websocket"
)

const defaultReadLimitBytes int64 = 1 << 20

type Frame struct {
	ConnectionID string
	RemoteAddr   string
	Payload      []byte
}

type FrameHandler interface {
	HandleFrame(context.Context, Frame) ([][]byte, error)
}

type FrameHandlerFunc func(context.Context, Frame) ([][]byte, error)

func (f FrameHandlerFunc) HandleFrame(ctx context.Context, frame Frame) ([][]byte, error) {
	if f == nil {
		return nil, fmt.Errorf("websocket frame handler function is nil")
	}
	return f(ctx, frame)
}

type Server struct {
	Handler        FrameHandler
	AcceptOptions  *websocket.AcceptOptions
	ReadLimitBytes int64

	nextConnectionID uint64
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.Handler == nil {
		http.Error(w, "websocket frame handler is required", http.StatusInternalServerError)
		return
	}

	conn, err := websocket.Accept(w, r, s.AcceptOptions)
	if err != nil {
		return
	}
	defer conn.CloseNow()

	readLimit := s.ReadLimitBytes
	if readLimit == 0 {
		readLimit = defaultReadLimitBytes
	}
	conn.SetReadLimit(readLimit)

	meta := connectionMetadata{
		id:         fmt.Sprintf("ws-%d", atomic.AddUint64(&s.nextConnectionID, 1)),
		remoteAddr: r.RemoteAddr,
	}
	_ = s.handleConnection(r.Context(), conn, meta)
}

type connectionMetadata struct {
	id         string
	remoteAddr string
}

func (s *Server) handleConnection(ctx context.Context, conn *websocket.Conn, meta connectionMetadata) error {
	for {
		messageType, payload, err := conn.Read(ctx)
		if err != nil {
			if closeStatus := websocket.CloseStatus(err); closeStatus == websocket.StatusNormalClosure || closeStatus == websocket.StatusGoingAway {
				return nil
			}
			return err
		}

		if messageType != websocket.MessageBinary {
			_ = conn.Close(websocket.StatusUnsupportedData, "binary frames required")
			return fmt.Errorf("unsupported websocket message type %d", messageType)
		}

		responses, err := s.Handler.HandleFrame(ctx, newFrame(meta, payload))
		if err != nil {
			_ = conn.Close(websocket.StatusInternalError, "frame handler failed")
			return err
		}

		for _, response := range responses {
			if err := conn.Write(ctx, websocket.MessageBinary, response); err != nil {
				return err
			}
		}
	}
}

func newFrame(meta connectionMetadata, payload []byte) Frame {
	return Frame{
		ConnectionID: meta.id,
		RemoteAddr:   meta.remoteAddr,
		Payload:      append([]byte(nil), payload...),
	}
}
