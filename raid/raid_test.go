package raid

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRaid0(t *testing.T) {

	files, err := os.ReadDir("test_input")
	assert.NoError(t, err)

	raid := NewRaid0()
	for _, file := range files {
		data, err := os.ReadFile(fmt.Sprintf("test_input/%s", file.Name()))
		assert.NoError(t, err)

		err = raid.Write(data)
		assert.NoError(t, err)

		s, err := raid.Read(len(data))
		assert.NoError(t, err)

		assert.Equal(t, string(data), s)

	}
}

func TestRaid1(t *testing.T) {

	files, err := os.ReadDir("test_input")
	assert.NoError(t, err)

	raid := NewRaid1()
	for _, file := range files {
		data, err := os.ReadFile(fmt.Sprintf("test_input/%s", file.Name()))
		assert.NoError(t, err)

		err = raid.Write(data)
		assert.NoError(t, err)

		s, err := raid.Read(len(data))
		assert.NoError(t, err)

		assert.Equal(t, string(data), s)

	}
}

func TestRaid10(t *testing.T) {

	files, err := os.ReadDir("test_input")
	assert.NoError(t, err)

	raid := NewRaid10()
	for _, file := range files {
		data, err := os.ReadFile(fmt.Sprintf("test_input/%s", file.Name()))
		assert.NoError(t, err)

		err = raid.Write(data)
		assert.NoError(t, err)

		s, err := raid.Read(len(data))
		assert.NoError(t, err)

		assert.Equal(t, string(data), s)

	}
}

func TestRaid5(t *testing.T) {

	files, err := os.ReadDir("test_input")
	assert.NoError(t, err)

	raid := NewRaid5()
	for _, file := range files {
		data, err := os.ReadFile(fmt.Sprintf("test_input/%s", file.Name()))
		assert.NoError(t, err)

		err = raid.Write(data)
		assert.NoError(t, err)

		s, err := raid.Read(len(data))
		assert.NoError(t, err)

		assert.Equal(t, string(data), s)

	}
}
