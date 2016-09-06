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
	// TODO(igm): add sums
	return a
}

func buildAggregate(base *apb.MatchSum, filtered *apb.MatchSum) *apb.MatchAggregate {
	// TODO(igm): implement
	return &apb.MatchAggregate{}
}
