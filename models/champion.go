package models

import (
	"github.com/golang/protobuf/ptypes"
	apb "github.com/simplyianm/apollo/gen-go/asuna"
	"golang.org/x/net/context"
	"math/rand"
	"time"
)

// ChampionDAO is a Champion DAO.
type ChampionDAO struct{}

// Get gets a champion.
func (dao *ChampionDAO) Get(ctx context.Context, req *apb.GetChampionRequest) (*apb.Champion, error) {
	rand.Seed(time.Now().UTC().Unix())

	const totalChampionsInRole = int32(42)
	patchStart, _ := ptypes.TimestampProto(time.Now().Add(-24 * time.Hour))
	patchEnd, _ := ptypes.TimestampProto(time.Now())

	return &apb.Champion{
		Metadata: &apb.Champion_Metadata{
			StaticInfo: &apb.ChampionInfo{
				Id:    uint32(64),
				Title: "The Blind Monk",
				Name:  "Lee Sin",
				Key:   "LeeSin",
			},
			PatchStart: patchStart,
			PatchEnd:   patchEnd,
		},
		MatchAggregate: &apb.MatchAggregate{
			Role: &apb.MatchAggregateRoles{
				Role:                 apb.Role_JUNGLE,
				TotalChampionsInRole: uint32(totalChampionsInRole),
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
				WinRate: &apb.MatchAggregateStatistics_Statistic{
					Rank:    uint32(rand.Int31n(totalChampionsInRole)),
					Change:  int32(-2),
					Value:   rand.Float64(),
					Average: rand.Float64(),
				},
			},
		},
	}, nil
}
