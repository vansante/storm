package index

import (
	"reflect"
	"testing"
)

func Test_generatePrefix(t *testing.T) {
	tests := []struct {
		name  string
		value []byte
		want  []byte
	}{
		{
			"normal",
			[]byte("hello"),
			[]byte("hello__"),
		},
		{
			"",
			[]byte(""),
			[]byte("__"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generatePrefix(tt.value)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generatePrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
