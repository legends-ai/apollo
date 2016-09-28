package models

import (
	"fmt"

	apb "github.com/asunaio/apollo/gen-go/asuna"
)

// Aggregator fetches MatchSums and derives aggregates.
type Aggregator interface {
	// Aggregate aggregates.
	Aggregate(
		championId uint32,
		enemyChampionId int32,
		patch *apb.PatchRange,
		tier *apb.TierRange,
		region apb.Region,
		role apb.Role,
	) (*apb.MatchAggregate, error)
}

// NewAggregator constructs a new Aggregator.
func NewAggregator() Aggregator {
	return &aggregatorImpl{}
}

// aggregatorImpl is an implementation of Aggregator.
type aggregatorImpl struct {
	MatchSumDAO MatchSumDAO `inject:"t"`
	Deriver     Deriver     `inject:"t"`
	Vulgate     Vulgate     `inject:"t"`
}

// Aggregate aggregates.
func (a *aggregatorImpl) Aggregate(
	aChampionId uint32,
	enemyChampionId int32,
	aPatch *apb.PatchRange,
	aTier *apb.TierRange,
	aRegion apb.Region,
	aRole apb.Role,
) (*apb.MatchAggregate, error) {
	champs, err := a.MatchSumDAO.SumsOfChampions(
		aPatch, enemyChampionId, aTier, aRegion, aRole,
	)
	if err != nil {
		return nil, fmt.Errorf("error finding champion sums: %v", err)
	}

	rolesSums, err := a.MatchSumDAO.SumsOfRoles(
		aPatch.Max, aChampionId, enemyChampionId, aTier, aRegion,
	)
	if err != nil {
		return nil, fmt.Errorf("error finding role sums: %v", err)
	}

	patches := map[string]map[uint32]*apb.MatchQuotient{}
	for id, champPatches := range champs {
		for patch, sum := range champPatches {
			if patches[patch] == nil {
				patches[patch] = map[uint32]*apb.MatchQuotient{}
			}
			if sum == nil {
				continue
			}
			patches[patch][id] = makeQuotient(sum)
		}
	}

	champions := map[uint32]*apb.MatchQuotient{}
	for _, id := range a.Vulgate.GetChampionIDs() {
		sum := &apb.MatchSum{}
		normalizeMatchSum(sum)
		for _, patch := range a.Vulgate.FindPatches(aPatch) {
			// Use existing fetched patches
			patchSum := champs[id][patch]

			// Retrieve patch if it does not exist
			if patchSum == nil {
				patchSum, err = a.MatchSumDAO.SumOfPatch(
					patch, aChampionId, enemyChampionId, aTier, aRegion, aRole,
				)
				if err != nil {
					return nil, err
				}
				if patchSum == nil {
					continue
				}
			}

			// Append sum
			sum = addMatchSums(sum, patchSum)
		}
		if sum == nil {
			continue
		}
		champions[id] = makeQuotient(sum)
	}

	roles := map[apb.Role]*apb.MatchQuotient{}
	for role, sum := range rolesSums {
		if sum == nil {
			continue
		}
		roles[role] = makeQuotient(sum)
	}

	// now let us build the match aggregate
	return a.Deriver.Derive(aRole, champions, roles, patches, aChampionId)
}

func addDelta(a *apb.MatchSum_Deltas_Delta, b *apb.MatchSum_Deltas_Delta) *apb.MatchSum_Deltas_Delta {
	return &apb.MatchSum_Deltas_Delta{
		ZeroToTen:      a.ZeroToTen + b.ZeroToTen,
		TenToTwenty:    a.TenToTwenty + b.TenToTwenty,
		TwentyToThirty: a.TwentyToThirty + b.TwentyToThirty,
		ThirtyToEnd:    a.ThirtyToEnd + b.ThirtyToEnd,
	}
}

func addStringSubscalarsMap(
	a, b map[string]*apb.MatchSum_Subscalars,
) map[string]*apb.MatchSum_Subscalars {
	for i, bv := range b {
		if av, exists := a[i]; exists {
			a[i] = &apb.MatchSum_Subscalars{
				Plays: av.Plays + bv.Plays,
				Wins:  av.Wins + bv.Wins,
			}
		} else {
			a[i] = bv
		}
	}
	return a
}

func addUint32SubscalarsMap(
	a, b map[uint32]*apb.MatchSum_Subscalars,
) map[uint32]*apb.MatchSum_Subscalars {
	for i, bv := range b {
		if av, exists := a[i]; exists {
			a[i] = &apb.MatchSum_Subscalars{
				Plays: av.Plays + bv.Plays,
				Wins:  av.Wins + bv.Wins,
			}
		} else {
			a[i] = bv
		}
	}
	return a
}

// makeQuotient creates a MatchQuotient from a MatchSum.
func makeQuotient(sum *apb.MatchSum) *apb.MatchQuotient {
	scalars := sum.Scalars
	plays := float64(scalars.Plays)
	dd := sum.DurationDistribution

	return &apb.MatchQuotient{
		Scalars: &apb.MatchQuotient_Scalars{
			Plays:                    plays,
			Wins:                     float64(scalars.Wins) / plays,
			GoldEarned:               float64(scalars.GoldEarned) / plays,
			Kills:                    float64(scalars.Kills) / plays,
			Deaths:                   float64(scalars.Deaths) / plays,
			Assists:                  float64(scalars.Assists) / plays,
			DamageDealt:              float64(scalars.DamageDealt) / plays,
			DamageTaken:              float64(scalars.DamageTaken) / plays,
			MinionsKilled:            float64(scalars.MinionsKilled) / plays,
			TeamJungleMinionsKilled:  float64(scalars.TeamJungleMinionsKilled) / plays,
			EnemyJungleMinionsKilled: float64(scalars.EnemyJungleMinionsKilled) / plays,
			StructureDamage:          float64(scalars.StructureDamage) / plays,
			KillingSpree:             float64(scalars.KillingSpree) / plays,
			WardsBought:              float64(scalars.WardsBought) / plays,
			WardsPlaced:              float64(scalars.WardsPlaced) / plays,
			CrowdControl:             float64(scalars.CrowdControl) / plays,
			FirstBlood:               float64(scalars.FirstBlood) / plays,
			FirstBloodAssist:         float64(scalars.FirstBloodAssist) / plays,
			Doublekills:              float64(scalars.Doublekills) / plays,
			Triplekills:              float64(scalars.Triplekills) / plays,
			Quadrakills:              float64(scalars.Quadrakills) / plays,
			Pentakills:               float64(scalars.Pentakills) / plays,
		},
		Deltas: &apb.MatchQuotient_Deltas{
			CsDiff:          makeQuotientDeltas(sum.Deltas.CsDiff, dd),
			XpDiff:          makeQuotientDeltas(sum.Deltas.XpDiff, dd),
			DamageTakenDiff: makeQuotientDeltas(sum.Deltas.DamageTakenDiff, dd),
			XpPerMin:        makeQuotientDeltas(sum.Deltas.XpPerMin, dd),
			GoldPerMin:      makeQuotientDeltas(sum.Deltas.GoldPerMin, dd),
			TowersPerMin:    makeQuotientDeltas(sum.Deltas.TowersPerMin, dd),
			WardsPlaced:     makeQuotientDeltas(sum.Deltas.WardsPlaced, dd),
			DamageTaken:     makeQuotientDeltas(sum.Deltas.DamageTaken, dd),
		},
		Masteries:     makeQuotientSubscalarStringMap(sum.Masteries, plays),
		Runes:         makeQuotientSubscalarStringMap(sum.Runes, plays),
		Keystones:     makeQuotientSubscalarStringMap(sum.Keystones, plays),
		Summoners:     makeQuotientSubscalarStringMap(sum.Summoners, plays),
		Trinkets:      makeQuotientSubscalarUint32Map(sum.Trinkets, plays),
		SkillOrders:   makeQuotientSubscalarStringMap(sum.SkillOrders, plays),
		Durations:     makeQuotientSubscalarUint32Map(sum.Durations, plays),
		Bans:          makeQuotientSubscalarUint32Map(sum.Bans, plays),
		Allies:        makeQuotientSubscalarUint32Map(sum.Allies, plays),
		Enemies:       makeQuotientSubscalarUint32Map(sum.Enemies, plays),
		StarterItems:  makeQuotientSubscalarStringMap(sum.StarterItems, plays),
		BuildPath:     makeQuotientSubscalarStringMap(sum.BuildPath, plays),
		CoreBuildList: makeQuotientSubscalarStringMap(sum.BuildPath, plays), // TODO(igm): implement
	}
}

// makeQuotientDeltas calculates deltas
func makeQuotientDeltas(deltas *apb.MatchSum_Deltas_Delta, dd *apb.MatchSum_DurationDistribution) *apb.MatchQuotient_Deltas_Delta {
	delta := &apb.MatchQuotient_Deltas_Delta{}
	if dd.ZeroToTen == 0 {
		delta.ZeroToTen = 0
	} else {
		delta.ZeroToTen = float64(deltas.ZeroToTen) / float64(dd.ZeroToTen)
	}
	if dd.TenToTwenty == 0 {
		delta.TenToTwenty = 0
	} else {
		delta.TenToTwenty = float64(deltas.TenToTwenty) / float64(dd.TenToTwenty)
	}
	if dd.TwentyToThirty == 0 {
		delta.TwentyToThirty = 0
	} else {
		delta.TwentyToThirty = float64(deltas.TwentyToThirty) / float64(dd.TwentyToThirty)
	}
	if dd.ThirtyToEnd == 0 {
		delta.ThirtyToEnd = 0
	} else {
		delta.ThirtyToEnd = float64(deltas.ThirtyToEnd) / float64(dd.ThirtyToEnd)
	}
	return delta
}

// makeQuotientSubscalars calculates subscalars
func makeQuotientSubscalars(ss *apb.MatchSum_Subscalars, plays float64) *apb.MatchQuotient_Subscalars {
	return &apb.MatchQuotient_Subscalars{
		Plays:     float64(ss.Plays) / plays,
		Wins:      float64(ss.Wins) / float64(ss.Plays),
		PlayCount: ss.Plays,
	}
}

func makeQuotientSubscalarStringMap(ss map[string]*apb.MatchSum_Subscalars, plays float64) map[string]*apb.MatchQuotient_Subscalars {
	// TODO(igm): truncate lower values
	ret := map[string]*apb.MatchQuotient_Subscalars{}
	for key, s := range ss {
		ret[key] = makeQuotientSubscalars(s, plays)
	}
	return ret
}

func makeQuotientSubscalarUint32Map(ss map[uint32]*apb.MatchSum_Subscalars, plays float64) map[uint32]*apb.MatchQuotient_Subscalars {
	// TODO(igm): truncate lower values
	ret := map[uint32]*apb.MatchQuotient_Subscalars{}
	for key, s := range ss {
		ret[key] = makeQuotientSubscalars(s, plays)
	}
	return ret
}
