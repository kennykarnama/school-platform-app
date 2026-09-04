package main

import (
	"github.com/kennykarnama/school-adminstration-api/domain/handler/interceptor"
	"github.com/kennykarnama/school-adminstration-api/domain/repository/user"
	user2 "github.com/kennykarnama/school-adminstration-api/domain/service/user"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kennykarnama/school-adminstration-api/config"
	"github.com/kennykarnama/school-adminstration-api/domain/handler"
	"github.com/kennykarnama/school-adminstration-api/domain/repository/academicyear"
	"github.com/kennykarnama/school-adminstration-api/domain/repository/attendancetype"
	"github.com/kennykarnama/school-adminstration-api/domain/repository/class"
	"github.com/kennykarnama/school-adminstration-api/domain/repository/core"
	setupRepo "github.com/kennykarnama/school-adminstration-api/domain/repository/setup"
	academicYearSvc "github.com/kennykarnama/school-adminstration-api/domain/service/academicyear"
	attendanceTypeSvc "github.com/kennykarnama/school-adminstration-api/domain/service/attendancetype"
	classSvc "github.com/kennykarnama/school-adminstration-api/domain/service/class"
	coreSvc "github.com/kennykarnama/school-adminstration-api/domain/service/core"
	setupSvc "github.com/kennykarnama/school-adminstration-api/domain/service/setup"
	"github.com/kennykarnama/school-adminstration-api/util/dbconn"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

type server struct {
	http.Server
}

func main() {
	cfg := config.Get()

	dbconn := dbconn.InitGorm()

	v := validator.New()

	coreRepo := core.NewSQLRepository(dbconn)

	academicYearRepo := academicyear.NewSQLRepository(dbconn)
	academicSvc := academicYearSvc.NewService(academicYearRepo)

	classRepo := class.NewEnumRepository()
	classSvc := classSvc.NewService(classRepo)

	attendanceTypeRepo := attendancetype.NewSQLRepository(dbconn)
	attendanceTypeSvc := attendanceTypeSvc.NewService(attendanceTypeRepo)

	coreSvc := coreSvc.NewService(coreRepo, attendanceTypeSvc)

	coreHandler := handler.NewHandler(coreSvc, academicSvc, classSvc, v, attendanceTypeSvc)

	userRepo := user.NewMySqlRepository(dbconn)
	userSvc := user2.NewService(userRepo, cfg)
	userHandler := handler.NewUserHandler(userSvc, v, cfg.SessionCookieSecure)
	setupRepository := setupRepo.NewSQLRepository(dbconn)
	setupService := setupSvc.NewService(setupRepository, classSvc)
	setupHandler := handler.NewSetupHandler(setupService)

	hostName, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	logrus.SetFormatter(&logrus.JSONFormatter{})
	standardFields := logrus.Fields{
		"appname":  cfg.ServiceName,
		"hostname": hostName,
	}

	httpServer := &server{
		Server: http.Server{
			Addr: ":" + cfg.RestPort,
		},
	}

	authMiddleware := interceptor.NewAuth(userSvc, []string{"/api/v1/teacher/login", "/healthz"})

	r := mux.NewRouter()
	r.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		handler.ResponseJson(w, map[string]string{"status": "ok"}, http.StatusOK)
	}).Methods(http.MethodGet)
	r.Handle("/api/v1/student/attendance/register", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(coreHandler.RegisterStudent))).Methods("POST")
	r.Handle("/api/v1/student/attendance/submit", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(coreHandler.SubmitAttendance))).Methods("POST")
	r.Handle("/api/v1/student/attendance/list", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(coreHandler.Students))).Methods("GET")
	r.Handle("/api/v1/academic-years", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(coreHandler.ListAcademicYear))).Methods("GET")
	r.Handle("/api/v1/classes", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(coreHandler.ListClasses))).Methods("GET")
	r.Handle("/api/v1/attendance/types", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(coreHandler.ListAttendanceTypes))).Methods("GET")
	r.Handle("/api/v1/student/class/{id}/deactivate", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(coreHandler.DeactivateStudentClass))).Methods("PATCH")
	r.Handle("/api/v1/student/class/transfer", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(coreHandler.TransferStudentClass))).Methods("POST")
	r.Handle("/api/v1/student/attendance/stats", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(coreHandler.Stats))).Methods("GET")

	r.Handle("/api/v1/teacher/login", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(userHandler.Login))).Methods("POST")
	r.Handle("/api/v1/teacher/session/validate", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(userHandler.Validate))).Methods("GET")
	r.Handle("/api/v1/setup/preview", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(setupHandler.Preview))).Methods("POST")
	r.Handle("/api/v1/setup/apply", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(setupHandler.Apply))).Methods("POST")
	r.Handle("/api/v1/setup/students/template", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(setupHandler.StudentTemplate))).Methods("GET")
	r.Handle("/api/v1/setup/students/import", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(setupHandler.ImportStudents))).Methods("POST")

	if cfg.EnableAuth {
		logrus.Infof("auth enabled")
		r.Use(authMiddleware.ValidateToken)
	}

	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
			"https://intense-sea-59889.herokuapp.com",
		},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug:          true,
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodHead, http.MethodPatch},
	})

	handler := c.Handler(r)

	httpServer.Handler = handler

	logrus.WithFields(standardFields).Infof("HTTP served on port: %v", cfg.RestPort)

	if err := httpServer.ListenAndServe(); err != nil {
		logrus.WithFields(standardFields).Fatalf("unable to serve. err: %v", err)
	}
}
