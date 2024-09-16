package main

import (
	_ "embed"
	"errors"
	"os"
	"os/exec"
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
	err := db.Create(s).Error
	if err != nil {
		return err
	}
	// trigger prometheus config update
	return triggerPrometheusConfigUpdate()
}

func DeleteServer(serverIP string) error {
	var server Server
	err := db.Where("ip = ?", serverIP).First(&server).Error
	if err != nil {
		return err
	}
	// delete all labs associated with this server
	db.Where("server_ip = ?", serverIP).Delete(&Lab{})
	err = db.Delete(&server).Error
	if err != nil {
		return err
	}
	// trigger prometheus config update
	return triggerPrometheusConfigUpdate()
}

func (s *Server) Enable() error {
	if !s.PodmanInstalled {
		return errors.New("Podman is not enabled on this server")
	}
	if !s.PrometheusExportedEnabled {
		return errors.New("Prometheus exporter is not enabled on this server")
	}
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

//go:embed scripts/install_prometheus_exporter.sh
var installPrometheusExporterScript string

func SetupPrometheusExporter(serverIP string) error {
	return runCommandOnServer(serverIP, installPrometheusExporterScript)
}

func triggerPrometheusConfigUpdate() error {
	config, err := generatePrometheusConfig()
	if err != nil {
		return err
	}
	// replace the /etc/prometheus/prometheus.yml file with the new config
	err = os.WriteFile("/etc/prometheus/prometheus.yml", []byte(config), 0644)
	if err != nil {
		return err
	}
	// trigger prometheus service restart
	return exec.Command("systemctl", "restart", "prometheus").Run()
}