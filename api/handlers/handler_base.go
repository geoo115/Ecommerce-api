package handlers

import (
	"net/http"
	"strconv"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HandlerBase provides common handler functionality
type HandlerBase struct{}

// GetUserID extracts user ID from context with validation
func (h *HandlerBase) GetUserID(c *gin.Context) (uint, error) {
	userID, exists := c.Get("userID")
	if !exists || userID == nil {
		utils.SendUnauthorized(c, "Unauthorized")
		return 0, utils.ErrUnauthorized
	}
	uid, ok := userID.(uint)
	if !ok {
		utils.SendUnauthorized(c, "Unauthorized")
		return 0, utils.ErrUnauthorized
	}
	return uid, nil
}

// ValidateIDParam validates and converts ID parameter from URL
func (h *HandlerBase) ValidateIDParam(c *gin.Context, paramName string) (uint, error) {
	idStr := c.Param(paramName)
	if idStr == "" {
		utils.SendValidationError(c, paramName+" is required")
		c.Abort()
		return 0, utils.ErrValidation
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		if paramName == "id" {
			utils.SendValidationError(c, "Invalid id")
		} else {
			utils.SendValidationError(c, "Invalid "+paramName)
		}
		c.Abort()
		return 0, utils.ErrValidation
	}

	return uint(id), nil
}

// BindJSON binds JSON input with validation
func (h *HandlerBase) BindJSON(c *gin.Context, input interface{}) error {
	if err := c.ShouldBindJSON(input); err != nil {
		utils.SendValidationError(c, "Invalid request payload")
		return err
	}
	return nil
}

// CheckOwnership verifies that a resource belongs to the authenticated user
func (h *HandlerBase) CheckOwnership(c *gin.Context, model interface{}, resourceID uint) error {
	uid, err := h.GetUserID(c)
	if err != nil {
		return err
	}

	err = db.DB.Where("id = ? AND user_id = ?", resourceID, uid).First(model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.SendNotFound(c, "Resource not found or access denied")
		} else {
			utils.SendInternalError(c, "Database error")
		}
		return err
	}
	return nil
}

// HandleDBError provides consistent database error handling
func (h *HandlerBase) HandleDBError(c *gin.Context, err error, notFoundMsg, internalMsg string) {
	if err == nil {
		return // No error, don't send any response
	}
	if err == gorm.ErrRecordNotFound {
		utils.SendNotFound(c, notFoundMsg)
	} else {
		utils.SendInternalError(c, internalMsg)
		utils.Error("Database error: %v", err)
	}
}

// PaginationParams holds pagination parameters
type PaginationParams struct {
	Page  int
	Limit int
}

// GetPaginationParams extracts pagination parameters from query string
func (h *HandlerBase) GetPaginationParams(c *gin.Context) PaginationParams {
	page := 1
	limit := 10

	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = p
	}

	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	return PaginationParams{Page: page, Limit: limit}
}

// ApplyPagination applies pagination to a GORM query
func (h *HandlerBase) ApplyPagination(query *gorm.DB, params PaginationParams) *gorm.DB {
	offset := (params.Page - 1) * params.Limit
	return query.Offset(offset).Limit(params.Limit)
}

// TransactionWrapper wraps operations in a database transaction
func (h *HandlerBase) TransactionWrapper(c *gin.Context, fn func(*gorm.DB) error) error {
	tx := db.DB.Begin()
	if tx.Error != nil {
		utils.SendInternalError(c, "Failed to start transaction")
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			utils.SendInternalError(c, "Transaction failed")
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		utils.SendInternalError(c, "Failed to commit transaction")
		return err
	}

	return nil
}

// SendCreatedResponse sends a 201 Created response
func (h *HandlerBase) SendCreatedResponse(c *gin.Context, message string, data interface{}) {
	utils.SendSuccess(c, http.StatusCreated, message, data)
}

// SendUpdatedResponse sends a 200 OK response for updates
func (h *HandlerBase) SendUpdatedResponse(c *gin.Context, message string, data interface{}) {
	utils.SendSuccess(c, http.StatusOK, message, data)
}

// SendDeletedResponse sends a 200 OK response for deletions
func (h *HandlerBase) SendDeletedResponse(c *gin.Context, message string) {
	utils.SendSuccess(c, http.StatusOK, message, nil)
}

// SendListResponse sends a 200 OK response for list operations
func (h *HandlerBase) SendListResponse(c *gin.Context, message string, data interface{}) {
	utils.SendSuccess(c, http.StatusOK, message, data)
}

// ResponseHelper provides standardized response patterns
type ResponseHelper struct{}

// SendCreatedResponse sends a 201 Created response
func (r *ResponseHelper) SendCreatedResponse(c *gin.Context, message string, data interface{}) {
	utils.SendSuccess(c, http.StatusCreated, message, data)
}

// SendUpdatedResponse sends a 200 OK response for updates
func (r *ResponseHelper) SendUpdatedResponse(c *gin.Context, message string, data interface{}) {
	utils.SendSuccess(c, http.StatusOK, message, data)
}

// SendDeletedResponse sends a 200 OK response for deletions
func (r *ResponseHelper) SendDeletedResponse(c *gin.Context, message string) {
	utils.SendSuccess(c, http.StatusOK, message, nil)
}

// SendListResponse sends a 200 OK response for list operations
func (r *ResponseHelper) SendListResponse(c *gin.Context, message string, data interface{}) {
	utils.SendSuccess(c, http.StatusOK, message, data)
}

// Global instances for convenience
var (
	Base     = &HandlerBase{}
	Response = &ResponseHelper{}
)
