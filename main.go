package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ----------------- MODELS -----------------

type Book struct {
	ID     uint   `gorm:"primaryKey"`
	Title  string `json:"title"`
	Author string `json:"author"`
	ISBN   string `json:"isbn"`
	Copies []BookCopy
}

type BookCopy struct {
	ID     uint   `gorm:"primaryKey"`
	BookID uint
	Status string `json:"status"` // available, checked_out
}

type User struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `json:"name"`
	Role string `json:"role"` // student, librarian
}

type Checkout struct {
	ID         uint `gorm:"primaryKey"`
	UserID     uint
	BookCopyID uint
	DueDate    time.Time
	Returned   bool
}

type Reservation struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	BookID    uint
	CreatedAt time.Time
	Active    bool
}

// ----------------- MAIN -----------------

func main() {

	var err error
	DB, err = gorm.Open(sqlite.Open("library.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	DB.AutoMigrate(&Book{}, &BookCopy{}, &User{}, &Checkout{}, &Reservation{})

	r := gin.Default()

	// ----------------- ROOT -----------------

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Library Management API"})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Library API running"})
	})

	// ----------------- ADD BOOK -----------------

	r.POST("/books", func(c *gin.Context) {

		var input struct {
			Title       string `json:"title"`
			Author      string `json:"author"`
			ISBN        string `json:"isbn"`
			NumOfCopies int    `json:"num_of_copies"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		book := Book{
			Title:  input.Title,
			Author: input.Author,
			ISBN:   input.ISBN,
		}

		DB.Create(&book)

		for i := 0; i < input.NumOfCopies; i++ {
			copy := BookCopy{
				BookID: book.ID,
				Status: "available",
			}
			DB.Create(&copy)
		}

		c.JSON(200, gin.H{
			"message": "Book added successfully",
			"book_id": book.ID,
		})
	})

	// ----------------- GET BOOKS -----------------

	r.GET("/books", func(c *gin.Context) {
		var books []Book
		DB.Preload("Copies").Find(&books)
		c.JSON(200, books)
	})

	// ----------------- REGISTER USER -----------------

	r.POST("/users", func(c *gin.Context) {

		var input struct {
			Name string `json:"name"`
			Role string `json:"role"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if input.Role != "student" && input.Role != "librarian" {
			c.JSON(400, gin.H{"error": "Role must be student or librarian"})
			return
		}

		user := User{
			Name: input.Name,
			Role: input.Role,
		}

		DB.Create(&user)

		c.JSON(200, gin.H{
			"message": "User created successfully",
			"user_id": user.ID,
		})
	})

	// ----------------- GET USERS -----------------

	r.GET("/users", func(c *gin.Context) {
		var users []User
		DB.Find(&users)
		c.JSON(200, users)
	})

	// ----------------- CHECKOUT -----------------

	r.POST("/checkout", func(c *gin.Context) {

		var input struct {
			UserID uint `json:"user_id"`
			BookID uint `json:"book_id"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		var user User
		if err := DB.First(&user, input.UserID).Error; err != nil {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}

		if user.Role != "student" {
			c.JSON(403, gin.H{"error": "Only students can checkout"})
			return
		}

		var copy BookCopy
		if err := DB.Where("book_id = ? AND status = ?", input.BookID, "available").
			First(&copy).Error; err != nil {

			c.JSON(400, gin.H{"error": "No copies available. You may reserve."})
			return
		}

		copy.Status = "checked_out"
		DB.Save(&copy)

		checkout := Checkout{
			UserID:     user.ID,
			BookCopyID: copy.ID,
			DueDate:    time.Now().AddDate(0, 0, 7),
			Returned:   false,
		}

		DB.Create(&checkout)

		c.JSON(200, gin.H{
			"message":     "Book checked out successfully",
			"checkout_id": checkout.ID,
			"due_date":    checkout.DueDate,
		})
	})

	// ----------------- GET CHECKOUTS -----------------

	r.GET("/checkouts", func(c *gin.Context) {
		var checkouts []Checkout
		DB.Find(&checkouts)
		c.JSON(200, checkouts)
	})

	// ----------------- RESERVE -----------------

	r.POST("/reserve", func(c *gin.Context) {

		var input struct {
			UserID uint `json:"user_id"`
			BookID uint `json:"book_id"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		reservation := Reservation{
			UserID:    input.UserID,
			BookID:    input.BookID,
			CreatedAt: time.Now(),
			Active:    true,
		}

		DB.Create(&reservation)

		c.JSON(200, gin.H{
			"message":        "Book reserved successfully",
			"reservation_id": reservation.ID,
		})
	})

	// ----------------- RETURN -----------------

	r.POST("/return", func(c *gin.Context) {

		var input struct {
			CheckoutID uint `json:"checkout_id"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		var checkout Checkout
		if err := DB.First(&checkout, input.CheckoutID).Error; err != nil {
			c.JSON(404, gin.H{"error": "Checkout not found"})
			return
		}

		if checkout.Returned {
			c.JSON(400, gin.H{"error": "Already returned"})
			return
		}

		checkout.Returned = true
		DB.Save(&checkout)

		var copy BookCopy
		DB.First(&copy, checkout.BookCopyID)
		copy.Status = "available"
		DB.Save(&copy)

		// Fine calculation
		var fine float64 = 0
		if time.Now().After(checkout.DueDate) {
			daysLate := int(time.Since(checkout.DueDate).Hours() / 24)
			fine = float64(daysLate * 10)
		}

		// FIFO Reservation Handling
		var reservation Reservation
		if err := DB.Where("book_id = ? AND active = ?", copy.BookID, true).
			Order("created_at ASC").
			First(&reservation).Error; err == nil {

			copy.Status = "checked_out"
			DB.Save(&copy)

			newCheckout := Checkout{
				UserID:     reservation.UserID,
				BookCopyID: copy.ID,
				DueDate:    time.Now().AddDate(0, 0, 7),
				Returned:   false,
			}
			DB.Create(&newCheckout)

			reservation.Active = false
			DB.Save(&reservation)
		}

		c.JSON(200, gin.H{
			"message": "Book returned successfully",
			"fine":    fine,
		})
	})

	r.Run(":8080")
}