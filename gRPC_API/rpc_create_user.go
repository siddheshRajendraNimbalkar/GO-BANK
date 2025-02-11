package grpcapi

import (
	"context"

	db "github.com/siddheshRajendraNimbalkar/GO-BANK/db/sqlc"
	"github.com/siddheshRajendraNimbalkar/GO-BANK/pb"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (service GRPCService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreatedUserResponse, error) {

	if req.GetPassword() == "" || req.GetEmail() == "" || req.GetFullName() == "" || req.Username == "" {
		return nil, status.Errorf(codes.Internal, "Bad Request Mention Every detals")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while hashing password: %s", err)
	}

	service.extractMetaData(ctx)

	// Prepare user parameters
	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: string(hashPassword),
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	// Create user in the database
	user, err := service.store.CreateUser(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while creating user: %s", err)
	}

	createUserRPC := &pb.CreatedUserResponse{
		User: &pb.User{
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			CreatedAt: timestamppb.New(user.CreatedAt),
		},
	}

	return createUserRPC, nil
}
