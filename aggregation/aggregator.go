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
	CQL *gocql.Session
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
	sanitizeMatchSum(a)
	sanitizeMatchSum(b)
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
	}
}

func sanitizeMatchSum(p *apb.MatchSum) {
	if p.Scalars == nil {
		p.Scalars = &apb.MatchSum_Scalars{}
	}

	if p.Deltas == nil {
		p.Deltas = &apb.MatchSum_Deltas{}
	}
}

func buildAggregate(base *apb.MatchSum, filtered *apb.MatchSum) *apb.MatchAggregate {
	// TODO(igm): implement
	return &apb.MatchAggregate{
		Statistics: &apb.MatchAggregateStatistics{
			WinRate:                  makeStatistic(float64(filtered.Wins) / filtered.Plays),
			GamesPlayed:              makeStatistic(float64(filtered.Plays)),
			GoldEarned:               makeStatistic(float64(filtered.GoldEarned) / filtered.Plays),
			Kills:                    makeStatistic(float64(filtered.Kills) / filtered.Plays),
			Deaths:                   makeStatistic(float64(filtered.Deaths) / filtered.Plays),
			Assists:                  makeStatistic(float64(filtered.Assists) / filtered.Plays),
			DamageDealt:              makeStatistic(float64(filtered.DamageDealt) / filtered.Plays),
			DamageTaken:              makeStatistic(float64(filtered.DamageTaken) / filtered.Plays),
			MinionsKilled:            makeStatistic(float64(filtered.MinionsKilled) / filtered.Plays),
			TeamJungleMinionsKilled:  makeStatistic(float64(filtered.TeamJungleMinionsKilled) / filtered.Plays),
			EnemyJungleMinionsKilled: makeStatistic(float64(filtered.EnemyJungleMinionsKilled) / filtered.Plays),
			StructureDamage:          makeStatistic(float64(filtered.StructureDamage) / filtered.Plays),
			KillingSpree:             makeStatistic(float64(filtered.KillingSpree) / filtered.Plays),
			WardsBought:              makeStatistic(float64(filtered.WardsBought) / filtered.Plays),
			WardsPlaced:              makeStatistic(float64(filtered.WardsPlaced) / filtered.Plays),
			CrowdControl:             makeStatistic(float64(filtered.CrowdControl) / filtered.Plays),
			FirstBlood:               makeStatistic(float64(filtered.FirstBlood) / filtered.Plays),
			FirstBloodAssist:         makeStatistic(float64(filtered.FirstBloodAssist) / filtered.Plays),
			Doubleills:               makeStatistic(float64(filtered.DoubleKills) / filtered.Plays),
			Triplekills:              makeStatistic(float64(filtered.TripleKills) / filtered.Plays),
			Quadrakills:              makeStatistic(float64(filtered.Quadrakills) / filtered.Plays),
			Pentakills:               makeStatistic(float64(filtered.Pentakills) / filtered.Plays),
		},
	}
}

func makeStatistic(val float64) *apb.MatchAggregateStatistics_Statistic {
	// TODO(igm): implement rank, change, average, percentile. Cache?
	return &apb.MatchAggregateStatistics_Statistic{
		Value: val,
	}
}
