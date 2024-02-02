package models

import (
	"github.com/erupshis/key_keeper/pb"
)

type User struct {
	ID       int64
	Login    string
	Password string
}

func ConvertUserFromGRPC(in *pb.Creds) *User {
	return &User{
		Login:    in.GetLogin(),
		Password: in.GetPassword(),
	}
}

func ConvertUserToGRPC(in *User) *pb.Creds {
	return &pb.Creds{
		Login:    in.Login,
		Password: in.Password,
	}
}
