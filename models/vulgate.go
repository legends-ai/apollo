package models

import (
	"io/ioutil"

	"github.com/golang/protobuf/proto"

	apb "github.com/simplyianm/apollo/gen-go/asuna"
)

const (
	TierChallenger = "CHALLENGER"
	TierMaster     = "MASTER"
	TierDiamond    = "DIAMOND"
	TierPlatinum   = "PLATINUM"
	TierGold       = "GOLD"
	TierSilver     = "SILVER"
	TierBronze     = "BRONZE"
)

// Vulgate defines the interface to the Vulgate.
type Vulgate interface {
	// FindPatches finds all patches within a patch range, inclusive.
	FindPatches(rg *apb.PatchRange) []string

	// FindTiers finds all tiers within a tier range, inclusive.
	FindTiers(rg *apb.TierRange) []int32

	// GetChampionInfo gets information about a champion.
	GetChampionInfo(id uint32) *apb.Vulgate_Champion

	// GetPatchTimes gets times for a patch.
	GetPatchTimes(rg *apb.PatchRange) *apb.Vulgate_PatchTime

	// GetChampionIDs gets a list of champion ids.
	GetChampionIDs() []uint32
}

// NewVulgate initializes the Vulgate.
func NewVulgate() (*vulgateImpl, error) {
	// TODO(igm): somehow make this not rely on binary location
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

	var tiers []int32
	for _, tier := range v.proto.Tiers {
		v := parseTier(tier)
		if rg.Min <= v && rg.Max >= v {
			tiers = append(tiers, int32(v))
		}
	}

	return tiers
}

func parseTier(s string) uint32 {
	switch s {
	case TierChallenger:
		return 0x70
	case TierMaster:
		return 0x60
	case TierDiamond:
		return 0x50
	case TierPlatinum:
		return 0x40
	case TierGold:
		return 0x30
	case TierSilver:
		return 0x20
	case TierBronze:
		return 0x10
	default:
		return 0
	}
}

func (v *vulgateImpl) GetChampionInfo(id uint32) *apb.Vulgate_Champion {
	// TODO(pradyuman): implement
	return nil
}

func (v *vulgateImpl) GetPatchTimes(rg *apb.PatchRange) *apb.Vulgate_PatchTime {
	// TODO(pradyuman): implement
	return nil
}

func (v *vulgateImpl) GetChampionIDs() []uint32 {
	// TODO(pradyuman): implement
	return nil
}
