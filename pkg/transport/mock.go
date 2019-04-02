package transport

import (
	"context"
	"errors"
	"io"
	"net"
	"time"

	"github.com/skycoin/skywire/pkg/cipher"
)

// ErrTransportCommunicationTimeout represent timeout error for a mock transport.
var ErrTransportCommunicationTimeout = errors.New("transport communication operation timed out")

type fConn struct {
	net.Conn
	cipher.PubKey
}

// MockFactory implements Factory over net.Pipe connections.
type MockFactory struct {
	local cipher.PubKey
	in    chan *fConn
	out   chan *fConn
	fType string
}

// NewMockFactoryPair constructs a pair of MockFactories.
func NewMockFactoryPair(local, remote cipher.PubKey) (*MockFactory, *MockFactory) {
	in := make(chan *fConn)
	out := make(chan *fConn)
	return &MockFactory{local, in, out, "mock"}, &MockFactory{remote, out, in, "mock"}
}

// SetType sets type of transport.
func (f *MockFactory) SetType(fType string) {
	f.fType = fType
}

// Accept waits for new net.Conn notification from another MockFactory.
func (f *MockFactory) Accept(ctx context.Context) (Transport, error) {
	conn, more := <-f.in
	if !more {
		return nil, errors.New("factory: closed")
	}
	return NewMockTransport(conn, f.local, conn.PubKey), nil
}

// Dial creates pair of net.Conn via net.Pipe and passes one end to another MockFactory.
func (f *MockFactory) Dial(ctx context.Context, remote cipher.PubKey) (Transport, error) {
	in, out := net.Pipe()
	f.out <- &fConn{in, f.local}
	return NewMockTransport(out, f.local, remote), nil
}

// Close closes notification channel between a pair of MockFactories.
func (f *MockFactory) Close() error {
	select {
	case <-f.in:
	default:
		close(f.in)
	}
	return nil
}

// Local returns a local PubKey of the Factory.
func (f *MockFactory) Local() cipher.PubKey {
	return f.local
}

// Type returns type of the Factory.
func (f *MockFactory) Type() string {
	return f.fType
}

// MockTransport is a transport that accepts custom writers and readers to use them in Read and Write
// operations
type MockTransport struct {
	rw    io.ReadWriteCloser
	edges [2]cipher.PubKey
	// local   cipher.PubKey
	// remote  cipher.PubKey
	context context.Context
}

// NewMockTransport creates a transport with the given secret key and remote public key, taking a writer
// and a reader that will be used in the Write and Read operation
func NewMockTransport(rw io.ReadWriteCloser, local, remote cipher.PubKey) *MockTransport {
	return &MockTransport{rw, [2]cipher.PubKey{local, remote}, context.Background()}
}

// Read implements reader for mock transport
func (m *MockTransport) Read(p []byte) (n int, err error) {
	select {
	case <-m.context.Done():
		return 0, ErrTransportCommunicationTimeout
	default:
		return m.rw.Read(p)
	}
}

// Write implements writer for mock transport
func (m *MockTransport) Write(p []byte) (n int, err error) {
	select {
	case <-m.context.Done():
		return 0, ErrTransportCommunicationTimeout
	default:
		return m.rw.Write(p)
	}
}

// Close implements closer for mock transport
func (m *MockTransport) Close() error {
	return m.rw.Close()
}

// Edges returns edges of MockTransport
func (m *MockTransport) Edges() [2]cipher.PubKey {
	return m.edges
}

// // Local returns the local static public key
// func (m *MockTransport) Local() cipher.PubKey {
// 	return m.local
// }

// // Remote returns the remote public key fo the mock transport
// func (m *MockTransport) Remote() cipher.PubKey {
// 	return m.remote
// }

// SetDeadline sets a deadline for the write/read operations of the mock transport
func (m *MockTransport) SetDeadline(t time.Time) error {
	// nolint
	ctx, cancel := context.WithDeadline(m.context, t)
	m.context = ctx

	go func(cancel context.CancelFunc) {
		time.Sleep(time.Until(t))
		cancel()
	}(cancel)

	return nil
}

// Type returns the type of the mock transport
func (m *MockTransport) Type() string {
	return "mock"
}

// MockTransportManagersPair constructs a pair of Transport Managers
func MockTransportManagersPair() (pk1, pk2 cipher.PubKey, m1, m2 *Manager, errCh chan error, err error) {
	discovery := NewDiscoveryMock()
	logs := InMemoryTransportLogStore()

	var sk1, sk2 cipher.SecKey
	pk1, sk1 = cipher.GenerateKeyPair()
	pk2, sk2 = cipher.GenerateKeyPair()

	c1 := &ManagerConfig{PubKey: pk1, SecKey: sk1, DiscoveryClient: discovery, LogStore: logs}
	c2 := &ManagerConfig{PubKey: pk2, SecKey: sk2, DiscoveryClient: discovery, LogStore: logs}

	f1, f2 := NewMockFactoryPair(pk1, pk2)

	if m1, err = NewManager(c1, f1); err != nil {
		return
	}
	if m2, err = NewManager(c2, f2); err != nil {
		return
	}

	errCh = make(chan error)
	go func() { errCh <- m1.Serve(context.TODO()) }()
	go func() { errCh <- m2.Serve(context.TODO()) }()

	return
}

// MockTransportManager creates Manager
func MockTransportManager() (cipher.PubKey, *Manager, error) {
	_, pkB, mgrA, _, _, err := MockTransportManagersPair()
	return pkB, mgrA, err
}
