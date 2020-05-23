package models

// Modelable represents the contract
// to be implemented by all models
type Modelable interface {
	// Indices returns a map of indices
	// names and indices properties
	Indices() map[string][]string

	// Table returns the table name
	Table() string
}
