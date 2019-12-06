package main

import (
	"context"
	"github.com/casbin/casbin"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	grpc2 "github.com/maykonlf/authorization-service/internal/server"
	"github.com/maykonlf/authorization-service/proto/authorization/v1"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"strings"
)

func main() {
	e := casbin.NewEnforcer("rbac_model.conf", "roles.csv")

	grpcServer := grpc.NewServer()
	authorization.RegisterAuthorizationServer(grpcServer, grpc2.NewAuthorizationService(e))

	go func() {
		lis, err := net.Listen("tcp", ":50000")
		if err != nil {
			panic(err)
		}

		err = grpcServer.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		"localhost:50000",
		grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	router := runtime.NewServeMux()
	if err = authorization.RegisterAuthorizationHandler(context.Background(), router, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	log.Println("listening on :9000")
	log.Fatalln(http.ListenAndServe(":9000", httpGrpcRouter(grpcServer, router)))
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