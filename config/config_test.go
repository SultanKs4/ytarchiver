package config

import (
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{name: "create empty config", want: &Config{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_LoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "load config", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{}
			if err := c.LoadConfig(); (err != nil) != tt.wantErr {
				t.Errorf("Config.LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
