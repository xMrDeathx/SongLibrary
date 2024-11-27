package repositories

import (
	"SongLibrary/songlibrary/impl/domain/model"
	"context"
	"github.com/google/uuid"
)

type SongLibraryRepository interface {
	GetNextID() uuid.UUID
	GetSongs(context context.Context, filters model.Filters) ([]model.Song, error)
	DeleteSong(context context.Context, uuid uuid.UUID) error
	UpdateSong(context context.Context, song model.Song, songDetail model.SongDetail) error
	GetSong(context context.Context, uuid uuid.UUID) (model.Song, model.SongDetail, error)
	AddSong(context context.Context, song model.Song, songDetail model.SongDetail) error
}
