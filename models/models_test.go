package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUserModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&User{})
	assert.NoError(t, err)

	user := User{
		Username: "testuser",
		Password: "password",
		Email:    "test@example.com",
		Phone:    "+1234567890",
		Role:     "customer",
	}

	err = db.Create(&user).Error
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestProductModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&Product{}, &Category{}, &Inventory{})
	assert.NoError(t, err)

	category := Category{Name: "Electronics"}
	err = db.Create(&category).Error
	assert.NoError(t, err)

	product := Product{
		Name:        "Laptop",
		Price:       999.99,
		CategoryID:  category.ID,
		Description: "A great laptop",
	}

	err = db.Create(&product).Error
	assert.NoError(t, err)
	assert.NotZero(t, product.ID)

	inventory := Inventory{
		ProductID: product.ID,
		Stock:     10,
	}
	err = db.Create(&inventory).Error
	assert.NoError(t, err)
}

// Add similar tests for other models: Category, Cart, Order, etc.
func TestCategoryModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&Category{})
	assert.NoError(t, err)

	category := Category{Name: "Books"}
	err = db.Create(&category).Error
	assert.NoError(t, err)
	assert.Equal(t, "Books", category.Name)
}

func TestCartModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&Cart{}, &User{}, &Product{})
	assert.NoError(t, err)

	user := User{Username: "user"}
	db.Create(&user)
	product := Product{Name: "Item"}
	db.Create(&product)

	cart := Cart{
		UserID:    user.ID,
		ProductID: product.ID,
		Quantity:  2,
	}
	err = db.Create(&cart).Error
	assert.NoError(t, err)
}

// Continue for Order, Payment, Review, Wishlist, Address
func TestOrderModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&Order{}, &User{})
	assert.NoError(t, err)

	user := User{Username: "user"}
	db.Create(&user)

	order := Order{
		UserID:      user.ID,
		TotalAmount: 100.0,
		Status:      "pending",
	}
	err = db.Create(&order).Error
	assert.NoError(t, err)
}

func TestPaymentModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&Payment{}, &Order{})
	assert.NoError(t, err)

	order := Order{TotalAmount: 100.0}
	db.Create(&order)

	payment := Payment{
		OrderID: order.ID,
		Status:  "completed",
		Amount:  100.0,
	}
	err = db.Create(&payment).Error
	assert.NoError(t, err)
}

func TestReviewModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&Review{}, &User{}, &Product{})
	assert.NoError(t, err)

	user := User{Username: "user"}
	db.Create(&user)
	product := Product{Name: "Item"}
	db.Create(&product)

	review := Review{
		UserID:    user.ID,
		ProductID: product.ID,
		Rating:    5,
		Comment:   "Great!",
	}
	err = db.Create(&review).Error
	assert.NoError(t, err)
}

func TestWishlistModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&Wishlist{}, &User{}, &Product{})
	assert.NoError(t, err)

	user := User{Username: "user"}
	db.Create(&user)
	product := Product{Name: "Item"}
	db.Create(&product)

	wishlist := Wishlist{
		UserID:    user.ID,
		ProductID: product.ID,
	}
	err = db.Create(&wishlist).Error
	assert.NoError(t, err)
}

func TestAddressModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&Address{}, &User{})
	assert.NoError(t, err)

	user := User{Username: "user"}
	db.Create(&user)

	address := Address{
		UserID:  user.ID,
		Address: "123 Main St",
		City:    "City",
		ZipCode: "12345",
	}
	err = db.Create(&address).Error
	assert.NoError(t, err)
}
