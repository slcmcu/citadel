package scheduler

type Dock interface {
	Info() (*Info, error)               // Info returns information about the current Dock
	Running() States                    // Running returns the State for all the containers currently running
	Execute(*Container) (*State, error) // Execute runs the container on the Dock and returns the state
}
