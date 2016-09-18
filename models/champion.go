package models

import (
	"golang.org/x/net/context"

	apb "github.com/asunaio/apollo/gen-go/asuna"
)

const (
	ANY_CHAMPION = -1
	ANY_ENEMY    = -1
)

// ChampionDAO is a Champion DAO.
type ChampionDAO interface {
	// Get gets a champion.
	Get(ctx context.Context, req *apb.GetChampionRequest) (*apb.Champion, error)

	// GetMatchup gets a matchup.
	GetMatchup(ctx context.Context, req *apb.GetMatchupRequest) (*apb.Matchup, error)
}

// NewChampionDAO returns a new ChampionDAO.
func NewChampionDAO() ChampionDAO {
	return &championDAOImpl{}
}

// championDAOImpl is an implementation of ChampionDAO.
type championDAOImpl struct {
	Aggregator Aggregator `inject:"t"`
	Vulgate    Vulgate    `inject:"t"`
}

// Get gets a champion.
func (c *championDAOImpl) Get(ctx context.Context, req *apb.GetChampionRequest) (*apb.Champion, error) {
	agg, err := c.Aggregator.Aggregate(req)
	if err != nil {
		return nil, err
	}

	// TODO(igm): implement

	patchTimes := c.Vulgate.GetPatchTimes(req.Patch)

	return &apb.Champion{
		Metadata: &apb.Champion_Metadata{
			StaticInfo: c.Vulgate.GetChampionInfo(req.ChampionId),
			PatchStart: patchTimes.Start,
			PatchEnd:   patchTimes.End,
		},
		MatchAggregate: agg,
	}, nil
}

func (c *championDAOImpl) GetMatchup(ctx context.Context, req *apb.GetMatchupRequest) (*apb.Matchup, error) {
	return &apb.Matchup{}, nil
}
