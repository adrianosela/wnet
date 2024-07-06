package wconn

import (
	"net"
	"time"
)

// ReadFunc represents a functional (net.Conn).Read implementation.
type ReadFunc func([]byte) (int, error)

// WriteFunc represents a functional (net.Conn).Write implementation.
type WriteFunc func([]byte) (int, error)

// CloseFunc represents a functional (net.Conn).Close implementation.
type CloseFunc func() error

// LocalAddrFunc represents a functional (net.Conn).LocalAddr implementation.
type LocalAddrFunc func() net.Addr

// RemoteAddrFunc represents a functional (net.Conn).RemoteAddr implementation.
type RemoteAddrFunc func() net.Addr

// SetDeadlineFunc represents a functional (net.Conn).SetDeadline implementation.
type SetDeadlineFunc func(time.Time) error

// SetReadDeadlineFunc represents a functional (net.Conn).SetReadDeadline implementation.
type SetReadDeadlineFunc func(time.Time) error

// SetWriteDeadlineFunc represents a functional (net.Conn).SetWriteDeadline implementation.
type SetWriteDeadlineFunc func(time.Time) error

// WrappedConn represents a wrapped net.Conn.
type WrappedConn struct {
	inner net.Conn

	onRead             ReadFunc
	onWrite            WriteFunc
	onClose            CloseFunc
	onLocalAddr        LocalAddrFunc
	onRemoteAddr       RemoteAddrFunc
	onSetDeadline      SetDeadlineFunc
	onSetReadDeadline  SetReadDeadlineFunc
	onSetWriteDeadline SetWriteDeadlineFunc
}

// Override represents an override for a net.Conn.
type Override func(*WrappedConn)

// OnRead sets the override for the inner net.Conn's Read method.
func OnRead(fn ReadFunc) Override {
	return func(wc *WrappedConn) { wc.onRead = fn }
}

// OnWrite sets the override for the inner net.Conn's Write method.
func OnWrite(fn WriteFunc) Override {
	return func(wc *WrappedConn) { wc.onWrite = fn }
}

// OnClose sets the override for the inner net.Conn's Close method.
func OnClose(fn CloseFunc) Override {
	return func(wc *WrappedConn) { wc.onClose = fn }
}

// OnLocalAddr sets the override for the inner net.Conn's LocalAddr method.
func OnLocalAddr(fn LocalAddrFunc) Override {
	return func(wc *WrappedConn) { wc.onLocalAddr = fn }
}

// OnRemoteAddr sets the override for the inner net.Conn's RemoteAddr method.
func OnRemoteAddr(fn RemoteAddrFunc) Override {
	return func(wc *WrappedConn) { wc.onRemoteAddr = fn }
}

// OnSetDeadline sets the override for the inner net.Conn's SetDeadline method.
func OnSetDeadline(fn SetDeadlineFunc) Override {
	return func(wc *WrappedConn) { wc.onSetDeadline = fn }
}

// OnSetReadDeadline sets the override for the inner net.Conn's SetReadDeadline method.
func OnSetReadDeadline(fn SetReadDeadlineFunc) Override {
	return func(wc *WrappedConn) { wc.onSetReadDeadline = fn }
}

// OnSetWriteDeadline sets the override for the inner net.Conn's SetWriteDeadline method.
func OnSetWriteDeadline(fn SetWriteDeadlineFunc) Override {
	return func(wc *WrappedConn) { wc.onSetWriteDeadline = fn }
}

// Wrap wraps a net.Conn with overriden methods.
func Wrap(conn net.Conn, overrides ...Override) net.Conn {
	wc := &WrappedConn{
		inner:              conn,
		onRead:             conn.Read,
		onWrite:            conn.Write,
		onClose:            conn.Close,
		onLocalAddr:        conn.LocalAddr,
		onRemoteAddr:       conn.RemoteAddr,
		onSetDeadline:      conn.SetDeadline,
		onSetReadDeadline:  conn.SetReadDeadline,
		onSetWriteDeadline: conn.SetWriteDeadline,
	}
	for _, o := range overrides {
		o(wc)
	}
	return wc
}

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (wc *WrappedConn) Read(b []byte) (int, error) { return wc.onRead(b) }

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (wc *WrappedConn) Write(b []byte) (int, error) { return wc.onWrite(b) }

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (wc *WrappedConn) Close() error { return wc.onClose() }

// LocalAddr returns the local network address, if known.
func (wc *WrappedConn) LocalAddr() net.Addr { return wc.onLocalAddr() }

// RemoteAddr returns the remote network address, if known.
func (wc *WrappedConn) RemoteAddr() net.Addr { return wc.onRemoteAddr() }

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail instead of blocking. The deadline applies to all future
// and pending I/O, not just the immediately following call to
// Read or Write. After a deadline has been exceeded, the
// connection can be refreshed by setting a deadline in the future.
//
// If the deadline is exceeded a call to Read or Write or to other
// I/O methods will return an error that wraps os.ErrDeadlineExceeded.
// This can be tested using errors.Is(err, os.ErrDeadlineExceeded).
// The error's Timeout method will return true, but note that there
// are other possible errors for which the Timeout method will
// return true even if the deadline has not been exceeded.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
func (wc *WrappedConn) SetDeadline(d time.Time) error { return wc.onSetDeadline(d) }

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (wc *WrappedConn) SetReadDeadline(d time.Time) error { return wc.onSetReadDeadline(d) }

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (wc *WrappedConn) SetWriteDeadline(d time.Time) error { return wc.onSetWriteDeadline(d) }
