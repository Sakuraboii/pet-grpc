package internal

import (
	"github.com/prometheus/client_golang/prometheus"
)

var RegCarCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "new_cars",
	Help: "New Car was created",
})

var DeletedCarCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "deleted_cars",
	Help: "Car that was deleted",
})

var GetCarCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "get_cars",
	Help: "Get car by id",
})

var UpdatedCarCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "updated_cars",
	Help: "Car that was updated",
})

var RegUserCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "new_users",
	Help: "New User was created",
})

var DeletedUserCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "deleted_users",
	Help: "User that was deleted",
})

var GetUserCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "get_users",
	Help: "Get user by id",
})

var UpdatedUserCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "updated_users",
	Help: "User that was updated",
})

func init() {
	prometheus.MustRegister(RegCarCounter)
	prometheus.MustRegister(DeletedCarCounter)
	prometheus.MustRegister(GetCarCounter)
	prometheus.MustRegister(UpdatedCarCounter)
	prometheus.MustRegister(DeletedUserCounter)
	prometheus.MustRegister(RegUserCounter)
	prometheus.MustRegister(GetUserCounter)
	prometheus.MustRegister(UpdatedUserCounter)
}
