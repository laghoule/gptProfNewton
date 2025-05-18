package config

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		error  bool
	}{
		{
			name:     "load good config",
			fileName: "testdata/good.yaml",
			error:  false,
		},
		{
			name:     "load bad config",
			fileName: "testdata/bad.yaml",
			error:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			conf := &Config{}
			err := conf.Load(test.fileName)
			if test.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}