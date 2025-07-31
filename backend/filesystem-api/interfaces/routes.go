package interfaces

import (
	"net/http"
)

// SetupFileSystemRoutes configura todas as rotas da API de sistema de arquivos
func SetupFileSystemRoutes(mux *http.ServeMux, handler *FileSystemHandler) {
	// Operações de arquivo
	mux.HandleFunc("/read-file/", handler.ReadFileHandler)
	mux.HandleFunc("/create-file/", handler.CreateFileHandler)
	mux.HandleFunc("/edit-file/", handler.EditFileHandler)
	mux.HandleFunc("/delete-file/", handler.DeleteFileHandler)
	
	// Operações de diretório
	mux.HandleFunc("/read-dir/", handler.ReadDirectoryHandler)
	mux.HandleFunc("/create-dir/", handler.CreateDirectoryHandler)
	mux.HandleFunc("/edit-dir/", handler.EditDirectoryHandler)
	mux.HandleFunc("/delete-dir/", handler.DeleteDirectoryHandler)
	
	// Execução de comandos
	mux.HandleFunc("/execute-command/", handler.ExecuteCommandHandler)
	
	// Endpoint de status/health check
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok", "service": "filesystem-api"}`))
	})
}