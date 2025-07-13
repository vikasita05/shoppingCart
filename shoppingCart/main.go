package main
import(
	"shoppingCart/database"
	"shoppingCart/model"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"net/http"
)
func main(){
	database.InitDB()
	router := gin.Default()
	router.Use(cors.Default())
	router.POST("/users", func(c *gin.Context){
		var user model.User
		if e := c.ShouldBindJSON(&user); e!=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid"})
			return
		}
		database.DB.Create(&user)
		c.JSON(http.StatusOK, gin.H{
			"msg": "User created",
			"user": user,
		})
	})

	router.POST("/users/login", func(c *gin.Context){
		var log struct{
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if e := c.ShouldBindJSON(&log); e!=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid"})
			return
		}
		var user model.User
		res := database.DB.Where("username = ? AND password = ?", log.Username, log.Password).First(&user)
		if res.Error != nil{
			c.JSON(http.StatusUnauthorized, gin.H{"error":"Invalid username and password"})
			return
		}
		token := log.Username+"token_1"
		user.Token=token
		database.DB.Save(&user)
		c.JSON(http.StatusOK,gin.H{
			"Msg":"Login Successful",
			"token": token,
		})
	})

	router.POST("/item",func(c *gin.Context){
		var item model.Item
		if e := c.ShouldBindJSON(&item); e != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid input"})
			return
		}
		database.DB.Create(&item)
		c.JSON(http.StatusOK, gin.H{"msg":"Item created","item":item})
	})
	//listing all items

	router.GET("/items",func(c *gin.Context){
	var items []model.Item
	database.DB.Find(&items)
	c.JSON(http.StatusOK,items)
	})

	router.POST("/carts",func(c *gin.Context){
		token := c.GetHeader("Authorization")
		if token == ""{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"Missing token"})
			return
		}
		var user model.User
		if e := database.DB.Where("token = ?",token).First(&user).Error; e != nil{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"Invalid token"})
			return
		}
		var itemData struct{
			ItemID uint `json:"item_id"`
			Name string `json:"name"`
		}
		if e := c.ShouldBindJSON(&itemData); e!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid item"})
			return
		}

		//find or create a cart for this user
		var cart model.Cart
		if e := database.DB.Where("user_id = ?",user.ID).First(&cart).Error; e != nil{
			cart=model.Cart{
				UserID: user.ID,
				Name: itemData.Name,
				Status: "active",
			}
			database.DB.Create(&cart)
		}

		//create a cart item linking item to cart
		CartItem := model.CartItem{
			CartID: cart.ID,
			ItemID: itemData.ItemID,
		}
		database.DB.Create(&CartItem)
		c.JSON(http.StatusOK,gin.H{
			"Msg":"Item added to cart",
			"cart": cart.ID,
		})
	})

	router.GET("/carts",func(c *gin.Context){
		token := c.GetHeader("Authorization")
		if token == ""{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"Missing token"})
			return 
		}
		var user model.User
		if e := database.DB.Where("token = ?",token).First(&user).Error; e != nil{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"Invalid"})
			return
		}
		var cart model.Cart
		if e := database.DB.Where("user_id = ?",user.ID).First(&cart).Error; e != nil{
			c.JSON(http.StatusNotFound,gin.H{"error":"Cart not found"})
			return
		}
		var cartItems []model.CartItem
		database.DB.Where("cart_id = ?",cart.ID).Find(&cartItems)
		c.JSON(http.StatusOK,gin.H{
			"cart":cart,
			"item":cartItems,
		})

	})

	router.POST("/orders",func(c *gin.Context){
		token := c.GetHeader("Authorization")
		if token == ""{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"Missing token"})
			return
		}
		var user model.User
		if e := database.DB.Where("token = ?",token).First(&user).Error; e != nil{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"Invalid token"})
			return
		}
		var cart model.Cart
		if e := database.DB.Where("user_id = ?",user.ID).First(&cart).Error; e != nil{
			c.JSON(http.StatusNotFound,gin.H{"error":"Cart not found"})
			return
		}
		order := model.Order{
			UserID: user.ID,
			CartID: cart.ID,
			Status: "Placed",
		}
		database.DB.Create(&order)
		cart.Status = "inactive"
		database.DB.Save(&cart)
		c.JSON(http.StatusOK,gin.H{"msg":"Order placed successfully","order_id":order.ID,})
	})

	router.GET("/orders",func(c *gin.Context){
		token := c.GetHeader("Authorization")
		if token == ""{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"Missing token"})
			return
		}
		var user model.User
		if e := database.DB.Where("token = ?",token).First(&user).Error; e != nil{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"Invalid token"})
			return
		}
		var orders []model.Order
		database.DB.Where("user_id = ?",user.ID).Find(&orders)
		c.JSON(http.StatusOK,gin.H{
			"orders":orders,
		})
		
	})
	router.Run(":8080")
}