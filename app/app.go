package app

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"

	"github.com/danielkraic/idmapper/app/handlers"
	// pq imported to use with database/sql
	_ "github.com/lib/pq"
	logrus "github.com/sirupsen/logrus"
)

// App consists of application configurationand resources
type App struct {
	log           *logrus.Logger
	version       *handlers.Version
	configuration *Configuration
	redisClient   *redis.Client
	db            *sql.DB
}

// NewApp creates new App with its configuration and resources (log, redis, pqsql)
func NewApp(version string, commit string, build string, configFile string) (*App, error) {
	configuration, err := readConfiguration(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration: %s", err)
	}

	printConfiguration(configuration)

	log := logrus.New()
	if configuration.Logger.JSON {
		log.Formatter = new(logrus.JSONFormatter)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     configuration.Redis.Addr,
		Password: configuration.Redis.Password,
	})

	db, err := sql.Open("postgres", configuration.PostgreSQL.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	return &App{
		log: log,
		version: &handlers.Version{
			Version: version,
			Commit:  commit,
			Build:   build,
		},
		configuration: configuration,
		redisClient:   redisClient,
		db:            db,
	}, nil
}

func printConfiguration(c *Configuration) {
	fmt.Println("CONFIGURATION:")

	data, err := json.Marshal(c)
	if err != nil {
		fmt.Println("failed to print configuration")
	}

	fmt.Printf("%s\n", data)
}
