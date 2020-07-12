package silk

import (
	"io/ioutil"
	"testing"

	"gotest.tools/assert"
)

func TestDecodeSilkBuffToWave(t *testing.T) {
	b, err := ioutil.ReadFile("../test.silk")
	assert.NilError(t, err)
	dst, err := DecodeSilkBuffToWave(b, 8000)
	assert.NilError(t, err)
	err = ioutil.WriteFile("../test.wav", dst, 0666)
	assert.NilError(t, err)
}
