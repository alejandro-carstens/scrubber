package configurations

// Configurable represents the contract to be implemented
// in order to comply with setting the configurations for
// the different notification channels
type Configurable interface {
	// Validate validates the configuration for a given channel
	Validate() (Configurable, error)

	// FillFromEnvs is responsible for setting the configuration
	// for the channel from the respective env variables
	FillFromEnvs() Configurable
}
