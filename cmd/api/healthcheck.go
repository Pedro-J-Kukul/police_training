package main

import (
	"context"
	"database/sql"
	"net/http"
	"time"
)

// healthCheckHandler returns application and database status metadata.
func (app *appDependencies) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := "available"

	var dbStatus string
	if err := pingDatabase(app.models.User.DB); err != nil {
		dbStatus = "unreachable"
		status = "degraded"
	} else {
		dbStatus = "ready"
	}

	payload := envelope{
		"status": status,
		"app": envelope{
			"version":     app.version(),
			"environment": app.config.env,
			"time":        time.Now().UTC(),
		},
		"database": dbStatus,
	}

	if err := app.writeJSON(w, http.StatusOK, payload, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func pingDatabase(db *sql.DB) error {
	if db == nil {
		return sql.ErrConnDone
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return db.PingContext(ctx)
}
