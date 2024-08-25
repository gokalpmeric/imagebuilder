package main

import (
	"archive/tar"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type BuildRequest struct {
	RepoURL        string `json:"repo_url"`
	DockerfilePath string `json:"dockerfile_path"`
}

func main() {
	r := gin.Default()
	r.POST("/build", buildHandler)
	fmt.Println("Server started at :8080")
	r.Run(":8080")
}

func buildHandler(c *gin.Context) {
	var req BuildRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repoDir := "repo"

	// Check if the directory exists and is not empty
	if dirExistsAndNotEmpty(repoDir) {
		// Remove the directory if it exists and is not empty
		if err := os.RemoveAll(repoDir); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Clone the repository
	if err := cloneRepo(req.RepoURL, repoDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Wait for the clone to complete
	if err := waitForClone(repoDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Build the Docker image
	if err := buildDockerImage(repoDir, req.DockerfilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := os.RemoveAll(repoDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Docker image built successfully"})
}

// dirExistsAndNotEmpty checks if a directory exists and is not empty
func dirExistsAndNotEmpty(dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	return err != io.EOF
}

// cloneRepo clones the repository from the given URL into the specified directory
func cloneRepo(url, dir string) error {
	cmd := exec.Command("git", "clone", url, dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// waitForClone waits for the clone operation to complete
func waitForClone(dir string) error {
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return nil
		}
	}
}

// buildDockerImage builds the Docker image from the Dockerfile
func buildDockerImage(dir, dockerfilePath string) error {
	cmd := exec.Command("docker", "build", "-t", "my-app:latest", "-f", dockerfilePath, ".")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// createTar creates a tar archive containing the Dockerfile
func createTar(dockerfile string) (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	// Add the Dockerfile to the tar archive
	err := tw.WriteHeader(&tar.Header{
		Name: "Dockerfile",
		Size: int64(len(dockerfile)),
	})
	if err != nil {
		return nil, err
	}

	_, err = tw.Write([]byte(dockerfile))
	if err != nil {
		return nil, err
	}

	// Close the tar writer
	err = tw.Close()
	if err != nil {
		return nil, err
	}

	return buf, nil
}
