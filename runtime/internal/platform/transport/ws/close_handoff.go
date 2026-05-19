package ws

import (
	"context"
	"strings"
	"sync"
	"time"
)

type CloseHandoffOutcome string

const (
	CloseHandoffOutcomeCloseRequested CloseHandoffOutcome = "close_requested"
	CloseHandoffOutcomeSocketNotFound CloseHandoffOutcome = "socket_not_found"
	CloseHandoffOutcomeEpochMismatch  CloseHandoffOutcome = "epoch_mismatch"
	CloseHandoffOutcomeAlreadyClosed  CloseHandoffOutcome = "already_closed"
	CloseHandoffOutcomeCloseFailed    CloseHandoffOutcome = "close_failed"
)

type CloseHandoffRequest struct {
	ConnectionID    string
	ConnectionEpoch uint64
	RequestedAt     time.Time
}

type CloseHandoffResult struct {
	ConnectionID    string
	ConnectionEpoch uint64
	Outcome         CloseHandoffOutcome
	ClosedAt        *time.Time
}

func (s *Server) RequestClose(ctx context.Context, request CloseHandoffRequest) CloseHandoffResult {
	request = normalizeCloseHandoffRequest(request)
	result := CloseHandoffResult{
		ConnectionID:    request.ConnectionID,
		ConnectionEpoch: request.ConnectionEpoch,
	}
	if s == nil || ctx == nil {
		result.Outcome = CloseHandoffOutcomeCloseFailed
		return result
	}
	if err := ctx.Err(); err != nil {
		result.Outcome = CloseHandoffOutcomeCloseFailed
		return result
	}
	if request.ConnectionID == "" || request.ConnectionEpoch == 0 {
		result.Outcome = CloseHandoffOutcomeSocketNotFound
		return result
	}

	entry, outcome := s.closeHandoff.find(request.ConnectionID, request.ConnectionEpoch)
	if outcome != CloseHandoffOutcomeCloseRequested {
		result.Outcome = outcome
		if entry != nil && outcome == CloseHandoffOutcomeAlreadyClosed {
			result.ClosedAt = entry.closedAt()
		}
		return result
	}

	if !entry.markCloseRequested() {
		result.Outcome = CloseHandoffOutcomeAlreadyClosed
		result.ClosedAt = entry.closedAt()
		return result
	}
	if err := entry.conn.CloseNow(); err != nil {
		entry.markCloseFailed()
		result.Outcome = CloseHandoffOutcomeCloseFailed
		return result
	}

	closedAt := request.RequestedAt
	if closedAt.IsZero() {
		closedAt = time.Now().UTC()
	} else {
		closedAt = closedAt.UTC()
	}
	entry.markClosedAt(closedAt)
	result.Outcome = CloseHandoffOutcomeCloseRequested
	result.ClosedAt = &closedAt
	return result
}

func (s *Server) registerCloseHandoffSocket(meta connectionMetadata, conn closeHandoffConn) {
	s.closeHandoff.register(meta, conn)
}

func (s *Server) unregisterCloseHandoffSocket(meta connectionMetadata) {
	s.closeHandoff.unregister(meta)
}

func normalizeCloseHandoffRequest(request CloseHandoffRequest) CloseHandoffRequest {
	request.ConnectionID = strings.TrimSpace(request.ConnectionID)
	if !request.RequestedAt.IsZero() {
		request.RequestedAt = request.RequestedAt.UTC()
	}
	return request
}

type closeHandoffSocketTable struct {
	mu      sync.RWMutex
	sockets map[string]*closeHandoffSocket
}

type closeHandoffConn interface {
	CloseNow() error
}

func (t *closeHandoffSocketTable) register(meta connectionMetadata, conn closeHandoffConn) {
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
		t.sockets = make(map[string]*closeHandoffSocket)
	}
	t.sockets[connectionID] = &closeHandoffSocket{
		connectionID:    connectionID,
		connectionEpoch: meta.epoch,
		conn:            conn,
	}
}

func (t *closeHandoffSocketTable) unregister(meta connectionMetadata) {
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
	entry.markClosedAt(time.Now().UTC())
}

func (t *closeHandoffSocketTable) find(connectionID string, epoch uint64) (*closeHandoffSocket, CloseHandoffOutcome) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	entry, ok := t.sockets[connectionID]
	if !ok {
		return nil, CloseHandoffOutcomeSocketNotFound
	}
	if entry.connectionEpoch != epoch {
		return nil, CloseHandoffOutcomeEpochMismatch
	}
	if entry.isClosed() {
		return entry, CloseHandoffOutcomeAlreadyClosed
	}
	return entry, CloseHandoffOutcomeCloseRequested
}

type closeHandoffSocket struct {
	mu              sync.Mutex
	connectionID    string
	connectionEpoch uint64
	conn            closeHandoffConn
	closeRequested  bool
	closedAtValue   *time.Time
}

func (s *closeHandoffSocket) markCloseRequested() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closeRequested || s.closedAtValue != nil {
		return false
	}
	s.closeRequested = true
	return true
}

func (s *closeHandoffSocket) markClosedAt(closedAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	closedAt = closedAt.UTC()
	s.closedAtValue = &closedAt
}

func (s *closeHandoffSocket) markCloseFailed() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closeRequested = false
}

func (s *closeHandoffSocket) isClosed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closeRequested || s.closedAtValue != nil
}

func (s *closeHandoffSocket) closedAt() *time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closedAtValue == nil {
		return nil
	}
	closedAt := *s.closedAtValue
	return &closedAt
}
