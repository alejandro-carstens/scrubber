package configurations

type Configurable interface {
	Validate() (Configurable, error)

	FillFromEnvs() Configurable
}
