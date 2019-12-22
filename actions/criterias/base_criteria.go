package criterias

type baseCriteria struct {
	Exclude bool `json:"exclude"`
}

// Include returns whether or not an index list
// be included on the actionable list
func (bc *baseCriteria) Include() bool {
	return !bc.Exclude
}
