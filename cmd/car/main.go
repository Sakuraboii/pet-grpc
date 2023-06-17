package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
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
	pb.UnimplementedCarServiceServer
	carsRepo repository.CarsRepo
}

func NewImplementation(db db.DBops) *Implementation {
	return &Implementation{
		carsRepo: postgresql.NewCars(db),
	}
}

const (
	service     = "car_api"
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

	lsn, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	_, database, _ := db.NewDB(ctx)

	// HTTP exporter для prometheus
	go http.ListenAndServe(":9091", promhttp.Handler())

	server := grpc.NewServer()
	pb.RegisterCarServiceServer(server, NewImplementation(database))

	log.Printf("starting server on %s", lsn.Addr().String())
	if err = server.Serve(lsn); err != nil {
		log.Fatal(err)
	}
}

func (s *Implementation) CreateCar(ctx context.Context, in *pb.CreateCarRequest) (*pb.CreateCarResponse, error) {
	tr := otel.Tracer("CreateCar")
	ctx, span := tr.Start(ctx, "received request")
	span.SetAttributes(attribute.Key("params").String(in.String()))
	defer span.End()

	car := &repository.Car{Model: in.Model, UserId: in.UserId}
	carId, err := s.carsRepo.Add(ctx, car)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	internal.RegCarCounter.Add(1)

	return &pb.CreateCarResponse{
		Id: carId,
	}, nil
}

func (s *Implementation) GetCar(ctx context.Context, in *pb.GetCarRequest) (*pb.GetCarResponse, error) {
	tr := otel.Tracer("GetCar")
	ctx, span := tr.Start(ctx, "received request")
	span.SetAttributes(attribute.Key("params").String(in.String()))
	defer span.End()

	car, err := s.carsRepo.GetById(ctx, in.Id)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	responseCar := &pb.Car{
		Id:     car.ID,
		Model:  car.Model,
		UserId: car.UserId,
	}

	internal.GetUserCounter.Add(1)

	return &pb.GetCarResponse{
		Car: responseCar,
	}, nil
}

func (s *Implementation) UpdateCar(ctx context.Context, in *pb.UpdateCarRequest) (*pb.UpdateCarResponse, error) {
	tr := otel.Tracer("UpdateCar")
	ctx, span := tr.Start(ctx, "received request")
	span.SetAttributes(attribute.Key("params").String(in.String()))
	defer span.End()

	car, err := s.carsRepo.GetById(ctx, in.Car.Id)
	if err != nil {
		fmt.Println(err)
	}
	car.Model = in.Car.Model

	complete, err := s.carsRepo.Update(ctx, car)

	if err != nil {
		fmt.Println(err)
	}

	internal.UpdatedUserCounter.Add(1)

	return &pb.UpdateCarResponse{
		Ok: complete,
	}, nil
}

func (s *Implementation) DeleteCar(ctx context.Context, in *pb.DeleteCarRequest) (*pb.DeleteCarResponse, error) {
	tr := otel.Tracer("DeleteCar")
	ctx, span := tr.Start(ctx, "received request")
	span.SetAttributes(attribute.Key("params").String(in.String()))
	defer span.End()

	err := s.carsRepo.Delete(ctx, in.Id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	internal.DeletedCarCounter.Add(1)

	return &pb.DeleteCarResponse{Ok: true}, nil
}
