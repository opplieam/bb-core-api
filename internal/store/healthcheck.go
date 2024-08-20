package store

import (
	"database/sql"
)

type HealthCheckStore struct {
	DB *sql.DB
}

func NewHealthCheckStore(db *sql.DB) *HealthCheckStore {
	return &HealthCheckStore{DB: db}
}

func (s *HealthCheckStore) HealthCheck() (bool, error) {
	tmp := false
	rows, err := s.DB.Query(`SELECT true`)
	if err != nil {
		return false, err
	}
	for rows.Next() {
		err = rows.Scan(&tmp)
		if err != nil {
			return false, err
		}
	}

	return tmp, nil
}
