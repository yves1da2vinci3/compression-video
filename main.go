package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/gofiber/fiber/v2"
)

func compressVideo(c *fiber.Ctx) error {
	// Get the video file from the request
	file, err := c.FormFile("video")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to get video file",
		})
	}

	// Create a temporary file to save the uploaded video
	tempFile, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open video file",
		})
	}
	defer tempFile.Close()

	// Create a destination file to save the uploaded video
	tempFileName := "temp_video.mp4"
	destFile, err := os.Create(tempFileName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create temporary file",
		})
	}
	defer destFile.Close()

	// Copy the uploaded video to the destination file
	_, err = io.Copy(destFile, tempFile)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save video file",
		})
	}

	// Perform video compression using FFmpeg
	compressedFileName := "compressed_video.mp4"
	cmd := exec.Command("ffmpeg", "-i", tempFileName, "-c:v", "libx264", "-crf", "23", "-preset", "fast", compressedFileName)
	err = cmd.Run()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to compress video",
		})
	}

	// Remove the temporary video file
	err = os.Remove(tempFileName)
	if err != nil {
		fmt.Printf("Failed to remove temporary file: %v\n", err)
	}

	// Return the compressed video as a response
	return c.SendFile(compressedFileName)
}

func main() {
	app := fiber.New()

	// Define the route for video compression
	app.Post("/compress", compressVideo)

	// Start the server
	err := app.Listen(":3000")
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
