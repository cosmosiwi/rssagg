package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/cosmosiwi/rssagg/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct{
	DB *database.Queries
}

func main(){
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == ""{
		log.Fatal("PORT environment variable is not set")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == ""{
		log.Fatal("DB_URL environment variable is not set")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Can't connect to the database")
	}

	apiCfg := apiConfig{
		DB: database.New(conn), 
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
	
	v1router.Post("/users", apiCfg.handlerCreateUser)
	v1router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	
	v1router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1router.Get("/feeds", apiCfg.handlerGetFeeds)

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