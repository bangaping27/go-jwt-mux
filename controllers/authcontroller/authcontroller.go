package authcontroller

import (
	"encoding/json"
	"go-jwt-mux/config"
	"go-jwt-mux/helper"
	"go-jwt-mux/models"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(w http.ResponseWriter, r *http.Request) {
	//mengambil inputan json
	var userInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error(), "status": "400"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	defer r.Body.Close()

	//cek apakah usernamenya ada
	var user models.User
	if err := models.DB.Where("username = ?", userInput.Username).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			response := map[string]string{"message": "username atau password salah", "status": "404"}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		default:
			response := map[string]string{"message": err.Error(), "status": "500"}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}
	}

	//cek password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password))
	if err != nil {
		response := map[string]string{"message": "username atau password salah", "status": "404"}
		helper.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

	//generate token jwt
	expTime := time.Now().Add(time.Hour * 1)

	claims := &config.JWTClaim{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-jwt-mux",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	//mendeklarasi algoritma yang akan digunakan untuk signing
	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//signing token
	token, err := tokenAlgo.SignedString([]byte(config.JWT_KEY))
	if err != nil {
		response := map[string]string{"message": err.Error(), "status": "500"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	//set token ke cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    token,
		HttpOnly: true,
	})

	//mengembalikan response
	response := map[string]string{"message": "Login Berhasil", "status": "200"}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func Register(w http.ResponseWriter, r *http.Request) {

	//mengambil inputan json
	var userInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error(), "status": "400"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	defer r.Body.Close()

	//hash pas menggunakan bcrypt
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	userInput.Password = string(hashPassword)

	// log.Fatal(userInput)

	//insert ke db
	if err := models.DB.Create(&userInput).Error; err != nil {
		response := map[string]string{"message": err.Error(), "status": "500"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}
	response := map[string]string{"message": "success", "status": "200"}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func Logout(w http.ResponseWriter, r *http.Request) {

	//menghapus cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	})

	//mengembalikan response
	response := map[string]string{"message": "Logout Berhasil", "status": "200"}
	helper.ResponseJSON(w, http.StatusOK, response)
}
