package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func generatePrometheusConfig() (string, error) {
	configContent := `
global:
  scrape_interval: 15s

scrape_configs:
- job_name: prometheus
  static_configs:
  - targets:
    - localhost:9090
- job_name: node
  static_configs:
  - targets:
#start_targets_list
#end_targets_list
`
	// Ensure that the tags exist in the config
	startTag := "#start_targets_list"
	endTag := "#end_targets_list"

	// if /etc/prometheus/prometheus.yml exist, read the file
	if _, err := os.Stat("/etc/prometheus/prometheus.yml"); err == nil {
		existingConfigContentBytes, err := os.ReadFile("/etc/prometheus/prometheus.yml")
		if err != nil {
			return "", errors.New("Failed to read /etc/prometheus/prometheus.yml > " + err.Error())
		}
		existingConfigContent := string(existingConfigContentBytes)
		if strings.Contains(existingConfigContent, "#start_targets_list") && strings.Contains(existingConfigContent, "#end_targets_list") {
			configContent = existingConfigContent
		}
	}

	servers, err := GetAllServers()
	if err != nil {
		return "", err
	}
	serverString := ""
	for _, server := range servers {
		serverString += fmt.Sprintf("    - %s:9100\n", server.IP)
	}

	re := regexp.MustCompile(fmt.Sprintf("(?s)%s(.*?)%s", regexp.QuoteMeta(startTag), regexp.QuoteMeta(endTag)))
	newConfig := re.ReplaceAllString(configContent, fmt.Sprintf("%s\n%s%s", startTag, serverString, endTag))
	return newConfig, nil
}
