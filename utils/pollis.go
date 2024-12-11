package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func Cattfusion(prompt string) string {
	text := fmt.Sprintf("%s except its a really cute cat", prompt)
	text = strings.TrimSpace(prompt)
	re := regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")
	text = strings.ReplaceAll(text, " ", "-")

	url := fmt.Sprintf("https://image.pollinations.ai/prompt/%s", text)

	client := &http.Client{}

	fmt.Println("catting...")
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Sprintf("bad query, (symbols should be avoided c: ) %s", err)
	}
	defer resp.Body.Close()

	// Check if response is an image
	if resp.StatusCode != http.StatusOK {
		return "failed to get some important stuff, might be a perms issue?"
	}
	if contentType := resp.Header.Get("Content-Type"); contentType != "image/jpeg" {
		return "content type mismatch"
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll("./catts", 0755); err != nil {
		return "failed to get some important stuff, might be a perms issue?"
	}

	// Generate unique filename using timestamp
	filename := fmt.Sprintf("image_%d.jpg", time.Now().Unix())
	cwd, err := os.Getwd()
	if err != nil {
		return "failed to get some important stuff, might be a perms issue?"
	}
	outputPath := filepath.Join(cwd, "./catts", filename)

	// Create output file
	out, err := os.Create(outputPath)
	if err != nil {
		return "failed to get some important stuff, might be a perms issue?"
	}
	defer out.Close()

	// Copy response body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "failed to get some important stuff, might be a perms issue?"
	}

	return outputPath
}
