package models

import "app/db"

type Repairman struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Password string
	Phone    string
	Email    string
	Address  string
	IsActive bool
}

type Donor struct {
	ID            uint `gorm:"primaryKey"`
	Name          string
	Phone         string
	Email         string
	Address       string
	IsAnonymous   bool
	IsLegalEntity bool
	AnonymousName string
	Items         []EquipmentItem //`gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
    Donations     []Donation
}

type Donation struct {
    ID        uint `gorm:"primaryKey"`
    DonorID   uint
    Donor     Donor
    Equipment []EquipmentItem //`gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
    Date      string `gorm:"autoCreateTime"`
}

type EquipmentItem struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Description string
	Category    string
	Quantity    int
	Tags        []Tag
    Donor      Donor
    Donation   Donation
}

type Tag struct {
    Code        string `gorm:"primaryKey"`
	Name        string
	Description string
}
type Recipient struct {
	ID            uint `gorm:"primaryKey"`
	Name          string
	Phone         string
	Email         string
	Address       string
	IsAnonymous   bool
	AnonymousName string
	Items         []ServicedMachine //`gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}


type ServicedMachine struct {
	ID           uint `gorm:"primaryKey"`
	SerialNumber string           `gorm:"unique"`
	Equipment    []EquipmentItem
    Status       ServiceStatus
	Comments     []ServiceComment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
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
	ID                uint `gorm:"primaryKey"`
	Comment           string
    Date              string `gorm:"autoCreateTime"`
	Person            Repairman
	ServicedMachineID uint
}


func Migrate() {
	db.DB.AutoMigrate(&Repairman{}, &Donor{}, &Recipient{}, &EquipmentItem{}, &Tag{}, &ServicedMachine{}, &ServiceComment{})
} 
