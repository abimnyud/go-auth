package routes

import (
	"fmt"
	"strings"

	"abimanyu.dev/go-auth/model"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewMiddleware() fiber.Handler {
	return AuthMiddleware
}

func AuthMiddleware(ctx *fiber.Ctx) error {
	sess, err := store.Get(ctx)

	if strings.Split(ctx.Path(), "/")[1] == "auth" {
		return ctx.Next()
	}

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	if sess.Get(AUTH_KEY) == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	return ctx.Next()
}

func Register(ctx *fiber.Ctx) error {
	var data User

	err := ctx.BodyParser(&data)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	password, err := bcrypt.GenerateFromPassword([]byte(data.Password), 14)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	user := model.User{
		Name:     data.Name,
		Email:    data.Email,
		Password: string(password),
	}

	err = model.CreateUser(&user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Create new user success",
	})
}

func Login(ctx *fiber.Ctx) error {
	var data User

	err := ctx.BodyParser(&data)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	var user model.User
	if !model.CheckUser(data.Email, &user) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized 1",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized 2",
		})
	}

	sess, err := store.Get(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	sess.Set(AUTH_KEY, true)
	sess.Set(USER_ID, user.ID)

	err = sess.Save()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logged in",
	})
}

func Logout(ctx *fiber.Ctx) error {
	sess, err := store.Get(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Logged out",
		})
	}

	err = sess.Destroy()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logged out",
	})
}

func HealtCheck(ctx *fiber.Ctx) error {
	sess, err := store.Get(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	auth := sess.Get(AUTH_KEY)
	if auth != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Authorized",
		})
	} else {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
}

func GetUser(ctx *fiber.Ctx) error {
	sess, err := store.Get(ctx);
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	if sess.Get(AUTH_KEY) == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	userId := sess.Get(USER_ID);

	if userId == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var user model.User;
	user, err = model.GetUser(fmt.Sprint(userId));
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(user);
}