package sepay

// String returns a pointer to the given string value.
func String(s string) *string {
	return &s
}

// Int returns a pointer to the given int value.
func Int(i int) *int {
	return &i
}

// Float64 returns a pointer to the given float64 value.
func Float64(f float64) *float64 {
	return &f
}
