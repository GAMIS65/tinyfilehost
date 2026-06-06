package server

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gamis65/tinyfilehost/internal/util"
)

type Handler struct {
	logger *slog.Logger
}

func NewRouter(logger *slog.Logger) *http.ServeMux {
	mux := http.NewServeMux()

	h := &Handler{logger: logger}

	mux.HandleFunc("GET /health", h.healthHandler)
	mux.HandleFunc("POST /upload", h.uploadHandler)

	fileServer := http.FileServer(http.Dir("./internal/static"))
	mux.Handle("/", fileServer)

	return mux
}

func (h *Handler) healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func (h *Handler) uploadHandler(w http.ResponseWriter, r *http.Request) {
	const fileSizeLimit = 10 * 1024 * 1024

	r.Body = http.MaxBytesReader(w, r.Body, fileSizeLimit)

	err := r.ParseMultipartForm(fileSizeLimit)
	if err != nil {
		h.logger.Error("Form was invalid or the file was too big", "error", err)
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		h.logger.Error("File not found in the form", "error", err)
		http.Error(w, "file missing", http.StatusBadRequest)
		return
	}
	defer file.Close()

	logger := h.logger.With("filename", handler.Filename, "headers", handler.Header, "size", handler.Size)
	
	logger.Info("Uploading a file")

	err = os.MkdirAll("uploads", 0755)
	if err != nil {
		logger.Error("Failed to create an uploads folder", "error", err)
		http.Error(w, "cannot create folder", http.StatusInternalServerError)
		return
	}

	fileId, err := util.GenerateId(6)
	if err != nil {
		logger.Error("Failed to generate an id for a file", "error", err)
		http.Error(w, "Failed to generate an id", http.StatusInternalServerError)
		return
	}

	newFileName := fileId + "." + filepath.Ext(handler.Filename)

	dstPath := filepath.Join("uploads", newFileName)
	dst, err := os.Create(dstPath)
	if err != nil {
		logger.Error("Failed to create a file", "error", err)
		http.Error(w, "cannot create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	hash := sha256.New()
	tee := io.TeeReader(file, hash)

	_, err = io.Copy(dst, tee)
	if err != nil {
		logger.Error("Failed to save a file", "error", err)
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}

	checksum := fmt.Sprintf("%x", hash.Sum(nil))

	logger.Info("File uploaded successfully", "checksum", checksum)

	fmt.Fprintln(w, "File uploaded successfully")
}