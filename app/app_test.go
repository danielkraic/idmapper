package app_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis"
	"github.com/danielkraic/idmapper/app"
	"github.com/danielkraic/idmapper/app/handlers"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const (
	AppVersion    = "testversion"
	AppCommit     = "testcommit"
	AppBuild      = "testbuild"
	AppConfigFile = ""
)

var (
	currencyCodes = map[string]interface{}{
		"usd": "dollar",
		"eur": "euro",
	}
	countryCodes = map[string]string{
		"us": "USA",
		"sk": "Slovakia",
	}
	languageCodes = map[string]string{
		"en": "English",
		"sk": "Slovak",
	}
)

func createHTTPTestServer(languages map[string]string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var response []handlers.IDMapperResponse
		for k, v := range languages {
			response = append(response, handlers.IDMapperResponse{ID: k, Name: v})
		}

		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
}

func createRedisMock(mr *miniredis.Miniredis, currenciesHashName string, currencies map[string]interface{}) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	redisClient.HMSet(currenciesHashName, currencies)
	return redisClient
}

func createPostgreSQLMock(countries map[string]string) (*sql.DB, *sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	rows := sqlmock.NewRows([]string{"id", "name"})
	for k, v := range countries {
		rows.AddRow(k, v)
	}
	mock.ExpectQuery("select id, name from country").WillReturnRows(rows)

	return db, &mock, err
}

type TestApp struct {
	App            *app.App
	Miniredis      *miniredis.Miniredis
	HTTPTestServer *httptest.Server
	SQL            *sql.DB
	SQLMock        *sqlmock.Sqlmock
}

func NewTestApp() (*TestApp, error) {
	app, err := app.NewApp(AppVersion, AppCommit, AppBuild, AppConfigFile)
	if err != nil {
		return nil, fmt.Errorf("failed to app: %s", err)
	}

	// http server
	server := createHTTPTestServer(languageCodes)
	app.Configuration.IDMappers.Loader.URLs.Language = server.URL

	// redis
	mr, err := miniredis.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to create miniredis: %s", err)
	}
	app.RedisClient = createRedisMock(mr, app.Configuration.IDMappers.Reloader.Currency.RedisHashName, currencyCodes)

	// pgsql
	db, mock, err := createPostgreSQLMock(countryCodes)
	if err != nil {
		return nil, fmt.Errorf("failed to create PostgreSQL mock: %s", err)
	}
	app.DB = db

	// idmappers
	err = app.SetupIDMappers()
	if err != nil {
		return nil, fmt.Errorf("failed to setup IDMappers for TestApp: %s", err)
	}

	return &TestApp{
		App:            app,
		Miniredis:      mr,
		HTTPTestServer: server,
		SQL:            db,
		SQLMock:        mock,
	}, nil
}

func (testApp *TestApp) Close() {
	testApp.Miniredis.Close()
	err := testApp.App.DB.Close()
	_ = err // ignore close error for DB.Close()
}

func TestAppVersion(t *testing.T) {
	testApp, err := NewTestApp()
	assert.Nil(t, err)
	defer testApp.Close()

	req, err := http.NewRequest(http.MethodGet, "/version", nil)
	assert.Nil(t, err)
	resp := httptest.NewRecorder()

	h := handlers.NewVersionHandler(testApp.App.Version)
	h.ServeHTTP(resp, req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)

	data, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	var version handlers.Version
	err = json.Unmarshal(data, &version)
	assert.Nil(t, err)
	assert.Equal(t, version.Version, AppVersion)
	assert.Equal(t, version.Commit, AppCommit)
	assert.Equal(t, version.Build, AppBuild)
}

func TestAppCountry(t *testing.T) {
	testApp, err := NewTestApp()
	assert.Nil(t, err)
	defer testApp.Close()

	req, err := http.NewRequest(http.MethodGet, "/v1/country/", nil)
	assert.Nil(t, err)
	req = mux.SetURLVars(req, map[string]string{"id": "sk"})

	resp := httptest.NewRecorder()

	h := handlers.NewIDMapperHandler(testApp.App.IDMappers.CountryCodes)
	h.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	data, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	var idMapperResponse handlers.IDMapperResponse
	err = json.Unmarshal(data, &idMapperResponse)
	assert.Nil(t, err)
	assert.Equal(t, idMapperResponse.ID, "sk")
	assert.Equal(t, idMapperResponse.Name, "Slovakia")
}

func TestAppCurrency(t *testing.T) {
	testApp, err := NewTestApp()
	assert.Nil(t, err)
	defer testApp.Close()

	req, err := http.NewRequest(http.MethodGet, "/v1/currency/", nil)
	assert.Nil(t, err)
	req = mux.SetURLVars(req, map[string]string{"id": "eur"})

	resp := httptest.NewRecorder()

	h := handlers.NewIDMapperHandler(testApp.App.IDMappers.CurrencyCodes)
	h.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	data, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	var idMapperResponse handlers.IDMapperResponse
	err = json.Unmarshal(data, &idMapperResponse)
	assert.Nil(t, err)
	assert.Equal(t, idMapperResponse.ID, "eur")
	assert.Equal(t, idMapperResponse.Name, "euro")
}

func TestAppLanguage(t *testing.T) {
	testApp, err := NewTestApp()
	assert.Nil(t, err)
	defer testApp.Close()

	req, err := http.NewRequest(http.MethodGet, "/v1/language/", nil)
	assert.Nil(t, err)
	req = mux.SetURLVars(req, map[string]string{"id": "sk"})

	resp := httptest.NewRecorder()

	h := handlers.NewIDMapperHandler(testApp.App.IDMappers.LanguageCodes)
	h.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	data, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	var idMapperResponse handlers.IDMapperResponse
	err = json.Unmarshal(data, &idMapperResponse)
	assert.Nil(t, err)
	assert.Equal(t, idMapperResponse.ID, "sk")
	assert.Equal(t, idMapperResponse.Name, "Slovak")
}

func TestAppLanguageNotFound(t *testing.T) {
	testApp, err := NewTestApp()
	assert.Nil(t, err)
	defer testApp.Close()

	req, err := http.NewRequest(http.MethodGet, "/v1/language/", nil)
	assert.Nil(t, err)
	req = mux.SetURLVars(req, map[string]string{"id": "non-existing"})

	resp := httptest.NewRecorder()

	h := handlers.NewIDMapperHandler(testApp.App.IDMappers.LanguageCodes)
	h.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestAppHealth(t *testing.T) {
	testApp, err := NewTestApp()
	assert.Nil(t, err)
	defer testApp.Close()

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	assert.Nil(t, err)
	resp := httptest.NewRecorder()

	h := handlers.NewVersionHandler(testApp.App.Version)
	h.ServeHTTP(resp, req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
}
