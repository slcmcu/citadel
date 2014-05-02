package citadel

import (
	"time"

	rethink "github.com/dancannon/gorethink"
)

func NewRethinkSession(addr string) (*rethink.Session, error) {
	return rethink.Connect(rethink.ConnectOpts{
		Address:     addr,
		Database:    "citadel",
		MaxIdle:     10,
		IdleTimeout: time.Second * 10,
	})
}
