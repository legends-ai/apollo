package server

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	apb "github.com/asunaio/apollo/gen-go/asuna"
	"github.com/asunaio/apollo/models"
)

type Server struct {
	Champions models.ChampionDAO `inject:"t"`
}

func (s *Server) GetChampion(ctx context.Context, in *apb.GetChampionRequest) (*apb.Champion, error) {
	champion, err := s.Champions.Get(ctx, in)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "could not get champion: %v", err)
	}
	return champion, nil
}

func (s *Server) GetMatchup(ctx context.Context, in *apb.GetMatchupRequest) (*apb.Matchup, error) {
	matchup, err := s.Champions.GetMatchup(ctx, in)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "could not get matchup: %v", err)
	}
	return matchup, nil
}

func (s *Server) GetProfile(ctx context.Context, in *apb.GetProfileRequest) (*apb.Profile, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "GetProfile unimplemented")
}
