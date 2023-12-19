package vm4_test

import (
	"bytes"
	"testing"

	"github.com/apoloval/scumm-go/scummtest"
	"github.com/apoloval/scumm-go/vm4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeIndexFile(t *testing.T) {
	input := bytes.NewReader(scummtest.MonkeyIsland["000.LFL"])
	index, err := vm4.DecodeIndex(input)

	require.NoError(t, err)
	assert.Equal(t, 83, len(index.Rooms))
	assert.Equal(t, 170, len(index.Scripts))
	assert.Equal(t, 109, len(index.Sounds))
	assert.Equal(t, 120, len(index.Costumes))
	assert.Equal(t, 1000, len(index.Objects))
}
