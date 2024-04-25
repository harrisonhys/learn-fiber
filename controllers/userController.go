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
	//Validate
	if err := user.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed Input new data",
			"error":   err.Error(),
		})
	}
	// Check Email Already Use
	var existingUser entities.User
	result := config.DB.Where("email = ?", user.Email).First(&existingUser)
	if result.Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Email already in use",
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

func UserControllerById(c *fiber.Ctx) error {
	userId := c.Params("userId")
	var user entities.User
	result := config.DB.First(&user, userId)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}
	return c.JSON(user)
}

func UserControllerUpdate(c *fiber.Ctx) error {
	userId := c.Params("userId")
	var user entities.User
	result := config.DB.First(&user, userId)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// Parse the request body into user request struct
	updateData := new(req.UserUpdateReq)
	if err := c.BodyParser(updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Validate the updated data
	if err := updateData.ValidateOnUpdate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"error":   err.Error(),
		})
	}

	// Check if the email is being updated to one that is already in use
	if updateData.Email != user.Email {
		var existingUser entities.User
		result := config.DB.Where("email = ?", updateData.Email).First(&existingUser)
		if result.Error == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "Email already in use",
			})
		}
	}

	// Hash the new password if it's provided
	if updateData.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateData.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to hash password",
			})
		}
		user.Password = string(hashedPassword)
	}

	// Update the user fields
	user.Name = updateData.Name
	user.Email = updateData.Email

	// Save the updated user
	if err := config.DB.Save(&user).Error; err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(user)
}

func UserControllerDelete(c *fiber.Ctx) error {
	userId := c.Params("userId")
	var user entities.User
	result := config.DB.First(&user, userId)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}
	config.DB.Delete(&user)
	return c.Status(fiber.StatusNoContent).JSON(nil)
}
