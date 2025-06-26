package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
	"urlshort/internal/database"
	"urlshort/internal/handlers"
	"urlshort/internal/middleware"
)

func main() {
	var fp string
	var useDB bool

	flag.StringVar(&fp, "f", "", "location of YAML or JSON file")
	flag.BoolVar(&useDB, "db", false, "use in-memory database")
	flag.Parse()

	var handler http.Handler
	if useDB {
		db := database.InitDB(getDataPath("redirects.db"))
		defer db.Close()
		err := database.SeedIfEmpty(db, getDataPath("redirects.json"))
		if err != nil {
			log.Printf("error seeding db: %v", err)
		}
		handler = handlers.BoltHandler(db, defaultMux())
	} else {
		fp = getDataPath(fp)
		if err := validatePath(fp); err != nil {
			log.Printf("error validating file: %v", err)
			log.Println("falling back to default config")
			fp = ""
		}
		format := inferExtension(fp)
		var err error
		handler, err = setupHandler(fp, format)
		if err != nil {
			log.Fatalf("failed to initialize handler: %v", err)
		}

	}

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Println("Server running on http://localhost:8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func setupHandler(path, format string) (http.Handler, error) {
	fallback := defaultMux()
	var handler http.HandlerFunc
	var err error

	switch format {
	case "yaml", "yml":
		handler, err = handlers.YamlHandler(path, fallback)
	case "json":
		handler, err = handlers.JsonHandler(path, fallback)
	default:
		log.Printf("unsupported format: %s, using fallback handler", format)
		handler = handlers.MapHandler(map[string]string{
			"/urlshort": "https://github.com/gophercises/urlshort",
			"/godoc":    "https://pkg.go.dev/std",
		}, fallback)
	}

	if err != nil {
		return nil, err
	}

	return middleware.LoggingMiddleware(handler), nil
}

func validatePath(path string) error {
	if path == "" {
		return nil
	}

	// validate extension is correct
	validExtensions := map[string]bool{
		".json": true,
		".yaml": true,
		".yml":  true,
	}

	if !validExtensions[filepath.Ext(path)] {
		return fmt.Errorf("malformed file path: %s contains incorrect extension", path)
	}

	//check that file exists
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %w", err)
	} else if err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}

	return nil
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, this is the fallback handler.")
	})
	return mux
}

func inferExtension(path string) string {
	switch ext := filepath.Ext(path); ext {
	case ".yaml", ".yml":
		return "yaml"
	case ".json":
		return "json"
	default:
		return ""
	}
}

// DataPath returns the absolute path to a file inside the /data directory.
func getDataPath(filename string) string {
	_, baseFile, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(baseFile)
	// Adjust this if your structure differs (e.g., utils is deeper)
	projectRoot := filepath.Join(baseDir, "..", "..")
	return filepath.Join(projectRoot, "data", filename)
}
