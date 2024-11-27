package services

import (
	"SongLibrary/songlibrary/impl/app/commands/songlibrarycommand"
	"SongLibrary/songlibrary/impl/app/mapper/songlibrarymapper"
	"SongLibrary/songlibrary/impl/domain/model"
	"SongLibrary/songlibrary/impl/domain/repositories"
	"SongLibrary/songlibrary/impl/domain/services"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"log"
	"strings"
)

type songLibraryService struct {
	repository repositories.SongLibraryRepository
}

func NewSongLibraryService(repository repositories.SongLibraryRepository) services.SongLibraryService {
	log.Println("INFO: Creating new service for song library")
	return &songLibraryService{
		repository: repository,
	}
}

func (service *songLibraryService) GetSongs(context context.Context, command songlibrarycommand.GetSongsCommand) ([]songlibrarycommand.SongResult, error) {
	filters := songlibrarymapper.NewSongFiltersToDomainSongFilters(command)
	songs, err := service.repository.GetSongs(context, filters)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	log.Println("INFO: Got all songs from library with filters and pagination")
	return songlibrarymapper.NewSongsResultFromEntity(songs), nil
}

func (service *songLibraryService) DeleteSong(context context.Context, songId uuid.UUID) error {
	return service.repository.DeleteSong(context, songId)
}

func (service *songLibraryService) UpdateSong(context context.Context, command songlibrarycommand.UpdateSongCommand) error {
	var newSong model.Song
	var newSongDetail model.SongDetail
	oldSong, oldSongDetail, err := service.repository.GetSong(context, command.ID)
	if err != nil {
		return err
	}

	newSong.ID = command.ID

	if command.Group == nil {
		newSong.Group = oldSong.Group
	} else {
		newSong.Group = *command.Group
	}
	if command.Song == nil {
		newSong.Song = oldSong.Song
	} else {
		newSong.Song = *command.Song
	}
	if command.ReleaseDate == nil {
		newSongDetail.ReleaseDate = oldSongDetail.ReleaseDate
	} else {
		newSongDetail.ReleaseDate = *command.ReleaseDate
	}
	if command.Text == nil {
		newSongDetail.Text = oldSongDetail.Text
	} else {
		newSongDetail.Text = *command.Text
	}
	if command.Link == nil {
		newSongDetail.Link = oldSongDetail.Link
	} else {
		newSongDetail.Link = *command.Link
	}

	return service.repository.UpdateSong(context, newSong, newSongDetail)
}

func (service *songLibraryService) GetSongText(context context.Context, songId uuid.UUID, page int, limit int) (songlibrarycommand.SongTextResult, error) {
	song, songDetail, err := service.repository.GetSong(context, songId)
	if err != nil {
		return songlibrarycommand.SongTextResult{}, err
	}

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 1
	}

	verses := strings.Split(songDetail.Text, "\n\n")

	totalVerses := len(verses)
	if totalVerses == 0 {
		return songlibrarycommand.SongTextResult{}, errors.New("no verses found")
	}

	startIndex := (page - 1) * limit
	endIndex := startIndex + limit

	if startIndex >= totalVerses {
		return songlibrarycommand.SongTextResult{}, errors.New("page out of range")
	}

	if endIndex > totalVerses {
		endIndex = totalVerses
	}

	selectedVerses := verses[startIndex:endIndex]

	log.Println("INFO: Got song text from library with pagination")
	return songlibrarymapper.NewSongTextResultFromEntity(song, songDetail, selectedVerses), nil
}

func (service *songLibraryService) AddSong(context context.Context, commandSong songlibrarycommand.AddSongCommand, commandSongDetail songlibrarycommand.AddSongDetailCommand) error {
	song := model.Song{
		ID:    service.repository.GetNextID(),
		Group: commandSong.Group,
		Song:  commandSong.Song,
	}

	songDetail := model.SongDetail{
		ReleaseDate: commandSongDetail.ReleaseDate,
		Text:        commandSongDetail.Text,
		Link:        commandSongDetail.Link,
	}

	return service.repository.AddSong(context, song, songDetail)
}
