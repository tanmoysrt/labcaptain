package main

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func initiateDB() {
	openDB, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to open database")
	}

	// Get the underlying SQL DB object
	sqlDB, err := openDB.DB()
	if err != nil {
		panic("failed to get SQL DB from GORM")
	}

	// Recommended connection pool settings for single backend use
	sqlDB.SetMaxOpenConns(1)    // Only one connection needed
	sqlDB.SetMaxIdleConns(1)    // Keep only one idle connection
	sqlDB.SetConnMaxLifetime(0) // No need to recycle connections

	// Enable WAL mode and set synchronous mode
	openDB.Exec("PRAGMA journal_mode = WAL;")
	openDB.Exec("PRAGMA synchronous = NORMAL;")

	// Migrate the schema
	openDB.AutoMigrate(Server{}, Lab{})

	db = openDB
}

/**
* In database operations, we will use sqlite as the database and move with default configuration
* so, assume there is no primary key, foreign key and default values
* So, make sure to handle primary key, foreign key and default values
* Decide only for getting max performance
 */

type Server struct {
	IP                        string `gorm:"primaryKey"`
	Enabled                   bool
	PodmanInstalled           bool
	PrometheusExportedEnabled bool
}

type LabStatus string

const (
	LabRequestedStatus   LabStatus = "requested"
	LabProvisionedStatus LabStatus = "provisioned"
	LabExpiredStatus     LabStatus = "expired"
)

type Lab struct {
	ID                   string `gorm:"primaryKey"`
	Image                string
	Status               LabStatus `gorm:"default:requested"`
	ExpiryTime           time.Time
	ServerIP             string
	ContainerPort        int
	MaxMemoryMB          int // 0 means no limit
	CPUs                 int // 0 means no limit
	WebTerminalEnabled   bool
	CodeServerEnabled    bool
	VNCEnabled           bool
	PortProxyEnabled     bool
	EnvironmentVariables string // json
}
