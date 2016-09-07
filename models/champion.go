package models

import (
	"golang.org/x/net/context"

	"github.com/simplyianm/apollo/aggregation"
	apb "github.com/simplyianm/apollo/gen-go/asuna"
)

const (
	ANY_CHAMPION = -1
	ANY_ENEMY    = -1
)

// ChampionDAO is a Champion DAO.
type ChampionDAO interface {

	// Get gets a champion.
	Get(ctx context.Context, req *apb.GetChampionRequest) (*apb.Champion, error)
}

// ChampionDAOImpl is an implementation of ChampionDAO.
type ChampionDAOImpl struct {
	Aggregator aggregation.Aggregator `inject:"t"`
	Vulgate    Vulgate                `inject:"t"`
}

// Get gets a champion.
func (c *ChampionDAOImpl) Get(ctx context.Context, req *apb.GetChampionRequest) (*apb.Champion, error) {
	base := c.buildBase(req)
	filters := c.buildFilters(req)

	_, err := c.Aggregator.Aggregate(base, filters)
	if err != nil {
		return nil, err
	}

	// TODO(igm): implement
	return nil, nil
}

func (c *ChampionDAOImpl) buildBase(req *apb.GetChampionRequest) []*apb.MatchFilters {
	var ret []*apb.MatchFilters
	for _, patch := range c.Vulgate.FindPatches(req.Patch) {
		for _, tier := range c.Vulgate.FindTiers(req.Tier) {
			ret = append(ret, &apb.MatchFilters{
				ChampionId: int32(ANY_CHAMPION),
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

func (c *ChampionDAOImpl) buildFilters(req *apb.GetChampionRequest) []*apb.MatchFilters {
	var ret []*apb.MatchFilters
	for _, patch := range c.Vulgate.FindPatches(req.Patch) {
		for _, tier := range c.Vulgate.FindTiers(req.Tier) {
			ret = append(ret, &apb.MatchFilters{
				ChampionId: int32(req.ChampionId),
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
