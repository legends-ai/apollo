package models

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/golang/protobuf/proto"

	apb "github.com/simplyianm/apollo/gen-go/asuna"
)

type MatchSumDAO interface {
	Get(f *apb.MatchFilters) (*apb.MatchSum, error)
}

// NewMatchSumDAO constructs a new MatchSumDAO.
func NewMatchSumDAO() MatchSumDAO {
	return &matchSumDAO{}
}

type matchSumDAO struct {
	CQL *gocql.Session `inject:"t"`
}

func (a *matchSumDAO) Get(f *apb.MatchFilters) (*apb.MatchSum, error) {
	var rawSum []byte
	if err := a.CQL.Query(
		stmtGetSum, f.ChampionId, f.EnemyId, f.Patch,
		f.Tier, int32(f.Region), int32(f.Role),
	).Scan(&rawSum); err != nil {
		return nil, fmt.Errorf("error fetching sum from Cassandra: %v", err)
	}

	var sum apb.MatchSum
	if err := proto.Unmarshal(rawSum, &sum); err != nil {
		return nil, fmt.Errorf("error unmarshaling sum: %v", err)
	}

	return &sum, nil
}
