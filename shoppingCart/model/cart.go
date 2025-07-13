package model
import "gorm.io/gorm"
type Cart struct{
	gorm.Model
	UserID uint `json:"user_id"`
	Name string `json:"name"`
	Status string `json:"status"`
}