package models

import (
	"fmt"
	"sync"

	"github.com/gocql/gocql"
	"github.com/golang/protobuf/proto"

	apb "github.com/asunaio/apollo/gen-go/asuna"
)

const (
	stmtGetSum = `SELECT match_sum
		FROM athena_out.match_sums
		WHERE
			champion_id = ? AND enemy_id = ? AND patch = ? AND
			tier = ? AND region = ? AND role = ?`
)

const (
	ANY_CHAMPION = -1

	// Number of previous patches to fetch.
	prevPatches = 5
)

type MatchSumDAO interface {
	// Get gets a MatchSum from MatchFilters.
	Get(f *apb.MatchFilters) (*apb.MatchSum, error)

	// Sum sums MatchSums derived from the given filters.
	Sum(filters []*apb.MatchFilters) (*apb.MatchSum, error)

	// SumsOfChampions gets the sums of champions per patch.
	SumsOfChampions(
		patchRange *apb.PatchRange, enemy int32,
		tiers *apb.TierRange, region apb.Region, role apb.Role,
	) (map[uint32]map[string]*apb.MatchSum, error)

	// SumsOfPatches gets the sums of a champion for a range of patches.
	SumsOfPatches(
		patchRange *apb.PatchRange, champion uint32, enemy int32,
		tiers *apb.TierRange, region apb.Region, role apb.Role,
	) (map[string]*apb.MatchSum, error)

	// SumOfPatch gets the sum of a champion for a patch.
	SumOfPatch(
		patch string, champion uint32, enemy int32,
		tiers *apb.TierRange, region apb.Region, role apb.Role,
	) (*apb.MatchSum, error)

	// SumsOfRoles gets the sums of a champion per role for a patch.
	SumsOfRoles(
		patch string, champion uint32, enemy int32,
		tiers *apb.TierRange, region apb.Region,
	) (map[apb.Role]*apb.MatchSum, error)
}

// NewMatchSumDAO constructs a new MatchSumDAO.
func NewMatchSumDAO() MatchSumDAO {
	return &matchSumDAO{}
}

type matchSumDAO struct {
	CQL     *gocql.Session `inject:"t"`
	Vulgate Vulgate        `inject:"t"`
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

	var wg sync.WaitGroup
	sums := make(chan *apb.MatchSum)

	var fetchErr error

	// Iterate over all filters
	for _, filter := range filters {
		wg.Add(1)
		go func(filter *apb.MatchFilters) {
			defer wg.Done()

			s, err := a.Get(filter)
			if err != nil {
				fetchErr = err
				return
			}
			if s == nil {
				return
			}
			sums <- s
		}(filter)
	}

	// Close when all are done
	go func() {
		wg.Wait()
		close(sums)
	}()

	for s := range sums {
		normalizeMatchSum(s)
		if sum == nil {
			sum = s
		} else {
			sum = addMatchSums(sum, s)
		}
	}

	if fetchErr != nil {
		return nil, fetchErr
	}

	// Return sum and error
	return sum, nil
}

func (m *matchSumDAO) SumsOfChampions(
	patchRange *apb.PatchRange, enemy int32,
	tiers *apb.TierRange, region apb.Region, role apb.Role,
) (map[uint32]map[string]*apb.MatchSum, error) {
	ret := map[uint32]map[string]*apb.MatchSum{}
	for _, id := range m.Vulgate.GetChampionIDs() {
		patches, err := m.SumsOfPatches(patchRange, id, enemy, tiers, region, role)
		if err != nil {
			return nil, err
		}
		ret[id] = patches
	}
	return ret, nil
}

func (m *matchSumDAO) SumsOfPatches(
	patchRange *apb.PatchRange, champion uint32, enemy int32,
	tiers *apb.TierRange, region apb.Region, role apb.Role,
) (map[string]*apb.MatchSum, error) {
	ret := map[string]*apb.MatchSum{}
	// TODO(igm): make prev patches configurable
	for _, patch := range m.Vulgate.FindNPreviousPatches(patchRange, prevPatches) {
		sum, err := m.SumOfPatch(patch, champion, enemy, tiers, region, role)
		if err != nil {
			return nil, err
		}
		ret[patch] = sum
	}
	return ret, nil
}

func (m *matchSumDAO) SumOfPatch(
	patch string, champion uint32, enemy int32,
	tiers *apb.TierRange, region apb.Region, role apb.Role,
) (*apb.MatchSum, error) {
	// TODO(igm): cache
	var filters []*apb.MatchFilters
	for _, tier := range m.Vulgate.FindTiers(tiers) {
		filters = append(filters, &apb.MatchFilters{
			ChampionId: int32(champion),
			EnemyId:    enemy,
			Patch:      patch,
			Tier:       tier,
			Region:     region,
			Role:       role,
		})
	}
	return m.Sum(filters)
}

func (m *matchSumDAO) SumsOfRoles(
	patch string, champion uint32, enemy int32,
	tiers *apb.TierRange, region apb.Region,
) (map[apb.Role]*apb.MatchSum, error) {
	ret := map[apb.Role]*apb.MatchSum{}
	for _, role := range []apb.Role{
		apb.Role_TOP,
		apb.Role_JUNGLE,
		apb.Role_MID,
		apb.Role_BOT,
		apb.Role_SUPPORT,
	} {
		sum, err := m.SumOfPatch(patch, champion, enemy, tiers, region, role)
		if err != nil {
			return nil, err
		}
		ret[role] = sum
	}
	return ret, nil
}

func addManyMatchSums(sums ...*apb.MatchSum) *apb.MatchSum {
	acc := sums[0]
	for _, sum := range sums[1:] {
		acc = addMatchSums(acc, sum)
	}
	return acc
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
		Durations:    addUint32SubscalarsMap(a.Durations, b.Durations),
		Bans:         addUint32SubscalarsMap(a.Bans, b.Bans),
		Allies:       addUint32SubscalarsMap(a.Allies, b.Allies),
		Enemies:      addUint32SubscalarsMap(a.Enemies, b.Enemies),
		StarterItems: addStringSubscalarsMap(a.StarterItems, b.StarterItems),
		BuildPath:    addStringSubscalarsMap(a.BuildPath, b.BuildPath),
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
	if p.StarterItems == nil {
		p.StarterItems = map[string]*apb.MatchSum_Subscalars{}
	}
	if p.BuildPath == nil {
		p.BuildPath = map[string]*apb.MatchSum_Subscalars{}
	}
}
