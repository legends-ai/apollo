package aggregation

import (
	"fmt"
	"sync"

	"github.com/gocql/gocql"
	"github.com/golang/protobuf/proto"

	apb "github.com/simplyianm/apollo/gen-go/asuna"
)

const (
	stmtGetSum = `SELECT match_sum
		FROM athena.matches
		WHERE
			champion_id = ? AND enemy_id = ? AND patch = ? AND
			tier = ? AND region = ? AND role = ?`
)

// Aggregator aggregates sums and derives aggregates.
type Aggregator struct {
	CQL *gocql.Session `inject:"t"`
}

// Aggregate derives a MatchAggregate from two sets of filters:
// - Base, the base stats to compare against
// - Filters, the filters to compare the specifics against
func (a *Aggregator) Aggregate(
	base []*apb.MatchFilters, filters []*apb.MatchFilters,
) (*apb.MatchAggregate, error) {
	baseSum, err := a.Sum(base)
	if err != nil {
		return nil, fmt.Errorf("error fetching base sum: %v", err)
	}

	filtersSum, err := a.Sum(filters)
	if err != nil {
		return nil, fmt.Errorf("error fetching filters sum: %v", err)
	}

	aggregate := buildAggregate(baseSum, filtersSum)
	return aggregate, nil
}

// Sum derives a sum from a set of filters.
func (a *Aggregator) Sum(filters []*apb.MatchFilters) (*apb.MatchSum, error) {
	// Channel containing sums
	sumsChan := make(chan *apb.MatchSum)

	// Error from fetching aggregates
	var fetchErr error = nil

	// Concurrently fetch all sums
	var wg sync.WaitGroup
	wg.Add(len(filters))

	// Iterate over all filters
	for _, filter := range filters {

		// Asynchronous get
		go func(filter *apb.MatchFilters) {
			// Error handling
			s, err := a.fetchSum(filter)
			if err != nil {
				fetchErr = err
			}

			// Process sum
			sumsChan <- s
			wg.Done()
		}(filter)
	}

	// Create aggregate sum
	sum := &apb.MatchSum{}
	go func() {
		for sumRow := range sumsChan {
			sum = addSums(sum, sumRow)
		}
	}()

	// Terminate when all sums are fetched
	wg.Wait()
	close(sumsChan)

	// Return sum and error
	return sum, fetchErr
}

func (a *Aggregator) fetchSum(f *apb.MatchFilters) (*apb.MatchSum, error) {
	var rawSum []byte
	if err := a.CQL.Query(
		stmtGetSum, f.ChampionId, f.EnemyId, f.Patch,
		f.Tier, int32(f.Region), int32(f.Role),
	).Scan(&rawSum); err != nil {
		return nil, fmt.Errorf("error fetching sum from Cassandra: %v", err)
	}

	var sum apb.MatchSum
	if err := proto.Unmarshal(rawSum, &sum); err != nil {
		return nil, fmt.Errorf("error unmarshaling sum: %v", err)
	}

	return &sum, nil
}

func addSums(a, b *apb.MatchSum) *apb.MatchSum {
	normalizeMatchSum(a)
	normalizeMatchSum(b)
	return &apb.MatchSum{
		Scalars: &apb.MatchSum_Scalars{
			Plays:                    a.Scalars.Plays + b.Scalars.Plays,
			Wins:                     a.Scalars.Wins + b.Scalars.Wins,
			GoldEarned:               a.Scalars.GoldEarned + b.Scalars.GoldEarned,
			Kills:                    a.Scalars.Kills + b.Scalars.Kills,
			Deaths:                   a.Scalars.Deaths + b.Scalars.Deaths,
			Assists:                  a.Scalars.Assists + b.Scalars.Assists,
			DamageDealt:              a.Scalars.DamageDealt + b.Scalars.DamageDealt,
			MinionsKilled:            a.Scalars.MinionsKilled + b.Scalars.MinionsKilled,
			TeamJungleMinionsKilled:  a.Scalars.TeamJungleMinionsKilled + b.Scalars.TeamJungleMinionsKilled,
			EnemyJungleMinionsKilled: a.Scalars.EnemyJungleMinionsKilled + b.Scalars.EnemyJungleMinionsKilled,
			StructureDamage:          a.Scalars.StructureDamage + b.Scalars.StructureDamage,
			KillingSpree:             a.Scalars.KillingSpree + b.Scalars.KillingSpree,
			WardsBought:              a.Scalars.WardsBought + b.Scalars.WardsBought,
			WardsPlaced:              a.Scalars.WardsPlaced + b.Scalars.WardsPlaced,
			WardsKilled:              a.Scalars.WardsKilled + b.Scalars.WardsKilled,
			CrowdControl:             a.Scalars.CrowdControl + b.Scalars.CrowdControl,
			FirstBlood:               a.Scalars.FirstBlood + b.Scalars.FirstBlood,
			FirstBloodAssist:         a.Scalars.FirstBloodAssist + b.Scalars.FirstBloodAssist,
			Doublekills:              a.Scalars.Doublekills + b.Scalars.Doublekills,
			Triplekills:              a.Scalars.Triplekills + b.Scalars.Triplekills,
			Quadrakills:              a.Scalars.Quadrakills + b.Scalars.Quadrakills,
			Pentakills:               a.Scalars.Pentakills + b.Scalars.Pentakills,
		},
		Deltas: &apb.MatchSum_Deltas{
			CsDiff:          addDelta(a.Deltas.CsDiff, b.Deltas.CsDiff),
			XpDiff:          addDelta(a.Deltas.XpDiff, b.Deltas.XpDiff),
			DamageTakenDiff: addDelta(a.Deltas.DamageTakenDiff, b.Deltas.DamageTakenDiff),
			XpPerMin:        addDelta(a.Deltas.XpPerMin, b.Deltas.XpPerMin),
			GoldPerMin:      addDelta(a.Deltas.GoldPerMin, b.Deltas.GoldPerMin),
			TowersPerMin:    addDelta(a.Deltas.TowersPerMin, b.Deltas.TowersPerMin),
			WardsPlaced:     addDelta(a.Deltas.WardsPlaced, b.Deltas.WardsPlaced),
			DamageTaken:     addDelta(a.Deltas.DamageTaken, b.Deltas.DamageTaken),
		},
		DurationDistribution: &apb.MatchSum_DurationDistribution{
			ZeroToTen:      a.DurationDistribution.ZeroToTen + b.DurationDistribution.ZeroToTen,
			TenToTwenty:    a.DurationDistribution.TenToTwenty + b.DurationDistribution.TenToTwenty,
			TwentyToThirty: a.DurationDistribution.TwentyToThirty + b.DurationDistribution.TwentyToThirty,
			ThirtyToEnd:    a.DurationDistribution.ThirtyToEnd + b.DurationDistribution.ThirtyToEnd,
		},
	}
}

func normalizeMatchSum(p *apb.MatchSum) {
	if p.Scalars == nil {
		p.Scalars = &apb.MatchSum_Scalars{}
	}

	if p.Deltas == nil {
		p.Deltas = &apb.MatchSum_Deltas{}
	}

	if p.Deltas.CsDiff == nil {
		p.Deltas.CsDiff = &apb.MatchSum_Deltas_Delta{}
	}

	if p.Deltas.XpDiff == nil {
		p.Deltas.XpDiff = &apb.MatchSum_Deltas_Delta{}
	}

	if p.Deltas.DamageTakenDiff == nil {
		p.Deltas.DamageTakenDiff = &apb.MatchSum_Deltas_Delta{}
	}

	if p.Deltas.XpPerMin == nil {
		p.Deltas.XpPerMin = &apb.MatchSum_Deltas_Delta{}
	}

	if p.Deltas.GoldPerMin == nil {
		p.Deltas.GoldPerMin = &apb.MatchSum_Deltas_Delta{}
	}

	if p.Deltas.TowersPerMin == nil {
		p.Deltas.TowersPerMin = &apb.MatchSum_Deltas_Delta{}
	}

	if p.Deltas.WardsPlaced == nil {
		p.Deltas.WardsPlaced = &apb.MatchSum_Deltas_Delta{}
	}

	if p.Deltas.DamageTaken == nil {
		p.Deltas.DamageTaken = &apb.MatchSum_Deltas_Delta{}
	}

	if p.DurationDistribution == nil {
		p.DurationDistribution = &apb.MatchSum_DurationDistribution{}
	}
}

func addDelta(a *apb.MatchSum_Deltas_Delta, b *apb.MatchSum_Deltas_Delta) *apb.MatchSum_Deltas_Delta {
	return &apb.MatchSum_Deltas_Delta{
		ZeroToTen:      a.ZeroToTen + b.ZeroToTen,
		TenToTwenty:    a.TenToTwenty + b.TenToTwenty,
		TwentyToThirty: a.TwentyToThirty + b.TwentyToThirty,
		ThirtyToEnd:    a.ThirtyToEnd + b.ThirtyToEnd,
	}
}

func buildAggregate(base *apb.MatchSum, filtered *apb.MatchSum) *apb.MatchAggregate {
	// TODO(igm): implement
	scalars := filtered.Scalars
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
