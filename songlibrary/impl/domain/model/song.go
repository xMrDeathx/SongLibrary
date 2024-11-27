package model

import "github.com/google/uuid"

type Song struct {
	ID    uuid.UUID
	Group string
	Song  string
}
