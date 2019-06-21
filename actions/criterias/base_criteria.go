package criterias

type baseCriteria struct {
	Exclude bool `json:"exclude"`
}

func (bc *baseCriteria) Include() bool {
	return !bc.Exclude
}
