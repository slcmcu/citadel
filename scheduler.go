package citadel

type Scheduler interface {
	// Schedule returns true if the engine can run the specified image
	Schedule(*Image, *Engine) (bool, error)
}
