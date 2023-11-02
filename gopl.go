package main

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// main CLI entry
func main() {
	var rootCmd = &cobra.Command{Use: "gopl"}

	var projectName string
	var useDocker bool

	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new Go API project structure",
		Run: func(cmd *cobra.Command, args []string) {
			createProjectStructure(projectName, useDocker)
		},
	}

	createCmd.Flags().StringVarP(&projectName, "project-name", "p", "", "Specify the project directory name")
	createCmd.Flags().BoolVarP(&useDocker, "use-docker", "d", false, "Configure Docker for the project")

	viper.BindPFlag("project-name", createCmd.Flags().Lookup("project-name"))
	viper.BindPFlag("use-docker", createCmd.Flags().Lookup("use-docker"))

	rootCmd.AddCommand(createCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// InstallLibrary installs a Go library and returns an error if the installation fails.
func InstallLibrary(library string) error {
	cmd := exec.Command("go", "get", "-u", library)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error installing %s: %v", library, err)
	}
	return nil
}

// creator function
func createProjectStructure(projectName string, useDocker bool) {

	// Check if the project directory already exists
	if _, err := os.Stat(projectName); !os.IsNotExist(err) {
		log.Fatalf("Directory %s already exists. Please choose a different project name.", projectName)
	}
	// Create the project directory and its parent directories
	if err := os.MkdirAll(projectName, os.ModePerm); err != nil {
		log.Fatalf("Failed to create project directory: %v", err)
	}
	// Change to the project directory
	if err := os.Chdir(projectName); err != nil {
		log.Fatalf("Failed to change to the project directory: %v", err)
	}

	// Defining the directory and file structure
	directories := []string{
		"cmd",
		"internal/api/handlers",
		"internal/app/config",
		"internal/app/database/postgres",
		"internal/app/middleware",
		"api/v1/routes",
		"scripts",
		"web",
	}

	// file structure definition
	files := []string{
		"Dockerfile",
		"go.mod",
		"go.sum",
		"README.md",
		"cmd/main.go",
		"internal/api/api.go",
		"internal/app/app.go",
		"internal/app/database/postgres/postgres.go",
		"api/v1/routes/routes.go",
	}

	// Creating subdirectories
	dirsProgressBar := pb.StartNew(len(directories))
	for _, dir := range directories {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}
		dirsProgressBar.Increment()
	}
	dirsProgressBar.Finish()

	// Creating necessary files
	filesProgressBar := pb.StartNew(len(files))
	for _, file := range files {
		if _, err := os.Create(file); err != nil {
			log.Fatalf("Failed to create file: %v", err)
		}
		filesProgressBar.Increment()
	}
	filesProgressBar.Finish()

	// installing needed dependencies for an API -  libraries of choice [ go-fiber, goorm, jwt-go, postgres ]
	// [
	// go get -u github.com/gofiber/fiber/v2
	// go get -u gorm.io/gorm
	// go get -u gorm.io/driver/postgres
	// go get -u github.com/dgrijalva/jwt-go
	// ]

	// List of library import paths to install
	librariesToInstall := []string{
		"github.com/gofiber/fiber/v2",
		"gorm.io/gorm",
		"gorm.io/driver/postgres",
		"github.com/dgrijalva/jwt-go",
	}

	// Create a progress bar
	libsInstallProgressBar := pb.StartNew(100)

	// Use a WaitGroup to wait for all installations to complete
	var wg sync.WaitGroup

	for _, library := range librariesToInstall {
		wg.Add(1)
		go func(lib string) {
			defer wg.Done()
			if err := InstallLibrary(lib); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Installed %s\n", lib)
			}
			libsInstallProgressBar.Increment()
		}(library)
	} // Wait for all installations to complete
	wg.Wait()

	// Close the progress bar
	libsInstallProgressBar.Finish()

	fmt.Println("All installations complete")

	// Prompting the user to create a .env file
	fmt.Print("Do you want to create a .env file? (Y/n): ")
	var createEnvFile string
	fmt.Scanln(&createEnvFile)

	if strings.ToLower(createEnvFile) == "y" {
		// Create a .env file
		envFilePath := filepath.Join(".env")
		envFileContent := `
# Environment Configuration
DATABASE_URL=your_database_url
SECRET_KEY=your_secret_key
# Add other environment variables here
`
		if err := os.WriteFile(envFilePath, []byte(strings.TrimSpace(envFileContent)), os.ModePerm); err != nil {
			log.Fatalf("Failed to create .env file: %v", err)
		}
		fmt.Printf("Created a .env file in the %s directory.\n", projectName)
	}

	if useDocker {
		// Create a Dockerfile
		dockerfileContent := `
FROM golang:latest
WORKDIR /app
COPY . .
RUN go build -o main cmd/myapi/main.go
EXPOSE 8080
CMD ["./main"]
`
		dockerfilePath := filepath.Join("Dockerfile")
		if err := os.WriteFile(dockerfilePath, []byte(dockerfileContent), os.ModePerm); err != nil {

			log.Fatalf("Failed to create Dockerfile: %v", err)
		}

		// Build a Docker image
		cmd := exec.Command("docker", "build", "-t", projectName, ".")
		cmd.Dir = projectName
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to build Docker image: %v", err)
		}

		fmt.Printf("API project structure created successfully in the %s directory, and a Docker image has been built with the same name.\n", projectName)
	} else {
		fmt.Printf("API project structure created successfully in the %s directory.\n", projectName)
	}
}
