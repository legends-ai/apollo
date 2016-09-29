package models

import (
	"reflect"
	"testing"

	apb "github.com/asunaio/apollo/gen-go/asuna"
)

func TestFindPatches(t *testing.T) {
	vimpl := &vulgateImpl{
		proto: &apb.Vulgate{
			Patches: []string{"6.13", "6.14", "6.15", "6.16", "6.17", "6.18"},
		},
	}

	for _, test := range []struct {
		Description string
		PatchRange  *apb.PatchRange
		Want        []string
	}{
		{
			Description: "Two patches",
			PatchRange: &apb.PatchRange{
				Min: "6.17",
				Max: "6.18",
			},
			Want: []string{"6.17", "6.18"},
		},
		{
			Description: "Six patches",
			PatchRange: &apb.PatchRange{
				Min: "6.13",
				Max: "6.18",
			},
			Want: []string{"6.13", "6.14", "6.15", "6.16", "6.17", "6.18"},
		},
		{
			Description: "One patch",
			PatchRange: &apb.PatchRange{
				Min: "6.18",
				Max: "6.18",
			},
			Want: []string{"6.18"},
		},
	} {
		patches := vimpl.FindPatches(test.PatchRange)
		if !reflect.DeepEqual(patches, test.Want) {
			t.Errorf("[%v] Got %v - Want %v", test.Description, patches, test.Want)
		}
	}
}

func TestFindNPreviousPatches(t *testing.T) {
	vimpl := &vulgateImpl{
		proto: &apb.Vulgate{
			Patches: []string{"6.13", "6.14", "6.15", "6.16", "6.17", "6.18"},
		},
	}

	for _, test := range []struct {
		Description string
		PatchRange  *apb.PatchRange
		Min         int
		Want        []string
	}{
		{
			Description: "[vulgate:FindNPreviousPatches] Minimum of N patches",
			PatchRange: &apb.PatchRange{
				Min: "6.17",
				Max: "6.18",
			},
			Min:  5,
			Want: []string{"6.14", "6.15", "6.16", "6.17", "6.18"},
		},
		{
			Description: "[vulgate:FindNPreviousPatches] Full Range if > N",
			PatchRange: &apb.PatchRange{
				Min: "6.13",
				Max: "6.18",
			},
			Min:  5,
			Want: []string{"6.13", "6.14", "6.15", "6.16", "6.17", "6.18"},
		},
		{
			Description: "[vulgate:FindNPreviousPatches] One Patch (Min = Max)",
			PatchRange: &apb.PatchRange{
				Min: "6.18",
				Max: "6.18",
			},
			Min:  5,
			Want: []string{"6.14", "6.15", "6.16", "6.17", "6.18"},
		},
	} {
		patches := vimpl.FindNPreviousPatches(test.PatchRange, test.Min)
		if !reflect.DeepEqual(patches, test.Want) {
			t.Errorf("[%v] Got %v - Want %v", test.Description, patches, test.Want)
		}
	}
}
