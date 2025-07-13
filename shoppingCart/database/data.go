package database
import(
	"log"
	"gorm.io/gorm"
	"github.com/glebarez/sqlite"
	"shoppingCart/model"
)
var DB *gorm.DB
func InitDB(){
	var e error
	DB, e=gorm.Open(sqlite.Open("cart.db"), &gorm.Config{})
	if e!=nil{
		log.Fatal("Failed to connect to database: ",e)
	}
	DB.AutoMigrate(&model.User{})
	DB.AutoMigrate(&model.Item{})
	DB.AutoMigrate(&model.Cart{})
	DB.AutoMigrate(&model.CartItem{})
	DB.AutoMigrate(&model.Order{})
}