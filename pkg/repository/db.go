package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"contrib.go.opencensus.io/integrations/ocsql"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type ConfigPostgres struct {
	Alias                string
	Address              string
	Database             string
	Username             string
	Password             string
	MaxConns             int
	MaxWaitForConnection time.Duration
	// Schema defaults to service env when value is ""
	Schema string `yaml:"schema" env:"PG_SCHEMA" env-default:""`
}

type Closer func()

type postgresConn struct {
	db     *sql.DB
	closer Closer
}

func DBConnection(ctx context.Context, dbCfg ConfigPostgres) (*sql.DB, Closer, error) {
	db, dbCloser, err := dbConnection(ctx, dbCfg)
	if err != nil {
		return nil, nil, err
	}
	if dbCfg.Schema != "" && dbCfg.Schema != "public" {
		_, err = db.ExecContext(ctx, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s AUTHORIZATION %s", dbCfg.Schema, dbCfg.Username))
		if err != nil {
			dbCloser()
			return nil, nil, err
		}
	}

	return db, dbCloser, err
}

func dbConnection(ctx context.Context, cfg ConfigPostgres) (*sql.DB, Closer, error) {
	dataSourceName := fmt.Sprintf(
		"postgres://%v:%v@%v/%v?sslmode=disable&search_path=%v",
		cfg.Username,
		cfg.Password,
		cfg.Address,
		cfg.Database,
		cfg.Schema,
	)
	// fallback to unwrapped driver
	driverName := "postgres"

	// driverName, err := otelsql.Register("postgres", semconv.DBSystemPostgreSQL)
	// if err != nil {
	// 	// fallback to unwrapped driver
	// 	driverName = "postgres"
	// }

	// Register SQL specific views, only needs to be done once per process
	ocsql.RegisterAllViews()

	// OpenCensusSQL wraps the actual driver, and returns a generated driver name. We do this once per connection.
	driverName, err := ocsql.Register(driverName, ocsql.WithAllTraceOptions(), ocsql.WithInstanceName(cfg.Alias))
	if err != nil {
		// fallback to unwrapped driver
		driverName = "postgres"
		fmt.Printf("Could not wrap sql driver: %+v\n", zap.Error(err))
	}

	// sql.Open for postgres doesn't do much, so we need to call Ping() to test the connection
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, nil, err
	}

	db.SetMaxOpenConns(cfg.MaxConns)

	start := time.Now()
	maxWait := cfg.MaxWaitForConnection
	fmt.Printf("pinging the db: %+v, %+v\n", time.Since(start), cfg.MaxWaitForConnection)
	for time.Since(start) < maxWait {
		// sometimes proxy needs a bit longer than server to start, so we just give it a few seconds
		err = db.Ping()
		if err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Printf("finished pinging db\n")
	if err != nil {
		return nil, nil, err
	}

	conn := postgresConn{
		db:     db,
		closer: ocsql.RecordStats(db, 5*time.Second),
	}
	return conn.db, conn.closer, nil
}
