package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/nhatquang342/chirpy/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform	   string
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	pf := os.Getenv("PLATFORM")
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	dbQueries := database.New(dbConn)
    apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db: 			dbQueries,
		platform:		pf,
	}

	// Step 1: Create a new ServeMux
	mux := http.NewServeMux()

	// Step 2: Create a FileServer handler for the current directory
	fs := http.FileServer(http.Dir(filepathRoot))

	// Step 3: Handle path (root "/", or sth else) with the file server
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", (fs))))
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerDeleteAllUsers)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	//mux.HandleFunc("POST /api/validate_chirp", handlerValidate)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpByID)

	// Step 4: Create the server struct
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Step 5: Start the server
	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}