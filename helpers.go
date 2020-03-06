package stream

// String is a helper function for making pointers to strings.
func String(i string) *string {
	return &i
}

// Bool is a helper function for making pointers to bools.
func Bool(i bool) *bool {
	return &i
}
