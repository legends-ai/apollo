package models

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/golang/protobuf/proto"

	apb "github.com/asunaio/apollo/gen-go/asuna"
)

const (
	stmtGetSum = `SELECT match_sum
		FROM athena.match_sums
		WHERE
			champion_id = ? AND enemy_id = ? AND patch = ? AND
			tier = ? AND region = ? AND role = ?`
)

type MatchSumDAO interface {
	// Get gets a MatchSum from MatchFilters.
	Get(f *apb.MatchFilters) (*apb.MatchSum, error)

	// Sum sums MatchSums derived from the given filters.
	Sum(filters []*apb.MatchFilters) (*apb.MatchSum, error)
}

// NewMatchSumDAO constructs a new MatchSumDAO.
func NewMatchSumDAO() MatchSumDAO {
	return &matchSumDAO{}
}

type matchSumDAO struct {
	CQL *gocql.Session `inject:"t"`
}

func (a *matchSumDAO) Get(f *apb.MatchFilters) (*apb.MatchSum, error) {
	var rawSum []byte
	if err := a.CQL.Query(
		stmtGetSum, f.ChampionId, f.EnemyId, f.Patch,
		f.Tier, int32(f.Region), int32(f.Role),
	).Scan(&rawSum); err != nil {
		if err == gocql.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching sum from Cassandra: %v", err)
	}

	var sum apb.MatchSum
	if err := proto.Unmarshal(rawSum, &sum); err != nil {
		return nil, fmt.Errorf("error unmarshaling sum: %v", err)
	}

	return &sum, nil
}

// Sum derives a sum from a set of filters.
func (a *matchSumDAO) Sum(filters []*apb.MatchFilters) (*apb.MatchSum, error) {
	// Create aggregate sum
	sum := (*apb.MatchSum)(nil)

	// Iterate over all filters
	for _, filter := range filters {
		// Error handling
		s, err := a.Get(filter)
		if err != nil {
			return nil, err
		}

		// Process sum
		if s != nil {
			normalizeMatchSum(s)
			if sum == nil {
				sum = s
			} else {
				sum = addMatchSums(sum, s)
			}
		}
	}

	// Return sum and error
	return sum, nil
}

func addMatchSums(a, b *apb.MatchSum) *apb.MatchSum {
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
		Durations: addUint32SubscalarsMap(a.Durations, b.Durations),
		Bans:      addUint32SubscalarsMap(a.Bans, b.Bans),
		Allies:    addUint32SubscalarsMap(a.Allies, b.Allies),
		Enemies:   addUint32SubscalarsMap(a.Enemies, b.Enemies),
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
	if p.Durations == nil {
		p.Durations = map[uint32]*apb.MatchSum_Subscalars{}
	}
	if p.Bans == nil {
		p.Bans = map[uint32]*apb.MatchSum_Subscalars{}
	}
	if p.Allies == nil {
		p.Allies = map[uint32]*apb.MatchSum_Subscalars{}
	}
	if p.Enemies == nil {
		p.Enemies = map[uint32]*apb.MatchSum_Subscalars{}
	}
}
