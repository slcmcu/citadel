package scheduler

type Dock interface {
	Info() (*Info, error)
	Containers() ([]*Container, error)
	Execute(*Container) (*State, error)
}
