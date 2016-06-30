package queue

import (
	"../communication"
	"../communication/protocols/http"
)

type Payload struct {
	method   string
	protocol *http.Protocol
	host     string
	path     string
}

func NewPayload(method string, protocol *http.Protocol, host, path string) *Payload {
	return &Payload{method, protocol, host, path}
}

func (self * Payload) Fetch(req *communication.Request) (*communication.Response, error) {
	return self.protocol.Send(req)
}
