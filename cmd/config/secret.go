package config

type SecretString string

func (s SecretString) MarshalYAML() (interface{}, error) {
	if len(s) == 0 {
		return "", nil
	}
	return "[REDACTED]", nil
}
