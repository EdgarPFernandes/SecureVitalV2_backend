package main

import (
	"log"
	"net/http"
	"time"

	repo "github.com/EdgarPFernandes/SecureVitalV2_backend/internal/adapters/postgresql/sqlc"
	"github.com/EdgarPFernandes/SecureVitalV2_backend/internal/alerts"
	"github.com/EdgarPFernandes/SecureVitalV2_backend/internal/devices"
	"github.com/EdgarPFernandes/SecureVitalV2_backend/internal/patients"
	"github.com/EdgarPFernandes/SecureVitalV2_backend/internal/type_alerts"
	"github.com/EdgarPFernandes/SecureVitalV2_backend/internal/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

type application struct {
	config config
	//logger
	db *pgx.Conn
}

type config struct {
	addr string //port
	db   dbConfig
}

type dbConfig struct {
	dsn string // user= password= dbname= etc
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID) //important for rate limiting
	r.Use(middleware.RealIP)    //important for rate limiting and analytics
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) //recover from crashes

	r.Use(middleware.Timeout(60 * time.Second)) //important for graceful shutdown and to prevent hanging requests

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	r.Route("/api", func(r chi.Router) {
		alertsService := alerts.NewService(repo.New(app.db), repo.New(app.db), app.db)
		alertsHandler := alerts.NewHandler(alertsService)
		r.Get("/alerts", alertsHandler.ListAlerts)

		r.Post("/alerts", alertsHandler.CreateAlert)

		typeAlertsService := type_alerts.NewService(repo.New(app.db), repo.New(app.db), app.db)
		typeAlertsHandler := type_alerts.NewHandler(typeAlertsService)
		r.Get("/type_alerts", typeAlertsHandler.ListAlertTypes)

		r.Post("/type_alerts", typeAlertsHandler.CreateAlertTypes)

		patientsService := patients.NewService(repo.New(app.db), repo.New(app.db), app.db)
		patientsHandler := patients.NewHandler(patientsService)
		r.Get("/patients", patientsHandler.ListPatients)

		r.Post("/patients", patientsHandler.CreatePatient)

		r.Get("/patients/{patient_id}/alerts", patientsHandler.ListAlertsByPatient)

		devicesService := devices.NewService(repo.New(app.db), repo.New(app.db), app.db)
		devicesHandler := devices.NewHandler(devicesService)
		r.Get("/devices", devicesHandler.ListDevices)

		r.Post("/devices", devicesHandler.CreateDevice)

		usersService := users.NewService(repo.New(app.db), repo.New(app.db), app.db)
		usersHandler := users.NewHandler(usersService)
		r.Get("/users", usersHandler.ListUsers)

		//r.Post("/users", usersHandler.CreateUser)
	})

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Minute,
	}

	log.Printf("listening on %s", srv.Addr)

	return srv.ListenAndServe()
}
