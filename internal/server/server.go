package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/tus/tusd/v2/pkg/gcsstore"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/tus/tusd/v2/pkg/memorylocker"
	"golang.org/x/exp/slog"
)

type Server struct {
	port int
}

var LOCALHOST = regexp.MustCompile(`/^https?:\/\/\w+(\.\w+)*(:[0-9]+)?(\/.*)?$/`)

func FileUploadHandler() *tusd.UnroutedHandler {
	bucket := "sample-bucket"
	absPath, _ := filepath.Abs("./")
	gcpCredentials := filepath.Join(absPath, "config", "credentials.json")

	gcsService, err := gcsstore.NewGCSService(gcpCredentials)
	if err != nil {
		log.Fatalf("Unable to create gcsService: %s\n", err.Error())
		return nil
	}
	locker := memorylocker.New()

	store := gcsstore.New(bucket, gcsService)
	store.ObjectPrefix = "store-products-csv"

	composer := tusd.NewStoreComposer()
	store.UseIn(composer)
	locker.UseIn(composer)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	handler, err := tusd.NewUnroutedHandler(tusd.Config{
		BasePath:              "/store-products-csv/",
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
		NotifyCreatedUploads:  true,
		Logger:                logger,
		Cors: &tusd.CorsConfig{
			Disable: true,
		},
	})

	if err != nil {
		log.Fatalf("could not create tusd handler: %s\n", err.Error())
		return nil
	}

	return handler

}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
