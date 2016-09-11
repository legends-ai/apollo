package models

import (
	"io/ioutil"

	"github.com/golang/protobuf/proto"

	apb "github.com/simplyianm/apollo/gen-go/asuna"
)

// Vulgate defines the interface to the Vulgate.
type Vulgate interface {
	// FindPatches finds all patches within a patch range, inclusive.
	FindPatches(rg *apb.PatchRange) []string

	// FindTiers finds all tiers within a tier range, inclusive.
	FindTiers(rg *apb.TierRange) []int32
}

// NewVulgate initializes the Vulgate.
func NewVulgate() (*vulgateImpl, error) {
	raw, err := ioutil.ReadFile("./vulgate/vulgate.textproto")
	if err != nil {
		return nil, err
	}

	vpb := &apb.Vulgate{}
	proto.UnmarshalText(string(raw), vpb)

	return &vulgateImpl{proto: vpb}, nil
}

// VulgateImpl is the implementation of Vulgate.
type vulgateImpl struct {
	proto *apb.Vulgate
}

// FindPatches implements FindPatches.
func (v *vulgateImpl) FindPatches(rg *apb.PatchRange) []string {
	if rg == nil {
		return []string{}
	}

	var start, end int = -1, -1
	for i, patch := range v.proto.Patches {
		if start == -1 && patch == rg.Min {
			start = i
		} else if end == -1 && patch == rg.Max {
			end = i + 1
			break
		}
	}

	if start == -1 || end == -1 {
		return []string{}
	}

	return v.proto.Patches[start:end]
}

// FindTiers implements FindTiers.
func (v *vulgateImpl) FindTiers(rg *apb.TierRange) []int32 {
	if rg == nil {
		return []int32{}
	}

	tiers := []int32{}
	var start, end int = -1, -1
	for i, _ := range v.proto.Tiers {
		if start == -1 && i+1 != int(rg.Min) {
			continue
		}

		start = i
		tiers = append(tiers, int32(i+1))
		if end == -1 && i+1 == int(rg.Max) {
			end = i + 1
			break
		}
	}

	if start == -1 || end == -1 {
		return []int32{}
	}

	return tiers
}
