package models

import (
	"reflect"
	"testing"
)

func TestDeserializeBonusSet(t *testing.T) {
	for _, test := range []struct {
		Description string
		In          string
		Want        map[uint32]uint32
	}{
		{
			Description: "One rune",
			In:          "123:1:4",
			Want: map[uint32]uint32{
				123: 4,
			},
		},
	} {
		got, _ := deserializeBonusSet(test.In)
		if !reflect.DeepEqual(got, test.Want) {
			t.Errorf("Error with test %q: got %v, want %v", test.Description, got, test.Want)
		}
	}
}
