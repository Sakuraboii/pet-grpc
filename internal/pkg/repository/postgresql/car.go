package postgresql

import (
	"context"
	"database/sql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"homework-8/internal/pkg/db"
	"homework-8/internal/pkg/repository"
)

type CarsRepo struct {
	db db.DBops
}

func NewCars(db db.DBops) *CarsRepo {
	return &CarsRepo{db: db}
}

func (r *CarsRepo) Add(ctx context.Context, car *repository.Car) (int64, error) {
	tr := otel.Tracer("AddCar")
	ctx, span := tr.Start(ctx, "database layer")
	span.SetAttributes(attribute.Key("params").String(car.Model))
	defer span.End()

	var id int64
	err := r.db.ExecQueryRow(ctx, `INSERT INTO cars(model,user_id) VALUES ($1,$2) RETURNING id`,
		car.Model, car.UserId).Scan(&id)
	return id, err
}

func (r *CarsRepo) GetById(ctx context.Context, id int64) (*repository.Car, error) {
	tr := otel.Tracer("GetByIdCar")
	ctx, span := tr.Start(ctx, "database layer")
	span.SetAttributes(attribute.Key("params").Int(int(id)))
	defer span.End()

	var car repository.Car
	err := r.db.Get(ctx, &car, "SELECT id,model,user_id,created_at,updated_at FROM cars WHERE id=$1", id)
	if err == sql.ErrNoRows {
		return nil, repository.ErrObjectNotFound
	}
	return &car, err
}

func (r *CarsRepo) List(ctx context.Context) ([]*repository.Car, error) {
	cars := make([]*repository.Car, 0)
	err := r.db.Select(ctx, &cars, "SELECT id, model, user_id, created_at, updated_at FROM cars")
	return cars, err
}

func (r *CarsRepo) Update(ctx context.Context, car *repository.Car) (bool, error) {
	tr := otel.Tracer("UpdateCar")
	ctx, span := tr.Start(ctx, "database layer")
	span.SetAttributes(attribute.Key("params").String(car.Model))
	defer span.End()

	result, err := r.db.Exec(ctx,
		"UPDATE cars SET model = $1 WHERE id = $2", car.Model, car.ID)
	return result.RowsAffected() > 0, err
}

func (r *CarsRepo) Delete(ctx context.Context, id int64) error {
	tr := otel.Tracer("DeleteCar")
	ctx, span := tr.Start(ctx, "database layer")
	span.SetAttributes(attribute.Key("params").Int(int(id)))
	defer span.End()

	_, err := r.db.Exec(ctx, "DELETE FROM cars WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
