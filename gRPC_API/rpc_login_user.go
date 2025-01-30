package grpcapi

import (
	"context"
	"database/sql"

	db "github.com/siddheshRajendraNimbalkar/GO-BANK/db/sqlc"
	"github.com/siddheshRajendraNimbalkar/GO-BANK/pb"
	"github.com/siddheshRajendraNimbalkar/GO-BANK/token"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (service GRPCService) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	user, err := service.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "User by %s name not found", req.Username)
		}
		return nil, status.Errorf(codes.NotFound, "User not found: %s", err)
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.GetPassword()))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "password: %s", err)
	}

	maker, err := token.NewPasetoMaker(service.config.Secret)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error creating token maker %s", err)
	}

	tokenStr, _, err := maker.CreateToken(user.Username, service.config.JwtDuration)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error creating token maker %s", err)
	}

	refesh_token, userPayload, err := service.tokenMaker.CreateToken(user.Username, service.config.SessionDuration)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error creating refesh token %s", err)
	}

	arg := db.CreateSessionParams{
		ID:           userPayload.ID,
		Username:     userPayload.Username,
		RefreshToken: refesh_token,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpireDate:   userPayload.ExpiresAt,
	}

	session, err := service.store.CreateSession(ctx, arg)

	if err != nil {

		return nil, status.Errorf(codes.Internal, "Error while creating sessioin %s", err)
	}

	// string session_id = 2;
	// string access_token = 3;
	// string refresh_token = 4;
	// google.protobuf.Timestamp token_expireAt = 5;

	loginUserResponce := &pb.LoginUserResponse{
		User: &pb.User{
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			CreatedAt: timestamppb.New(user.CreatedAt),
		},
		SessionId:     session.ID.String(),
		AccessToken:   tokenStr,
		RefreshToken:  refesh_token,
		TokenExpireAt: timestamppb.New(session.ExpireDate),
	}

	return loginUserResponce, nil
}
