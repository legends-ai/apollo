package server

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	apb "github.com/simplyianm/apollo/gen-go/apollo"
)

type Server struct{}

func (s *Server) GetChampion(ctx context.Context, in *apb.GetChampionRequest) (*apb.Champion, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "GetChampion unimplemented")
}

func (s *Server) GetMatchup(ctx context.Context, in *apb.GetMatchupRequest) (*apb.Matchup, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "GetMatchup unimplemented")
}

func (s *Server) GetProfile(ctx context.Context, in *apb.GetProfileRequest) (*apb.Profile, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "GetProfile unimplemented")
}
