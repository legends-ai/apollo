package models

import (
	"sort"

	apb "github.com/simplyianm/apollo/gen-go/asuna"
)

func makeMatchAggregate(quots map[uint32]*apb.MatchQuotient, id uint32) *apb.MatchAggregate {
	return nil
}

type groupedQuotients struct {
	scalars groupedScalarsQuotients
	deltas  groupedDeltasQuotients
}

type groupedScalarsQuotients struct {
	winRate                  []float64
	pickRate                 []float64
	banRate                  []float64
	gamesPlayed              []float64
	goldEarned               []float64
	kills                    []float64
	deaths                   []float64
	assists                  []float64
	damageDealt              []float64
	damageTaken              []float64
	minionsKilled            []float64
	teamJungleMinionsKilled  []float64
	enemyJungleMinionsKilled []float64
	structureDamage          []float64
	killingSpree             []float64
	wardsBought              []float64
	wardsPlaced              []float64
	wardsKilled              []float64
	crowdControl             []float64
	firstBlood               []float64
	firstBloodAssist         []float64
	doubleKills              []float64
	tripleKills              []float64
	quadrakills              []float64
	pentakills               []float64
}

type groupedDeltaQuotients struct {
	zeroToTen      []float64
	tenToTwenty    []float64
	twentyToThirty []float64
	thirtyToEnd    []float64
}

type groupedDeltasQuotients struct {
	csDiff          groupedDeltaQuotients
	xpDiff          groupedDeltaQuotients
	damageTakenDiff groupedDeltaQuotients
	xpPerMin        groupedDeltaQuotients
	goldPerMin      groupedDeltaQuotients
	towersPerMin    groupedDeltaQuotients
	wardsPlaced     groupedDeltaQuotients
	damageTaken     groupedDeltaQuotients
}

func makeMatchAggregateStatistics(quots map[uint32]*apb.MatchQuotient, id uint32) *apb.MatchAggregateStatistics {
	// grouped quotient aggregates
	var gs groupedQuotients
	self := quots[id]

	for _, quot := range quots {
		// Scalars
		gs.scalars.winRate = append(gs.scalars.winRate, quot.Scalars.Wins)
		gs.scalars.gamesPlayed = append(gs.scalars.gamesPlayed, quot.Scalars.Plays)
		gs.scalars.goldEarned = append(gs.scalars.goldEarned, quot.Scalars.GoldEarned)
		gs.scalars.kills = append(gs.scalars.kills, quot.Scalars.Kills)
		gs.scalars.deaths = append(gs.scalars.deaths, quot.Scalars.Deaths)
		gs.scalars.assists = append(gs.scalars.assists, quot.Scalars.Assists)
		gs.scalars.damageDealt = append(gs.scalars.damageDealt, quot.Scalars.DamageDealt)
		gs.scalars.damageTaken = append(gs.scalars.damageTaken, quot.Scalars.DamageTaken)
		gs.scalars.minionsKilled = append(gs.scalars.minionsKilled, quot.Scalars.MinionsKilled)
		gs.scalars.teamJungleMinionsKilled = append(gs.scalars.teamJungleMinionsKilled, quot.Scalars.TeamJungleMinionsKilled)
		gs.scalars.enemyJungleMinionsKilled = append(gs.scalars.enemyJungleMinionsKilled, quot.Scalars.EnemyJungleMinionsKilled)
		gs.scalars.structureDamage = append(gs.scalars.structureDamage, quot.Scalars.StructureDamage)
		gs.scalars.killingSpree = append(gs.scalars.killingSpree, quot.Scalars.KillingSpree)
		gs.scalars.wardsBought = append(gs.scalars.wardsBought, quot.Scalars.WardsBought)
		gs.scalars.wardsPlaced = append(gs.scalars.wardsPlaced, quot.Scalars.WardsPlaced)
		gs.scalars.wardsKilled = append(gs.scalars.wardsKilled, quot.Scalars.WardsKilled)
		gs.scalars.crowdControl = append(gs.scalars.crowdControl, quot.Scalars.CrowdControl)
		gs.scalars.firstBlood = append(gs.scalars.firstBlood, quot.Scalars.FirstBlood)
		gs.scalars.firstBloodAssist = append(gs.scalars.firstBloodAssist, quot.Scalars.FirstBloodAssist)
		gs.scalars.doubleKills = append(gs.scalars.doubleKills, quot.Scalars.Doublekills)
		gs.scalars.tripleKills = append(gs.scalars.tripleKills, quot.Scalars.Triplekills)
		gs.scalars.quadrakills = append(gs.scalars.quadrakills, quot.Scalars.Quadrakills)
		gs.scalars.pentakills = append(gs.scalars.pentakills, quot.Scalars.Pentakills)

		// Deltas
		gs.deltas.csDiff = appendDeltas(gs.deltas.csDiff, quot.Deltas.CsDiff)
		gs.deltas.xpDiff = appendDeltas(gs.deltas.xpDiff, quot.Deltas.XpDiff)
		gs.deltas.damageTakenDiff = appendDeltas(gs.deltas.damageTakenDiff, quot.Deltas.DamageTakenDiff)
		gs.deltas.xpPerMin = appendDeltas(gs.deltas.xpPerMin, quot.Deltas.XpPerMin)
		gs.deltas.goldPerMin = appendDeltas(gs.deltas.goldPerMin, quot.Deltas.GoldPerMin)
		gs.deltas.towersPerMin = appendDeltas(gs.deltas.towersPerMin, quot.Deltas.TowersPerMin)
		gs.deltas.wardsPlaced = appendDeltas(gs.deltas.wardsPlaced, quot.Deltas.WardsPlaced)
		gs.deltas.damageTaken = appendDeltas(gs.deltas.damageTaken, quot.Deltas.DamageTaken)
	}

	return &apb.MatchAggregateStatistics{
		Scalars: &apb.MatchAggregateStatistics_Scalars{
			WinRate:                  deriveStatistic(gs.scalars.winRate, self.Scalars.Wins),
			GamesPlayed:              deriveStatistic(gs.scalars.gamesPlayed, self.Scalars.Plays),
			GoldEarned:               deriveStatistic(gs.scalars.goldEarned, self.Scalars.GoldEarned),
			Kills:                    deriveStatistic(gs.scalars.kills, self.Scalars.Kills),
			Deaths:                   deriveStatistic(gs.scalars.deaths, self.Scalars.Deaths),
			Assists:                  deriveStatistic(gs.scalars.assists, self.Scalars.Assists),
			DamageDealt:              deriveStatistic(gs.scalars.damageDealt, self.Scalars.DamageDealt),
			MinionsKilled:            deriveStatistic(gs.scalars.minionsKilled, self.Scalars.MinionsKilled),
			TeamJungleMinionsKilled:  deriveStatistic(gs.scalars.teamJungleMinionsKilled, self.Scalars.TeamJungleMinionsKilled),
			EnemyJungleMinionsKilled: deriveStatistic(gs.scalars.enemyJungleMinionsKilled, self.Scalars.EnemyJungleMinionsKilled),
			StructureDamage:          deriveStatistic(gs.scalars.structureDamage, self.Scalars.StructureDamage),
			KillingSpree:             deriveStatistic(gs.scalars.killingSpree, self.Scalars.KillingSpree),
			WardsBought:              deriveStatistic(gs.scalars.wardsBought, self.Scalars.WardsBought),
			WardsPlaced:              deriveStatistic(gs.scalars.wardsPlaced, self.Scalars.WardsPlaced),
			WardsKilled:              deriveStatistic(gs.scalars.wardsKilled, self.Scalars.WardsKilled),
			CrowdControl:             deriveStatistic(gs.scalars.crowdControl, self.Scalars.CrowdControl),
			FirstBlood:               deriveStatistic(gs.scalars.firstBlood, self.Scalars.FirstBlood),
			FirstBloodAssist:         deriveStatistic(gs.scalars.firstBloodAssist, self.Scalars.FirstBloodAssist),
			DoubleKills:              deriveStatistic(gs.scalars.doubleKills, self.Scalars.Doublekills),
			TripleKills:              deriveStatistic(gs.scalars.tripleKills, self.Scalars.Triplekills),
			Quadrakills:              deriveStatistic(gs.scalars.quadrakills, self.Scalars.Quadrakills),
			Pentakills:               deriveStatistic(gs.scalars.pentakills, self.Scalars.Pentakills),
		},

		Deltas: &apb.MatchAggregateStatistics_Deltas{
			CsDiff:          deriveDeltaStatistic(gs.deltas.csDiff, self.Deltas.CsDiff),
			XpDiff:          deriveDeltaStatistic(gs.deltas.xpDiff, self.Deltas.XpDiff),
			DamageTakenDiff: deriveDeltaStatistic(gs.deltas.damageTakenDiff, self.Deltas.DamageTakenDiff),
			XpPerMin:        deriveDeltaStatistic(gs.deltas.xpPerMin, self.Deltas.XpPerMin),
			GoldPerMin:      deriveDeltaStatistic(gs.deltas.goldPerMin, self.Deltas.GoldPerMin),
			TowersPerMin:    deriveDeltaStatistic(gs.deltas.towersPerMin, self.Deltas.TowersPerMin),
			WardsPlaced:     deriveDeltaStatistic(gs.deltas.wardsPlaced, self.Deltas.WardsPlaced),
			DamageTaken:     deriveDeltaStatistic(gs.deltas.damageTaken, self.Deltas.DamageTaken),
		},
	}
}

func deriveDeltaStatistic(dqs groupedDeltaQuotients, ds *apb.MatchQuotient_Deltas_Delta) *apb.MatchAggregateStatistics_Deltas_Delta {
	return &apb.MatchAggregateStatistics_Deltas_Delta{
		ZeroToTen:      deriveStatistic(dqs.zeroToTen, ds.ZeroToTen),
		TenToTwenty:    deriveStatistic(dqs.tenToTwenty, ds.TenToTwenty),
		TwentyToThirty: deriveStatistic(dqs.twentyToThirty, ds.TwentyToThirty),
		ThirtyToEnd:    deriveStatistic(dqs.thirtyToEnd, ds.ThirtyToEnd),
	}
}

func deriveStatistic(vals []float64, val float64) *apb.MatchAggregateStatistics_Statistic {
	var sum float64
	for _, v := range vals {
		sum += v
	}
	avg := sum / float64(len(vals))

	// sort desc so we can get the rank
	sort.Sort(sort.Reverse(sort.Float64Slice(vals)))
	var rank int
	for idx, v := range vals {
		if v == val {
			rank = idx + 1
			break
		}
	}

	percentile := 1.0 - float64(rank)/float64(len(vals))

	// TODO(igm): implement change
	return &apb.MatchAggregateStatistics_Statistic{
		Rank:       uint32(rank),
		Value:      val,
		Average:    avg,
		Percentile: percentile,
	}
}

func appendDeltas(dqs groupedDeltaQuotients, ds *apb.MatchQuotient_Deltas_Delta) groupedDeltaQuotients {
	return groupedDeltaQuotients{
		zeroToTen:      append(dqs.zeroToTen, ds.ZeroToTen),
		tenToTwenty:    append(dqs.tenToTwenty, ds.TenToTwenty),
		twentyToThirty: append(dqs.twentyToThirty, ds.TwentyToThirty),
		thirtyToEnd:    append(dqs.thirtyToEnd, ds.ThirtyToEnd),
	}
}
