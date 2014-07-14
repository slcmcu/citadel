package citadel

import (
	"github.com/citadel/citadel/utils"
)

type TransactionType string

const (
	RunTransaction    TransactionType = "run"
	StopTransaction   TransactionType = "stop"
	LoadTransaction   TransactionType = "load"
	DeleteTransaction TransactionType = "delete"
)

type Transaction struct {
	// ID is the uuid of a specific transaction
	ID string `json:"id,omitempty"`
	// Type is the transaction type, run, stop, register
	Type TransactionType `json:"type,omitempty"`
	// Containers is a list of containers affected for the given trasnaction
	Containers []*Container `json:"containers,omitempty"`
	// Error is the encountered if any
	Err error `json:"error,omitempty"`
	// Children are transactions run to complete this parent
	Children []*Transaction `json:"children,omitempty"`
}

func NewTransaction(t TransactionType) *Transaction {
	return &Transaction{
		ID:   utils.GenerateUUID(32),
		Type: t,
	}
}

func (t *Transaction) Error(err error) *Transaction {
	t.Err = err
	return t
}
