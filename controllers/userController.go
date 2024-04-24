package controllers

import (
	"log"
	"strconv"

	"math"

	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
	"github.com/harrisonhys/learn-fiber/config"
	"github.com/harrisonhys/learn-fiber/models/entities"
	"github.com/harrisonhys/learn-fiber/models/req"
)

func UserControllerShow(c *fiber.Ctx) error {
	// Default values for pagination and filtering
	page := c.Query("page", "1")
	perPage := c.Query("perpage", "10")
	active := c.Query("active", "") // Filter by "active" column

	// Convert page and perPage to integers
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	perPageInt, err := strconv.Atoi(perPage)
	if err != nil || perPageInt < 1 {
		perPageInt = 10
	}

	// Calculate the offset
	offset := (pageInt - 1) * perPageInt

	var totalData int64
	var users []entities.User

	// Create a query builder
	query := config.DB.Model(&entities.User{})

	// Apply filter by "active" column if provided
	if active != "" {
		activeBool, err := strconv.ParseBool(active)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid value for 'active' parameter",
			})
		}
		query = query.Where("active = ?", activeBool)
	}

	// Count total number of data
	if err := query.Count(&totalData).Error; err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalData) / float64(perPageInt)))

	// Fetch users for the current page with filter
	if err := query.Offset(offset).Limit(perPageInt).Find(&users).Error; err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	// Response payload
	response := fiber.Map{
		"users":       users,
		"page":        pageInt,
		"per_page":    perPageInt,
		"total_data":  totalData,
		"total_pages": totalPages,
	}

	return c.JSON(response)
}

func UserControllerCreate(c *fiber.Ctx) error {

	user := new(req.UserReq)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to hash password",
		})
	}

	newUser := entities.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: string(hashedPassword),
	}

	if err := config.DB.Create(&newUser).Error; err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(newUser)
}
