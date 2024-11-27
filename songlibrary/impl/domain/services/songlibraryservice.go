package services

import (
	"SongLibrary/songlibrary/impl/app/commands/songlibrarycommand"
	"context"
	"github.com/google/uuid"
)

type SongLibraryService interface {
	GetSongs(context context.Context, command songlibrarycommand.GetSongsCommand) ([]songlibrarycommand.SongResult, error)
	DeleteSong(context context.Context, songId uuid.UUID) error
	UpdateSong(context context.Context, command songlibrarycommand.UpdateSongCommand) error
	GetSongText(context context.Context, songId uuid.UUID, page int, limit int) (songlibrarycommand.SongTextResult, error)
	AddSong(context context.Context, commandSong songlibrarycommand.AddSongCommand, commandSongDetail songlibrarycommand.AddSongDetailCommand) error
}
