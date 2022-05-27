package config

type SecretString string

func (s SecretString) MarshalYAML() (interface{}, error) {
	if len(s) == 0 {
		return "", nil
	}
	return "[REDACTED]", nil
}

func (s SecretString) MarshalJSON() ([]byte, error) {
	if len(s) == 0 {
		return []byte(`""`), nil
	}
	return []byte(`"[REDACTED]"`), nil
}
