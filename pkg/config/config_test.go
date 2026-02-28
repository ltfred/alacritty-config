package config

import "testing"

func TestConfig(t *testing.T) {
	config := GetConfigStruct()
	t.Log(config)
}
