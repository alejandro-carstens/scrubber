package criterias

type baseCriteria struct {
	Exclude bool `json:"exclude"`
}

// Include returns whether or not an index should
// be included on the actionable list
func (bc *baseCriteria) Include() bool {
	return !bc.Exclude
}
