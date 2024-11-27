package transport

import (
	"SongLibrary/core/utils"
	frontendapi "SongLibrary/songlibrary/api/frontend"
	"SongLibrary/songlibrary/impl/app/commands/songlibrarycommand"
	"SongLibrary/songlibrary/impl/domain/services"
	"SongLibrary/songlibrary/impl/infrastructure/transport/mapper"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"io"
	"log"
	"net/http"
	"net/url"
)

type songLibraryServer struct {
	songLibraryService services.SongLibraryService
}

func NewSongLibraryServer(
	songLibraryService services.SongLibraryService,
) frontendapi.ServerInterface {
	return &songLibraryServer{
		songLibraryService: songLibraryService,
	}
}

func (server *songLibraryServer) AddSong(w http.ResponseWriter, r *http.Request, params frontendapi.AddSongParams) {
	song := songlibrarycommand.AddSongCommand{
		Group: params.Group,
		Song:  params.Song,
	}

	log.Println("INFO: Getting song details from external API")
	songDetail, err := server.getSongDetailFromApi(song.Group, song.Song)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("INFO: Adding new song into library")
	err = server.songLibraryService.AddSong(r.Context(), song, songDetail)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (server *songLibraryServer) GetSongs(w http.ResponseWriter, r *http.Request, params frontendapi.GetSongsParams) {
	filters := songlibrarycommand.GetSongsCommand{
		Group:       params.Group,
		Song:        params.Song,
		ReleaseDate: params.ReleaseDate,
		Text:        params.Text,
		Link:        params.Link,
		Page:        params.Page,
		Limit:       params.Limit,
	}

	log.Println("INFO: Getting songs from library")
	songs, err := server.songLibraryService.GetSongs(r.Context(), filters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	songsJson := utils.ListMaker(songs, mapper.MapSongToJson)
	response, err := json.Marshal(frontendapi.SongsResponse{Songs: songsJson})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (server *songLibraryServer) GetSongText(w http.ResponseWriter, r *http.Request, songId openapi_types.UUID, params frontendapi.GetSongTextParams) {
	if songId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("INFO: Getting song text from library")
	songTextResponse, err := server.songLibraryService.GetSongText(r.Context(), songId, params.Page, params.Limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(frontendapi.GetSongTextResponse{
		Song:        mapper.MapSongToJson(songTextResponse.SongInfo),
		ReleaseDate: songTextResponse.ReleaseDate,
		Text:        songTextResponse.Text,
		Link:        songTextResponse.Link,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (server *songLibraryServer) UpdateSong(w http.ResponseWriter, r *http.Request, songId openapi_types.UUID, params frontendapi.UpdateSongParams) {
	if songId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	song := songlibrarycommand.UpdateSongCommand{
		ID:          songId,
		Group:       params.Group,
		Song:        params.Song,
		ReleaseDate: params.ReleaseDate,
		Text:        params.Text,
		Link:        params.Link,
	}

	log.Println("INFO: Updating song info in library")
	err := server.songLibraryService.UpdateSong(r.Context(), song)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (server *songLibraryServer) DeleteSong(w http.ResponseWriter, r *http.Request, songId openapi_types.UUID) {
	if songId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("INFO: Deleting song from library")
	err := server.songLibraryService.DeleteSong(r.Context(), songId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (server *songLibraryServer) getSongDetailFromApi(group string, song string) (songlibrarycommand.AddSongDetailCommand, error) {
	encodedGroup := url.QueryEscape(group)
	encodedSong := url.QueryEscape(song)
	apiURL := fmt.Sprintf("http://localhost:8080/info?group=%s&song=%s", encodedGroup, encodedSong)

	response, err := http.Get(apiURL)
	if err != nil {
		return songlibrarycommand.AddSongDetailCommand{}, errors.New("failed to request external API")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return songlibrarycommand.AddSongDetailCommand{}, errors.New("failed to get song details from external API")
	}

	var songDetail songlibrarycommand.AddSongDetailCommand
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return songlibrarycommand.AddSongDetailCommand{}, errors.New("failed to read external API response")
	}

	if err = json.Unmarshal(responseBody, &songDetail); err != nil {
		return songlibrarycommand.AddSongDetailCommand{}, err
	}

	return songDetail, nil
}
