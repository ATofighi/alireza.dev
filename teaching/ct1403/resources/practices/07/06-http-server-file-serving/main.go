package main

import (
	"fmt"
	"html"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const rootDirectory = "../files"

func main() {
	if err := os.MkdirAll(rootDirectory, 0755); err != nil {
		log.Fatalf("Could not create/access root directory: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGet(w, r)
		case http.MethodPost:
			handlePost(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Serving on http://localhost:9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	// Clean the URL path and ensure we don't allow directory traversal.
	requestPath := path.Clean(r.URL.Path)
	safePath, err := getSafePath(requestPath)
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Check if safePath is a directory or a file
	info, err := os.Stat(safePath)
	if os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, "Error accessing the file", http.StatusInternalServerError)
		return
	}

	// If it's a directory, list contents
	if info.IsDir() {
		err = listDirectory(w, safePath, requestPath)
		if err != nil {
			http.Error(w, "Error listing directory", http.StatusInternalServerError)
		}
		return
	}

	// Otherwise, serve the file with the correct MIME type
	serveFile(w, r, safePath)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	// We expect the URL path to be the location of the file to write.
	requestPath := path.Clean(r.URL.Path)
	safePath, err := getSafePath(requestPath)
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Read the request body (file content).
	content, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Create or overwrite the file at the safePath.
	// For more security, consider checking if file already exists or restricting overwrites.
	err = os.WriteFile(safePath, content, 0644)
	if err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "File %q successfully created/updated.\n", html.EscapeString(requestPath))
}

// getSafePath joins the user request path with the rootDirectory
// and ensures the final path is still within the rootDirectory
// (helps prevent path traversal: `../` etc.).
func getSafePath(requestPath string) (string, error) {
	// Join with the root directory
	fullPath := filepath.Join(rootDirectory, requestPath)

	// Resolve any symlinks or relative paths
	absFullPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", err
	}

	absRoot, err := filepath.Abs(rootDirectory)
	if err != nil {
		return "", err
	}

	// Ensure the resolved path starts with the root directory path
	if !strings.HasPrefix(absFullPath, absRoot) {
		return "", fmt.Errorf("path outside of root directory")
	}

	return absFullPath, nil
}

// listDirectory reads the specified directory and displays an HTML index of its contents.
func listDirectory(w http.ResponseWriter, dirPath string, requestPath string) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<html><body>")
	fmt.Fprintf(w, "<h1>Index of %s</h1>", html.EscapeString(requestPath))
	fmt.Fprintf(w, "<ul>")

	// Add a link to go up one level (unless we are at root)
	if requestPath != "/" {
		parent := path.Dir(requestPath)
		if parent == "." {
			parent = "/"
		}
		fmt.Fprintf(w, `<li><a href="%s">..</a></li>`, html.EscapeString(parent))
	}

	for _, entry := range entries {
		name := entry.Name()
		linkPath := path.Join(requestPath, name)
		displayName := html.EscapeString(name)
		if entry.IsDir() {
			fmt.Fprintf(w, `<li><a href="%s">%s/</a></li>`, html.EscapeString(linkPath), displayName)
		} else {
			fmt.Fprintf(w, `<li><a href="%s">%s</a></li>`, html.EscapeString(linkPath), displayName)
		}
	}

	fmt.Fprintf(w, "</ul></body></html>")
	return nil
}

// serveFile serves the file located at safePath with appropriate MIME type detection.
func serveFile(w http.ResponseWriter, r *http.Request, safePath string) {
	file, err := os.Open(safePath)
	if err != nil {
		http.Error(w, "Cannot open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Read a small portion of the file to detect MIME type
	const sniffLen = 512
	buf := make([]byte, sniffLen)
	n, _ := file.Read(buf)
	fileType := http.DetectContentType(buf[:n])

	// If no known detection, fallback to extension-based
	if fileType == "application/octet-stream" {
		ext := filepath.Ext(safePath)
		if ext != "" {
			mt := mime.TypeByExtension(ext)
			if mt != "" {
				fileType = mt
			}
		}
	}

	// Reset file offset to serve the entire file
	if _, err := file.Seek(0, 0); err != nil {
		http.Error(w, "Failed to seek file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", fileType)
	// Serve the file content
	if _, err := io.Copy(w, file); err != nil {
		http.Error(w, "Failed to send file", http.StatusInternalServerError)
		return
	}
}
