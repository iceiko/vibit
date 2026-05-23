package ws

import (
	"bytes"
	"context"
	"errors"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/coder/websocket"
)

func TestDeliverBinaryFrameWritesToServerObservedConnectionAndEpoch(t *testing.T) {
	server := &Server{}
	conn := &fakeOutboundConn{}
	server.registerOutboundSocket(connectionMetadata{id: " ws-1 ", epoch: 7}, conn)

	requestedAt := time.Date(2026, 5, 23, 12, 0, 0, 0, time.FixedZone("test", 8*60*60))
	payload := []byte("encoded-envelope")
	result := server.DeliverBinaryFrame(context.Background(), OutboundDeliveryRequest{
		ConnectionID:    " ws-1 ",
		ConnectionEpoch: 7,
		Payload:         payload,
		RequestedAt:     requestedAt,
	})

	if result.Outcome != OutboundDeliveryOutcomeDelivered {
		t.Fatalf("DeliverBinaryFrame() outcome = %s, want %s", result.Outcome, OutboundDeliveryOutcomeDelivered)
	}
	if result.ConnectionID != "ws-1" || result.ConnectionEpoch != 7 {
		t.Fatalf("DeliverBinaryFrame() target = %q/%d, want ws-1/7", result.ConnectionID, result.ConnectionEpoch)
	}
	if result.DeliveredAt == nil || !result.DeliveredAt.Equal(requestedAt.UTC()) {
		t.Fatalf("DeliveredAt = %v, want %v", result.DeliveredAt, requestedAt.UTC())
	}
	if conn.writeCount() != 1 {
		t.Fatalf("Write calls = %d, want 1", conn.writeCount())
	}
	if got := conn.payload(0); !bytes.Equal(got, []byte("encoded-envelope")) {
		t.Fatalf("payload = %q, want encoded-envelope", string(got))
	}
	if conn.messageType(0) != websocket.MessageBinary {
		t.Fatalf("message type = %v, want binary", conn.messageType(0))
	}

	payload[0] = 'x'
	if got := conn.payload(0); !bytes.Equal(got, []byte("encoded-envelope")) {
		t.Fatalf("written payload was mutated to %q", string(got))
	}
}

func TestDeliverBinaryFrameRejectsMissingOrStaleSocketWithoutWriting(t *testing.T) {
	server := &Server{}
	current := &fakeOutboundConn{}
	server.registerOutboundSocket(connectionMetadata{id: "ws-1", epoch: 8}, current)

	stale := server.DeliverBinaryFrame(context.Background(), OutboundDeliveryRequest{
		ConnectionID:    "ws-1",
		ConnectionEpoch: 7,
		Payload:         []byte("payload"),
	})
	if stale.Outcome != OutboundDeliveryOutcomeEpochMismatch {
		t.Fatalf("stale outcome = %s, want %s", stale.Outcome, OutboundDeliveryOutcomeEpochMismatch)
	}
	if current.writeCount() != 0 {
		t.Fatalf("stale epoch wrote to current socket %d times, want 0", current.writeCount())
	}

	missing := server.DeliverBinaryFrame(context.Background(), OutboundDeliveryRequest{
		ConnectionID:    "missing",
		ConnectionEpoch: 1,
		Payload:         []byte("payload"),
	})
	if missing.Outcome != OutboundDeliveryOutcomeSocketNotFound || missing.DeliveredAt != nil {
		t.Fatalf("missing result = %#v, want redacted socket-not-found", missing)
	}
}

func TestDeliverBinaryFrameReturnsAlreadyClosedForClosedEntry(t *testing.T) {
	server := &Server{}
	conn := &fakeOutboundConn{}
	meta := connectionMetadata{id: "ws-1", epoch: 1}
	server.registerOutboundSocket(meta, conn)

	entry, outcome := server.outbound.find("ws-1", 1)
	if outcome != OutboundDeliveryOutcomeDelivered || entry == nil {
		t.Fatalf("find before unregister outcome = %s entry = %#v, want delivered entry", outcome, entry)
	}
	entry.markClosed()

	result := server.DeliverBinaryFrame(context.Background(), OutboundDeliveryRequest{
		ConnectionID:    "ws-1",
		ConnectionEpoch: 1,
		Payload:         []byte("payload"),
	})
	if result.Outcome != OutboundDeliveryOutcomeAlreadyClosed {
		t.Fatalf("DeliverBinaryFrame() outcome = %s, want already_closed", result.Outcome)
	}
	if conn.writeCount() != 0 {
		t.Fatalf("closed socket write calls = %d, want 0", conn.writeCount())
	}
	server.unregisterOutboundSocket(meta)
}

func TestDeliverBinaryFrameReturnsRedactedWriteFailedOutcome(t *testing.T) {
	server := &Server{}
	conn := &fakeOutboundConn{err: errors.New("raw socket detail")}
	server.registerOutboundSocket(connectionMetadata{id: "ws-1", epoch: 1}, conn)

	result := server.DeliverBinaryFrame(context.Background(), OutboundDeliveryRequest{
		ConnectionID:    "ws-1",
		ConnectionEpoch: 1,
		Payload:         []byte("payload"),
	})

	if result.Outcome != OutboundDeliveryOutcomeWriteFailed {
		t.Fatalf("DeliverBinaryFrame() outcome = %s, want %s", result.Outcome, OutboundDeliveryOutcomeWriteFailed)
	}
	if result.DeliveredAt != nil {
		t.Fatalf("DeliveredAt = %v, want nil after failed write", result.DeliveredAt)
	}
	if conn.writeCount() != 1 {
		t.Fatalf("Write calls = %d, want 1", conn.writeCount())
	}
}

func TestDeliverBinaryFrameSerializesSocketWrites(t *testing.T) {
	server := &Server{}
	conn := &fakeOutboundConn{block: make(chan struct{})}
	server.registerOutboundSocket(connectionMetadata{id: "ws-1", epoch: 1}, conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	firstDone := make(chan OutboundDeliveryResult, 1)
	secondDone := make(chan OutboundDeliveryResult, 1)
	go func() {
		firstDone <- server.DeliverBinaryFrame(ctx, OutboundDeliveryRequest{
			ConnectionID:    "ws-1",
			ConnectionEpoch: 1,
			Payload:         []byte("first"),
		})
	}()
	if !eventually(time.Second, func() bool { return conn.inFlightCount() == 1 }) {
		t.Fatal("first write did not start")
	}
	go func() {
		secondDone <- server.DeliverBinaryFrame(ctx, OutboundDeliveryRequest{
			ConnectionID:    "ws-1",
			ConnectionEpoch: 1,
			Payload:         []byte("second"),
		})
	}()
	time.Sleep(50 * time.Millisecond)
	if conn.maxInFlightCount() != 1 {
		t.Fatalf("max concurrent writes = %d, want serialized writes", conn.maxInFlightCount())
	}
	conn.releaseOne()
	if result := <-firstDone; result.Outcome != OutboundDeliveryOutcomeDelivered {
		t.Fatalf("first result = %#v, want delivered", result)
	}
	if !eventually(time.Second, func() bool { return conn.inFlightCount() == 1 }) {
		t.Fatal("second write did not start after first completed")
	}
	conn.releaseOne()
	if result := <-secondDone; result.Outcome != OutboundDeliveryOutcomeDelivered {
		t.Fatalf("second result = %#v, want delivered", result)
	}
	if got := conn.maxInFlightCount(); got != 1 {
		t.Fatalf("max concurrent writes = %d, want 1", got)
	}
}

func TestServerDeliversOutboundFrameToAcceptedSocket(t *testing.T) {
	accepted := make(chan Frame, 1)
	wsServer := &Server{
		AcceptOptions: &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		},
		Handler: FrameHandlerFunc(func(ctx context.Context, frame Frame) ([][]byte, error) {
			accepted <- frame
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(5 * time.Second):
				return nil, nil
			}
		}),
	}
	server := httptest.NewServer(wsServer)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, websocketURL(server.URL), nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer conn.CloseNow()

	if err := conn.Write(ctx, websocket.MessageBinary, []byte("register")); err != nil {
		t.Fatalf("write initial frame: %v", err)
	}
	var frame Frame
	select {
	case frame = <-accepted:
	case <-ctx.Done():
		t.Fatal("handler did not receive accepted frame")
	}

	result := wsServer.DeliverBinaryFrame(ctx, OutboundDeliveryRequest{
		ConnectionID:    frame.ConnectionID,
		ConnectionEpoch: frame.ConnectionEpoch,
		Payload:         []byte("server-push"),
	})
	if result.Outcome != OutboundDeliveryOutcomeDelivered {
		t.Fatalf("DeliverBinaryFrame() result = %#v, want delivered", result)
	}

	messageType, payload, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("read outbound frame: %v", err)
	}
	if messageType != websocket.MessageBinary || string(payload) != "server-push" {
		t.Fatalf("outbound frame type=%v payload=%q, want binary server-push", messageType, string(payload))
	}
}

func TestOutboundDeliveryRequestAndResultRemainCredentialNeutral(t *testing.T) {
	for _, typ := range []reflect.Type{
		reflect.TypeOf(OutboundDeliveryRequest{}),
		reflect.TypeOf(OutboundDeliveryResult{}),
	} {
		for _, forbidden := range []string{
			"PlayerID",
			"RuntimeSessionID",
			"AccessToken",
			"AccessTokenRecordID",
			"Credential",
			"RawToken",
			"RawCredential",
			"LookupDigest",
			"VerifierDigest",
			"VerifierKeyID",
			"Authorization",
			"Cookie",
			"QueryString",
			"Subprotocol",
			"RemoteAddr",
			"PayloadType",
			"MessageKind",
			"Nakama",
			"Pitaya",
		} {
			if _, ok := typ.FieldByName(forbidden); ok {
				t.Fatalf("%s contains forbidden policy, proof, or compatibility field %s", typ.Name(), forbidden)
			}
		}
	}
}

type fakeOutboundConn struct {
	mu          sync.Mutex
	types       []websocket.MessageType
	payloads    [][]byte
	err         error
	block       chan struct{}
	inFlight    int
	maxInFlight int
}

func (f *fakeOutboundConn) Write(_ context.Context, messageType websocket.MessageType, payload []byte) error {
	f.mu.Lock()
	f.inFlight++
	if f.inFlight > f.maxInFlight {
		f.maxInFlight = f.inFlight
	}
	f.mu.Unlock()

	if f.block != nil {
		<-f.block
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	f.inFlight--
	f.types = append(f.types, messageType)
	f.payloads = append(f.payloads, append([]byte(nil), payload...))
	return f.err
}

func (f *fakeOutboundConn) writeCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.payloads)
}

func (f *fakeOutboundConn) payload(index int) []byte {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]byte(nil), f.payloads[index]...)
}

func (f *fakeOutboundConn) messageType(index int) websocket.MessageType {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.types[index]
}

func (f *fakeOutboundConn) inFlightCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.inFlight
}

func (f *fakeOutboundConn) maxInFlightCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.maxInFlight
}

func (f *fakeOutboundConn) releaseOne() {
	f.block <- struct{}{}
}
