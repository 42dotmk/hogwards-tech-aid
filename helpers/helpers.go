package helpers

import (
	"app/db"
	"app/models"
)

type Stats struct {
	ServicedMachines int64
	Donations        int64
	Donors           int64
	EquipmentItems   int64
}

func CountStats() Stats {
	var stats Stats
	db.DB.Model(&models.ServicedMachine{}).Count(&stats.ServicedMachines)
	db.DB.Model(&models.Donation{}).Count(&stats.Donations)
	db.DB.Model(&models.Donor{}).Count(&stats.Donors)
	db.DB.Model(&models.EquipmentItem{}).Count(&stats.EquipmentItems)
	return stats
}
