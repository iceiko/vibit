package ws

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
)

type OutboundDeliveryOutcome string

const (
	OutboundDeliveryOutcomeDelivered      OutboundDeliveryOutcome = "delivered"
	OutboundDeliveryOutcomeSocketNotFound OutboundDeliveryOutcome = "socket_not_found"
	OutboundDeliveryOutcomeEpochMismatch  OutboundDeliveryOutcome = "epoch_mismatch"
	OutboundDeliveryOutcomeAlreadyClosed  OutboundDeliveryOutcome = "already_closed"
	OutboundDeliveryOutcomeWriteFailed    OutboundDeliveryOutcome = "write_failed"
)

type OutboundDeliveryRequest struct {
	ConnectionID    string
	ConnectionEpoch uint64
	Payload         []byte
	RequestedAt     time.Time
}

type OutboundDeliveryResult struct {
	ConnectionID    string
	ConnectionEpoch uint64
	Outcome         OutboundDeliveryOutcome
	DeliveredAt     *time.Time
}

func (s *Server) DeliverBinaryFrame(ctx context.Context, request OutboundDeliveryRequest) OutboundDeliveryResult {
	request = normalizeOutboundDeliveryRequest(request)
	result := OutboundDeliveryResult{
		ConnectionID:    request.ConnectionID,
		ConnectionEpoch: request.ConnectionEpoch,
	}
	if s == nil || ctx == nil {
		result.Outcome = OutboundDeliveryOutcomeWriteFailed
		return result
	}
	if err := ctx.Err(); err != nil {
		result.Outcome = OutboundDeliveryOutcomeWriteFailed
		return result
	}
	if request.ConnectionID == "" || request.ConnectionEpoch == 0 || len(request.Payload) == 0 {
		result.Outcome = OutboundDeliveryOutcomeSocketNotFound
		return result
	}

	entry, outcome := s.outbound.find(request.ConnectionID, request.ConnectionEpoch)
	if outcome != OutboundDeliveryOutcomeDelivered {
		result.Outcome = outcome
		return result
	}

	if err := entry.write(ctx, request.Payload); err != nil {
		result.Outcome = OutboundDeliveryOutcomeWriteFailed
		return result
	}

	deliveredAt := request.RequestedAt
	if deliveredAt.IsZero() {
		deliveredAt = time.Now().UTC()
	} else {
		deliveredAt = deliveredAt.UTC()
	}
	result.Outcome = OutboundDeliveryOutcomeDelivered
	result.DeliveredAt = &deliveredAt
	return result
}

func (s *Server) registerOutboundSocket(meta connectionMetadata, conn outboundConn) {
	s.outbound.register(meta, conn)
}

func (s *Server) unregisterOutboundSocket(meta connectionMetadata) {
	s.outbound.unregister(meta)
}

func (s *Server) writeAcceptedBinaryFrame(ctx context.Context, meta connectionMetadata, payload []byte) error {
	entry, outcome := s.outbound.find(strings.TrimSpace(meta.id), meta.epoch)
	if outcome != OutboundDeliveryOutcomeDelivered || entry == nil {
		return errOutboundSocketClosed
	}
	return entry.write(ctx, payload)
}

func normalizeOutboundDeliveryRequest(request OutboundDeliveryRequest) OutboundDeliveryRequest {
	request.ConnectionID = strings.TrimSpace(request.ConnectionID)
	request.Payload = append([]byte(nil), request.Payload...)
	if !request.RequestedAt.IsZero() {
		request.RequestedAt = request.RequestedAt.UTC()
	}
	return request
}

type outboundConn interface {
	Write(context.Context, websocket.MessageType, []byte) error
}

type outboundSocketTable struct {
	mu      sync.RWMutex
	sockets map[string]*outboundSocket
}

func (t *outboundSocketTable) register(meta connectionMetadata, conn outboundConn) {
	if conn == nil {
		return
	}
	connectionID := strings.TrimSpace(meta.id)
	if connectionID == "" || meta.epoch == 0 {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	if t.sockets == nil {
		t.sockets = make(map[string]*outboundSocket)
	}
	t.sockets[connectionID] = &outboundSocket{
		connectionID:    connectionID,
		connectionEpoch: meta.epoch,
		conn:            conn,
	}
}

func (t *outboundSocketTable) unregister(meta connectionMetadata) {
	connectionID := strings.TrimSpace(meta.id)
	if connectionID == "" {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	entry, ok := t.sockets[connectionID]
	if !ok || entry.connectionEpoch != meta.epoch {
		return
	}
	delete(t.sockets, connectionID)
	entry.markClosed()
}

func (t *outboundSocketTable) find(connectionID string, epoch uint64) (*outboundSocket, OutboundDeliveryOutcome) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	entry, ok := t.sockets[connectionID]
	if !ok {
		return nil, OutboundDeliveryOutcomeSocketNotFound
	}
	if entry.connectionEpoch != epoch {
		return nil, OutboundDeliveryOutcomeEpochMismatch
	}
	if entry.isClosed() {
		return entry, OutboundDeliveryOutcomeAlreadyClosed
	}
	return entry, OutboundDeliveryOutcomeDelivered
}

type outboundSocket struct {
	mu              sync.Mutex
	connectionID    string
	connectionEpoch uint64
	conn            outboundConn
	closed          bool
}

func (s *outboundSocket) write(ctx context.Context, payload []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return errOutboundSocketClosed
	}
	return s.conn.Write(ctx, websocket.MessageBinary, append([]byte(nil), payload...))
}

func (s *outboundSocket) markClosed() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
}

func (s *outboundSocket) isClosed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closed
}

type outboundSocketClosedError struct{}

func (outboundSocketClosedError) Error() string {
	return "websocket outbound delivery: socket is closed"
}

var errOutboundSocketClosed error = outboundSocketClosedError{}
