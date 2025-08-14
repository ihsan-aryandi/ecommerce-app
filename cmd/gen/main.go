package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go make:<type> <Name>")
		return
	}

	cmd := os.Args[1]
	name := os.Args[2]

	switch cmd {
	case "make:handler", "make:service", "make:repository":
		createSingle(cmd, name)
	case "make:module":
		createModule(name)
	default:
		fmt.Println("Unknown command:", cmd)
	}
}

// Single file generation
func createSingle(cmd, name string) {
	nameLower := strings.ToLower(name)
	var folder, fileName, tmpl, providerSet string

	switch cmd {
	case "make:handler":
		folder = "internal/api/handler"
		fileName = fmt.Sprintf("%s-handler.go", nameLower)
		tmpl = handlerTemplateSingle
		providerSet = fmt.Sprintf("handler.New%sHandler,", name)
	case "make:service":
		folder = "internal/api/service"
		fileName = fmt.Sprintf("%s-service.go", nameLower)
		tmpl = serviceTemplateSingle
		providerSet = fmt.Sprintf("service.New%sService,", name)
	case "make:repository":
		folder = "internal/api/repository"
		fileName = fmt.Sprintf("%s-repository.go", nameLower)
		tmpl = repositoryTemplate
		providerSet = fmt.Sprintf("repository.New%sRepository,", name)
	}

	createFile(folder, fileName, tmpl, name, nameLower)
	appendProvider("internal/provider/providers.go", cmd, providerSet)
}

// Module: generate handler+service+repository dengan DI lengkap
func createModule(name string) {
	fmt.Println("Generating module:", name)
	createModuleFile("repository", name)
	createModuleFile("service", name)
	createModuleFile("handler", name)
}

func createModuleFile(kind, name string) {
	nameLower := strings.ToLower(name)
	var folder, fileName, tmpl, providerSet string

	switch kind {
	case "handler":
		folder = "internal/api/handler"
		fileName = fmt.Sprintf("%s-handler.go", nameLower)
		tmpl = handlerTemplateModule
		providerSet = fmt.Sprintf("handler.New%sHandler,", name)
	case "service":
		folder = "internal/api/service"
		fileName = fmt.Sprintf("%s-service.go", nameLower)
		tmpl = serviceTemplateModule
		providerSet = fmt.Sprintf("service.New%sService,", name)
	case "repository":
		folder = "internal/api/repository"
		fileName = fmt.Sprintf("%s-repository.go", nameLower)
		tmpl = repositoryTemplate
		providerSet = fmt.Sprintf("repository.New%sRepository,", name)
	}

	createFile(folder, fileName, tmpl, name, nameLower)
	appendProvider("internal/provider/providers.go", "make:"+kind, providerSet)
}

// Utility: create file from template
func createFile(folder, fileName, tmpl, name, nameLower string) {
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	filePath := filepath.Join(folder, fileName)
	if _, err := os.Stat(filePath); err == nil {
		fmt.Println("File already exists:", filePath)
		return
	}

	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	t := template.Must(template.New("file").Parse(tmpl))
	err = t.Execute(f, struct {
		Name      string
		NameLower string
	}{
		Name:      name,
		NameLower: nameLower,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created:", filePath)
}

// Append provider to providers.go
func appendProvider(filePath, cmd, line string) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", filePath, err)
	}

	strContent := string(content)
	var setName string
	switch cmd {
	case "make:handler":
		setName = "HandlersSet"
	case "make:service":
		setName = "ServicesSet"
	case "make:repository":
		setName = "RepositoriesSet"
	}

	marker := fmt.Sprintf("var %s = wire.NewSet(", setName)
	if !strings.Contains(strContent, line) {
		idx := strings.Index(strContent, marker)
		if idx == -1 {
			log.Fatalf("Cannot find marker '%s' in %s", marker, filePath)
		}
		before := strContent[:idx+len(marker)]
		after := strContent[idx+len(marker):]
		newContent := before + "\n\t" + line + after
		err = ioutil.WriteFile(filePath, []byte(newContent), 0644)
		if err != nil {
			log.Fatalf("Failed to write to %s: %v", filePath, err)
		}
		fmt.Println("Updated providers.go:", setName)
	} else {
		fmt.Println("Provider already exists in providers.go")
	}
}

// Templates

// Single file templates (no DI)
var handlerTemplateSingle = `package handler

type {{.Name}}Handler struct {}

func New{{.Name}}Handler() *{{.Name}}Handler {
	return &{{.Name}}Handler{}
}
`

var serviceTemplateSingle = `package service

type {{.Name}}Service struct {}

func New{{.Name}}Service() *{{.Name}}Service {
	return &{{.Name}}Service{}
}
`

// Module templates (with DI)
var handlerTemplateModule = `package handler

import "ecommerce-app/internal/api/service"

type {{.Name}}Handler struct {
	{{.NameLower}}Service *service.{{.Name}}Service
}

func New{{.Name}}Handler({{.NameLower}}Service *service.{{.Name}}Service) *{{.Name}}Handler {
	return &{{.Name}}Handler{
		{{.NameLower}}Service: {{.NameLower}}Service,
	}
}
`

var serviceTemplateModule = `package service

import "ecommerce-app/internal/api/repository"

type {{.Name}}Service struct {
	{{.NameLower}}Repository *repository.{{.Name}}Repository
}

func New{{.Name}}Service({{.NameLower}}Repository *repository.{{.Name}}Repository) *{{.Name}}Service {
	return &{{.Name}}Service{
		{{.NameLower}}Repository: {{.NameLower}}Repository,
	}
}
`

var repositoryTemplate = `package repository

import "database/sql"

type {{.Name}}Repository struct {
	db *sql.DB
}

func New{{.Name}}Repository(db *sql.DB) *{{.Name}}Repository {
	return &{{.Name}}Repository{db: db}
}
`
