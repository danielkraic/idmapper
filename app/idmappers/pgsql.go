package idmappers

import (
	"database/sql"
	"fmt"

	"github.com/danielkraic/idmapper/idmapper"
	"github.com/sirupsen/logrus"
)

// NewPgSQLIDMapper creates IDMapper that reads data from sql database
func NewPgSQLIDMapper(log *logrus.Logger, db *sql.DB, query string) (*idmapper.IDMapper, error) {
	if db == nil {
		return nil, fmt.Errorf("failed to create PgSQL IDMapper: sql.DB is nil")
	}
	return idmapper.NewIDMapper(&pgSQLSource{
		log:   log,
		query: query,
		db:    db,
	})
}

type pgSQLSource struct {
	log   *logrus.Logger
	query string
	db    *sql.DB
}

func (source *pgSQLSource) Read() (idmapper.ValuesMap, error) {
	rows, err := source.db.Query(source.query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query '%s': %s", source.query, err)
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			source.log.Errorf("failed to close pq rows")
		}
	}()

	var id, name string
	result := make(idmapper.ValuesMap)

	for rows.Next() {
		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows: %s", err)
		}

		result[id] = name
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows scan failed: %s", err)
	}

	return result, nil
}
