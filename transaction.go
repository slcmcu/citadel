package citadel

import "time"

type Transaction struct {
	// Started is the time the transaction began
	Started time.Time `json:"started,omitempty"`

	// Ended is the time that the tranaction finished
	Ended time.Time `json:"ended,omitempty"`

	// Container is the current container that is being scheduled
	Container *Container `json:"container,omitempty"`

	// Placement is the selection from the cluster that is able to run the container
	Placement *Placement `json:"placement,omitempty"`
}

func newTransaction(c *Container) *Transaction {
	return &Transaction{
		Started:   time.Now(),
		Container: c,
	}
}

func (t *Transaction) place(e *Engine) {
	t.Placement = &Placement{
		Engine: e,
	}
}

func (t *Transaction) Close() error {
	t.Ended = time.Now()

	return nil
}
