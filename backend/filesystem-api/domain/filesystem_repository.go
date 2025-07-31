package domain

// FileSystemRepository define as operações de sistema de arquivos
// Interface para facilitar testes e desacoplamento (Clean Architecture)
type FileSystemRepository interface {
	// Operações de arquivo
	ReadFile(path string) (*FileContent, error)
	CreateFile(path, content string) error
	UpdateFile(path, content string) error
	DeleteFile(path string) error
	
	// Operações de diretório
	ReadDirectory(path string) (*DirectoryContent, error)
	CreateDirectory(path string) error
	RenameDirectory(oldPath, newPath string) error
	DeleteDirectory(path string) error
	
	// Operações de sistema
	ExecuteCommand(command string) (*CommandResult, error)
	
	// Utilitários
	FileExists(path string) bool
	DirectoryExists(path string) bool
}