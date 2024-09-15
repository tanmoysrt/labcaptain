package main

import (
	"github.com/google/uuid"
)

func (l *Lab) Create() error {
	l.ID = uuid.NewString()
	l.Status = LabRequestedStatus
	return db.Create(l).Error
}

func UpdateLabStatus(id string, status LabStatus) error {
	return db.Model(&Lab{}).Where("id = ?", id).Update("status", status).Error
}

func GetAllLabs() ([]Lab, error) {
	var labs []Lab
	err := db.Find(&labs).Error
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
