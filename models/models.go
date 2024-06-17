package models

import "gorm.io/gorm"

// Ensures that the models and the tables are in sync with taking the models as a baseline
func Migrate(db *gorm.DB) {
	db.AutoMigrate(&Repairman{}, &Donor{}, &Recipient{}, &EquipmentItem{}, &Tag{}, &ServicedMachine{}, &ServiceComment{})
}

type IModel interface {
	GetID() uint
	SetID(id uint)
}

func (self BaseModel) GetID() uint {
	return self.ID
}

func (self BaseModel) SetID(id uint) {
	self.ID = id
}

type BaseModel struct {
	gorm.Model
}

type Repairman struct {
	BaseModel
	Name     string
	Password string
	Phone    string
	Email    string
	Address  string
	IsActive bool
}

type Donor struct {
	BaseModel
	Name          string
	Phone         string
	Email         string
	Address       string
	IsAnonymous   bool
	IsLegalEntity bool
	AnonymousName string

	Donations []Donation
	Items     []EquipmentItem //`gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Donation struct {
	BaseModel
	DonorID uint
	Donor   Donor
	Date    string `gorm:"autoCreateTime"`

	Equipment []EquipmentItem //`gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type EquipmentItem struct {
	BaseModel
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

	Tags []Tag `gorm:"many2many:equipment_tags;"`
}

type Tag struct {
	BaseModel
	Code        string `gorm:"primaryKey"`
	Name        string
	Description string
}

type Recipient struct {
	BaseModel
	Name          string
	Phone         string
	Email         string
	Address       string
	IsAnonymous   bool
	AnonymousName string

	Items []ServicedMachine //`gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type ServicedMachine struct {
	BaseModel
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
	BaseModel
	Comment string

	PersonID uint
	Person   Repairman

	ServicedMachineID uint
	ServicedMachine   ServicedMachine
}

