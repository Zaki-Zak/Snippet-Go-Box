package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql" // New import
	"github.com/joho/godotenv"
)

type application struct {
	logger *slog.Logger
}
type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return f, nil
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to datebase: %w", err)
	}
	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)
	// Verify the connction
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping datebase: %w", err)
	}
	return db, nil
}

func getDefaultDSN() string {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbParseTime := os.Getenv("DB_PARSE_TIME")
	return fmt.Sprintf("%s:%s@/%s?parseTime=%s", dbUser, dbPassword, dbName, dbParseTime)
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP nerwork address")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := flag.String("dsn", getDefaultDSN(), "MySQL data source name")

	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	app := &application{
		logger: logger,
	}

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error("failed to open database", "dbError", err)
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("starting server on", "addr", *addr)
	err = http.ListenAndServe(*addr, app.routes())
	logger.Error("server error", "error", err)
	os.Exit(1)
}
