package mapper

import (
	frontendapi "SongLibrary/songlibrary/api/frontend"
	"SongLibrary/songlibrary/impl/app/commands/songlibrarycommand"
)

func MapSongToJson(song songlibrarycommand.SongResult) frontendapi.Song {
	return frontendapi.Song{
		Id:    song.ID,
		Group: song.Group,
		Song:  song.Song,
	}
}
