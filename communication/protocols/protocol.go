package protocols

import (
	".."
)

type Protocol interface {
	Send(*communication.Request) (*communication.Response, error)
}
