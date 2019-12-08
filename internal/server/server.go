package server

import (
	"context"
	"github.com/casbin/casbin"
	"github.com/golang/protobuf/ptypes"
	v1 "github.com/maykonlf/authorization-service/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type AuthorizationService struct {
	enforcer *casbin.Enforcer
}

func NewAuthorizationService(enforcer *casbin.Enforcer) v1.AuthorizationServer {
	return &AuthorizationService{
		enforcer: enforcer,
	}
}

func (service *AuthorizationService) CreatePolicy(ctx context.Context, policy *v1.PolicyRequest) (*v1.PolicyResponse, error) {
	_ = service.enforcer.LoadPolicy()
	service.enforcer.AddPolicy(policy.Role, policy.Tenant, policy.Resource, policy.Action)
	err := service.enforcer.SavePolicy()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	t, _ := ptypes.TimestampProto(time.Now())
	return &v1.PolicyResponse{
		When: t,
	}, err
}
