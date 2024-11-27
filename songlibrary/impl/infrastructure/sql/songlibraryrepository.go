package sql

import (
	"SongLibrary/songlibrary/impl/domain/model"
	"SongLibrary/songlibrary/impl/domain/repositories"
	"SongLibrary/songlibrary/impl/infrastructure/sql/transactionwrapper"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"strconv"
	"time"
)

type songLibraryRepository struct {
	conn *pgxpool.Pool
}

func NewSongLibraryRepository(conn *pgxpool.Pool) repositories.SongLibraryRepository {
	log.Println("INFO: Creating new repository for song library")
	return &songLibraryRepository{
		conn: conn,
	}
}

func (repo *songLibraryRepository) GetNextID() uuid.UUID {
	return uuid.New()
}

func (repo *songLibraryRepository) GetSongs(ctx context.Context, filters model.Filters) ([]model.Song, error) {
	query := repo.buildQueryWithFilters(filters)

	log.Println("DEBUG: Select songs from database with filters")
	rows, err := repo.conn.Query(ctx, query)
	defer rows.Close()
	if errors.Is(err, pgx.ErrNoRows) {
		return []model.Song{}, nil
	} else if err != nil {
		return []model.Song{}, err
	}

	var songs []model.Song
	for rows.Next() {
		var song model.Song
		if err := rows.Scan(&song.ID, &song.Group, &song.Song); err != nil {
			return songs, err
		}

		songs = append(songs, song)
	}
	return songs, nil
}

func (repo *songLibraryRepository) GetSong(ctx context.Context, songId uuid.UUID) (model.Song, model.SongDetail, error) {
	var song model.Song
	var songDetail model.SongDetail

	log.Printf("DEBUG: Select song with ID: '%s' from database", songId)
	err := repo.conn.QueryRow(ctx, `
		SELECT id, band, song, release_date, text, link
		FROM song
		WHERE id=$1`, songId).Scan(
		&song.ID, &song.Group, &song.Song,
		&songDetail.ReleaseDate, &songDetail.Text, &songDetail.Link)
	if err != nil {
		return model.Song{}, model.SongDetail{}, err
	}

	return song, songDetail, nil
}

func (repo *songLibraryRepository) DeleteSong(ctx context.Context, songId uuid.UUID) error {
	wrapper := transactionwrapper.NewTransactionWrapper(repo.conn)
	err := wrapper.ExecuteWithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		now := time.Now()
		log.Printf("DEBUG: Delete song with ID: '%s' from database")
		_, err := tx.Exec(ctx, `
			UPDATE song
			SET deleted_at=$1
			WHERE id=$2`, now, songId)
		if err != nil {
			return err
		}

		log.Printf("INFO: Deleted song with ID: %s\n", songId)
		return nil
	})
	return err
}

func (repo *songLibraryRepository) UpdateSong(ctx context.Context, song model.Song, songDetail model.SongDetail) error {
	wrapper := transactionwrapper.NewTransactionWrapper(repo.conn)
	err := wrapper.ExecuteWithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		now := time.Now()
		log.Printf("DEBUG: Update song with ID: '%s' in database")
		_, err := tx.Exec(ctx, `
			UPDATE song
			SET band=$1, song=$2, release_date=$3, text=$4, link=$5, updated_at=$6
			WHERE id=$7`, song.Group, song.Song,
			songDetail.ReleaseDate, songDetail.Text, songDetail.Link, now, song.ID)
		if err != nil {
			return err
		}

		log.Printf("INFO: Updated song with ID: %s\n", song.ID)
		return nil
	})
	return err
}

func (repo *songLibraryRepository) AddSong(ctx context.Context, song model.Song, songDetail model.SongDetail) error {
	wrapper := transactionwrapper.NewTransactionWrapper(repo.conn)
	err := wrapper.ExecuteWithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		now := time.Now()
		log.Printf("DEBUG: Add new song in database")
		_, err := tx.Exec(ctx, `
			INSERT INTO song (id, band, song, release_date, text, link, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			song.ID, song.Group, song.Song, songDetail.ReleaseDate, songDetail.Text, songDetail.Link, now)
		if err != nil {
			return err
		}

		log.Println("INFO: Added new song into library")
		return nil
	})
	return err
}

func (repo *songLibraryRepository) buildQueryWithFilters(filters model.Filters) string {
	query := `
		SELECT id, band, song
		FROM song`

	filtersOn := false
	if *filters.Group != "" {
		if !filtersOn {
			filtersOn = true
			query += `
				WHERE band='` + *filters.Group + `'`
		} else {
			query += ` AND band='` + *filters.Group + `'`
		}
	}

	if *filters.Song != "" {
		if !filtersOn {
			filtersOn = true
			query += `
				WHERE song='` + *filters.Song + `'`
		} else {
			query += ` AND song='` + *filters.Song + `'`
		}
	}

	if *filters.ReleaseDate != "" {
		if !filtersOn {
			filtersOn = true
			query += `
				WHERE release_date='` + *filters.ReleaseDate + `'`
		} else {
			query += ` AND release_date='` + *filters.ReleaseDate + `'`
		}
	}

	if *filters.Text != "" {
		if !filtersOn {
			filtersOn = true
			query += `
				WHERE strpos(text, '` + *filters.Text + `') > 0`
		} else {
			query += ` AND strpos(text, '` + *filters.Text + `') > 0`
		}
	}

	if *filters.Link != "" {
		if !filtersOn {
			filtersOn = true
			query += `
				WHERE link='` + *filters.Link + `'`
		} else {
			query += ` AND link='` + *filters.Link + `'`
		}
	}

	pageNumber := *filters.Page
	if pageNumber < 1 {
		pageNumber = 1
	}

	limitNumber := *filters.Limit
	if limitNumber < 1 {
		limitNumber = 10
	}

	offset := (pageNumber - 1) * limitNumber

	query += `
		LIMIT ` + strconv.Itoa(limitNumber) + `
		OFFSET ` + strconv.Itoa(offset)

	return query
}
