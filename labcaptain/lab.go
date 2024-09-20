package main

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	_ "embed"

	"github.com/jaevor/go-nanoid"
	"gorm.io/gorm"
)

var nanoidGenerator func() string

func init() {
	gen, err := nanoid.CustomASCII("abcdefghijklmnopqrstuvwxyz0123456789", 15)
	if err != nil {
		panic(err)
	}
	nanoidGenerator = gen
}

func (l *Lab) Create() error {
	l.ID = nanoidGenerator()
	l.Status = LabRequestedStatus
	return db.Create(l).Error
}

func GetLabByID(id string) (*Lab, error) {
	var lab Lab
	err := db.Where("id = ?", id).First(&lab).Error
	if err != nil {
		return nil, err
	}
	return &lab, err
}

func GetAllLabs() ([]Lab, error) {
	var labs []Lab
	err := db.Find(&labs).Error
	if err != nil {
		return labs, err
	}
	return labs, nil
}

func GetPendingLabs() ([]Lab, error) {
	var labs []Lab
	err := db.Where("status = ?", LabRequestedStatus).Find(&labs).Error
	if err != nil {
		return labs, err
	}
	return labs, nil
}

func GetExpiredLabs() ([]Lab, error) {
	var labs []Lab
	err := db.Where("status = ?", LabExpiredStatus).Find(&labs).Error
	if err != nil {
		return labs, err
	}
	return labs, nil
}

//go:embed scripts/deploy_lab.sh
var deployLabScript string

func DeployLab(labID string) error {
	var lab Lab
	err := db.Where("id = ?", labID).First(&lab).Error
	if err != nil {
		return err
	}
	if lab.Status != LabRequestedStatus {
		return errors.New("Lab is not in requested status")
	}
	server, err := GetRandomDeployableServer()
	if err != nil {
		return err
	}
	environmentVariables := fmt.Sprintf("ENABLE_WEB_TERMINAL=%d ENABLE_CODE_SERVER=%d ENABLE_VNC=%d ENABLE_PORT_PROXY=%d %s", boolToInt(lab.WebTerminalEnabled), boolToInt(lab.CodeServerEnabled), boolToInt(lab.VNCEnabled), boolToInt(lab.PortProxyEnabled), lab.EnvironmentVariables)
	// deploy the lab on server
	stdoutBuffer := new(bytes.Buffer)
	stderrBuffer := new(bytes.Buffer)
	// replace variables in the script
	deployLabScriptCopy := string(deployLabScript)
	deployLabScriptCopy = strings.ReplaceAll(deployLabScriptCopy, "{{lab_id}}", lab.ID)
	deployLabScriptCopy = strings.ReplaceAll(deployLabScriptCopy, "{{lab_image}}", lab.Image)
	deployLabScriptCopy = strings.ReplaceAll(deployLabScriptCopy, "{{lab_environment_variables}}", environmentVariables)
	// run the script
	err = runCommandOnServerWithBuffer(server+":22", deployLabScriptCopy, stdoutBuffer, stderrBuffer)
	if err != nil {
		return err
	}
	if stderrBuffer.Len() > 0 {
		return errors.New(stderrBuffer.String())
	}
	// parse the output using regex match to get the port > format > assigned_${port}_port
	regex_match := regexp.MustCompile(`assigned_(\d+)_port`)
	matches := regex_match.FindStringSubmatch(stdoutBuffer.String())
	if len(matches) != 2 {
		return errors.New("Failed to parse output from deploy_lab.sh script")
	}
	port, err := strconv.Atoi(matches[1])
	if err != nil {
		return err
	}
	lab.ServerIP = server
	lab.ContainerPort = port
	// run the deploy script
	lab.Status = LabProvisionedStatus
	err = db.Save(&lab).Error
	if err != nil {
		return err
	}
	generateNginxConfig(lab.ID, server, port)
	return nil
}

//go:embed scripts/destroy_lab.sh
var destroyLabScript string

func DestroyLab(labID string) error {
	var lab Lab
	err := db.Where("id = ?", labID).First(&lab).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	if lab.Status != LabExpiredStatus {
		return errors.New("Lab is not in expired status")
	}
	// destroy from server
	err = runCommandOnServer(lab.ServerIP+":22", strings.ReplaceAll(destroyLabScript, "{{lab_id}}", lab.ID))
	if err != nil {
		return err
	}
	// delete lab from db
	db.Delete(&lab)
	// remove nginx config
	removeNginxConfig(lab.ID)
	return nil
}
