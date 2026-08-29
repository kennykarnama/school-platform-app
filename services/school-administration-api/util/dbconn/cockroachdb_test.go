package dbconn

import (
	"net/url"
	"testing"
)

func TestBuildDSNForLocalCockroach(t *testing.T) {
	dsn, err := url.Parse(buildDSN(Config{
		Host:     "localhost",
		Port:     "26257",
		Username: "root",
		DBName:   "school_administration",
		SSLMode:  "disable",
	}))
	if err != nil {
		t.Fatalf("parse DSN: %v", err)
	}

	if got, want := dsn.Host, "localhost:26257"; got != want {
		t.Fatalf("host = %q, want %q", got, want)
	}
	if got, want := dsn.Path, "/school_administration"; got != want {
		t.Fatalf("path = %q, want %q", got, want)
	}
	if got, want := dsn.Query().Get("sslmode"), "disable"; got != want {
		t.Fatalf("sslmode = %q, want %q", got, want)
	}
	if got := dsn.Query().Get("options"); got != "" {
		t.Fatalf("local options = %q, want empty", got)
	}
}

func TestBuildDSNForHostedCockroach(t *testing.T) {
	dsn, err := url.Parse(buildDSN(Config{
		Host:     "example.cockroachlabs.cloud",
		Port:     "26257",
		Username: "school",
		Password: "p@ss word",
		DBName:   "defaultdb",
		Cluster:  "routing-id",
		SSLMode:  "verify-full",
	}))
	if err != nil {
		t.Fatalf("parse DSN: %v", err)
	}

	if got, want := dsn.Query().Get("sslmode"), "verify-full"; got != want {
		t.Fatalf("sslmode = %q, want %q", got, want)
	}
	if got, want := dsn.Query().Get("options"), "--cluster=routing-id"; got != want {
		t.Fatalf("options = %q, want %q", got, want)
	}
	if password, ok := dsn.User.Password(); !ok || password != "p@ss word" {
		t.Fatalf("password was not preserved")
	}
}
