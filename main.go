package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

const storagePath = "/app/storage"

func uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10 MB
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error Retrieving the File", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		os.MkdirAll(storagePath, os.ModePerm)
	}

	filePath := filepath.Join(storagePath, handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error Uploading the File", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error Uploading the File", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Successfully Uploaded File: %s\n", handler.Filename)
}

func ListFiles(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(storagePath)
	if err != nil {
		http.Error(w, "Error Listing Files", http.StatusInternalServerError)
		return
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Files: %v\n", strings.Join(fileNames, ", "))
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileName := vars["filename"]
	filepath := filepath.Join(storagePath, fileName)

	http.ServeFile(w, r, filepath)
}

func deleteFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileName := vars["filename"]
	filepath := filepath.Join(storagePath, fileName)

	err := os.Remove(filepath)
	if err != nil {
		http.Error(w, "Error Deleting the File", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Successfully Deleted File: %s\n", fileName)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/upload", uploadFile).Methods("POST")
	r.HandleFunc("/files", ListFiles).Methods("GET")
	r.HandleFunc("/files/{filename}", DownloadFile).Methods("GET")
	r.HandleFunc("/files/{filename}", deleteFile).Methods("DELETE")

	fmt.Println("Starting server on port 8080...")
	http.ListenAndServe(":8080", r)
}
