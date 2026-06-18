package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

// Upload handles file upload (multipart or base64 JSON).
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		h.uploadMultipart(w, r)
	} else {
		h.uploadJSON(w, r)
	}
}

func (h *Handler) uploadJSON(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Filename string `json:"filename"`
		Content  string `json:"content"`
		Encoding string `json:"encoding"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "invalid json body")
		return
	}

	var data []byte
	if req.Encoding == "base64" {
		var err error
		data, err = base64.StdEncoding.DecodeString(req.Content)
		if err != nil {
			writeError(w, 400, "invalid base64")
			return
		}
	} else {
		data = []byte(req.Content)
	}

	text := extractText(data, req.Filename)
	if strings.TrimSpace(text) == "" {
		writeError(w, 400, "could not extract text from file")
		return
	}
	if len(text) > 50000 {
		text = text[:50000]
	}

	writeJSON(w, 200, map[string]interface{}{
		"text":     text,
		"filename": req.Filename,
		"size":     len(text),
	})
}

func (h *Handler) uploadMultipart(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

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

	text := extractText(content, header.Filename)
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

func extractText(content []byte, filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".txt":
		return string(content)
	case ".pdf", ".docx", ".png", ".jpg", ".jpeg":
		return extractViaAI(content, filename)
	default:
		return string(content)
	}
}

func extractViaAI(content []byte, filename string) string {
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	fw, _ := mw.CreateFormFile("file", filename)
	fw.Write(content)
	mw.Close()

	resp, err := http.Post("http://localhost:8001/extract", mw.FormDataContentType(), &b)
	if err != nil {
		return fmt.Sprintf("[Error: AI service unreachable]")
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
