package delivery

import (
	"fmt"
	"net/http"

	"github.com/daugminas/kyc/app/domain"

	"github.com/labstack/echo/v4"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
)

func (s *server) createUser(c echo.Context) error {

	// validate request body
	var user domain.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, domain.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	// validate required fields
	validate = validator.New()
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.JSON(http.StatusBadRequest, domain.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &echo.Map{"data": validationErr.Error()},
		})
	}

	userId, err := s.a.CreateUser(&user, true)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	return c.JSON(http.StatusCreated, domain.UserResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data:    &echo.Map{"data": userId},
	})
}

func (s *server) getUser(c echo.Context) error {
	userId := c.Param("userId") // 'userId' is the param in the URL route
	user, err := s.a.GetUser(userId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	return c.JSON(http.StatusOK, domain.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &echo.Map{"data": user},
	})
}

func (s *server) editUser(c echo.Context) error {

	// validate request body
	var user domain.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, domain.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	// validate required fields
	validate = validator.New()
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.JSON(http.StatusBadRequest, domain.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &echo.Map{"data": validationErr.Error()},
		})
	}

	userId := c.Param("userId")
	updatedUser, err := s.a.UpdateUser(userId, &user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	return c.JSON(http.StatusOK, domain.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &echo.Map{"data": updatedUser},
	})
}

func (s *server) deleteUser(c echo.Context) error {
	userId := c.Param("userId")

	err := s.a.DeleteUser(userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.UserResponse{
			Status:  http.StatusNotFound,
			Message: "error",
			Data:    &echo.Map{"data": fmt.Sprintf("user with id '%s' not found!", userId)},
		})
	}

	return c.JSON(http.StatusOK, domain.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &echo.Map{"data": fmt.Sprintf("user with id '%s' successfully deleted!", userId)},
	})
}
