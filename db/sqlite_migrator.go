package db

import (
	"database/sql"
	"errors"
)

type SqliteMigrator struct {
	migrations []SqliteMigration
}

type SqliteMigration struct {
	statements []string
}

func New_migrator() SqliteMigrator {
	return SqliteMigrator{
		migrations: make([]SqliteMigration, 0),
	}
}

func New_migration(statements []string) SqliteMigration {
	return SqliteMigration{
		statements: statements,
	}
}

func (s *SqliteMigrator) AddMigration(migration SqliteMigration) {
	s.migrations = append(s.migrations, migration)
}

func (s *SqliteMigrator) Run(db *sql.DB) error {
	var err error

	row := db.QueryRow("PRAGMA user_version")
	if row == nil {
		return errors.New("could not retrieve user version")
	}
	var starting_version int
	if err = row.Scan(&starting_version); err != nil {
		return err
	}

	for i, migration := range s.migrations {
		version := i + 1
		if version <= starting_version {
			continue
		}

		for _, statement := range migration.statements {
			if _, err = db.Exec(statement); err != nil {
				return err
			}
		}

		if _, err = db.Exec("PRAGMA user_version=?", version); err != nil {
			return err
		}
	}

	return nil
}
