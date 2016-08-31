package models

import (
	"math/rand"
	"time"

	"golang.org/x/net/context"

	"github.com/golang/protobuf/ptypes"
	apb "github.com/simplyianm/apollo/gen-go/apollo"
)

// ChampionDAO is a Champion DAO.
type ChampionDAO struct{}

// Get gets a champion.
func (dao *ChampionDAO) Get(ctx context.Context, req *apb.GetChampionRequest) (*apb.Champion, error) {
	// TODO(igm): generate real data
	patchStart, _ := ptypes.TimestampProto(time.Now().Add(-24 * time.Hour))
	patchEnd, _ := ptypes.TimestampProto(time.Now())

	return &apb.Champion{
		Metadata: &apb.Champion_Metadata{
			StaticInfo: &apb.ChampionStatic{
				Id:    uint32(64),
				Title: "the Blind Monk",
				Name:  "Lee Sin",
				Key:   "LeeSin",
			},
			PatchStart:           patchStart,
			PatchEnd:             patchEnd,
			Role:                 apb.Role_JUNGLE,
			TotalChampionsInRole: 42,
			RoleStats: []*apb.Champion_Metadata_RoleStats{
				{
					Role:       apb.Role_TOP,
					PickRate:   0.25,
					NumMatches: 3,
				},
				{
					Role:       apb.Role_JUNGLE,
					PickRate:   0.05,
					NumMatches: rand.Uint32(),
				},
			},
		},

		// Random statistics
		Statistics: &apb.Champion_Statistics{
			WinRate:                  randStatistic(),
			PickRate:                 randStatistic(),
			BanRate:                  randStatistic(),
			GamesPlayed:              randStatistic(),
			GoldEarned:               randStatistic(),
			Kills:                    randStatistic(),
			Deaths:                   randStatistic(),
			Assists:                  randStatistic(),
			DamageDealt:              randStatistic(),
			DamageTaken:              randStatistic(),
			MinionsKilled:            randStatistic(),
			TeamJungleMinionsKilled:  randStatistic(),
			EnemyJungleMinionsKilled: randStatistic(),
			StructureDamage:          randStatistic(),
			KillingSpree:             randStatistic(),
			WardsPlaced:              randStatistic(),
			WardsKilled:              randStatistic(),
			CrowdControl:             randStatistic(),
			FirstBlood:               randStatistic(),
			FirstBloodAssist:         randStatistic(),
			Diffs: &apb.Champion_Statistics_Differentials{
				Cs:          randStatistic(),
				DamageTaken: randStatistic(),
				Xp:          randStatistic(),
			},
			MultikillStats: &apb.Champion_Statistics_MultikillStats{
				DoubleKill: randStatistic(),
				TripleKill: randStatistic(),
				Quadrakill: randStatistic(),
				Pentakill:  randStatistic(),
			},
		},

		Graphs: &apb.Champion_Graphs{
			Distribution: &apb.Champion_Graphs_Distribution{
				WinRate: map[string]float64{
					"MonkeyKing": 3.14,
				},
				PickRate: map[string]float64{
					"MonkeyKing": 3.14,
				},
				BanRate: map[string]float64{
					"MonkeyKing": 3.14,
				},
			},

			ByPatch: []*apb.Champion_Graphs_ByPatch{
				{
					Patch:    "3.14",
					WinRate:  0.5,
					PickRate: 0.2,
					BanRate:  0.1,
				},
				{
					Patch:    "3.15",
					WinRate:  0.2,
					PickRate: 0.1,
					BanRate:  0.6,
				},
				{
					Patch:    "3.16",
					WinRate:  0.2,
					PickRate: 0.6,
					BanRate:  0.02,
				},
			},

			ByGameLength: []*apb.Champion_Graphs_ByGameLength{
				{
					GameLength: &apb.IntRange{
						Min: 1,
						Max: 5,
					},
					WinRate: 0.24,
				},
				{
					GameLength: &apb.IntRange{
						Min: 6,
						Max: 10,
					},
					WinRate: 0.54,
				},
			},

			ByGamesPlayed: []*apb.Champion_Graphs_ByGamesPlayed{
				{
					GamesPlayed: &apb.IntRange{
						Min: 1,
						Max: 5,
					},
					WinRate: 0.24,
				},
				{
					GamesPlayed: &apb.IntRange{
						Min: 6,
						Max: 10,
					},
					WinRate: 0.54,
				},
			},

			PhysicalDamage: 0.23,
			MagicDamage:    0.37,
			TrueDamage:     0.4,

			ByExperience: []*apb.Champion_Graphs_ByExperience{
				{
					Experience: &apb.IntRange{
						Min: 1,
						Max: 5,
					},
					Count: 23,
				},
			},
		},
	}, nil
}

func randStatistic() *apb.Champion_Statistics_Statistic {
	return &apb.Champion_Statistics_Statistic{
		Rank:       rand.Uint32(),
		Change:     rand.Uint32(),
		Value:      rand.Float64(),
		Average:    rand.Float64(),
		Percentile: rand.Float64(),
	}
}
