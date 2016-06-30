package http

import (
	"net/http"
	"bytes"
	"../.."
)

type Protocol struct {
	client *http.Client
}

func NewProtocol(client *http.Client) *Protocol {
	return &Protocol{client}
}

func (self *Protocol) Client() *http.Client {
	return self.client
}

func (self *Protocol) Send(req *communication.Request) (*communication.Response, error){

	httpRequest, err := NewHttpRequest(req)

	if (nil != err) {
		return nil, err
	}

	nativeRequest, err := http.NewRequest(
		httpRequest.Method(),
		httpRequest.Protocol() + httpRequest.Host() + httpRequest.path,
		bytes.NewBuffer([]byte(httpRequest.Body())),
	)

	if (nil != err) {
		return nil, err
	}

	nativeResponse, err := self.Client().Do(nativeRequest)

	if (nil != err) {
		return nil, err
	}

	httpResponse, err := NewHttpResponseFromNative(nativeResponse)

	if (nil != err) {
		return nil, err
	}

	return httpResponse.toBaseResponse(), nil
}