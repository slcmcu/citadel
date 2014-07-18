package citadel

import "errors"

var (
	ErrNoValidHost = errors.New("no valid host for application")
)

type Scheduler interface {
	Schedule(*Application, []*Host) (*Host, error)
}
