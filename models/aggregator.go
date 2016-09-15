package models

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

// Aggregator fetches MatchSums and derives aggregates.
type Aggregator interface {
	// Aggregate aggregates.
	Aggregate(req *apb.GetChampionRequest) (*apb.MatchAggregate, error)
}

// AggregatorImpl is an implementation of Aggregator.
type AggregatorImpl struct {
	CQL     *gocql.Session `inject:"t"`
	Vulgate Vulgate        `inject:"t"`
}

// Aggregate aggregates.
func (a *AggregatorImpl) Aggregate(req *apb.GetChampionRequest) (*apb.MatchAggregate, error) {
	quots := map[uint32]*apb.MatchQuotient{}
	for _, id := range a.Vulgate.GetChampionIDs() {
		quot, err := a.findChampionQuotient(req, id)
		if err != nil {
			return nil, err
		}
		quots[id] = quot
	}

	// now let us build the match aggregate
	return makeMatchAggregate(quots, req.ChampionId), nil
}

func (a *AggregatorImpl) findChampionQuotient(req *apb.GetChampionRequest, cid uint32) (*apb.MatchQuotient, error) {
	f := a.buildFilters(req, cid)
	return a.deriveQuotient(f)
}

// buildFilters builds a list of filters for a given champion.
func (a *AggregatorImpl) buildFilters(req *apb.GetChampionRequest, cid uint32) []*apb.MatchFilters {
	var ret []*apb.MatchFilters
	for _, patch := range a.Vulgate.FindPatches(req.Patch) {
		for _, tier := range a.Vulgate.FindTiers(req.Tier) {
			ret = append(ret, &apb.MatchFilters{
				ChampionId: int32(cid),
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

func (a *AggregatorImpl) deriveQuotient(filters []*apb.MatchFilters) (*apb.MatchQuotient, error) {
	sum, err := a.Sum(filters)
	if err != nil {
		return nil, fmt.Errorf("error fetching sum: %v", err)
	}
	return makeQuotient(sum), nil
}

// Sum derives a sum from a set of filters.
func (a *AggregatorImpl) Sum(filters []*apb.MatchFilters) (*apb.MatchSum, error) {
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
			sum = addMatchSums(sum, sumRow)
		}
	}()

	// Terminate when all sums are fetched
	wg.Wait()
	close(sumsChan)

	// Return sum and error
	return sum, fetchErr
}

func (a *AggregatorImpl) fetchSum(f *apb.MatchFilters) (*apb.MatchSum, error) {
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

func addMatchSums(a, b *apb.MatchSum) *apb.MatchSum {
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
		Masteries:   addStringSubscalarsMap(a.Masteries, b.Masteries),
		Runes:       addStringSubscalarsMap(a.Runes, b.Runes),
		Keystones:   addStringSubscalarsMap(a.Keystones, b.Keystones),
		Summoners:   addStringSubscalarsMap(a.Summoners, b.Summoners),
		Trinkets:    addUint32SubscalarsMap(a.Trinkets, b.Trinkets),
		SkillOrders: addStringSubscalarsMap(a.SkillOrders, b.SkillOrders),
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

	if p.Masteries == nil {
		p.Masteries = map[string]*apb.MatchSum_Subscalars{}
	}

	if p.Runes == nil {
		p.Runes = map[string]*apb.MatchSum_Subscalars{}
	}

	if p.Keystones == nil {
		p.Keystones = map[string]*apb.MatchSum_Subscalars{}
	}

	if p.Summoners == nil {
		p.Summoners = map[string]*apb.MatchSum_Subscalars{}
	}

	if p.Trinkets == nil {
		p.Trinkets = map[uint32]*apb.MatchSum_Subscalars{}
	}

	if p.SkillOrders == nil {
		p.SkillOrders = map[string]*apb.MatchSum_Subscalars{}
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