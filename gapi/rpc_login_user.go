package gapi

import (
	"context"

	db "github.com/bbsemih/gobank/internal/db/sqlc"
	"github.com/bbsemih/gobank/pb"
	"github.com/bbsemih/gobank/pkg/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == db.ErrRecordNotFound {
			return nil, status.Errorf(codes.NotFound, "user not found in the database: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %s", err)
	}
	err = util.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "password is not correct!")
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %s", err)
	}

	rsp := &pb.LoginUserResponse{
		User: &pb.User{
			Username:          user.Username,
			FullName:          user.FullName,
			Email:             user.Email,
			PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
			CreatedAt:         timestamppb.New(user.CreatedAt),
		},
		AccessToken:          accessToken,
		AccessTokenExpiresAt: timestamppb.New(accessPayload.ExpiredAt),
	}
	return rsp, nil
}
