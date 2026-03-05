package commands

import (
	"github.com/Murchoid/iwashere/internal/domain/models"
	"github.com/Murchoid/iwashere/internal/repository"
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
	Usage() string      // Detailed usage example
	Examples() []string // Multiple examples
	Execute(ctx *Context) error
}

type BaseCommand struct {
	NameStr      string
	DescStr      string
	UsageStr     string
	ExamplesList []string
}

func (c *BaseCommand) Name() string        { return c.NameStr }
func (c *BaseCommand) Description() string { return c.DescStr }
func (c *BaseCommand) Usage() string       { return c.UsageStr }
func (c *BaseCommand) Examples() []string  { return c.ExamplesList }
