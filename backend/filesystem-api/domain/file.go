package domain

import "time"

// FileInfo representa informações sobre um arquivo ou diretório
type FileInfo struct {
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	Size    int64     `json:"size"`
	IsDir   bool      `json:"is_dir"`
	ModTime time.Time `json:"mod_time"`
	Mode    string    `json:"mode"`
}

// FileContent representa o conteúdo de um arquivo
type FileContent struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Size    int64  `json:"size"`
}

// CommandResult representa o resultado de um comando executado
type CommandResult struct {
	Command  string `json:"command"`
	Output   string `json:"output"`
	Error    string `json:"error,omitempty"`
	ExitCode int    `json:"exit_code"`
}

// DirectoryContent representa o conteúdo de um diretório
type DirectoryContent struct {
	Path  string      `json:"path"`
	Files []*FileInfo `json:"files"`
	Total int         `json:"total"`
}