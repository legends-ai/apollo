package models

import (
	"golang.org/x/net/context"

	apb "github.com/simplyianm/apollo/gen-go/asuna"
)

// ChampionDAO is a Champion DAO.
type ChampionDAO struct{}

// Get gets a champion.
func (dao *ChampionDAO) Get(ctx context.Context, req *apb.GetChampionRequest) (*apb.Champion, error) {
	return nil, nil
}
