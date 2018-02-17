package persistence

import (
	"database/sql"
	_ "github.com/lib/pq"  // it's pq; shut it, golint
	log "github.com/sirupsen/logrus"
    "fmt"
)

type Storage struct {
}

// SaveUnknownResolvedQuery saves a resolved query into the Postgres DB
func (s *Storage) SaveUnknownResolvedQuery(rq string) int {
	db, err := sql.Open("postgres", connectionString)
    if err != nil {
        log.Errorf("Could not connect to Postgres! Error: %s", err)
        return -1
    }
    var insertedID int
    query := fmt.Sprintf("INSERT INTO beardbot.missed_queries(resolved_query) VALUES(%s) RETURNING id", rq)
    err = db.QueryRow(query).Scan(&insertedID)
    if err != nil {
        return -1
    }
    return insertedID
}
