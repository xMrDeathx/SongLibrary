package songlibrarymapper

import (
	"SongLibrary/songlibrary/impl/app/commands/songlibrarycommand"
	"SongLibrary/songlibrary/impl/domain/model"
)

func NewSongDetailResultFromEntity(songDetail model.SongDetail) songlibrarycommand.SongDetailResult {
	return songlibrarycommand.SongDetailResult{
		ReleaseDate: songDetail.ReleaseDate,
		Text:        songDetail.Text,
		Link:        songDetail.Link,
	}
}

func NewSongFiltersToDomainSongFilters(filters songlibrarycommand.GetSongsCommand) model.Filters {
	return model.Filters{
		Group:       filters.Group,
		Song:        filters.Song,
		ReleaseDate: filters.ReleaseDate,
		Text:        filters.Text,
		Link:        filters.Link,
		Page:        filters.Page,
		Limit:       filters.Limit,
	}
}

func NewSongResultFromEntity(song model.Song) songlibrarycommand.SongResult {
	return songlibrarycommand.SongResult{
		ID:    song.ID,
		Group: song.Group,
		Song:  song.Song,
	}
}

func NewSongsResultFromEntity(songs []model.Song) []songlibrarycommand.SongResult {
	result := make([]songlibrarycommand.SongResult, 0, len(songs))

	for _, v := range songs {
		songResult := NewSongResultFromEntity(v)
		result = append(result, songResult)
	}

	return result
}

func NewVersesResultFromEntity(verses []string) []string {
	var result []string

	for _, v := range verses {
		result = append(result, v)
	}

	return result
}

func NewSongTextResultFromEntity(song model.Song, songDetail model.SongDetail, verses []string) songlibrarycommand.SongTextResult {
	return songlibrarycommand.SongTextResult{
		SongInfo:    NewSongResultFromEntity(song),
		ReleaseDate: songDetail.ReleaseDate,
		Text:        NewVersesResultFromEntity(verses),
		Link:        songDetail.Link,
	}
}
