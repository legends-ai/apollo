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
		role apb.Role,
		champions map[uint32]*apb.MatchQuotient,
		roles map[apb.Role]*apb.MatchQuotient,
		patches map[string]map[uint32]*apb.MatchQuotient,
		id uint32,
		minPlayRate float64,
	) (*apb.MatchAggregate, error)
}

// NewDeriver constructs a new Deriver.
func NewDeriver() Deriver {
	return &deriverImpl{}
}

type deriverImpl struct{}

func (d *deriverImpl) Derive(
	role apb.Role,
	champions map[uint32]*apb.MatchQuotient,
	roles map[apb.Role]*apb.MatchQuotient,
	patches map[string]map[uint32]*apb.MatchQuotient,
	id uint32,
	minPlayRate float64,
) (*apb.MatchAggregate, error) {
	// precondition -- champ must exist
	if champions[id] == nil {
		return nil, fmt.Errorf("champion %d does not exist in quotient map", id)
	}

	collections, err := makeMatchAggregateCollections(champions[id], minPlayRate)
	if err != nil {
		return nil, fmt.Errorf("error parsing collections: %v", err)
	}

	return &apb.MatchAggregate{
		Role:        makeMatchAggregateRoles(champions, roles, role, id),
		Statistics:  makeMatchAggregateStatistics(champions, id),
		Graphs:      makeMatchAggregateGraphs(champions, patches, id),
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
	selfPick := calculatePickRate(quots, id)
	selfBan := calculateBanRate(quots, id)

	for cid, quot := range quots {
		// Scalars
		gs.scalars.winRate = append(gs.scalars.winRate, quot.Scalars.Wins)
		// TODO(igm): optimize this
		gs.scalars.pickRate = append(gs.scalars.pickRate, calculatePickRate(quots, cid))
		gs.scalars.banRate = append(gs.scalars.banRate, calculateBanRate(quots, cid))
		gs.scalars.gamesPlayed = append(gs.scalars.gamesPlayed, float64(quot.Scalars.Plays))
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
			PickRate:                 deriveStatistic(gs.scalars.pickRate, selfPick),
			BanRate:                  deriveStatistic(gs.scalars.pickRate, selfBan),
			GamesPlayed:              deriveStatistic(gs.scalars.gamesPlayed, float64(self.Scalars.Plays)),
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

func makeMatchAggregateRoles(
	champions map[uint32]*apb.MatchQuotient,
	roles map[apb.Role]*apb.MatchQuotient,
	role apb.Role, id uint32,
) *apb.MatchAggregateRoles {
	var total uint32
	for _, champ := range champions {
		if champ.Scalars.Plays != 0 {
			total++
		}
	}

	totalForChamp := 0.0
	for _, roleQuotient := range roles {
		totalForChamp += float64(roleQuotient.Scalars.Plays)
	}

	var roleStats []*apb.MatchAggregateRoles_RoleStats
	for role, roleQuotient := range roles {
		roleStats = append(roleStats, &apb.MatchAggregateRoles_RoleStats{
			Role:       role,
			PickRate:   float64(roleQuotient.Scalars.Plays) / float64(totalForChamp),
			NumMatches: uint32(roleQuotient.Scalars.Plays),
		})
	}

	return &apb.MatchAggregateRoles{
		Role:                 role,
		TotalChampionsInRole: total,
		RoleStats:            roleStats,
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
			return nil, fmt.Errorf("unknown skill: %q", r)
		}
	}
	return ret, nil
}

func deserializeBuild(s string) ([]uint32, error) {
	var ret []uint32
	if s == "" {
		return ret, nil
	}
	for _, el := range strings.Split(s, "|") {
		eli, err := strconv.ParseUint(el, 10, 0)
		if err != nil {
			return nil, fmt.Errorf("unknown item: %q", el)
		}
		ret = append(ret, uint32(eli))
	}
	return ret, nil
}

func makeMatchAggregateGraphs(
	champions map[uint32]*apb.MatchQuotient,
	patches map[string]map[uint32]*apb.MatchQuotient, id uint32,
) *apb.MatchAggregateGraphs {
	winRate := map[uint32]float64{}
	pickRate := map[uint32]float64{}
	banRate := map[uint32]float64{}

	for cid, champ := range champions {
		winRate[cid] = champ.Scalars.Wins
		pickRate[cid] = calculatePickRate(champions, cid)
		banRate[cid] = calculateBanRate(champions, cid)
	}

	distribution := &apb.MatchAggregateGraphs_Distribution{
		WinRate:  winRate,
		PickRate: pickRate,
		BanRate:  banRate,
	}

	var byPatch []*apb.MatchAggregateGraphs_ByPatch
	for patch, championsOfPatch := range patches {
		self := championsOfPatch[id]
		winRate := 0.0
		if self != nil {
			winRate = self.Scalars.Wins
		}
		byPatch = append(byPatch, &apb.MatchAggregateGraphs_ByPatch{
			Patch:    patch,
			WinRate:  winRate,
			PickRate: calculatePickRate(championsOfPatch, id),
			BanRate:  calculateBanRate(championsOfPatch, id),
		})
	}

	quot := champions[id]
	var byGameLength []*apb.MatchAggregateGraphs_ByGameLength
	for duration, stats := range quot.Durations {
		byGameLength = append(byGameLength, &apb.MatchAggregateGraphs_ByGameLength{
			GameLength: &apb.IntRange{
				Min: duration,
				Max: duration,
			},
			WinRate: stats.Wins,
		})
	}

	return &apb.MatchAggregateGraphs{
		Distribution:   distribution,
		ByPatch:        byPatch,
		ByGameLength:   byGameLength,
		PhysicalDamage: quot.Scalars.PhysicalDamage,
		MagicDamage:    quot.Scalars.MagicDamage,
		TrueDamage:     quot.Scalars.TrueDamage,
	}
}

func makeMatchAggregateCollections(quot *apb.MatchQuotient, minPlayRate float64) (*apb.MatchAggregateCollections, error) {
	// derive runes
	var runes []*apb.MatchAggregateCollections_RuneSet
	for rs, rstats := range quot.Runes {
		if rstats.Plays < minPlayRate {
			continue
		}

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
		if mstats.Plays < minPlayRate {
			continue
		}

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
		if kstats.Plays < minPlayRate {
			continue
		}

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
		if sstats.Plays < minPlayRate {
			continue
		}

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
		if tstats.Plays < minPlayRate {
			continue
		}

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
		if sostats.Plays < minPlayRate {
			continue
		}

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

	// derive starter items
	var starterItems []*apb.MatchAggregateCollections_Build
	for sis, sistats := range quot.StarterItems {
		if sistats.Plays < minPlayRate {
			continue
		}

		si, err := deserializeBuild(sis)
		if err != nil {
			return nil, fmt.Errorf("could not deserialize starter items: %v", err)
		}
		// sistats is skill order subscalars
		starterItems = append(starterItems, &apb.MatchAggregateCollections_Build{
			Build:      si,
			PickRate:   sistats.Plays,
			WinRate:    sistats.Wins,
			NumMatches: uint32(sistats.PlayCount),
		})
	}

	// derive build path
	var buildPath []*apb.MatchAggregateCollections_Build
	for bps, bpstats := range quot.BuildPath {
		if bpstats.Plays < minPlayRate {
			continue
		}

		bp, err := deserializeBuild(bps)
		if err != nil {
			return nil, fmt.Errorf("could not deserialize build path: %v", err)
		}
		// bpstats is skill order subscalars
		buildPath = append(buildPath, &apb.MatchAggregateCollections_Build{
			Build:      bp,
			PickRate:   bpstats.Plays,
			WinRate:    bpstats.Wins,
			NumMatches: uint32(bpstats.PlayCount),
		})
	}
	buildPath = groupBuildPaths(buildPath)

	// derive core build list
	var coreBuildList []*apb.MatchAggregateCollections_Build
	for cbs, cbstats := range quot.CoreBuildList {
		if cbstats.Plays < minPlayRate {
			continue
		}

		cb, err := deserializeBuild(cbs)
		if err != nil {
			return nil, fmt.Errorf("could not deserialize build path: %v", err)
		}
		// cbstats is skill order subscalars
		coreBuildList = append(coreBuildList, &apb.MatchAggregateCollections_Build{
			Build:      cb,
			PickRate:   cbstats.Plays,
			WinRate:    cbstats.Wins,
			NumMatches: uint32(cbstats.PlayCount),
		})
	}

	return &apb.MatchAggregateCollections{
		Runes:          runes,
		Masteries:      masteries,
		Keystones:      keystones,
		SummonerSpells: summonerSpells,
		Trinkets:       trinkets,
		SkillOrders:    skillOrders,
		StarterItems:   starterItems,
		BuildPath:      buildPath,
		CoreBuildList:  coreBuildList,
	}, nil
}

func calculatePickRate(champions map[uint32]*apb.MatchQuotient, id uint32) float64 {
	var plays float64
	var champPlays float64
	for _, quot := range champions {
		plays += float64(quot.Scalars.Plays)
		if allied := quot.Allies[id]; allied != nil {
			champPlays += float64(allied.PlayCount)
		}
	}
	// 5 on a team
	champPlays /= 5
	// 10 in a game
	plays /= 10
	return champPlays / plays
}

func calculateBanRate(champions map[uint32]*apb.MatchQuotient, id uint32) float64 {
	var bans float64
	var champBans float64
	for _, quot := range champions {
		bans += float64(quot.Scalars.Plays)
		if banned := quot.Bans[id]; banned != nil {
			champBans += float64(banned.PlayCount)
		}
	}
	return champBans / bans
}

func groupBuildPaths(in []*apb.MatchAggregateCollections_Build) []*apb.MatchAggregateCollections_Build {
	// TODO(igm)
	return in
}
