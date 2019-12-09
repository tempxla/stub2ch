package types

type Maybe []string // optional parameter

func Nothing() Maybe {
	return nil
}

func Just(s string) Maybe {
	return []string{s}
}

func FromMaybe(m Maybe, def string) string {
	if m != nil {
		return m[0]
	} else {
		return def
	}
}

func FromJust(m Maybe) string {
	return m[0]
}

func IsJust(m Maybe) bool {
	return m != nil
}

func IsNothing(m Maybe) bool {
	return m == nil
}
