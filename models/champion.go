package models

import (
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
	patchStart, _ := ptypes.TimestampProto(time.Now().Sub(24 * time.Hour))
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
			RoleStats: []*apb.RoleStats{
				{
					Role:       apb.Role_TOP,
					PickRate:   0.25,
					NumMatches: 3,
				},
				{
					Role:       apb.Role_JUNGLE,
					PickRate:   0.05,
					NumMatches: 4,
				},
			},
		},
	}, nil
}
