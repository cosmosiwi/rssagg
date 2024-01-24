package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main(){
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == ""{
		log.Fatal("PORT environment variable is not set")
	}

	router := chi.NewRouter()
	
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"PUT", "POST", "GET", "DELETE", "OPTION"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300,
	}))

	v1router := chi.NewRouter()
	v1router.Get("/healthz", handlerReadiness)
	v1router.Get("/err", handlerErr)
	

	router.Mount("/v1", v1router)
	
	svc := &http.Server{
		Addr: ":" + portString,
		Handler: router,
	}

	log.Printf("Serving on port: %s", portString)
	if err := svc.ListenAndServe(); err != nil{
		log.Fatal(err)
	}
}