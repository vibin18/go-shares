package config

import (
	"github.com/vibin18/go-shares/internal/models"
	"html/template"
	"log"
)

type AppConfig struct {
	TemplateCache     map[string]*template.Template
	UseCache          bool
	SessionLifetime   float64
	InfoLog           *log.Logger
	ProdMode          bool
	ShareCache        *[]models.Stock
	DashShareList     []string
	DashShareCodeList []string
}
