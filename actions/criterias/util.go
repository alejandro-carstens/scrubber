package criterias

import "errors"

// NewBuilder returns a *Builder instance
func NewBuilder() *Builder {
	return new(Builder)
}

func validateTimestring(timestring string) error {
	if len(timestring) == 0 {
		return errors.New("Since the Source is name you need to specify a timestring")
	}

	switch timestring {
	case "Y.m.d":
		break
	case "m.d.Y":
		break
	case "Y.m":
		break
	case "Y-m-d":
		break
	case "Y-m-d H:M":
		break
	case "Y-m-d H:M:S":
		break
	case "m-d-Y":
		break
	case "Y-m":
		break
	default:
		return errors.New("Invalid timestring please use Y.m.d, m.d.Y, Y.m, Y-m-d, Y-m-d H:M, Y-m-d H:M:S, m-d-Y, or Y-m")
	}

	return nil
}
