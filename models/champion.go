package models

import (
	"math/rand"
	"time"

	"github.com/golang/protobuf/ptypes"
	apb "github.com/simplyianm/apollo/gen-go/asuna"
	"golang.org/x/net/context"
)

// ChampionDAO is a Champion DAO.
type ChampionDAO struct{}

// Get gets a champion.
func (dao *ChampionDAO) Get(ctx context.Context, req *apb.GetChampionRequest) (*apb.Champion, error) {
	rand.Seed(time.Now().UTC().Unix())

	const total = int32(42)
	patchStart, _ := ptypes.TimestampProto(time.Now().Add(-24 * time.Hour))
	patchEnd, _ := ptypes.TimestampProto(time.Now())

	return &apb.Champion{
		Metadata: &apb.Champion_Metadata{
			StaticInfo: &apb.Vulgate_Champion{
				Base: &apb.Vulgate_Champion_Base{
					Id:    uint32(64),
					Title: "The Blind Monk",
					Name:  "Lee Sin",
					Key:   "LeeSin",
				},
			},
			PatchStart: patchStart,
			PatchEnd:   patchEnd,
		},
		MatchAggregate: &apb.MatchAggregate{
			Role: &apb.MatchAggregateRoles{
				Role:                 apb.Role_JUNGLE,
				TotalChampionsInRole: uint32(total),
				RoleStats: []*apb.MatchAggregateRoles_RoleStats{
					{
						Role:       apb.Role_TOP,
						PickRate:   rand.Float64(),
						NumMatches: rand.Uint32(),
					},
					{
						Role:       apb.Role_JUNGLE,
						PickRate:   rand.Float64(),
						NumMatches: rand.Uint32(),
					},
				},
			},
			Statistics: &apb.MatchAggregateStatistics{
				Scalars: &apb.MatchAggregateStatistics_Scalars{
					WinRate:                  generateRandomStatistic(total),
					PickRate:                 generateRandomStatistic(total),
					BanRate:                  generateRandomStatistic(total),
					GamesPlayed:              generateRandomStatistic(total),
					GoldEarned:               generateRandomStatistic(total),
					Kills:                    generateRandomStatistic(total),
					Deaths:                   generateRandomStatistic(total),
					Assists:                  generateRandomStatistic(total),
					DamageDealt:              generateRandomStatistic(total),
					DamageTaken:              generateRandomStatistic(total),
					MinionsKilled:            generateRandomStatistic(total),
					TeamJungleMinionsKilled:  generateRandomStatistic(total),
					EnemyJungleMinionsKilled: generateRandomStatistic(total),
					StructureDamage:          generateRandomStatistic(total),
					KillingSpree:             generateRandomStatistic(total),
					WardsBought:              generateRandomStatistic(total),
					WardsPlaced:              generateRandomStatistic(total),
					WardsKilled:              generateRandomStatistic(total),
					CrowdControl:             generateRandomStatistic(total),
					FirstBlood:               generateRandomStatistic(total),
					FirstBloodAssist:         generateRandomStatistic(total),
					DoubleKills:              generateRandomStatistic(total),
					TripleKills:              generateRandomStatistic(total),
					Quadrakills:              generateRandomStatistic(total),
					Pentakills:               generateRandomStatistic(total),
				},
				Deltas: &apb.MatchAggregateStatistics_Deltas{
					CsDiff:          generateRandomDelta(total),
					XpDiff:          generateRandomDelta(total),
					DamageTakenDiff: generateRandomDelta(total),
					XpPerMin:        generateRandomDelta(total),
					GoldPerMin:      generateRandomDelta(total),
					TowersPerMin:    generateRandomDelta(total),
					WardsPlaced:     generateRandomDelta(total),
					DamageTaken:     generateRandomDelta(total),
				},
			},
			Graphs: &apb.MatchAggregateGraphs{
				/*Distribution: &apb.MatchAggregateGraphs_Distribution{
					WinRate:  generateRandomDistributionGraph(),
					PickRate: generateRandomDistributionGraph(),
					BanRate:  generateRandomDistributionGraph(),
				},*/
				ByPatch: []*apb.MatchAggregateGraphs_ByPatch{
					{
						Patch:    "6.17",
						WinRate:  rand.Float64(),
						PickRate: rand.Float64(),
						BanRate:  rand.Float64(),
					},
					{
						Patch:    "6.17",
						WinRate:  rand.Float64(),
						PickRate: rand.Float64(),
						BanRate:  rand.Float64(),
					},
				},
				ByGameLength: []*apb.MatchAggregateGraphs_ByGameLength{
					{
						GameLength: &apb.IntRange{Min: rand.Uint32(), Max: rand.Uint32()},
						WinRate:    rand.Float64(),
					},
					{
						GameLength: &apb.IntRange{Min: rand.Uint32(), Max: rand.Uint32()},
						WinRate:    rand.Float64(),
					},
				},
				ByGamesPlayed: []*apb.MatchAggregateGraphs_ByGamesPlayed{
					{
						GamesPlayed: &apb.IntRange{Min: rand.Uint32(), Max: rand.Uint32()},
						WinRate:     rand.Float64(),
					},
					{
						GamesPlayed: &apb.IntRange{Min: rand.Uint32(), Max: rand.Uint32()},
						WinRate:     rand.Float64(),
					},
				},
				PhysicalDamage: rand.Float64(),
				MagicDamage:    rand.Float64(),
				TrueDamage:     rand.Float64(),
				ByExperience: []*apb.MatchAggregateGraphs_ByExperience{
					{
						Experience: &apb.IntRange{Min: rand.Uint32(), Max: rand.Uint32()},
						Count:      rand.Uint32(),
					},
					{
						Experience: &apb.IntRange{Min: rand.Uint32(), Max: rand.Uint32()},
						Count:      rand.Uint32(),
					},
				},
				GoldOverTime: []*apb.MatchAggregateGraphs_GoldPerTime{
					{
						Gold: rand.Float64(),
						Time: &apb.IntRange{Min: rand.Uint32(), Max: rand.Uint32()},
					},
					{
						Gold: rand.Float64(),
						Time: &apb.IntRange{Min: rand.Uint32(), Max: rand.Uint32()},
					},
				},
			},
			Collections: &apb.MatchAggregateCollections{
				Runes:          generateRandomRuneSet(),
				Masteries:      generateRandomMasterySet(),
				Keystones:      generateRandomKeystones(),
				SummonerSpells: generateRandomSummonerSets(),
				Trinkets:       generateRandomTrinkets(),
				SkillOrders:    generateRandomSkillOrder(),
				StarterItems:   generateRandomBuild(2),
				BuildPath:      generateRandomBuild(10),
				CoreBuildList:  generateRandomBuild(5),
			},
		},
	}, nil
}

func generateRandomStatistic(num int32) *apb.MatchAggregateStatistics_Statistic {
	return &apb.MatchAggregateStatistics_Statistic{
		Rank:    uint32(rand.Int31n(num)),
		Change:  int32(rand.Int31n(num*2) - num),
		Value:   rand.Float64(),
		Average: rand.Float64(),
	}
}

func generateRandomDelta(num int32) *apb.MatchAggregateStatistics_Deltas_Delta {
	return &apb.MatchAggregateStatistics_Deltas_Delta{
		ZeroToTen:      generateRandomStatistic(num),
		TenToTwenty:    generateRandomStatistic(num),
		TwentyToThirty: generateRandomStatistic(num),
		ThirtyToEnd:    generateRandomStatistic(num),
	}
}

func generateRandomDistributionGraph() map[string]float64 {
	graph := map[string]float64{}
	for _, x := range championList() {
		graph[x] = rand.Float64()
	}
	return graph
}

func generateRandomRuneSet() []*apb.MatchAggregateCollections_RuneSet {
	runeSet := []*apb.MatchAggregateCollections_RuneSet{}
	for i := 0; i < 5; i++ {
		runeSet = append(runeSet, &apb.MatchAggregateCollections_RuneSet{
			// Runes:      generateRandomUint32Map(),
			PickRate:   rand.Float64(),
			WinRate:    rand.Float64(),
			NumMatches: rand.Uint32(),
		})
	}
	return runeSet
}

func generateRandomMasterySet() []*apb.MatchAggregateCollections_MasterySet {
	masterySet := []*apb.MatchAggregateCollections_MasterySet{}
	for i := 0; i < 5; i++ {
		masterySet = append(masterySet, &apb.MatchAggregateCollections_MasterySet{
			// Masteries:      generateRandomUint32Map(),
			PickRate:   rand.Float64(),
			WinRate:    rand.Float64(),
			NumMatches: rand.Uint32(),
		})
	}
	return masterySet
}

func generateRandomKeystones() []*apb.MatchAggregateCollections_Keystone {
	keystones := []*apb.MatchAggregateCollections_Keystone{}
	for i := 0; i < 5; i++ {
		keystones = append(keystones, &apb.MatchAggregateCollections_Keystone{
			Keystone:   rand.Uint32(),
			PickRate:   rand.Float64(),
			WinRate:    rand.Float64(),
			NumMatches: rand.Uint32(),
		})
	}
	return keystones
}

func generateRandomSummonerSets() []*apb.MatchAggregateCollections_SummonerSet {
	summonerSet := []*apb.MatchAggregateCollections_SummonerSet{}
	for i := 0; i < 5; i++ {
		summonerSet = append(summonerSet, &apb.MatchAggregateCollections_SummonerSet{
			Spell1:     rand.Uint32(),
			Spell2:     rand.Uint32(),
			PickRate:   rand.Float64(),
			WinRate:    rand.Float64(),
			NumMatches: rand.Uint32(),
		})
	}
	return summonerSet
}

func generateRandomTrinkets() []*apb.MatchAggregateCollections_Trinket {
	trinkets := []*apb.MatchAggregateCollections_Trinket{}
	for i := 0; i < 5; i++ {
		trinkets = append(trinkets, &apb.MatchAggregateCollections_Trinket{
			Trinket:    rand.Uint32(),
			PickRate:   rand.Float64(),
			WinRate:    rand.Float64(),
			NumMatches: rand.Uint32(),
		})
	}
	return trinkets
}

func generateRandomSkillOrder() []*apb.MatchAggregateCollections_SkillOrder {
	skillOrder := []*apb.MatchAggregateCollections_SkillOrder{}
	for i := 0; i < 5; i++ {
		rand.Seed(time.Now().UTC().Unix())
		var skills = []apb.Ability{}
		for j := 0; j < int(rand.Int31n(19)); j++ {
			skills = append(skills, apb.Ability(rand.Int31n(4)))
		}
		skillOrder = append(skillOrder, &apb.MatchAggregateCollections_SkillOrder{
			SkillOrder: skills,
			PickRate:   rand.Float64(),
			WinRate:    rand.Float64(),
			NumMatches: rand.Uint32(),
		})
	}
	return skillOrder
}

func generateRandomBuild(num int) []*apb.MatchAggregateCollections_Build {
	build := []*apb.MatchAggregateCollections_Build{}
	for i := 0; i < 5; i++ {
		var items = []uint32{}
		for j := 0; j < num; j++ {
			items = append(items, rand.Uint32())
		}
		build = append(build, &apb.MatchAggregateCollections_Build{
			Build:      items,
			PickRate:   rand.Float64(),
			WinRate:    rand.Float64(),
			NumMatches: rand.Uint32(),
		})
	}
	return build
}

func generateRandomUint32Map() map[uint32]uint32 {
	uint32map := map[uint32]uint32{}
	for i := uint32(0); i < 30; i++ {
		uint32map[i] = uint32(rand.Int31n(5))
	}
	return uint32map
}

func championList() []string {
	return []string{"Aatrox", "Ahri", "Akali", "Alistar", "Amumu", "Anivia", "Annie", "Ashe", "Aurelion Sol", "Azir", "Bard", "Blitzcrank", "Brand", "Braum", "Caitlyn", "Cassiopeia", "Cho'Gath", "Corki", "Darius", "Diana", "Dr. Mundo", "Draven", "Ekko", "Elise", "Evelynn", "Ezreal", "Fiddlesticks", "Fiora", "Fizz", "Galio", "Gangplank", "Garen", "Gnar", "Gragas", "Graves", "Hecarim", "Heimerdinger", "Illaoi", "Irelia", "Janna", "Jarvan", "Jax", "Jayce", "Jhin", "Jinx", "Kalista", "Karma", "Karthus", "Kassadin", "Katarina", "Kayle", "Kennen", "Kha'Zix", "Kindred", "Kled", "Kog'Maw", "LeBlanc", "Lee Sin", "Leona", "Lissandra", "Lucian", "Lulu", "Lux", "Malphite", "Malzahar", "Maokai", "Master", "Miss", "Mordekaiser", "Morgana", "Nami", "Nasus", "Nautilus", "Nidalee", "Nocturne", "Nunu", "Olaf", "Orianna", "Pantheon", "Poppy", "Quinn", "Rammus", "Rek'Sai", "Rengar", "Riven", "Rumble", "Ryze", "Sejuani", "Shaco", "Shen", "Shyvana", "Singed", "Sion", "Sivir", "Skarner", "Sona", "Soraka", "Swain", "Syndra", "Tahm", "Taliyah", "Talon", "Taric", "Teemo", "Thresh", "Tristana", "Trundle", "Tryndamere", "Twisted", "Twitch", "Udyr", "Urgot", "Varus", "Vayne", "Veigar", "Vel'Koz", "Vi", "Viktor", "Vladimir", "Volibear", "Warwick", "Wukong", "Xerath", "Xin", "Yasuo", "Yorick", "Zac", "Zed", "Ziggs", "Zilean", "Zyra"}
}
