package models

import (
	"fmt"

	apb "github.com/asunaio/apollo/gen-go/asuna"
)

// Aggregator fetches MatchSums and derives aggregates.
type Aggregator interface {
	// Aggregate aggregates.
	Aggregate(req *apb.GetChampionRequest) (*apb.MatchAggregate, error)
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
func (a *aggregatorImpl) Aggregate(req *apb.GetChampionRequest) (*apb.MatchAggregate, error) {
	// TODO(igm): optimize these calls. We have quite a few duplicates.
	champions, err := a.findChampionQuotients(req)
	if err != nil {
		return nil, fmt.Errorf("error finding champion quotients: %v", err)
	}

	roles, err := a.findRoleQuotients(req)
	if err != nil {
		return nil, fmt.Errorf("error finding role quotients: %v", err)
	}

	patches, err := a.findPatchQuotients(req)
	if err != nil {
		return nil, fmt.Errorf("error finding patch quotients: %v", err)
	}

	// now let us build the match aggregate
	return a.Deriver.Derive(req.Role, champions, roles, patches, req.ChampionId)
}

func (a *aggregatorImpl) findChampionQuotients(req *apb.GetChampionRequest) (map[uint32]*apb.MatchQuotient, error) {
	champions := map[uint32]*apb.MatchQuotient{}
	for _, id := range a.Vulgate.GetChampionIDs() {
		copy := *req
		copy.ChampionId = id
		f := a.buildFilters(&copy)
		champ, err := a.deriveQuotient(f)
		if err != nil {
			return nil, err
		}

		// no nil champs
		if champ == nil {
			continue
		}

		champions[id] = champ
	}
	return champions, nil
}

func (a *aggregatorImpl) findRoleQuotients(req *apb.GetChampionRequest) (map[apb.Role]*apb.MatchQuotient, error) {
	roles := map[apb.Role]*apb.MatchQuotient{}

	for _, role := range []apb.Role{
		apb.Role_TOP,
		apb.Role_JUNGLE,
		apb.Role_MID,
		apb.Role_BOT,
		apb.Role_SUPPORT,
	} {
		copy := *req
		copy.Role = role
		f := a.buildFilters(&copy)
		champ, err := a.deriveQuotient(f)
		if err != nil {
			return nil, err
		}

		// no nil champs
		if champ == nil {
			continue
		}

		roles[role] = champ
	}

	return roles, nil
}

func (a *aggregatorImpl) findPatchQuotients(
	req *apb.GetChampionRequest,
) (map[string]map[uint32]*apb.MatchQuotient, error) {
	patches := map[string]map[uint32]*apb.MatchQuotient{}

	// build patch filters
	for _, patch := range a.Vulgate.FindNPreviousPatches(req.Patch, 5) {

		// build champion filters
		champions := map[uint32]*apb.MatchQuotient{}
		for _, id := range a.Vulgate.GetChampionIDs() {

			// Construct filters
			var f []*apb.MatchFilters
			for _, tier := range a.Vulgate.FindTiers(req.Tier) {
				f = append(f, &apb.MatchFilters{
					ChampionId: int32(id),
					EnemyId:    ANY_ENEMY,
					Patch:      patch,
					Tier:       tier,
					Region:     req.Region,
					Role:       req.Role,
				})
			}

			champ, err := a.deriveQuotient(f)
			if err != nil {
				return nil, err
			}

			// no nil champs
			if champ == nil {
				continue
			}

			champions[id] = champ
		}

		patches[patch] = champions
	}

	return patches, nil
}

// buildFilters builds a list of filters for a given champion.
func (a *aggregatorImpl) buildFilters(req *apb.GetChampionRequest) []*apb.MatchFilters {
	var ret []*apb.MatchFilters
	for _, patch := range a.Vulgate.FindPatches(req.Patch) {
		for _, tier := range a.Vulgate.FindTiers(req.Tier) {
			ret = append(ret, &apb.MatchFilters{
				ChampionId: int32(req.ChampionId),
				EnemyId:    ANY_ENEMY,
				Patch:      patch,
				Tier:       tier,
				Region:     req.Region,
				Role:       req.Role,
			})
		}
	}
	return ret
}

func (a *aggregatorImpl) deriveQuotient(filters []*apb.MatchFilters) (*apb.MatchQuotient, error) {
	sum, err := a.MatchSumDAO.Sum(filters)
	if err != nil {
		return nil, fmt.Errorf("error fetching sum: %v", err)
	}
	if sum == nil {
		return nil, nil
	}
	return makeQuotient(sum), nil
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
		Masteries:   makeQuotientSubscalarStringMap(sum.Masteries, plays),
		Runes:       makeQuotientSubscalarStringMap(sum.Runes, plays),
		Keystones:   makeQuotientSubscalarStringMap(sum.Keystones, plays),
		Summoners:   makeQuotientSubscalarStringMap(sum.Summoners, plays),
		Trinkets:    makeQuotientSubscalarUint32Map(sum.Trinkets, plays),
		SkillOrders: makeQuotientSubscalarStringMap(sum.SkillOrders, plays),
		Durations:   makeQuotientSubscalarUint32Map(sum.Durations, plays),
		Bans:        makeQuotientSubscalarUint32Map(sum.Bans, plays),
		Allies:      makeQuotientSubscalarUint32Map(sum.Allies, plays),
		Enemies:     makeQuotientSubscalarUint32Map(sum.Enemies, plays),
	}
}

// makeQuotientDeltas calculates deltas
func makeQuotientDeltas(deltas *apb.MatchSum_Deltas_Delta, dd *apb.MatchSum_DurationDistribution) *apb.MatchQuotient_Deltas_Delta {
	return &apb.MatchQuotient_Deltas_Delta{
		ZeroToTen:      float64(deltas.ZeroToTen) / float64(dd.ZeroToTen),
		TenToTwenty:    float64(deltas.TenToTwenty) / float64(dd.TenToTwenty),
		TwentyToThirty: float64(deltas.TwentyToThirty) / float64(dd.TwentyToThirty),
		ThirtyToEnd:    float64(deltas.ThirtyToEnd) / float64(dd.ThirtyToEnd),
	}
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
