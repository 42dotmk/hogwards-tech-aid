package models

import "gorm.io/gorm"
import "app/db"

type Repairman struct {
	gorm.Model
	Name     string
	Password string
	Phone    string
	Email    string
	Address  string
	IsActive bool
}

type Donor struct {
	gorm.Model
	Name          string
	Phone         string
	Email         string
	Address       string
	IsAnonymous   bool
	IsLegalEntity bool
	AnonymousName string

	Items     []EquipmentItem //`gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Donations []Donation
}

type Donation struct {
	gorm.Model
	DonorID uint
	Donor   Donor
	Date    string `gorm:"autoCreateTime"`

	Equipment []EquipmentItem //`gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type EquipmentItem struct {
	gorm.Model
	Name        string
	Description string
	Category    string
	Quantity    int

	DonorID uint
	Donor   Donor

	DonationID uint
	Donation   Donation

	ServicedMachineID uint
	ServicedMachine   ServicedMachine

	Tags []Tag
}

type Tag struct {
	Code            string `gorm:"primaryKey"`
	Name            string
	Description     string
	EquipmentItemID uint
}
type Recipient struct {
	gorm.Model
	Name          string
	Phone         string
	Email         string
	Address       string
	IsAnonymous   bool
	AnonymousName string

	Items []ServicedMachine //`gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type ServicedMachine struct {
	gorm.Model
	SerialNumber string `gorm:"unique"`
	Status       ServiceStatus

	RecipientID uint
	Recipient   Recipient

	Equipment []EquipmentItem
	Comments  []ServiceComment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type ServiceStatus int

const (
	Pending ServiceStatus = iota
	InProgress
	Paused
	Completed
	Cancelled
)

type ServiceComment struct {
	gorm.Model
	Comment string

	PersonID uint
	Person   Repairman

	ServicedMachineID uint
	ServicedMachine   ServicedMachine
}

func Migrate() {
	db.DB.AutoMigrate(&Repairman{}, &Donor{}, &Recipient{}, &EquipmentItem{}, &Tag{}, &ServicedMachine{}, &ServiceComment{})
}
