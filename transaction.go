package citadel

import "time"

type Placement struct {
	Engine *Docker `json:"engine",omitempty"`
}

type Transaction struct {
	// Started is the time the transaction began
	Started time.Time `json:"started,omitempty"`

	// Ended is the time that the tranaction finished
	Ended time.Time `json:"ended,omitempty"`

	// Container is the current container that is being scheduled
	Container *Container `json:"container,omitempty"`

	// Placement is the selection from the cluster that is able to run the container
	Placement *Placement `json:"placement,omitempty"`

	// engines are the current engines in the bidding the run the container
	engines []*Docker
}

func newTransaction(c *Container, engines []*Docker) (*Transaction, error) {
	t := &Transaction{
		Started:   time.Now(),
		Container: c,
		engines:   engines,
	}

	for _, e := range engines {
		if err := e.loadContainers(); err != nil {
			return t, err
		}
	}

	return t, nil
}

func (t *Transaction) GetEngines() []*Docker {
	return t.engines
}

func (t *Transaction) Reduce(engines []*Docker) {
	t.engines = engines
}

func (t *Transaction) Close() error {
	t.Ended = time.Now()

	for _, e := range t.engines {
		e.cleanContainers()
	}

	return nil
}
