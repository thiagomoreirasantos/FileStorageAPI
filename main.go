package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/thiagomoreirasantos/FileStorageAPI/docs" // Importação para o Swagger

	"github.com/gorilla/mux"
)

const storagePath = "/app/storage"

// uploadFile recebe um arquivo via formulário e o armazena no servidor.
//
// @Summary Faz upload de um arquivo
// @Description Envia um arquivo e salva no volume do Docker
// @Accept multipart/form-data
// @Produce text/plain
// @Param file formData file true "Arquivo a ser enviado"
// @Success 200 {string} string "Arquivo salvo com sucesso"
// @Failure 400 {string} string "Erro ao receber o arquivo"
// @Failure 500 {string} string "Erro ao salvar arquivo"
// @Router /upload [post]
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

// uploadFile recebe um arquivo via formulário e o armazena no servidor.
//
// @Summary Faz upload de um arquivo
// @Description Envia um arquivo e salva no volume do Docker
// @Accept multipart/form-data
// @Produce text/plain
// @Param file formData file true "Arquivo a ser enviado"
// @Success 200 {string} string "Arquivo salvo com sucesso"
// @Failure 400 {string} string "Erro ao receber o arquivo"
// @Failure 500 {string} string "Erro ao salvar arquivo"
// @Router /upload [post]
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

// downloadFile baixa um arquivo específico.
//
// @Summary Faz o download de um arquivo
// @Description Baixa um arquivo pelo nome
// @Produce application/octet-stream
// @Param filename path string true "Nome do arquivo"
// @Success 200 {file} file "Arquivo baixado com sucesso"
// @Failure 404 {string} string "Arquivo não encontrado"
// @Router /files/{filename} [get]
func DownloadFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileName := vars["filename"]
	filepath := filepath.Join(storagePath, fileName)

	http.ServeFile(w, r, filepath)
}

// deleteFile exclui um arquivo específico.
//
// @Summary Exclui um arquivo
// @Description Remove permanentemente um arquivo armazenado
// @Param filename path string true "Nome do arquivo"
// @Success 200 {string} string "Arquivo excluído com sucesso"
// @Failure 404 {string} string "Arquivo não encontrado"
// @Router /files/{filename} [delete]
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

// deleteFile exclui um arquivo específico.
//
// @Summary Exclui um arquivo
// @Description Remove permanentemente um arquivo armazenado
// @Param filename path string true "Nome do arquivo"
// @Success 200 {string} string "Arquivo excluído com sucesso"
// @Failure 404 {string} string "Arquivo não encontrado"
// @Router /files/{filename} [delete]
func main() {
	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r.HandleFunc("/upload", uploadFile).Methods("POST")
	r.HandleFunc("/files", ListFiles).Methods("GET")
	r.HandleFunc("/files/{filename}", DownloadFile).Methods("GET")
	r.HandleFunc("/files/{filename}", deleteFile).Methods("DELETE")

	fmt.Println("Starting server on port 8080...")
	http.ListenAndServe(":8080", r)
}
