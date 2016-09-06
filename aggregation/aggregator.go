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

// GenerateAggregate generates an aggregate for a champion.
func (a *Aggregator) AggregateFromFilters(filters []*apb.MatchFilters) (*apb.MatchAggregate, error) {
	// Channel containing sums
	sumsChan := make(chan *apb.MatchSum)

	// Concurrently fetch all sums
	var wg sync.WaitGroup
	wg.Add(len(filters))

	// Iterate over all filters
	for _, filter := range filters {

		// Asynchronous get
		go func(filter *apb.MatchFilters) {
			sumsChan <- a.fetchSum(filter)
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

	return buildAggregate(sum)
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

func buildAggregate(sum *apb.MatchSum) (*apb.MatchAggregate, error) {
	return nil, nil
}
