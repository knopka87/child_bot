package store

import "database/sql"

type Store struct {
	DB       *sql.DB
	Attempts *AttemptStore
	Villains *VillainStore
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		DB:       db,
		Attempts: NewAttemptStore(db),
		Villains: NewVillainStore(db),
	}
}
