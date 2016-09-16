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

func TestDeserializeSummoners(t *testing.T) {
	for _, test := range []struct {
		Description string
		In          string
		Want1       uint32
		Want2       uint32
	}{
		{
			Description: "Normal summoners",
			In:          "123|456",
			Want1:       123,
			Want2:       456,
		},
	} {
		got1, got2, _ := deserializeSummoners(test.In)
		if got1 != test.Want1 || got2 != test.Want2 {
			t.Errorf("Error with test %q: got (%v, %v), want (%v, %v)", test.Description, got1, got2, test.Want1, test.Want2)
		}
	}
}
