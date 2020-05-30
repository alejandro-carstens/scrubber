package controllers

type Controllerable interface {
	new() Controllerable

	Routes() []*Route
}
