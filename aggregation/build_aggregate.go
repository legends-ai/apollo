package aggregation

import (
	apb "github.com/simplyianm/apollo/gen-go/asuna"
)

func buildAggregate(sum *apb.MatchSum) *apb.MatchAggregate {
	// TODO(igm): implement
	scalars := sum.Scalars
	return &apb.MatchAggregate{
		Statistics: &apb.MatchAggregateStatistics{
			Scalars: &apb.MatchAggregateStatistics_Scalars{
				WinRate:                  makeStatistic(float64(scalars.Wins) / float64(scalars.Plays)),
				GamesPlayed:              makeStatistic(float64(scalars.Plays)),
				GoldEarned:               makeStatistic(float64(scalars.GoldEarned) / float64(scalars.Plays)),
				Kills:                    makeStatistic(float64(scalars.Kills) / float64(scalars.Plays)),
				Deaths:                   makeStatistic(float64(scalars.Deaths) / float64(scalars.Plays)),
				Assists:                  makeStatistic(float64(scalars.Assists) / float64(scalars.Plays)),
				DamageDealt:              makeStatistic(float64(scalars.DamageDealt) / float64(scalars.Plays)),
				DamageTaken:              makeStatistic(float64(scalars.DamageTaken) / float64(scalars.Plays)),
				MinionsKilled:            makeStatistic(float64(scalars.MinionsKilled) / float64(scalars.Plays)),
				TeamJungleMinionsKilled:  makeStatistic(float64(scalars.TeamJungleMinionsKilled) / float64(scalars.Plays)),
				EnemyJungleMinionsKilled: makeStatistic(float64(scalars.EnemyJungleMinionsKilled) / float64(scalars.Plays)),
				StructureDamage:          makeStatistic(float64(scalars.StructureDamage) / float64(scalars.Plays)),
				KillingSpree:             makeStatistic(float64(scalars.KillingSpree) / float64(scalars.Plays)),
				WardsBought:              makeStatistic(float64(scalars.WardsBought) / float64(scalars.Plays)),
				WardsPlaced:              makeStatistic(float64(scalars.WardsPlaced) / float64(scalars.Plays)),
				CrowdControl:             makeStatistic(float64(scalars.CrowdControl) / float64(scalars.Plays)),
				FirstBlood:               makeStatistic(float64(scalars.FirstBlood) / float64(scalars.Plays)),
				FirstBloodAssist:         makeStatistic(float64(scalars.FirstBloodAssist) / float64(scalars.Plays)),
				DoubleKills:              makeStatistic(float64(scalars.Doublekills) / float64(scalars.Plays)),
				TripleKills:              makeStatistic(float64(scalars.Triplekills) / float64(scalars.Plays)),
				Quadrakills:              makeStatistic(float64(scalars.Quadrakills) / float64(scalars.Plays)),
				Pentakills:               makeStatistic(float64(scalars.Pentakills) / float64(scalars.Plays)),
			},
		},
	}
}

func makeStatistic(val float64) *apb.MatchAggregateStatistics_Statistic {
	// TODO(igm): implement rank, change, average, percentile. Cache?
	return &apb.MatchAggregateStatistics_Statistic{
		Value: val,
	}
}
