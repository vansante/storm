package gob

import (
	"testing"

	"github.com/vansante/storm/v3/codec/internal"
)

func TestGob(t *testing.T) {
	internal.RoundtripTester(t, Codec)
}
