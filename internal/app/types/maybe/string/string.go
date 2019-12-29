package string

type Maybe struct {
	a         string
	isNothing bool
}

func Nothing() *Maybe {
	return &Maybe{a: "", isNothing: true}
}

func Just(s string) *Maybe {
	return &Maybe{a: s, isNothing: false}
}

func (m *Maybe) FromMaybe(defaultString string) string {
	if m.IsNothing() {
		return defaultString
	} else {
		return m.a
	}
}

func (m *Maybe) IsJust() bool {
	return !m.isNothing
}

func (m *Maybe) IsNothing() bool {
	return m.isNothing
}
