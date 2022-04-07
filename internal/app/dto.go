package app

import (
	"io"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/config"
)

type DTO struct {
	Config  *config.Config
	Storage storage.Storage
	Stdout  io.WriteCloser
	Stdin   io.ReadCloser
}
