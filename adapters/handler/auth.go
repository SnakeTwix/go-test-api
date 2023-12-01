package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"server/adapters/tools/jwt"
	"server/core/domain"
	"server/core/ports"
)

type AuthHandler struct {
	serviceUser ports.ServiceUser
}

var authHandler *AuthHandler

func GetAuthHandler(serviceUser ports.ServiceUser) *AuthHandler {
	if authHandler != nil {
		return authHandler
	}

	authHandler = &AuthHandler{
		serviceUser: serviceUser,
	}

	return authHandler
}

func (h *AuthHandler) RegisterRoutes(e *echo.Group) {
	e.POST("/auth/register", h.Register)
	e.POST("/auth/login", h.Login)
}

// Register godoc
// @Summary      Registers a user
// @Description  Registers a user with default permissions
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param user body domain.CreateUser true "The user data to use when registering"
// @Success      200  {array}   domain.Token
// @Failure      400  {object}  error "Should only happen when there is already a user with the same username"
// @Failure      500  {object}  error "Shouldn't happen, but will if something fails"
// @Router       /api/v1/auth/register [POST]
func (h *AuthHandler) Register(ctx echo.Context) error {
	var user domain.CreateUser

	if err := ctx.Bind(&user); err != nil {
		return err
	}
	if err := ctx.Validate(&user); err != nil {
		return err
	}

	createdUser, err := h.serviceUser.SaveUser(ctx.Request().Context(), &user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tokener := &jwt.TokenGenerator{
		User: createdUser,
	}

	token, err := tokener.Token()
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, token)
}

func (h *AuthHandler) Login(ctx echo.Context) error {
	var user domain.LoginUser

	if err := ctx.Bind(&user); err != nil {
		return err
	}
	if err := ctx.Validate(&user); err != nil {
		return err
	}

	token, err := h.serviceUser.LoginUser(ctx.Request().Context(), &user)

	fmt.Println(err)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, token)
}
