# wnet

Go library that makes it easy to override net package interfaces with minimal code.

[![Go Report Card](https://goreportcard.com/badge/github.com/adrianosela/wnet)](https://goreportcard.com/report/github.com/adrianosela/wnet)
[![Documentation](https://godoc.org/github.com/adrianosela/wnet?status.svg)](https://godoc.org/github.com/adrianosela/wnet)
[![GitHub issues](https://img.shields.io/github/issues/adrianosela/wnet.svg)](https://github.com/adrianosela/wnet/issues)
[![license](https://img.shields.io/github/license/adrianosela/wnet.svg)](https://github.com/adrianosela/wnet/blob/master/LICENSE)

### Usage

```
listener, err := net.Listen("tcp", ":4242")
if err != nil {
    log.Fatalf("failed to start tcp listener: %v", err)
}

wrapped := wlistener.Wrap(
    listener,
    wlistener.OnAccept(func() (net.Conn, error) {
        log.Println("waiting for connections...")
        return listener.Accept()
    }),
)
defer wrapped.Close()

// use wrapped as your usual listener...
```