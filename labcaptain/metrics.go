package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ServerMetrics struct {
	CPUUsage    int
	MemoryUsage int
}

var query = `
(100 - (avg by (instance) (rate(node_cpu_seconds_total{mode="idle", instance="{{server_ip}}:9100"}[2m])) * 100))
or
(100 * (node_memory_MemTotal_bytes{instance="{{server_ip}}:9100"} - node_memory_MemAvailable_bytes{instance="{{server_ip}}:9100"}) / node_memory_MemTotal_bytes{instance="{{server_ip}}:9100"})
`
var prometheusServer = "http://localhost:9090"

func GetServerMetrics(serverIP string) (*ServerMetrics, error) {
	encodedQuery := url.QueryEscape(strings.ReplaceAll(query, "{{server_ip}}", serverIP))
	req, err := http.NewRequest("GET", prometheusServer+"/api/v1/query?query="+encodedQuery, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	var data struct {
		Data struct {
			Result []struct {
				Value []interface{} `json:"value"`
			} `json:"result"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	if len(data.Data.Result) != 2 {
		return nil, errors.New("Failed to get server metrics")
	}
	if len(data.Data.Result[0].Value) != 2 {
		return nil, errors.New("Failed to get server metrics")
	}
	cpuUsage, err := strconv.ParseFloat(data.Data.Result[0].Value[1].(string), 64)
	if err != nil {
		return nil, err
	}
	memoryUsage, err := strconv.ParseFloat(data.Data.Result[1].Value[1].(string), 64)
	if err != nil {
		return nil, err
	}
	if cpuUsage < 0 || memoryUsage < 0 {
		return nil, errors.New("Failed to get server metrics")
	}

	return &ServerMetrics{
		CPUUsage:    int(cpuUsage),
		MemoryUsage: int(memoryUsage),
	}, nil
}
