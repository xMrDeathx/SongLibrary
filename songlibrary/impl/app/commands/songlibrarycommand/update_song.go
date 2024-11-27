package songlibrarycommand

import "github.com/google/uuid"

type UpdateSongCommand struct {
	ID          uuid.UUID
	Group       *string
	Song        *string
	ReleaseDate *string
	Text        *string
	Link        *string
}
