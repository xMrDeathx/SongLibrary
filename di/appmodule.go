package di

import (
	"SongLibrary/songlibrary/cmd"
	"SongLibrary/songlibrary/impl/di"
	"log"
	"net/http"
	"os"
)

func InitAppModule() {
	songLibraryConn, err := di.InitSongLibraryModule(cmd.NewConfig())
	defer songLibraryConn.Close()

	appPort := os.Getenv("BACKEND_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	log.Println("INFO: Start song library server")
	err = http.ListenAndServe(":"+appPort, nil)
	if err != nil {
		log.Panic("ListenAndServe: ", err)
	}
}

func Migrate() {
	err := cmd.Migrate(cmd.NewConfig())
	if err != nil {
		log.Fatal("Failed to migrate song library module:", err)
	}

	log.Println("INFO: Database migrations completed")
}
