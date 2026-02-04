package storage

import(
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

//Opens a SQLite DB with sane defaults 
func Open(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("[DB] open db: %w", err)
	}

	//Verify connection early
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("[DB] ping db: %w", err)
	}

	//Enforce strict pragmas for correctness and concurrency behaviour
	pragmas := []string{
		"PRAGMA journal_mode = WAL;", 	//For better concurrency {It enables concurrent READS & WRITES}
		"PRAGMA foreign_keys = ON;", 	//For ForeignKey constraints
		"PRAGMA synchronous = NORMAL;", //For Good Durability/Perfor tradeoff
		"PRAGMA busy_timeout = 5000;",	//For waiting before DB is locked
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			return nil, fmt.Errorf("[DB] execution of pragma %q: %w", pragma, err)
		}
	}

	return db, nil
}