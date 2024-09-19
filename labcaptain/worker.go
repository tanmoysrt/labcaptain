package main

import (
	"log"
	"os"
	"time"
)

func processPendingLabsDeployment() {
	successLogger := log.New(os.Stdout, "[SUCCESS] [ProcessPendingLabsDeployment] : ", log.LstdFlags)
	errorLogger := log.New(os.Stderr, "[ERROR] [ProcessPendingLabsDeployment] : ", log.LstdFlags)

	for {
		pendingLabs, err := GetPendingLabs()
		if err != nil {
			errorLogger.Println(err.Error())
			continue
		}
		for _, lab := range pendingLabs {
			err := DeployLab(lab.ID)
			if err != nil {
				errorLogger.Println(lab.ID + " > " + err.Error())
			} else {
				successLogger.Println(lab.ID + " > " + "Lab deployed successfully")
			}
		}
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func processExpiredLabsDeletion() {
	successLogger := log.New(os.Stdout, "[SUCCESS] [ProcessExpiredLabsDeletion] : ", log.LstdFlags)
	errorLogger := log.New(os.Stderr, "[ERROR] [ProcessExpiredLabsDeletion] : ", log.LstdFlags)
	for {
		expiredLabs, err := GetExpiredLabs()
		if err != nil {
			errorLogger.Println(err.Error())
			continue
		}
		for _, lab := range expiredLabs {
			err := DestroyLab(lab.ID)
			if err != nil {
				errorLogger.Println(lab.ID + " > " + err.Error())
			} else {
				successLogger.Println(lab.ID + " > " + "Lab destroyed successfully")
			}
		}

		time.Sleep(time.Duration(1) * time.Second)
	}
}

func processAutoExpiratonOfLabs() {
	errorLogger := log.New(os.Stderr, "[ERROR] [ProcessAutoExpiratonOfLabs] : ", log.LstdFlags)
	for {
		// mark deployed and expiry time < now as expired
		err := db.Model(&Lab{}).Where("status = ? AND expiry_time < ?", LabProvisionedStatus, time.Now()).Update("status", LabExpiredStatus).Error
		if err != nil {
			errorLogger.Println(err.Error())
		}
		// delete records with pending and expiry time < now
		err = db.Model(&Lab{}).Where("status = ? AND expiry_time < ?", LabRequestedStatus, time.Now()).Delete(&Lab{}).Error
		if err != nil {
			errorLogger.Println(err.Error())
		}
		time.Sleep(time.Duration(10) * time.Second)
	}
}
