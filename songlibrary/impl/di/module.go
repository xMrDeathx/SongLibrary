package di

import (
	frontendapi "SongLibrary/songlibrary/api/frontend"
	"SongLibrary/songlibrary/cmd"
	"SongLibrary/songlibrary/impl/app/services"
	"SongLibrary/songlibrary/impl/infrastructure/sql"
	"SongLibrary/songlibrary/impl/infrastructure/transport"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
	"time"
)

func InitSongLibraryModule(config cmd.Config) (
	*pgxpool.Pool,
	error,
) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)

	conn, _ := ConnectLoop(connStr, 30*time.Second)

	songLibraryRepository := sql.NewSongLibraryRepository(conn)
	songLibraryService := services.NewSongLibraryService(songLibraryRepository)
	songLibraryServer := transport.NewSongLibraryServer(songLibraryService)

	router := mux.NewRouter()

	options := frontendapi.GorillaServerOptions{
		BaseRouter: router,
		Middlewares: []frontendapi.MiddlewareFunc{func(handler http.Handler) http.Handler {
			return http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.ServeHTTP(w, r)
			}))
		}},
	}
	r := frontendapi.HandlerWithOptions(songLibraryServer, options)
	http.Handle("/songlibrary/", r)

	return conn, nil
}

func ConnectLoop(connStr string, timeout time.Duration) (*pgxpool.Pool, error) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeoutExceeded := time.After(timeout)
	for {
		select {
		case <-timeoutExceeded:
			return nil, fmt.Errorf("db connection failed after %s timeout", timeout)

		case <-ticker.C:
			db, err := pgxpool.Connect(context.Background(), connStr)
			if err == nil {
				return db, nil
			}
		}
	}
}
