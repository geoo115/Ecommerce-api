package services

import (
	"errors"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"gorm.io/gorm"
)

// CartService interface defines cart business logic
type CartService interface {
	AddToCart(userID uint, productID uint, quantity int) error
	GetUserCart(userID uint) ([]models.Cart, error)
	UpdateCartItem(userID uint, productID uint, quantity int) error
	RemoveFromCart(userID uint, productID uint) error
	CheckStock(productID uint, quantity int) (bool, error)
	CalculateCartTotal(cartItems []models.Cart) float64
}

// cartService implements CartService interface
type cartService struct {
	db *gorm.DB
}

// NewCartService creates a new cart service instance
func NewCartService() CartService {
	return &cartService{
		db: db.DB,
	}
}

// AddToCart adds a product to user's cart
func (s *cartService) AddToCart(userID uint, productID uint, quantity int) error {
	// Check if product exists and has sufficient stock
	if available, err := s.CheckStock(productID, quantity); err != nil {
		return err
	} else if !available {
		return errors.New("insufficient stock")
	}

	// Check if item already exists in cart
	var existingCart models.Cart
	err := s.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&existingCart).Error

	if err == nil {
		// Item exists, update quantity
		existingCart.Quantity += quantity
		return s.db.Save(&existingCart).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Item doesn't exist, create new cart entry
		cart := models.Cart{
			UserID:    userID,
			ProductID: productID,
			Quantity:  quantity,
		}
		return s.db.Create(&cart).Error
	}

	return err
}

// GetUserCart retrieves all cart items for a user
func (s *cartService) GetUserCart(userID uint) ([]models.Cart, error) {
	var cartItems []models.Cart
	err := s.db.Where("user_id = ?", userID).
		Preload("Product").
		Find(&cartItems).Error
	return cartItems, err
}

// UpdateCartItem updates the quantity of a cart item
func (s *cartService) UpdateCartItem(userID uint, productID uint, quantity int) error {
	if quantity <= 0 {
		return s.RemoveFromCart(userID, productID)
	}

	// Check stock availability
	if available, err := s.CheckStock(productID, quantity); err != nil {
		return err
	} else if !available {
		return errors.New("insufficient stock")
	}

	return s.db.Model(&models.Cart{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Update("quantity", quantity).Error
}

// RemoveFromCart removes a product from user's cart
func (s *cartService) RemoveFromCart(userID uint, productID uint) error {
	return s.db.Where("user_id = ? AND product_id = ?", userID, productID).
		Delete(&models.Cart{}).Error
}

// CheckStock verifies if product has sufficient stock
func (s *cartService) CheckStock(productID uint, quantity int) (bool, error) {
	var inventory models.Inventory
	err := s.db.Where("product_id = ?", productID).First(&inventory).Error
	if err != nil {
		return false, err
	}
	return inventory.Stock >= quantity, nil
}

// CalculateCartTotal calculates the total price of cart items
func (s *cartService) CalculateCartTotal(cartItems []models.Cart) float64 {
	var total float64
	for _, item := range cartItems {
		if item.Product.ID != 0 {
			total += item.Product.Price * float64(item.Quantity)
		}
	}
	return total
}
