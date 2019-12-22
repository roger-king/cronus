// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package tasker

import (
	"github.com/google/wire"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron/v3"
	"github.com/roger-king/tasker/config"
	"github.com/roger-king/tasker/services"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"strings"
)

import (
	_ "github.com/joho/godotenv/autoload"
)

// Injectors from tasker.go:

func New(tc *config.TaskerConfig) (*Tasker, error) {
	cron := ProvideCron()
	db := services.NewDBConnection(tc)
	tasker := ProivdeTasker(tc, cron, db)
	return tasker, nil
}

// tasker.go:

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

// Tasker -
type Tasker struct {
	Config    *config.TaskerConfig
	DB        *mongo.Client
	Scheduler *cron.Cron
	Router    *mux.Router
}

var TaskerSet = wire.NewSet(services.ServiceSet, ProvideCron, ProivdeTasker)

func ProvideCron() *cron.Cron {
	return cron.New()
}

func ProivdeTasker(tc *config.TaskerConfig, c *cron.Cron, db *sqlx.DB) *Tasker {
	if tc.Auth {
		if len(tc.GithubClientID) == 0 && len(tc.GithubClientSecret) == 0 {
			logrus.Fatal("Authentication is enabled. Please provide the github client id and secret.")
			os.Exit(1)
		}
	}

	if len(tc.DBConnectionURL) > 0 {
		if !strings.Contains(tc.DBConnectionURL, "postgres") {
			logrus.Fatal("Please provide a valid postgres db connection")
			os.Exit(1)
		}
	} else {
		logrus.Fatal("DBConnectionURL is required")
		os.Exit(1)
	}

	return &Tasker{
		Config:    tc,
		Scheduler: c,
	}
}

// Start - returns a mux router instance
func (t *Tasker) Start() *mux.Router {
	logrus.Info("Starting Tasker application")
	t.Scheduler.Start()

	return t.Router
}
