package server

import (
	"context"
	"github.com/casbin/casbin"
	"github.com/golang/protobuf/ptypes"
	"github.com/maykonlf/authorization-service/proto/authorization/v1"
	"time"
)

type AuthorizationService struct {
	enforcer *casbin.Enforcer
}

func NewAuthorizationService(enforcer *casbin.Enforcer) authorization.AuthorizationServer {
	return &AuthorizationService{
		enforcer: enforcer,
	}
}

func (service *AuthorizationService) CreatePolicy(ctx context.Context, policy *authorization.PolicyRequest) (*authorization.PolicyResponse, error) {
	_ = service.enforcer.LoadPolicy()
	service.enforcer.AddPolicy(policy.Role, policy.Tenant, policy.Resource, policy.Action)
	err := service.enforcer.SavePolicy()

	t, _ := ptypes.TimestampProto(time.Now())
	return &authorization.PolicyResponse{
		When: t,
	}, err
}
