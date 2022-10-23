package config

import (
	"github.com/alexedwards/scs/v2"
	"html/template"
	"log"
)

type AppConfig struct {
	UseCache      bool
	InProduction  bool
	TemplateCache map[string]*template.Template
	Sessions      *scs.SessionManager
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
}
