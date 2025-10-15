package main

import (
	"context"
	"database/sql"
	"expvar"
	"flag"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Pedro-J-Kukul/police_training/internal/data"
	_ "github.com/lib/pq"
)

const AppVersion = "1.0.0"

type serverConfig struct {
	port int    // server port
	env  string // environment (development, staging, production)
	db   struct {
		dsn          string        // database source name
		maxOpenConns int           // maximum number of open connections
		maxIdleConns int           // maximum number of idle connections
		maxIdleTime  time.Duration // maximum idle time for connections
	}
	cors struct {
		trustedOrigins []string // list of trusted CORS origins
	}
	limiter struct {
		rps     float64 // requests per second
		burst   int     // burst size
		enabled bool    // whether the limiter is enabled
	}
	smtp struct {
		host     string // SMTP host
		port     int    // SMTP port
		username string // SMTP username
		password string // SMTP password
		sender   string // SMTP sender address
	}
}

type appDependencies struct {
	config serverConfig   // application configuration settings
	logger *slog.Logger   // logger for structured logging
	wg     sync.WaitGroup // wait group for managing goroutines
	models data.Models
}

func (app *appDependencies) version() string {
	return AppVersion // return the application version
}

/************************************************************************************************************/
// Main Application Entry Point
/************************************************************************************************************/
func main() {
	// For application setup
	cfg := loadConfig()            // load the application configuration
	logger := setUpLogger(cfg.env) // set up the logger
	db, err := openDB(cfg)         // open the database connection
	if err != nil {
		logger.Error("unable to connect to database", slog.Any("error", err)) // log any error connecting to the database
		os.Exit(1)                                                            // exit if there is a database connection error
	}
	defer db.Close()                                    // ensure the database connection is closed when main() exits
	logger.Info("database connection pool established") // log successful database connection

	// For metrics
	expvar.NewString("version").Set(AppVersion) // publish the application version
	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine() // publish the number of active goroutines
	}))
	expvar.Publish("database", expvar.Func(func() interface{} {
		return db.Stats() // publish database connection pool statistics
	}))
	expvar.Publish("timestamp", expvar.Func(func() interface{} {
		return time.Now().Unix() // publish the current Unix timestamp
	}))

	// Initialize the application dependencies
	app := &appDependencies{
		config: cfg,
		logger: logger,
		wg:     sync.WaitGroup{},
		models: data.NewModels(db),
		// userModel:       data.NewUserModel(db),
		// tokenModel:      data.NewTokenModel(db),
		// permissionModel: data.NewPermissionModel(db),
		// mailer:          mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	err = app.serve() // start the HTTP server
	if err != nil {
		logger.Error("error starting server", slog.Any("error", err)) // log any error starting the server
		os.Exit(1)                                                    // exit if there is a server error
	}
}

/************************************************************************************************************/
// API setup functions
/************************************************************************************************************/
// loadConfig loads the application configuration settings from environment variables
func loadConfig() serverConfig {
	var cfg serverConfig // create a new serverConfig instance

	flag.IntVar(&cfg.port, "port", 4000, "API server port")                                        // server port
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)") // environment

	// Database settings
	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "PostgreSQL DSN")                                                   // database source name
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")                 // max open connections
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")                 // max idle connections
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", time.Minute, "PostgreSQL max connection idle time") // max idle time

	// CORS settings
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(s string) error {
		cfg.cors.trustedOrigins = strings.Fields(s) // split the input string by spaces and assign to trustedOrigins
		return nil
	})

	// Rate limiter settings
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second") // requests per second
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")               // burst size
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")              // whether the limiter is enabled

	// SMTP settings
	flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.mailtrap.io", "SMTP host")                             // SMTP host
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2, "SMTP port")                                                 // SMTP port
	flag.StringVar(&cfg.smtp.username, "smtp-username", "", "SMTP username")                                 // SMTP username
	flag.StringVar(&cfg.smtp.password, "smtp-password", "", "SMTP password")                                 // SMTP password
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Training <noreply@example.com>", "SMTP sender address") // SMTP sender address

	flag.Parse() // parse the command-line flags
	return cfg   // return the populated configuration
}

// setUpLogger initializes and returns a structured logger
func setUpLogger(env string) *slog.Logger {
	var logger *slog.Logger                                     // declare a logger variable
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))      // default to text handler
	logger = logger.With("app_version", AppVersion, "env", env) // add default fields to the logger
	return logger                                               // return the configured logger
}

// openDB opens a database connection pool
func openDB(cfg serverConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn) // open a new database connection
	if err != nil {
		return nil, err // return any error encountered while opening the connection
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)   // set the maximum number of open connections
	db.SetMaxIdleConns(cfg.db.maxIdleConns)   // set the maximum number of idle connections
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime) // set the maximum idle time for connections

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // create a context with a 5-second timeout
	defer cancel()                                                          // ensure the context is cancelled to free resources

	err = db.PingContext(ctx) // ping the database to verify the connection
	if err != nil {
		db.Close()      // close the database connection if ping fails
		return nil, err // return any error encountered while pinging the database
	}

	return db, nil // return the database connection pool
}
