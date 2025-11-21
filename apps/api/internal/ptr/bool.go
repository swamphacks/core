package ptr

// Takes a boolean and returns a pointer to that boolean
func BoolToPtr(b bool) *bool {
	return &b
}
