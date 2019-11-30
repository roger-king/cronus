// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package tasker

import (
	"github.com/google/wire"
	"github.com/gorilla/mux"
	"github.com/robfig/cron/v3"
	"github.com/roger-king/tasker/handlers"
	"github.com/roger-king/tasker/models"
	"github.com/roger-king/tasker/services"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

import (
	_ "github.com/joho/godotenv/autoload"
)

// Injectors from tasker.go:

func New(tc models.TaskerConfig) (*Tasker, error) {
	client, err := services.NewMongoConnection(tc)
	if err != nil {
		return nil, err
	}
	cron := ProvideCron()
	taskService := services.NewTaskService(client, cron)
	settingService := services.NewSettingService(client)
	githubAuthService := services.NewGithubAuthService()
	router := handlers.NewRouter(taskService, settingService, githubAuthService, client)
	tasker := ProivdeTasker(tc, router, cron)
	return tasker, nil
}

// tasker.go:

// Tasker -
type Tasker struct {
	Config    models.TaskerConfig
	DB        *mongo.Client
	Scheduler *cron.Cron
	Router    *mux.Router
}

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

var TaskerSet = wire.NewSet(services.ServiceSet, handlers.RouterSet, ProvideCron, ProivdeTasker)

func ProvideCron() *cron.Cron {
	return cron.New()
}

func ProivdeTasker(tc models.TaskerConfig, r *mux.Router, c *cron.Cron) *Tasker {
	return &Tasker{
		Config:    tc,
		Scheduler: c,
		Router:    r,
	}
}

// Start - returns a mux router instance
func (t *Tasker) Start() *mux.Router {
	logrus.Info("Starting Tasker application")
	t.Scheduler.Start()

	return t.Router
}
