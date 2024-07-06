package wlistener

import (
	"net"
)

// AcceptFunc represents a functional (net.Listener).Accept implementation.
type AcceptFunc func() (net.Conn, error)

// CloseFunc represents a functional (net.Listener).Close implementation.
type CloseFunc func() error

// AddrFunc represents a functional (net.Listener).Addr implementation.
type AddrFunc func() net.Addr

// WrappedListener represents a wrapped net.Listener.
type WrappedListener struct {
	inner net.Listener

	onAccept AcceptFunc
	onClose  CloseFunc
	onAddr   AddrFunc
}

// Override represents an override for a net.Listener.
type Override func(*WrappedListener)

// OnAccept sets the override for the inner net.Listener's Accept method.
func OnAccept(fn AcceptFunc) Override { return func(wl *WrappedListener) { wl.onAccept = fn } }

// OnClose sets the override for the inner net.Listener's Close method.
func OnClose(fn CloseFunc) Override { return func(wl *WrappedListener) { wl.onClose = fn } }

// OnAddr sets the override for the inner net.Listener's Addr method.
func OnAddr(fn AddrFunc) Override { return func(wl *WrappedListener) { wl.onAddr = fn } }

// Wrap wraps a net.Listener with overriden methods.
func Wrap(listener net.Listener, overrides ...Override) net.Listener {
	wl := &WrappedListener{
		inner: listener,

		onAccept: listener.Accept,
		onClose:  listener.Close,
		onAddr:   listener.Addr,
	}
	for _, o := range overrides {
		o(wl)
	}
	return wl
}

// Accept waits for and returns the next connection to the listener.
func (wl *WrappedListener) Accept() (net.Conn, error) { return wl.onAccept() }

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (wl *WrappedListener) Close() error { return wl.onClose() }

// Addr returns the listener's network address.
func (wl *WrappedListener) Addr() net.Addr { return wl.onAddr() }
