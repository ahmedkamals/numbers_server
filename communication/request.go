package communication

type Request struct {
	body         string
	meta         map[string]string
}

func NewRequest(body string, meta map[string]string) *Request {
	return &Request{body, meta}
}

func (self *Request) Body() string {
	return self.body
}

func (self *Request) Meta() map[string]string {
	return self.meta
}