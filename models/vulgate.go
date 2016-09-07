package models

import apb "github.com/simplyianm/apollo/gen-go/asuna"

// Vulgate defines the interface to the Vulgate.
type Vulgate interface {
	// FindPatches finds all patches within a patch range, inclusive.
	FindPatches(rg *apb.PatchRange) []string

	// FindTiers finds all tiers within a tier range, inclusive.
	FindTiers(rg *apb.TierRange) []int32
}

// NewVulgate initializes the Vulgate.
func NewVulgate() (*vulgateImpl, error) {
	// TODO(pradyuman): implement
	return &vulgateImpl{}, nil
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
	// TODO(pradyuman): implement from vulgate
	return []string{rg.Min, rg.Max}
}

// FindTiers implements FindTiers.
func (v *vulgateImpl) FindTiers(rg *apb.TierRange) []int32 {
	if rg == nil {
		return []int32{}
	}
	// TODO(pradyuman): implement from vulgate
	return []int32{int32(rg.Min), int32(rg.Max)}
}
