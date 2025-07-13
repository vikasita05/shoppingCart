package model
import "gorm.io/gorm"
type CartItem struct{
	gorm.Model
	CartID uint `json:"cartID"`
	ItemID uint `json:"itemId"`
}