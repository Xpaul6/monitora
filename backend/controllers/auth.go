package controllers

import (
	"net/http"

	. "github.com/XPaul6/monitora/models"
	authutils "github.com/XPaul6/monitora/utils/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody AuthRequest
		if err := c.Bind(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var num int64
		db.Model(&User{}).Where("email = ?", reqBody.Email).Count(&num)
		if num > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
			return
		}

		passwordHash, err := authutils.HashPassword(reqBody.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}

		newUser := User{
			Email:        reqBody.Email,
			PasswordHash: passwordHash,
		}

		result := db.Create(&newUser)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register new user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	}
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody AuthRequest
		if err := c.Bind(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user User
		result := db.Where("email = ?", reqBody.Email).First(&user)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
			return
		}

		isValid := authutils.VerificatePassword(user.PasswordHash, reqBody.Password)
		if !isValid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
			return
		}

		token, err := authutils.GenerateToken(reqBody.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generations error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
