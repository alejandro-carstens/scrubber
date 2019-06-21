package options

type baseSnapshotOptions struct {
	defaultOptions
	Repository string `json:"repository"`
}

func (bso *baseSnapshotOptions) IsSnapshot() bool {
	return true
}
