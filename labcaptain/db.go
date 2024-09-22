package main

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func initiateDB() {
	openDB, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		fmt.Println(err.Error())
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
	CpuUsage                  int
	MemoryUsage               int
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
	BaseEndpoint         string
	ExpiryTime           time.Time
	ServerIP             string
	ContainerPort        int
	WebTerminalEnabled   bool
	CodeServerEnabled    bool
	VNCEnabled           bool
	PortProxyEnabled     bool
	EnvironmentVariables string // <variable> = <value>
}
