package main

import (
	"context"
	"github.com/casbin/casbin"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/maykonlf/authorization-service/internal/server"
	v1 "github.com/maykonlf/authorization-service/pkg/api/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net/http"
	"strings"
)

func main() {
	serverCert, err := credentials.NewServerTLSFromFile("server.crt", "server.key")
	if err != nil {
		log.Fatalln("failed to create server cert", err)
	}

	e := casbin.NewEnforcer("rbac_model.conf", "roles.csv")

	grpcServer := grpc.NewServer(grpc.Creds(serverCert))
	v1.RegisterAuthorizationServer(grpcServer, server.NewAuthorizationService(e))

	clientCert, err := credentials.NewClientTLSFromFile("server.crt", "")
	if err != nil {
		log.Fatalln("failed to create client cert", err)
	}

	conn, err := grpc.DialContext(context.Background(), "localhost:9000", grpc.WithTransportCredentials(clientCert))
	if err != nil {
		log.Fatalln("failed to dial gRPC server", err)
	}

	router := runtime.NewServeMux()
	if err = v1.RegisterAuthorizationHandler(context.Background(), router, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	log.Info("serving on port :9000")
	log.Fatalln(http.ListenAndServeTLS(":9000", "server.crt", "server.key", httpGrpcRouter(grpcServer, router)))
}

func httpGrpcRouter(grpcServer *grpc.Server, httpHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			httpHandler.ServeHTTP(w, r)
		}
	})
}
