package main

import (
	"context"
	"fmt"
	"github.com/pressly/goose/v3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"google.golang.org/grpc"
	"homework-8/internal"
	"homework-8/internal/app/homework/pb"
	"homework-8/internal/pkg/db"
	"homework-8/internal/pkg/repository"
	"homework-8/internal/pkg/repository/postgresql"
	"log"
	"net"
	"net/http"
	"time"
)

type Implementation struct {
	pb.UnimplementedUserServiceServer
	userRepo repository.UsersRepo
}

func NewImplementation(db db.DBops) *Implementation {
	return &Implementation{
		userRepo: postgresql.NewUsers(db),
	}
}

const (
	service     = "user_api"
	environment = "development"
)

func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(service),
			attribute.String("environment", environment),
		)),
	)
	return tp, nil
}

func main() {
	tp, err := tracerProvider("http://localhost:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cleanly shutdown and flush telemetry when the application exits.
	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err = tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	lsn, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal(err)
	}

	sqlDb, database, _ := db.NewDB(ctx)

	err = goose.Up(sqlDb.DB, "./internal/pkg/db/migrations")

	if err != nil {
		log.Fatalf("невозможно накатить миграции: %v", err)
	}

	// HTTP exporter для prometheus
	go http.ListenAndServe(":9091", promhttp.Handler())

	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, NewImplementation(database))

	log.Printf("starting server on %s", lsn.Addr().String())
	if err = server.Serve(lsn); err != nil {
		log.Fatal(err)
	}
}

func (s *Implementation) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	tr := otel.Tracer("CreateUser")
	ctx, span := tr.Start(ctx, "received request")
	span.SetAttributes(attribute.Key("params").String(in.String()))
	defer span.End()

	userId, err := s.userRepo.Add(ctx, &repository.User{Name: in.Name})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	internal.RegUserCounter.Add(1)

	return &pb.CreateUserResponse{Id: userId}, nil
}

func (s *Implementation) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	tr := otel.Tracer("GetUser")
	ctx, span := tr.Start(ctx, "received request")
	span.SetAttributes(attribute.Key("params").String(in.String()))
	defer span.End()

	user, err := s.userRepo.GetById(ctx, in.Id)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	internal.GetUserCounter.Add(1)

	responseUser := &pb.User{
		Id:   user.ID,
		Name: user.Name,
	}

	return &pb.GetUserResponse{User: responseUser}, nil
}

func (s *Implementation) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	tr := otel.Tracer("UpdateUser")
	ctx, span := tr.Start(ctx, "received request")
	span.SetAttributes(attribute.Key("params").String(in.String()))
	defer span.End()

	user, err := s.userRepo.GetById(ctx, in.User.Id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	user.Name = in.User.Name

	complete, err := s.userRepo.Update(ctx, user)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	internal.UpdatedUserCounter.Add(1)

	return &pb.UpdateUserResponse{Ok: complete}, nil
}

func (s *Implementation) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	tr := otel.Tracer("UpdateUser")
	ctx, span := tr.Start(ctx, "received request")
	span.SetAttributes(attribute.Key("params").String(in.String()))
	defer span.End()

	err := s.userRepo.Delete(ctx, in.Id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	internal.DeletedUserCounter.Add(1)

	return &pb.DeleteUserResponse{Ok: true}, nil
}
