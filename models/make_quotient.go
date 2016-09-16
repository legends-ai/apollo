package models

import apb "github.com/asunaio/apollo/gen-go/asuna"

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
		Plays: float64(ss.Plays) / plays,
		Wins:  float64(ss.Wins) / float64(ss.Plays),
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
