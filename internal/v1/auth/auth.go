package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"github.com/opplieam/bb-core-api/.gen/buy-better-core/public/model"
	"github.com/opplieam/bb-core-api/internal/store"
	"github.com/opplieam/bb-core-api/internal/utils"
)

var refreshTokenDuration = 730 * time.Hour // 1 month
var tokenDuration = 15 * time.Minute

type Store interface {
	InsertOrUpdateUser(email, firstName, lastName, role string) error
	FindUserByEmail(email string) (model.Users, error)
	IsValidUser(userID int32) error
}

type Handler struct {
	Store Store
}

func NewHandler(store Store) *Handler {
	return &Handler{
		Store: store,
	}
}

type ReqURI struct {
	Provider string `uri:"provider" binding:"required,oneof=google"`
}

func (h *Handler) ProviderHandler(c *gin.Context) {
	var reqURI ReqURI
	if err := c.ShouldBindUri(&reqURI); err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	q := c.Request.URL.Query()
	q.Add("provider", reqURI.Provider)
	c.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (h *Handler) CallbackHandler(c *gin.Context) {
	var reqURI ReqURI
	if err := c.ShouldBindUri(&reqURI); err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	q := c.Request.URL.Query()
	q.Add("provider", reqURI.Provider)
	c.Request.URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		_ = c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	err = h.Store.InsertOrUpdateUser(user.Email, user.FirstName, user.LastName, "basic")
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	token, err := utils.GenerateToken(1*time.Minute, user.UserID, "guest")
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	redirectURL := fmt.Sprintf("http://localhost:5174/callback/%s/%s", user.Email, token)
	c.Redirect(http.StatusFound, redirectURL)
}

type GetTokenReq struct {
	Email string `json:"email" binding:"required"`
	Token string `json:"token" binding:"required"`
}

func (h *Handler) GetTokenHandler(c *gin.Context) {
	var getTokenReq GetTokenReq
	if err := c.ShouldBindJSON(&getTokenReq); err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	_, err := utils.VerifyToken(getTokenReq.Token)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	user, err := h.Store.FindUserByEmail(getTokenReq.Email)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			c.AbortWithStatus(http.StatusForbidden)
			return
		default:
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	accessToken, err := utils.GenerateToken(tokenDuration, string(user.ID), user.Role)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	refreshToken, err := utils.GenerateToken(refreshTokenDuration, string(user.ID), user.Role)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	//TODO: change domain name and secure depend on environment
	c.SetCookie(
		"refresh_token",
		refreshToken,
		2629800,
		"/",
		"localhost",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "email": user.Email})

}

func (h *Handler) LogoutHandler(c *gin.Context) {
	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"msg": "logged out"})
}

func (h *Handler) RefreshTokenHandler(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		c.JSON(-1, gin.H{"msg": "no token"})
		return
	}
	token, err := utils.VerifyToken(refreshToken)
	if err != nil {
		_ = c.AbortWithError(http.StatusUnauthorized, err)
		c.JSON(-1, gin.H{"msg": "invalid token"})
		return
	}
	userIdString, err := token.GetString("user_id")
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		c.JSON(-1, gin.H{"msg": "no user id"})
		return
	}
	userId, err := strconv.ParseInt(userIdString, 10, 32)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	roleString, err := token.GetString("role")

	err = h.Store.IsValidUser(int32(userId))
	if err != nil {
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			c.AbortWithStatus(http.StatusForbidden)
			return
		default:
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	accessToken, err := utils.GenerateToken(tokenDuration, userIdString, roleString)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	newRefreshToken, err := utils.GenerateToken(refreshTokenDuration, userIdString, roleString)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	//TODO: change domain name and secure depend on environment
	c.SetCookie(
		"refresh_token",
		newRefreshToken,
		2629800,
		"/",
		"localhost",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})

}
