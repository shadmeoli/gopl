package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	var rootCmd = &cobra.Command{Use: "api-generator"}

	var projectName string

	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new Go API project structure",
		Run: func(cmd *cobra.Command, args []string) {
			createProjectStructure(projectName)
		},
	}

	createCmd.Flags().StringVarP(&projectName, "project-name", "p", "", "Specify the project directory name")

	viper.BindPFlag("project-name", createCmd.Flags().Lookup("project-name"))

	rootCmd.AddCommand(createCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func createProjectStructure(projectName string) {
	// Check if the project directory already exists
	if _, err := os.Stat(projectName); !os.IsNotExist(err) {
		log.Fatalf("Directory %s already exists. Please choose a different project name.", projectName)
	}

	// Create the project directory
	if err := os.MkdirAll(projectName, os.ModePerm); err != nil {
		log.Fatalf("Failed to create project directory: %v", err)
	}

	// Change to the project directory
	if err := os.Chdir(projectName); err != nil {
		log.Fatalf("Failed to change to the project directory: %v", err)
	}

	// Define the directory and file structure
	directories := []string{
		"cmd/myapi",
		"internal/api/handlers",
		"internal/app/config",
		"internal/app/database/postgres",
		"internal/app/middleware",
		"api/v1/routes",
		"scripts",
		"web",
	}
	files := []string{
		"Dockerfile",
		"go.mod",
		"go.sum",
		"README.md",
		"cmd/myapi/main.go",
		"internal/api/api.go",
		"internal/app/app.go",
		"internal/app/database/postgres/postgres.go",
		"api/v1/routes/routes.go",
	}

	// Create subdirectories
	for _, dir := range directories {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}
	}

	// Create necessary files
	for _, file := range files {
		if _, err := os.Create(file); err != nil {
			log.Fatalf("Failed to create file: %v", err)
		}
	}

	// Install Go dependencies
	cmd := exec.Command("go", "get", "-u", "github.com/gofiber/fiber/v2")
	cmd.Dir = projectName
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to install Fiber: %v", err)
	}

	cmd = exec.Command("go", "get", "-u", "gorm.io/gorm")
	cmd.Dir = projectName
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to install GORM: %v", err)
	}

	cmd = exec.Command("go", "get", "-u", "gorm.io/driver/postgres")
	cmd.Dir = projectName
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to install PostgreSQL driver: %v", err)
	}

	cmd = exec.Command("go", "get", "-u", "github.com/dgrijalva/jwt-go")
	cmd.Dir = projectName
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to install JWT library: %v", err)
	}

	// Create a Dockerfile
	dockerfileContent := `
# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the local code to the container
COPY . .

# Build the Go application
RUN go build -o main cmd/myapi/main.go

# Expose the port the application will run on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
`
	dockerfilePath := filepath.Join(projectName, "Dockerfile")
	if err := os.WriteFile(dockerfilePath, []byte(dockerfileContent), os.ModePerm); err != nil {
		log.Fatalf("Failed to create Dockerfile: %v", err)
	}

	// Build a Docker image
	cmd = exec.Command("docker", "build", "-t", projectName, ".")
	cmd.Dir = projectName
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to build Docker image: %v", err)
	}

	fmt.Printf("API project structure created successfully in the %s directory, and a Docker image has been built with the same name.\n", projectName)
}
