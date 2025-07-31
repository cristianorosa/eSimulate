package interfaces

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/cristianorosa/eSimulate/backend/filesystem-api/usecase"
)

// FileSystemHandler lida com requisições HTTP relacionadas ao sistema de arquivos
type FileSystemHandler struct {
	UC *usecase.FileSystemUsecase
}

// NewFileSystemHandler cria uma nova instância do handler
func NewFileSystemHandler(uc *usecase.FileSystemUsecase) *FileSystemHandler {
	return &FileSystemHandler{
		UC: uc,
	}
}

// extractPathFromURL extrai o caminho da URL
func extractPathFromURL(requestPath, prefix string) (string, error) {
	// Remove o prefixo da URL
	path := strings.TrimPrefix(requestPath, prefix)
	
	// Decodifica URL encoding
	decodedPath, err := url.QueryUnescape(path)
	if err != nil {
		return "", err
	}
	
	// Remove barra inicial se existir
	if strings.HasPrefix(decodedPath, "/") {
		decodedPath = decodedPath[1:]
	}
	
	return decodedPath, nil
}

// ReadFileHandler trata o endpoint GET /read-file/{path}
func (h *FileSystemHandler) ReadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	
	path, err := extractPathFromURL(r.URL.Path, "/read-file/")
	if err != nil {
		http.Error(w, "Caminho inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	fileContent, err := h.UC.ReadFile(context.Background(), path)
	if err != nil {
		http.Error(w, "Erro ao ler arquivo: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fileContent)
}

// CreateFileHandler trata o endpoint PUT /create-file/{path}
func (h *FileSystemHandler) CreateFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	
	path, err := extractPathFromURL(r.URL.Path, "/create-file/")
	if err != nil {
		http.Error(w, "Caminho inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	var req struct {
		Content string `json:"content"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := h.UC.CreateFile(context.Background(), path, req.Content); err != nil {
		http.Error(w, "Erro ao criar arquivo: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Arquivo criado com sucesso",
		"path":    path,
	})
}

// EditFileHandler trata o endpoint PATCH /edit-file/{path}
func (h *FileSystemHandler) EditFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	
	path, err := extractPathFromURL(r.URL.Path, "/edit-file/")
	if err != nil {
		http.Error(w, "Caminho inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	var req struct {
		Content string `json:"content"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := h.UC.UpdateFile(context.Background(), path, req.Content); err != nil {
		http.Error(w, "Erro ao editar arquivo: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Arquivo editado com sucesso",
		"path":    path,
	})
}

// DeleteFileHandler trata o endpoint DELETE /delete-file/{path}
func (h *FileSystemHandler) DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	
	path, err := extractPathFromURL(r.URL.Path, "/delete-file/")
	if err != nil {
		http.Error(w, "Caminho inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := h.UC.DeleteFile(context.Background(), path); err != nil {
		http.Error(w, "Erro ao deletar arquivo: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Arquivo deletado com sucesso",
		"path":    path,
	})
}

// ReadDirectoryHandler trata o endpoint GET /read-dir/{path}
func (h *FileSystemHandler) ReadDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	
	path, err := extractPathFromURL(r.URL.Path, "/read-dir/")
	if err != nil {
		http.Error(w, "Caminho inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	// Se o caminho estiver vazio, usa o diretório atual
	if path == "" {
		path = "."
	}
	
	dirContent, err := h.UC.ReadDirectory(context.Background(), path)
	if err != nil {
		http.Error(w, "Erro ao ler diretório: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dirContent)
}

// CreateDirectoryHandler trata o endpoint PUT /create-dir/{path}
func (h *FileSystemHandler) CreateDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	
	path, err := extractPathFromURL(r.URL.Path, "/create-dir/")
	if err != nil {
		http.Error(w, "Caminho inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := h.UC.CreateDirectory(context.Background(), path); err != nil {
		http.Error(w, "Erro ao criar diretório: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Diretório criado com sucesso",
		"path":    path,
	})
}

// EditDirectoryHandler trata o endpoint PATCH /edit-dir/{path}
func (h *FileSystemHandler) EditDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	
	oldPath, err := extractPathFromURL(r.URL.Path, "/edit-dir/")
	if err != nil {
		http.Error(w, "Caminho inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	var req struct {
		NewPath string `json:"new_path"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	if req.NewPath == "" {
		http.Error(w, "Novo caminho é obrigatório", http.StatusBadRequest)
		return
	}
	
	if err := h.UC.RenameDirectory(context.Background(), oldPath, req.NewPath); err != nil {
		http.Error(w, "Erro ao renomear diretório: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Diretório renomeado com sucesso",
		"old_path": oldPath,
		"new_path": req.NewPath,
	})
}

// DeleteDirectoryHandler trata o endpoint DELETE /delete-dir/{path}
func (h *FileSystemHandler) DeleteDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	
	path, err := extractPathFromURL(r.URL.Path, "/delete-dir/")
	if err != nil {
		http.Error(w, "Caminho inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := h.UC.DeleteDirectory(context.Background(), path); err != nil {
		http.Error(w, "Erro ao deletar diretório: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Diretório deletado com sucesso",
		"path":    path,
	})
}

// ExecuteCommandHandler trata o endpoint POST /execute-command/{command}
func (h *FileSystemHandler) ExecuteCommandHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	
	command, err := extractPathFromURL(r.URL.Path, "/execute-command/")
	if err != nil {
		http.Error(w, "Comando inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	if command == "" {
		http.Error(w, "Comando é obrigatório", http.StatusBadRequest)
		return
	}
	
	result, err := h.UC.ExecuteCommand(context.Background(), command)
	if err != nil {
		http.Error(w, "Erro ao executar comando: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}