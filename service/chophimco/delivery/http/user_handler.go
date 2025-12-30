package http

import (
	"github.com/gin-gonic/gin"
	"github.com/leehai1107/chophimco-server/pkg/apiwrapper"
	"github.com/leehai1107/chophimco-server/pkg/logger"
	"github.com/leehai1107/chophimco-server/pkg/middleware/auth"
	"github.com/leehai1107/chophimco-server/service/chophimco/model/request"
)

type IUserHandler interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
	GetProfile(ctx *gin.Context)
}

// Login godoc
// @Summary User login
// @Description Authenticate a user and return a token
// @Tags user
// @Accept json
// @Produce json
// @Param request body request.Login true "Login credentials"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Failure 401 {object} apiwrapper.APIResponse
// @Router /api/v1/user/login [post]
func (h *Handler) Login(ctx *gin.Context) {
	var req request.Login
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	response, err := h.userUsecase.Login(ctx, req)
	if err != nil {
		logger.EnhanceWith(ctx).Errorw("Login failed", "error", err, "email", req.Email)
		apiwrapper.SendUnauthorized(ctx, "Login failed")
		return
	}

	apiwrapper.SendSuccess(ctx, response)
}

// Register godoc
// @Summary Register new user
// @Description Register a new user account
// @Tags user
// @Accept json
// @Produce json
// @Param request body request.Register true "Registration information"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Failure 500 {object} apiwrapper.APIResponse
// @Router /api/v1/user/register [post]
func (h *Handler) Register(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	var req request.Register
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorw("Invalid request format", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	err := h.userUsecase.Register(ctx, req)
	if err != nil {
		log.Errorw("Registration failed", "error", err, "email", req.Email)
		if err.Error() == "email already registered" {
			apiwrapper.SendBadRequest(ctx, "Email already exists")
			return
		}
		apiwrapper.SendInternalError(ctx, "Registration failed")
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Registration successful"})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get user profile information
// @Tags user
// @Produce json
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 401 {object} apiwrapper.APIResponse
// @Router /api/v1/user/profile [get]
func (h *Handler) GetProfile(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		apiwrapper.SendUnauthorized(ctx, "Unauthorized")
		return
	}

	profile, err := h.userUsecase.GetUserProfile(ctx, userID)
	if err != nil {
		apiwrapper.SendInternalError(ctx, "Failed to get profile")
		return
	}

	apiwrapper.SendSuccess(ctx, profile)
}
