package app

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"

	"github.com/danielkraic/idmapper/app/handlers"
	"github.com/danielkraic/idmapper/app/idmappers"

	// pq imported to use with database/sql
	_ "github.com/lib/pq"
	logrus "github.com/sirupsen/logrus"
)

// App consists of application configurationand resources
type App struct {
	log           *logrus.Logger
	Version       *handlers.Version
	Configuration *Configuration
	RedisClient   *redis.Client
	DB            *sql.DB
	IDMappers     *idmappers.IDMappers
}

// NewApp creates new App with its configuration and logger
func NewApp(version string, commit string, build string, configFile string) (*App, error) {
	configuration, err := readConfiguration(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration: %s", err)
	}

	log := logrus.New()
	if configuration.Logger.JSON {
		log.Formatter = new(logrus.JSONFormatter)
	}

	return &App{
		log: log,
		Version: &handlers.Version{
			Version: version,
			Commit:  commit,
			Build:   build,
		},
		Configuration: configuration,
	}, nil
}

// SetupRedis creates redis client (useful during testing when using redis mock)
func (app *App) SetupRedis() {
	app.RedisClient = redis.NewClient(&redis.Options{
		Addr:     app.Configuration.Redis.Addr,
		Password: app.Configuration.Redis.Password,
	})
}

// SetupPostgreSQL opens connection to PostgresSQL (useful during testing when using sql mock)
func (app *App) SetupPostgreSQL() error {
	db, err := sql.Open("postgres", app.Configuration.PostgreSQL.ConnectionString)
	if err != nil {
		return err
	}

	app.DB = db
	return nil
}

// SetupIDMappers creates IDMappers
func (app *App) SetupIDMappers() error {
	idMappers, err := idmappers.NewIDMappers(app.log, app.RedisClient, app.DB, &app.Configuration.IDMappers)
	if err != nil {
		return fmt.Errorf("failed to create IDMappers: %s", err)
	}

	app.IDMappers = idMappers
	return nil
}

// PrintConfiguration prints configuration to stdout
func (app App) PrintConfiguration() {
	data, err := json.MarshalIndent(app.Configuration, "", " ")
	if err != nil {
		fmt.Println("failed to print configuration")
	}

	fmt.Printf("%s\n", data)
}
