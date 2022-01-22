package axudp

import (
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBooleansBytes(t *testing.T) {

	mask := []bool{true, false, false, true, true}
	maskBytes := booleansToBytes(mask)
	assert.Equal(t, len(maskBytes), 1)
	assert.EqualValues(t, mask, bytesToBooleans(maskBytes)[:5])

	mask = []bool{true, false, false, true, true, false, false, false, true}
	maskBytes = booleansToBytes(mask)
	assert.Equal(t, len(maskBytes), 2)
	assert.EqualValues(t, mask, bytesToBooleans(maskBytes)[:9])
}

func TestLength(t *testing.T) {
	size := 65_535
	lenBytes := getBytesFromInt(size)
	log.Info().Hex("size", lenBytes).Msg("bytes")
	assert.Equal(t, size, getIntFromBytes(lenBytes))
}

func TestMergeList(t *testing.T) {
	a := []bool{false, false, true, false}
	b := []bool{true, false, false, true}
	mergeBoolLists(&a, b)
	assert.EqualValues(t, a, []bool{true, false, true, true})
}

func TestChunk(t *testing.T) {
	a := []byte{1, 2, 3, 4, 5}
	s := getChunk(a, 3)
	assert.Equal(t, len(s), 2)
	assert.EqualValues(t, s[0], []byte{1, 2, 3})
	assert.EqualValues(t, s[1], []byte{4, 5})
	s = getChunk(a, 5)
	assert.Equal(t, len(s), 1)
	assert.EqualValues(t, s[0], []byte{1, 2, 3, 4, 5})
	s = getChunk(a, 10)
	assert.Equal(t, len(s), 1)
	assert.EqualValues(t, s[0], []byte{1, 2, 3, 4, 5})

	a = []byte{}
	s = getChunk(a, 10)
	assert.Equal(t, len(s), 0)

}
