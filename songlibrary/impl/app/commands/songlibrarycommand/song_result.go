package songlibrarycommand

import "github.com/google/uuid"

type SongResult struct {
	ID    uuid.UUID
	Group string
	Song  string
}
