package main

import (
	"coinChat/pkg/httpgateway"
	"coinChat/pkg/repository"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	postgresMigration "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/cors"

	"github.com/jmoiron/sqlx"
)

func main() {
	cfg := Config{
		ConfigPostgres: repository.ConfigPostgres{
			Address:              ":5432",
			Username:             "postgres",
			Database:             "coinchat",
			Password:             "password",
			MaxConns:             10,
			MaxWaitForConnection: time.Second * 5,
			Schema:               "dev",
		},
	}
	ctx := context.Background()
	db, dbCloser, err := initDB(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to initiate the db: %+v", err)
		return
	}
	defer dbCloser()

	err = performDBMigration(db, true)
	if err != nil {
		log.Fatalf("failed to perform db migration: %+v", err)
		return
	}

	dbx := sqlx.NewDb(db, "postgres")
	userRepo := repository.NewUserRepo(dbx)
	messageRepo := repository.NewMessageRepo(dbx)
	serverHandler := httpgateway.NewHandler(userRepo, messageRepo)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})
	s := &http.Server{
		Addr:           ":8080",
		Handler:        c.Handler(serverHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	quit := InitGracefulShutdown()
	go func() {
		s.ListenAndServe()
	}()
	<-quit
	fmt.Printf("ended with error: %+v", performDBMigration(db, false))
}

func InitGracefulShutdown() chan struct{} {
	quit := make(chan struct{})
	go func() {
		// We initiate a graceful shutdowns when quit via
		//  SIGINT  (Ctrl+C)
		//  SIGTERM (kill -15)
		cSig := make(chan os.Signal, 1)
		signal.Notify(cSig, os.Interrupt)
		signal.Notify(cSig, syscall.SIGTERM)

		// we block until we get a signal
		<-cSig

		// signal quit
		close(quit)
	}()
	return quit
}

func initDB(ctx context.Context, cfg Config) (*sql.DB, repository.Closer, error) {
	db, dbCloser, err := repository.DBConnection(ctx, cfg.ConfigPostgres)
	if err != nil {
		return nil, nil, err
	}

	return db, dbCloser, nil
}

func performDBMigration(db *sql.DB, up bool) error {
	driver, err := postgresMigration.WithInstance(db, &postgresMigration.Config{})
	if err != nil {
		return err
	}

	dbMigrator, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s/db/migrations", RootDir()),
		"postgres", driver)

	if err != nil {
		return err
	}
	log.Printf("error is: %+v", err)
	if up {
		err = dbMigrator.Up()
		if err != nil && err != migrate.ErrNoChange {
			return err
		}
	} else {
		err = dbMigrator.Down()
		if err != nil && err != migrate.ErrNoChange {
			return err
		}
	}

	return nil
}

// RootDir returns the current projects root directory.
// It can be useful when we need some file pathes based on project root.
func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "../..")
}
