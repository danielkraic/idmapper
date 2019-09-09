package idmappers

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/danielkraic/idmapper/idmapper"
	"github.com/danielkraic/idmapper/scheduler"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

// Config configuration of IDMappers
type Config struct {
	Reloader struct {
		Currency struct {
			Interval      time.Duration `mapstructure:"interval"`
			RedisHashName string        `mapstructure:"redis_hash_name"`
		} `mapstructure:"currency"`
		Country struct {
			Interval time.Duration `mapstructure:"interval"`
		} `mapstructure:"country"`
		Language struct {
			Interval time.Duration `mapstructure:"interval"`
		} `mapstructure:"language"`
	} `mapstructure:"reloader"`
	Loader struct {
		Timeout time.Duration `mapstructure:"timeout"`
		URLs    struct {
			Currency string `mapstructure:"currency"`
			Country  string `mapstructure:"country"`
			Language string `mapstructure:"language"`
		} `mapstructure:"urls"`
	} `mapstructure:"loader"`
}

// IDMappers consists of available IDMapper objects
type IDMappers struct {
	config        *Config
	CurrencyCodes *idmapper.IDMapper
	CountryCodes  *idmapper.IDMapper
	LanguageCodes *idmapper.IDMapper
	// this mutex will prevent multiple IDMappers to be reloaded at the same time
	mtx      sync.Mutex
	reloader *scheduler.Scheduler
}

// NewIDMappers creates IDMappers with available IDMapper objects
func NewIDMappers(log *logrus.Logger, client *redis.Client, db *sql.DB, config *Config) (*IDMappers, error) {
	currencyCodes, err := NewRedisIDMapper(client, config.Reloader.Currency.RedisHashName)
	if err != nil {
		return nil, fmt.Errorf("failed to create IDMapper for currency codes: %s", err)
	}

	countryCodes, err := NewPgSQLIDMapper(log, db, "select id, name from country")
	if err != nil {
		return nil, fmt.Errorf("failed to create IDMapper for country codes: %s", err)
	}

	languageCodes, err := NewHTTPIDMapper(log, config.Loader.URLs.Language, config.Loader.Timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to create IDMapper for language codes: %s", err)
	}

	return &IDMappers{
		config:        config,
		CurrencyCodes: currencyCodes,
		CountryCodes:  countryCodes,
		LanguageCodes: languageCodes,
		reloader:      &scheduler.Scheduler{},
	}, nil
}

// RunReloader starts scheduler for automatic reloading of IDMapper objects
func (idMappers *IDMappers) RunReloader(log *logrus.Logger) {
	logOperation := func(description string, err error) {
		if err != nil {
			log.Errorf("%s failed: %s", description, err)
		} else {
			log.Infof("%s was successful", description)
		}
	}

	logOperation("setup of CurrencyCodes reloading", idMappers.reloader.AddFunc(func() {
		idMappers.mtx.Lock()
		defer idMappers.mtx.Unlock()
		logOperation("reload of CurrencyCodes", idMappers.CurrencyCodes.Reload())
	}, idMappers.config.Reloader.Currency.Interval))

	logOperation("setup of CountryCodes reloading", idMappers.reloader.AddFunc(func() {
		idMappers.mtx.Lock()
		defer idMappers.mtx.Unlock()
		logOperation("reload of CountryCodes", idMappers.CountryCodes.Reload())
	}, idMappers.config.Reloader.Country.Interval))

	logOperation("setup of LanguageCodes reloading", idMappers.reloader.AddFunc(func() {
		idMappers.mtx.Lock()
		defer idMappers.mtx.Unlock()
		logOperation("reload of LanguageCodes", idMappers.LanguageCodes.Reload())
	}, idMappers.config.Reloader.Language.Interval))

	go idMappers.reloader.Start()
}

// StopReloader stops scheduler for automatic reloading of IDMapper objects
func (idMappers *IDMappers) StopReloader() {
	idMappers.reloader.Stop()
}
