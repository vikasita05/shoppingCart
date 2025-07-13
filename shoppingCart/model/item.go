package model
import "gorm.io/gorm"
type Item struct{
	gorm.Model
	Name string `json:"name`
	Status string `json:"status`
}