package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://127.0.0.1:3000", "http://0.0.0.0:3000"},
		AllowedMethods: []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE", "HEAD"},
		AllowedHeaders: []string{"*"},

		AllowCredentials: true,
	}))
	fileRoutes := FileUploadHandler()

	go func() {
		for {
			event := <-fileRoutes.CompleteUploads
			fmt.Printf("\nSuccessfully uploaded a file with an id: %v", event.Upload.ID)
			createdE := <-fileRoutes.CreatedUploads
			fmt.Printf("\nCreated an upload with an id: %v", createdE.Upload.ID)
		}
	}()

	r.Post("/store-products-csv/", fileRoutes.PostFile)
	r.Head("/store-products-csv/{id}", fileRoutes.HeadFile)
	r.Patch("/store-products-csv/{id}", fileRoutes.PatchFile)
	r.Get("/store-products-csv/{id}", fileRoutes.GetFile)

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
