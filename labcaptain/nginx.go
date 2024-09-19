package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

//go:embed nginx.conf.template
var nginxConfigTemplate string

var successLogger = log.New(os.Stdout, "[SUCCESS] [NGINX] : ", log.LstdFlags)
var errorLogger = log.New(os.Stderr, "[ERROR] [NGINX] : ", log.LstdFlags)

func generateNginxConfig(labID string, serverIP string, containerPort int) {
	nginxCnfig := string(nginxConfigTemplate)
	nginxCnfig = strings.ReplaceAll(nginxCnfig, "{{lab_id}}", labID)
	nginxCnfig = strings.ReplaceAll(nginxCnfig, "{{server_ip}}", serverIP)
	nginxCnfig = strings.ReplaceAll(nginxCnfig, "{{container_port}}", strconv.Itoa(containerPort))
	err := os.WriteFile(fmt.Sprintf("/etc/nginx/sites-enabled/%s", labID), []byte(nginxCnfig), 0777)
	if err != nil {
		errorLogger.Println(fmt.Sprintf("%s > %s", labID, err.Error()))
	} else {
		successLogger.Println(fmt.Sprintf("%s > Nginx config generated successfully", labID))
	}
}

func removeNginxConfig(labID string) {
	err := os.Remove(fmt.Sprintf("/etc/nginx/sites-enabled/%s", labID))
	if err != nil {
		errorLogger.Println(fmt.Sprintf("%s > %s", labID, err.Error()))
	} else {
		successLogger.Println(fmt.Sprintf("%s > Nginx config removed successfully", labID))
	}
}

func reloadNginxConfig() {
	// run /etc/init.d/nginx reload
	err := exec.Command("/etc/init.d/nginx", "reload").Run()
	if err != nil {
		errorLogger.Println(err.Error())
	} else {
		successLogger.Println("Nginx config reloaded successfully")
	}
}
