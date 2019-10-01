package natsclient

import (
	nats "github.com/nats-io/nats.go"
)

// NewNATSServerConnection creates a new NATS server encoded connection
func NewNATSServerConnection(server string) (*nats.EncodedConn, error) {

	nc, err := nats.Connect(server)
	c, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}
	return c, nil

}
