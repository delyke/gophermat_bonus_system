package flag

import "time"

func parseDurationOrDefault(val *string, def string, fallback time.Duration) time.Duration {
	var (
		d   time.Duration
		err error
	)

	if val == nil {
		d, err = time.ParseDuration(def)
	} else {
		d, err = time.ParseDuration(*val)
	}

	if err != nil {
		return fallback
	}
	return d
}

type flagValueSetter func() bool

func newFlagConfig[T any](raw *T, setters ...flagValueSetter) *T {
	flagWasSet := false

	for _, set := range setters {
		if set() {
			flagWasSet = true
		}
	}

	if !flagWasSet {
		return nil
	}
	return raw
}
