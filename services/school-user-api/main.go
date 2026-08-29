package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kennykarnama/school-user-api/config"
	"github.com/kennykarnama/school-user-api/domain/api/user"
	"github.com/kennykarnama/school-user-api/domain/api/userauth"
	"github.com/kennykarnama/school-user-api/domain/repository/user/sql"
	"github.com/kennykarnama/school-user-api/domain/repository/userauth/redis"
	userService "github.com/kennykarnama/school-user-api/domain/service/user"
	userAuthService "github.com/kennykarnama/school-user-api/domain/service/userauth"
	"github.com/kennykarnama/school-user-api/util"
	"github.com/kennykarnama/school-user-api/util/dbconn"
	"github.com/sirupsen/logrus"
)

type server struct {
	http.Server
}

func main() {

	ctx := context.Background()
	cfg := config.Get()

	hostName, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	logrus.SetFormatter(&logrus.JSONFormatter{})
	standardFields := logrus.Fields{
		"appname":  cfg.ServiceName,
		"hostname": hostName,
	}

	db := dbconn.InitGorm()
	redisPool := dbconn.InitRedis()

	redisWrapper := util.NewRedisWrapper(redisPool)

	userSqlRepository := sql.NewMysqlRepository(db)
	userAuthRedisRepository := redis.NewRepository(redisWrapper)

	userService := userService.NewService(userSqlRepository)
	userAuthService := userAuthService.NewService(cfg, userService, userAuthRedisRepository)

	v := validator.New()
	userHandler := user.NewHandler(ctx, userService, v)
	userAuthHandler := userauth.NewHandler(ctx, v, userAuthService)

	httpServer := &server{
		Server: http.Server{
			Addr: ":" + cfg.RestPort,
		},
	}
	r := mux.NewRouter()
	r.Handle("/api/v1/user", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(userHandler.RegisterUser))).Methods("POST")
	r.Handle("/api/v1/user/auth", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(userAuthHandler.Login))).Methods("POST")
	r.Handle("/api/v1/user/auth/token/validate", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(userAuthHandler.ValidateToken))).Methods("POST")
	r.Handle("/api/v1/user/auth/logout", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(userAuthHandler.Logout))).Methods("POST")
	r.Handle("/api/v1/user/auth/token/refresh", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(userAuthHandler.RefreshToken))).Methods("POST")

	httpServer.Handler = r

	logrus.WithFields(standardFields).Infof("HTTP served on port: %v", cfg.RestPort)

	if err := httpServer.ListenAndServe(); err != nil {
		logrus.WithFields(standardFields).Fatalf("unable to serve. err: %v", err)
	}
}
