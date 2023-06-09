package handler

import (
	"net/http"
	"os"
	"strconv"
	"path/filepath"
	"time"


	"github.com/KhoirulAziz99/final_project_e_wallet/internal/app"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase app.UserUsecase
}

func NewUserHandler(userUsecase app.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}
func (h *UserHandler) InsertUser(c *gin.Context) {
	var user domain.User
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mengambil data dari form-data
	idStr := c.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	user.ID = id
	user.Name = c.PostForm("name")
	user.Email = c.PostForm("email")
	user.Password = c.PostForm("password")

	// Mengambil file dari form-data
	file, err := c.FormFile("profile_picture")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload profile picture"})
		return
	}

	// Simpan file ke lokasi yang diinginkan
	dstPath := filepath.Join("cmd", file.Filename)
	err = c.SaveUploadedFile(file, dstPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save profile picture"})
		return
	}
	user.ProfilePicture = dstPath

	if err := h.userUsecase.InsertUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Mengambil data dari form-data
	var user domain.User
	user.ID = userID
	user.Name = c.PostForm("name")
	user.Email = c.PostForm("email")
	user.Password = c.PostForm("password")

	// Mengambil file dari form-data
	file, err := c.FormFile("profile_picture")
	if err == nil {
		// Jika ada file gambar baru diunggah, simpan ke lokasi yang diinginkan
		dstPath := filepath.Join("cmd", file.Filename)
		err = c.SaveUploadedFile(file, dstPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save profile picture"})
			return
		}
		user.ProfilePicture = dstPath
	}

	if err := h.userUsecase.UpdateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (h *UserHandler) FindOneUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userUsecase.FindOne(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *UserHandler) FindAllUsers(c *gin.Context) {
	users, err := h.userUsecase.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.userUsecase.Delete(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandler) Login(c *gin.Context) {

	var user domain.LoginUser

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userInDb, err := h.userUsecase.FindByUsername(user.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
		return
	}

	if userInDb.Password != user.Password {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})

}

func ProfileHandler(c *gin.Context) {

	claims := c.MustGet("claims").(jwt.MapClaims)
	username := claims["username"].(string)
	c.JSON(http.StatusOK, gin.H{"username": "Heloo, " + username})
}
