package main

import (
	_ "embed"
	"errors"
	"strings"
)

func GetAllServers() ([]Server, error) {
	var servers []Server
	err := db.Find(&servers).Error
	if err != nil {
		return nil, err
	}
	return servers, nil
}

func GetServerByIP(ip string) (*Server, error) {
	var server *Server
	err := db.Where("ip = ?", ip).First(&server).Error
	if err != nil {
		return nil, err
	}
	return server, nil
}

func (s *Server) Create() error {
	if strings.Compare(s.IP, "") == 0 {
		return errors.New("IP is required")
	}
	// check if server is already present
	var servercount int64
	count := db.Table("servers").Where("ip = ?", s.IP).Count(&servercount)
	if count.Error != nil {
		return errors.New("Failed to check if server is already present")
	}
	if servercount > 0 {
		return errors.New("Server already exists with IP " + s.IP)
	}
	// create server
	s.Enabled = false
	s.PrometheusExportedEnabled = false
	s.PodmanInstalled = false
	return db.Create(s).Error
}

func DeleteServer(serverIP string) error {
	var server Server
	err := db.Where("ip = ?", serverIP).First(&server).Error
	if err != nil {
		return err
	}
	// delete all labs associated with this server
	db.Where("server_ip = ?", serverIP).Delete(&Lab{})
	return db.Delete(&server).Error
}

func (s *Server) Enable() error {
	s.Enabled = true
	return db.Save(s).Error
}

func (s *Server) Disable() error {
	s.Enabled = false
	s.PodmanInstalled = false
	s.PrometheusExportedEnabled = false
	return db.Save(s).Error
}

func (s *Server) SetPodmanInstalled(installed bool) error {
	s.PodmanInstalled = installed
	return db.Save(s).Error
}

func (s *Server) SetPrometheusExportedEnabled(enabled bool) error {
	s.PrometheusExportedEnabled = enabled
	return db.Save(s).Error
}

//go:embed scripts/install_podmain.sh
var installPodmanScript string

func SetupPodman(serverIP string) error {
	return runCommandOnServer(serverIP, installPodmanScript)
}
