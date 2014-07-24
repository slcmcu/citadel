package citadel

// Executor is able to run the provided container on the specific resource
type Executor interface {
	Run(*Resource, *Container) error
}
