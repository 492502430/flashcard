package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

// Upload handles file upload and text extraction.
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10 MB max

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, 400, "file is required")
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		writeError(w, 500, "failed to read file")
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	var text string

	switch ext {
	case ".txt":
		text = string(content)
	case ".pdf", ".docx":
		text = extractViaAI(content, header.Filename)
	default:
		writeError(w, 400, "unsupported format: "+ext)
		return
	}

	if strings.TrimSpace(text) == "" {
		writeError(w, 400, "could not extract text from file")
		return
	}

	if len(text) > 50000 {
		text = text[:50000]
	}

	writeJSON(w, 200, map[string]interface{}{
		"text":     text,
		"filename": header.Filename,
		"size":     len(text),
	})
}

// extractViaAI forwards binary file to Python AI service for text extraction.
func extractViaAI(content []byte, filename string) string {
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)

	fw, _ := mw.CreateFormFile("file", filename)
	fw.Write(content)
	mw.Close()

	resp, err := http.Post("http://localhost:8001/extract", mw.FormDataContentType(), &b)
	if err != nil {
		return fmt.Sprintf("[Error: AI service unreachable — %v]", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Sprintf("[Error: extraction failed — %s]", string(body))
	}

	var result struct {
		Text string `json:"text"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Text
}
