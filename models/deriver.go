package models

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	apb "github.com/asunaio/apollo/gen-go/asuna"
)

// Deriver derives MatchAggregates from quotients.
// This does not need its own class, but there's so much code that this makes sense.
type Deriver interface {
	// Derive derives a MatchAggregate from a map of MatchQuotients and a champion id.
	// - Champions is a map of all champions to their match quotient for the current role
	// - Roles is a map of all roles of the current champion
	Derive(
		champions map[uint32]*apb.MatchQuotient, roles map[apb.Role]*apb.MatchQuotient, id uint32,
	) (*apb.MatchAggregate, error)
}

// NewDeriver constructs a new Deriver.
func NewDeriver() Deriver {
	return &deriverImpl{}
}

type deriverImpl struct{}

func (d *deriverImpl) Derive(
	champions map[uint32]*apb.MatchQuotient, roles map[apb.Role]*apb.MatchQuotient, id uint32,
) (*apb.MatchAggregate, error) {
	// precondition -- champ must exist
	if champions[id] == nil {
		return nil, fmt.Errorf("champion %d does not exist in quotient map", id)
	}

	collections, err := makeMatchAggregateCollections(champions[id])
	if err != nil {
		return nil, fmt.Errorf("error parsing collections: %v", err)
	}

	return &apb.MatchAggregate{
		Statistics:  makeMatchAggregateStatistics(champions, id),
		Graphs:      makeMatchAggregateGraphs(champions[id]),
		Collections: collections,
	}, nil
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

// deserializeBonusSet works for masteries, runes, and keystones.
func deserializeBonusSet(s string) (map[uint32]uint32, error) {
	ret := map[uint32]uint32{}

	// no runes
	if len(s) == 0 {
		return ret, nil
	}

	// get rune counts
	rs := strings.Split(s, "|")
	for _, r := range rs {
		id, ct, err := deserializeBonusSetElement(r)
		if err != nil {
			return nil, err
		}
		// assign
		ret[uint32(id)] = uint32(ct)
	}
	return ret, nil
}

func deserializeBonusSetElement(s string) (uint32, uint32, error) {
	// No element
	if s == "" {
		return 0, 0, nil
	}

	// get data of rune count
	ps := strings.Split(s, ":")

	// rune id
	id, err := strconv.Atoi(ps[0])

	// rune count
	ct, err := strconv.Atoi(ps[2])

	// check for strconv errors
	if err != nil {
		return 0, 0, err
	}

	return uint32(id), uint32(ct), nil
}

// deserializeSummoners deserializes the summoners key
func deserializeSummoners(s string) (uint32, uint32, error) {
	rs := strings.Split(s, "|")
	a, err := strconv.Atoi(rs[0])
	b, err := strconv.Atoi(rs[1])
	if err != nil {
		return 0, 0, err
	}
	return uint32(a), uint32(b), nil
}

// deserializeSkillOrder converts a skill order string to a list of abilities.
func deserializeSkillOrder(s string) ([]apb.Ability, error) {
	var ret []apb.Ability
	for _, r := range s {
		switch r {
		case 'Q':
			ret = append(ret, apb.Ability_Q)
		case 'W':
			ret = append(ret, apb.Ability_W)
		case 'E':
			ret = append(ret, apb.Ability_E)
		case 'R':
			ret = append(ret, apb.Ability_R)
		default:
			return nil, fmt.Errorf("Unknown skill: %q", r)
		}
	}
	return ret, nil
}
func makeMatchAggregateGraphs(quot *apb.MatchQuotient) *apb.MatchAggregateGraphs {
	return &apb.MatchAggregateGraphs{}
}

func makeMatchAggregateCollections(quot *apb.MatchQuotient) (*apb.MatchAggregateCollections, error) {
	// derive runes
	var runes []*apb.MatchAggregateCollections_RuneSet
	for rs, rstats := range quot.Runes {
		// rs is rune set string
		// rstats is rune set subscalars
		runeSet, err := deserializeBonusSet(rs)
		if err != nil {
			return nil, fmt.Errorf("could not deserialize rune set: %v", err)
		}
		runes = append(runes, &apb.MatchAggregateCollections_RuneSet{
			Runes:      runeSet,
			PickRate:   rstats.Plays,
			WinRate:    rstats.Wins,
			NumMatches: uint32(rstats.PlayCount),
		})
	}

	// derive masteries
	var masteries []*apb.MatchAggregateCollections_MasterySet
	for ms, mstats := range quot.Masteries {
		// ms is mastery set string
		// mstats is mastery set subscalars
		masterySet, err := deserializeBonusSet(ms)
		if err != nil {
			return nil, fmt.Errorf("could not deserialize mastery set: %v", err)
		}
		masteries = append(masteries, &apb.MatchAggregateCollections_MasterySet{
			Masteries:  masterySet,
			PickRate:   mstats.Plays,
			WinRate:    mstats.Wins,
			NumMatches: uint32(mstats.PlayCount),
		})
	}

	// derive keystones
	var keystones []*apb.MatchAggregateCollections_Keystone
	for ks, kstats := range quot.Keystones {
		// ks is keystone string
		// kstats is keystone subscalars
		keystone, ct, err := deserializeBonusSetElement(ks)
		if err != nil {
			return nil, fmt.Errorf("could not deserialize keystone: %v", err)
		}
		if ct == 0 {
			// check for nil keystone
			continue
		}
		keystones = append(keystones, &apb.MatchAggregateCollections_Keystone{
			Keystone:   keystone,
			PickRate:   kstats.Plays,
			WinRate:    kstats.Wins,
			NumMatches: uint32(kstats.PlayCount),
		})
	}

	// derive summoners
	var summonerSpells []*apb.MatchAggregateCollections_SummonerSet
	for ss, sstats := range quot.Summoners {
		// ss is summoner string
		// sstats is summoner subscalars
		spell1, spell2, err := deserializeSummoners(ss)
		if err != nil {
			return nil, fmt.Errorf("could not deserialize summoners: %v", err)
		}
		summonerSpells = append(summonerSpells, &apb.MatchAggregateCollections_SummonerSet{
			Spell1:     spell1,
			Spell2:     spell2,
			PickRate:   sstats.Plays,
			WinRate:    sstats.Wins,
			NumMatches: uint32(sstats.PlayCount),
		})
	}

	// derive trinkets
	var trinkets []*apb.MatchAggregateCollections_Trinket
	for trinket, tstats := range quot.Trinkets {
		// tstats is trinket subscalars
		trinkets = append(trinkets, &apb.MatchAggregateCollections_Trinket{
			Trinket:    trinket,
			PickRate:   tstats.Plays,
			WinRate:    tstats.Wins,
			NumMatches: uint32(tstats.PlayCount),
		})
	}

	// derive skill orders
	var skillOrders []*apb.MatchAggregateCollections_SkillOrder
	for sos, sostats := range quot.SkillOrders {
		so, err := deserializeSkillOrder(sos)
		if err != nil {
			return nil, fmt.Errorf("could not deserialize skill order: %v", err)
		}
		// sostats is skill order subscalars
		skillOrders = append(skillOrders, &apb.MatchAggregateCollections_SkillOrder{
			SkillOrder: so,
			PickRate:   sostats.Plays,
			WinRate:    sostats.Wins,
			NumMatches: uint32(sostats.PlayCount),
		})
	}

	// TODO(igm ^ pradyuman): builds

	return &apb.MatchAggregateCollections{
		Runes:          runes,
		Masteries:      masteries,
		Keystones:      keystones,
		SummonerSpells: summonerSpells,
		Trinkets:       trinkets,
		SkillOrders:    skillOrders,
	}, nil
}
