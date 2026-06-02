package protobuf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vansante/storm/v3"
	"github.com/vansante/storm/v3/codec/internal"
)

func TestProtobuf(t *testing.T) {
	u1 := SimpleUser{Id: 1, Name: "John"}
	u2 := SimpleUser{}
	internal.RoundtripTester(t, Codec, &u1, &u2)
	require.True(t, u1.Id == u2.Id)
}

func TestSave(t *testing.T) {
	dir, _ := os.MkdirTemp(os.TempDir(), "storm")
	defer func() { _ = os.RemoveAll(dir) }()
	db, _ := storm.Open(filepath.Join(dir, "storm.db"), storm.Codec(Codec))
	u1 := SimpleUser{Id: 1, Name: "John"}
	err := db.Save(&u1)
	require.NoError(t, err)
	u2 := SimpleUser{}
	err = db.One("Id", uint64(1), &u2)
	require.NoError(t, err)
	require.Equal(t, u2.Name, u1.Name)
}

func TestGetSet(t *testing.T) {
	dir, _ := os.MkdirTemp(os.TempDir(), "storm")
	defer func() { _ = os.RemoveAll(dir) }()
	db, _ := storm.Open(filepath.Join(dir, "storm.db"), storm.Codec(Codec))
	err := db.Set("bucket", "key", "value")
	require.NoError(t, err)
	var s string
	err = db.Get("bucket", "key", &s)
	require.NoError(t, err)
	require.Equal(t, "value", s)
}
