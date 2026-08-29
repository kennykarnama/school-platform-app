package dbconn

import (
	"database/sql"
	"log"
	"net"
	"net/url"
	"time"

	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	_ "github.com/lib/pq"
)

type Config struct {
	// MaxIdleConnection to set max idle connection pooling
	MaxIdleConnection int `envconfig:"DB_MAX_IDLE_CONNECTION" default:"5"`
	// MaxOpenConnection to set max open connection pooling
	MaxOpenConnection int `envconfig:"DB_MAX_OPEN_CONNECTION" default:"10"`
	// MaxLifetimeConnectionn to set max lifetime of pooling | minutes unit
	MaxLifetimeConnection int `envconfig:"DB_MAX_LIFETIME_CONNECTION" default:"10"`
	// Host is host of mysql service
	Host string `envconfig:"DB_HOST" required:"true"`
	// Port is port of mysql service
	Port string `envconfig:"DB_PORT" required:"true" default:"26257"`
	// Username is name of registered user in mysql service
	Username string `envconfig:"DB_USERNAME" required:"true"`
	// DBName is name of registered database in mysql service
	DBName string `envconfig:"DB_NAME" required:"true"`
	// Password is password of used Username in mysql service
	Password string `envconfig:"DB_PASSWORD" default:""`
	Cluster  string `envconfig:"DB_CLUSTER" default:""`
	SSLMode  string `envconfig:"DB_SSL_MODE" default:"disable"`
	// LogMode is toggle to enable/disable log query in your service by default false
	LogMode bool `envconfig:"DB_LOG_MODE" default:"true"`
}

func buildDSN(cfg Config) string {
	user := url.User(cfg.Username)
	if cfg.Password != "" {
		user = url.UserPassword(cfg.Username, cfg.Password)
	}

	dsn := &url.URL{
		Scheme: "postgresql",
		User:   user,
		Host:   net.JoinHostPort(cfg.Host, cfg.Port),
		Path:   cfg.DBName,
	}

	query := url.Values{}
	if cfg.SSLMode != "" {
		query.Set("sslmode", cfg.SSLMode)
	}
	if cfg.Cluster != "" {
		query.Set("options", "--cluster="+cfg.Cluster)
	}
	dsn.RawQuery = query.Encode()

	return dsn.String()
}

// InitGorm is helper function to init gorm database from envar values
// it will panic it cannot find required keys or failed to open database connection
func InitGorm() *gorm.DB {
	var cfg Config
	envconfig.MustProcess("", &cfg)
	dsn := buildDSN(cfg)
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("cannot open postgres connection: %v", err)
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnection)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifetimeConnection) * time.Minute)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnection)

	logMode := logger.Error
	if cfg.LogMode {
		logMode = logger.Info
	}
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: false},
		Logger:         logger.Default.LogMode(logMode),
	})
	if err != nil {
		log.Fatalf("cannot initialize gorm database: %v", err)
	}

	return db
}
