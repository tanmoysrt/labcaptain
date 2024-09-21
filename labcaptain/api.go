package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type LabCreateRequest struct {
	Image                string    `json:"image"`
	ExpiryTime           time.Time `json:"expiry_time"`
	WebTerminalEnabled   bool      `json:"web_terminal_enabled"`
	CodeServerEnabled    bool      `json:"code_server_enabled"`
	VNCEnabled           bool      `json:"vnc_enabled"`
	PortProxyEnabled     bool      `json:"port_proxy_enabled"`
	EnvironmentVariables string    `json:"environment_variables"`
}

type LabInfo struct {
	ID         string    `json:"id"`
	Status     LabStatus `json:"status"`
	ExpiryTime time.Time `json:"expiry_time"`
}

func startAPIServer() {
	e := echo.New()
	authtoken := os.Getenv("LABCAPTAIN_API_TOKEN")
	if strings.Compare(authtoken, "") == 0 {
		fmt.Println("No API token provided. API access will be open to all. Pass LABCAPTAIN_API_TOKEN environment variable to enable authentication.")
	} else {
		fmt.Println("API token provided. API access will be restricted to authorized users")
		e.Use(middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
			return strings.Compare(key, authtoken) == 0, nil
		}))
	}
	e.HideBanner = true
	e.GET("/status/:lab_id", func(c echo.Context) error {
		lab_id := c.Param("lab_id")
		lab, err := GetLabByID(lab_id)
		if err != nil {
			return c.JSON(200, LabInfo{
				ID:         lab_id,
				Status:     LabExpiredStatus,
				ExpiryTime: time.Now(),
			})
		}
		return c.JSON(200, LabInfo{
			ID:         lab.ID,
			Status:     lab.Status,
			ExpiryTime: lab.ExpiryTime,
		})
	})
	e.POST("/start", func(c echo.Context) error {
		var labCreateRequest LabCreateRequest
		if err := c.Bind(&labCreateRequest); err != nil {
			return c.String(http.StatusBadRequest, "Invalid request")
		}
		lab := Lab{
			Image:                labCreateRequest.Image,
			ExpiryTime:           labCreateRequest.ExpiryTime,
			WebTerminalEnabled:   labCreateRequest.WebTerminalEnabled,
			CodeServerEnabled:    labCreateRequest.CodeServerEnabled,
			VNCEnabled:           labCreateRequest.VNCEnabled,
			PortProxyEnabled:     labCreateRequest.PortProxyEnabled,
			EnvironmentVariables: labCreateRequest.EnvironmentVariables,
		}
		// check if ExpiryTime is in the future
		if lab.ExpiryTime.Before(time.Now()) {
			return c.String(http.StatusBadRequest, "Expiry time is in the past")
		}
		err := lab.Create()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to create lab")
		}
		return c.JSON(200, LabInfo{
			ID:         lab.ID,
			Status:     lab.Status,
			ExpiryTime: lab.ExpiryTime,
		})
	})
	e.POST("/stop/:lab_id", func(c echo.Context) error {
		labID := c.Param("lab_id")
		if labID == "" {
			return c.String(http.StatusBadRequest, "Lab ID is required")
		}
		err := db.Model(&Lab{}).Where("id = ?", labID).Update("status", LabExpiredStatus).Error
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to destroy lab")
		}
		return c.String(http.StatusOK, "Lab destroyed successfully")
	})
	e.Logger.Fatal(e.Start(":8888"))
}
