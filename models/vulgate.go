package models

import apb "github.com/simplyianm/apollo/gen-go/asuna"

// Vulgate defines the interface to the Vulgate.
type Vulgate interface {
	// FindPatches finds all patches within a patch range, inclusive.
	FindPatches(rg *apb.PatchRange) []string

	// FindTiers finds all tiers within a tier range, inclusive.
	FindTiers(rg *apb.TierRange) []uint32
}

// NewVulgate initializes the Vulgate.
func NewVulgate() (Vulgate, error) {
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
func (v *vulgateImpl) FindTiers(rg *apb.TierRange) []uint32 {
	if rg == nil {
		return []uint32{}
	}
	// TODO(pradyuman): implement from vulgate
	return []uint32{rg.Min, rg.Max}
}
