package commands

import (
	"githum.com/Murchoid/iwashere/internal/domain/models"
	"githum.com/Murchoid/iwashere/internal/repository"
)

type Context struct {
	Args        []string
	Flags       map[string]string
	WorkDir     string
	ProjectPath string
	Config      *models.Config
	Repo        repository.Repository
}

type Command interface {
	Name() string
	Description() string
	Execute(ctx *Context) error
}
