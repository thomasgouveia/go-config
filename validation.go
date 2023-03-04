package config

// isValidFormat checks if the configuration format is valid or not,
// and it return an error if it is invalid.
func isValidFormat(f format) bool {
	return f == JSON || f == YAML
}
