package auth

import (
	"github.com/erupshis/key_keeper/internal/common/auth"
	"github.com/erupshis/key_keeper/pb"
)

type Controller struct {
	pb.UnimplementedAuthServer

	authManager *auth.Manager
}

func NewController(authManager *auth.Manager) *Controller {
	return &Controller{
		authManager: authManager,
	}
}
