package httpdelivery

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-tutorial/internal/usecase"
)

type AuthHandler struct {
	authUC usecase.AuthUsecase
}

func NewAuthHandler(r *gin.Engine, authUC usecase.AuthUsecase) {
	h := &AuthHandler{authUC: authUC}

	api := r.Group("/api/v1")
	auth := api.Group("/auth")

	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)

	protected := auth.Group("/")
	protected.Use(AuthMiddleware(authUC))
	protected.GET("/me", h.Me)
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Register godoc
// @Tags Auth
// @Accept json
// @Param payload body registerRequest true "Register request"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authUC.Register(req.Email, req.Password)
	if err != nil {
		if err == usecase.ErrEmailExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"email": user.Email,
	})
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login godoc
// @Summary Login and obtain JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body loginRequest true "Login request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authUC.Login(req.Email, req.Password)
	if err != nil {
		if err == usecase.ErrInvalidCredential {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Me godoc
// @Summary Get current user info
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	uid, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user in context"})
		return
	}

	user, err := h.authUC.GetUserByID(uid.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"email": user.Email,
	})
}
