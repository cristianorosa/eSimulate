package usecase

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	"github.com/cristianorosa/eSimulate/backend/filesystem-api/domain"
)

// FileSystemUsecase implementa as regras de negócio para operações de sistema de arquivos
// Inclui validação, sanitização e controle de acesso
type FileSystemUsecase struct {
	Repo domain.FileSystemRepository
}

// NewFileSystemUsecase cria uma nova instância do usecase
func NewFileSystemUsecase(repo domain.FileSystemRepository) *FileSystemUsecase {
	return &FileSystemUsecase{
		Repo: repo,
	}
}

// validatePath valida e sanitiza caminhos de entrada
func (uc *FileSystemUsecase) validatePath(path string) error {
	if path == "" {
		return errors.New("caminho não pode estar vazio")
	}
	
	// Remove espaços em branco
	path = strings.TrimSpace(path)
	
	// Verifica caracteres perigosos
	dangerousChars := []string{"../", "..\\", "<", ">", "|", "&", ";", "`", "$"}
	for _, char := range dangerousChars {
		if strings.Contains(path, char) {
			return errors.New("caminho contém caracteres não permitidos")
		}
	}
	
	return nil
}

// ReadFile lê o conteúdo de um arquivo com validações
func (uc *FileSystemUsecase) ReadFile(ctx context.Context, path string) (*domain.FileContent, error) {
	if err := uc.validatePath(path); err != nil {
		return nil, err
	}
	
	if !uc.Repo.FileExists(path) {
		return nil, errors.New("arquivo não encontrado")
	}
	
	return uc.Repo.ReadFile(path)
}

// CreateFile cria um novo arquivo com validações
func (uc *FileSystemUsecase) CreateFile(ctx context.Context, path, content string) error {
	if err := uc.validatePath(path); err != nil {
		return err
	}
	
	// Verifica se já existe
	if uc.Repo.FileExists(path) {
		return errors.New("arquivo já existe")
	}
	
	// Verifica se o diretório pai existe ou pode ser criado
	dir := filepath.Dir(path)
	if dir != "." && !uc.Repo.DirectoryExists(dir) {
		// Tenta criar o diretório pai
		if err := uc.Repo.CreateDirectory(dir); err != nil {
			return errors.New("não foi possível criar diretório pai: " + err.Error())
		}
	}
	
	return uc.Repo.CreateFile(path, content)
}

// UpdateFile atualiza um arquivo existente
func (uc *FileSystemUsecase) UpdateFile(ctx context.Context, path, content string) error {
	if err := uc.validatePath(path); err != nil {
		return err
	}
	
	if !uc.Repo.FileExists(path) {
		return errors.New("arquivo não encontrado")
	}
	
	return uc.Repo.UpdateFile(path, content)
}

// DeleteFile remove um arquivo
func (uc *FileSystemUsecase) DeleteFile(ctx context.Context, path string) error {
	if err := uc.validatePath(path); err != nil {
		return err
	}
	
	if !uc.Repo.FileExists(path) {
		return errors.New("arquivo não encontrado")
	}
	
	return uc.Repo.DeleteFile(path)
}

// ReadDirectory lê o conteúdo de um diretório
func (uc *FileSystemUsecase) ReadDirectory(ctx context.Context, path string) (*domain.DirectoryContent, error) {
	if err := uc.validatePath(path); err != nil {
		return nil, err
	}
	
	if !uc.Repo.DirectoryExists(path) {
		return nil, errors.New("diretório não encontrado")
	}
	
	return uc.Repo.ReadDirectory(path)
}

// CreateDirectory cria um novo diretório
func (uc *FileSystemUsecase) CreateDirectory(ctx context.Context, path string) error {
	if err := uc.validatePath(path); err != nil {
		return err
	}
	
	if uc.Repo.DirectoryExists(path) {
		return errors.New("diretório já existe")
	}
	
	return uc.Repo.CreateDirectory(path)
}

// RenameDirectory renomeia/move um diretório
func (uc *FileSystemUsecase) RenameDirectory(ctx context.Context, oldPath, newPath string) error {
	if err := uc.validatePath(oldPath); err != nil {
		return errors.New("caminho origem inválido: " + err.Error())
	}
	
	if err := uc.validatePath(newPath); err != nil {
		return errors.New("caminho destino inválido: " + err.Error())
	}
	
	if !uc.Repo.DirectoryExists(oldPath) {
		return errors.New("diretório origem não encontrado")
	}
	
	if uc.Repo.DirectoryExists(newPath) {
		return errors.New("diretório destino já existe")
	}
	
	return uc.Repo.RenameDirectory(oldPath, newPath)
}

// DeleteDirectory remove um diretório
func (uc *FileSystemUsecase) DeleteDirectory(ctx context.Context, path string) error {
	if err := uc.validatePath(path); err != nil {
		return err
	}
	
	if !uc.Repo.DirectoryExists(path) {
		return errors.New("diretório não encontrado")
	}
	
	// Proteção: não permite deletar diretório raiz
	if path == "/" || path == "." || path == ".." {
		return errors.New("não é possível deletar diretório raiz")
	}
	
	return uc.Repo.DeleteDirectory(path)
}

// ExecuteCommand executa um comando com validações de segurança
func (uc *FileSystemUsecase) ExecuteCommand(ctx context.Context, command string) (*domain.CommandResult, error) {
	if command == "" {
		return nil, errors.New("comando não pode estar vazio")
	}
	
	// Remove espaços extras
	command = strings.TrimSpace(command)
	
	// Verifica comandos perigosos
	dangerousCommands := []string{"rm", "rmdir", "del", "format", "fdisk", "mkfs", "dd", "sudo", "su", "chmod", "chown"}
	commandParts := strings.Fields(command)
	if len(commandParts) > 0 {
		for _, dangerous := range dangerousCommands {
			if commandParts[0] == dangerous {
				return &domain.CommandResult{
					Command:  command,
					Output:   "",
					Error:    "Comando não permitido por motivos de segurança",
					ExitCode: 1,
				}, nil
			}
		}
	}
	
	return uc.Repo.ExecuteCommand(command)
}

// GetFileInfo retorna informações sobre um arquivo ou diretório
func (uc *FileSystemUsecase) GetFileInfo(ctx context.Context, path string) (*domain.FileInfo, error) {
	if err := uc.validatePath(path); err != nil {
		return nil, err
	}
	
	// Tenta ler como diretório primeiro
	if uc.Repo.DirectoryExists(path) {
		_, err := uc.Repo.ReadDirectory(path)
		if err != nil {
			return nil, err
		}
		
		return &domain.FileInfo{
			Name:  filepath.Base(path),
			Path:  path,
			IsDir: true,
		}, nil
	}
	
	// Tenta ler como arquivo
	if uc.Repo.FileExists(path) {
		fileContent, err := uc.Repo.ReadFile(path)
		if err != nil {
			return nil, err
		}
		
		return &domain.FileInfo{
			Name:  filepath.Base(path),
			Path:  path,
			Size:  fileContent.Size,
			IsDir: false,
		}, nil
	}
	
	return nil, errors.New("arquivo ou diretório não encontrado")
}