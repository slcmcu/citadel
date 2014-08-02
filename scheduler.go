package citadel

type Scheduler interface {
	// Schedule returns true if the engine can run the specified container
	Schedule(*Container, *Engine) (bool, error)
}
