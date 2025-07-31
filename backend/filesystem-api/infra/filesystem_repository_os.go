package infra

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cristianorosa/eSimulate/backend/filesystem-api/domain"
)

// FileSystemRepositoryOS implementa FileSystemRepository usando o sistema operacional
type FileSystemRepositoryOS struct {
	// Diretório base para operações (por segurança)
	BaseDir string
}

// NewFileSystemRepositoryOS cria uma nova instância do repository
func NewFileSystemRepositoryOS(baseDir string) *FileSystemRepositoryOS {
	return &FileSystemRepositoryOS{
		BaseDir: baseDir,
	}
}

// validatePath valida se o caminho está dentro do diretório base (segurança)
func (r *FileSystemRepositoryOS) validatePath(path string) (string, error) {
	// Limpa o caminho
	cleanPath := filepath.Clean(path)
	
	// Se não for absoluto, torna relativo ao BaseDir
	if !filepath.IsAbs(cleanPath) {
		cleanPath = filepath.Join(r.BaseDir, cleanPath)
	}
	
	// Verifica se está dentro do BaseDir
	absBaseDir, err := filepath.Abs(r.BaseDir)
	if err != nil {
		return "", fmt.Errorf("erro ao resolver diretório base: %v", err)
	}
	
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", fmt.Errorf("erro ao resolver caminho: %v", err)
	}
	
	if !strings.HasPrefix(absPath, absBaseDir) {
		return "", fmt.Errorf("acesso negado: caminho fora do diretório permitido")
	}
	
	return absPath, nil
}

// ReadFile lê o conteúdo de um arquivo
func (r *FileSystemRepositoryOS) ReadFile(path string) (*domain.FileContent, error) {
	validPath, err := r.validatePath(path)
	if err != nil {
		return nil, err
	}
	
	content, err := os.ReadFile(validPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo: %v", err)
	}
	
	stat, err := os.Stat(validPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter informações do arquivo: %v", err)
	}
	
	return &domain.FileContent{
		Path:    path,
		Content: string(content),
		Size:    stat.Size(),
	}, nil
}

// CreateFile cria um novo arquivo
func (r *FileSystemRepositoryOS) CreateFile(path, content string) error {
	validPath, err := r.validatePath(path)
	if err != nil {
		return err
	}
	
	// Verifica se o arquivo já existe
	if _, err := os.Stat(validPath); err == nil {
		return fmt.Errorf("arquivo já existe: %s", path)
	}
	
	// Cria diretórios pai se necessário
	dir := filepath.Dir(validPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretórios pai: %v", err)
	}
	
	// Cria o arquivo
	if err := os.WriteFile(validPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("erro ao criar arquivo: %v", err)
	}
	
	return nil
}

// UpdateFile atualiza o conteúdo de um arquivo existente
func (r *FileSystemRepositoryOS) UpdateFile(path, content string) error {
	validPath, err := r.validatePath(path)
	if err != nil {
		return err
	}
	
	// Verifica se o arquivo existe
	if _, err := os.Stat(validPath); os.IsNotExist(err) {
		return fmt.Errorf("arquivo não encontrado: %s", path)
	}
	
	// Atualiza o arquivo
	if err := os.WriteFile(validPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("erro ao atualizar arquivo: %v", err)
	}
	
	return nil
}

// DeleteFile remove um arquivo
func (r *FileSystemRepositoryOS) DeleteFile(path string) error {
	validPath, err := r.validatePath(path)
	if err != nil {
		return err
	}
	
	// Verifica se é um arquivo (não diretório)
	stat, err := os.Stat(validPath)
	if err != nil {
		return fmt.Errorf("arquivo não encontrado: %s", path)
	}
	
	if stat.IsDir() {
		return fmt.Errorf("caminho é um diretório, não um arquivo: %s", path)
	}
	
	if err := os.Remove(validPath); err != nil {
		return fmt.Errorf("erro ao deletar arquivo: %v", err)
	}
	
	return nil
}

// ReadDirectory lê o conteúdo de um diretório
func (r *FileSystemRepositoryOS) ReadDirectory(path string) (*domain.DirectoryContent, error) {
	validPath, err := r.validatePath(path)
	if err != nil {
		return nil, err
	}
	
	entries, err := os.ReadDir(validPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler diretório: %v", err)
	}
	
	var files []*domain.FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue // Pula arquivos com erro
		}
		
		fileInfo := &domain.FileInfo{
			Name:    info.Name(),
			Path:    filepath.Join(path, info.Name()),
			Size:    info.Size(),
			IsDir:   info.IsDir(),
			ModTime: info.ModTime(),
			Mode:    info.Mode().String(),
		}
		files = append(files, fileInfo)
	}
	
	return &domain.DirectoryContent{
		Path:  path,
		Files: files,
		Total: len(files),
	}, nil
}

// CreateDirectory cria um novo diretório
func (r *FileSystemRepositoryOS) CreateDirectory(path string) error {
	validPath, err := r.validatePath(path)
	if err != nil {
		return err
	}
	
	// Verifica se já existe
	if _, err := os.Stat(validPath); err == nil {
		return fmt.Errorf("diretório já existe: %s", path)
	}
	
	if err := os.MkdirAll(validPath, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório: %v", err)
	}
	
	return nil
}

// RenameDirectory renomeia/move um diretório
func (r *FileSystemRepositoryOS) RenameDirectory(oldPath, newPath string) error {
	validOldPath, err := r.validatePath(oldPath)
	if err != nil {
		return err
	}
	
	validNewPath, err := r.validatePath(newPath)
	if err != nil {
		return err
	}
	
	// Verifica se o diretório origem existe
	stat, err := os.Stat(validOldPath)
	if err != nil {
		return fmt.Errorf("diretório não encontrado: %s", oldPath)
	}
	
	if !stat.IsDir() {
		return fmt.Errorf("caminho não é um diretório: %s", oldPath)
	}
	
	// Verifica se o destino já existe
	if _, err := os.Stat(validNewPath); err == nil {
		return fmt.Errorf("destino já existe: %s", newPath)
	}
	
	if err := os.Rename(validOldPath, validNewPath); err != nil {
		return fmt.Errorf("erro ao renomear diretório: %v", err)
	}
	
	return nil
}

// DeleteDirectory remove um diretório
func (r *FileSystemRepositoryOS) DeleteDirectory(path string) error {
	validPath, err := r.validatePath(path)
	if err != nil {
		return err
	}
	
	// Verifica se é um diretório
	stat, err := os.Stat(validPath)
	if err != nil {
		return fmt.Errorf("diretório não encontrado: %s", path)
	}
	
	if !stat.IsDir() {
		return fmt.Errorf("caminho não é um diretório: %s", path)
	}
	
	if err := os.RemoveAll(validPath); err != nil {
		return fmt.Errorf("erro ao deletar diretório: %v", err)
	}
	
	return nil
}

// ExecuteCommand executa um comando no sistema
func (r *FileSystemRepositoryOS) ExecuteCommand(command string) (*domain.CommandResult, error) {
	// Por segurança, limita comandos permitidos
	allowedCommands := []string{"ls", "pwd", "echo", "cat", "head", "tail", "grep", "find"}
	
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return nil, fmt.Errorf("comando vazio")
	}
	
	cmdName := parts[0]
	allowed := false
	for _, allowedCmd := range allowedCommands {
		if cmdName == allowedCmd {
			allowed = true
			break
		}
	}
	
	if !allowed {
		return &domain.CommandResult{
			Command:  command,
			Output:   "",
			Error:    "Comando não permitido por motivos de segurança",
			ExitCode: 1,
		}, nil
	}
	
	// Executa o comando no diretório base
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = r.BaseDir
	
	output, err := cmd.CombinedOutput()
	
	result := &domain.CommandResult{
		Command:  command,
		Output:   string(output),
		ExitCode: 0,
	}
	
	if err != nil {
		result.Error = err.Error()
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
		} else {
			result.ExitCode = 1
		}
	}
	
	return result, nil
}

// FileExists verifica se um arquivo existe
func (r *FileSystemRepositoryOS) FileExists(path string) bool {
	validPath, err := r.validatePath(path)
	if err != nil {
		return false
	}
	
	stat, err := os.Stat(validPath)
	return err == nil && !stat.IsDir()
}

// DirectoryExists verifica se um diretório existe
func (r *FileSystemRepositoryOS) DirectoryExists(path string) bool {
	validPath, err := r.validatePath(path)
	if err != nil {
		return false
	}
	
	stat, err := os.Stat(validPath)
	return err == nil && stat.IsDir()
}