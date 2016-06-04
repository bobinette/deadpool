package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

type Generator interface {
	Generate(name string) error
}

type generator struct {
	lister Lister
}

func NewGenerator() Generator {
	return &generator{
		lister: NewLister(),
	}
}

func (g *generator) Generate(name string) error {
	lower := strings.ToLower(name)

	// Create the architecture
	// {{ name }}
	// |_ clients
	// |_ protos
	// |_ server
	dirs := []string{
		lower, fmt.Sprintf("%s/clients", lower), fmt.Sprintf("%s/protos", lower), fmt.Sprintf("%s/server", lower),
	}
	for _, dir := range dirs {
		if err := os.Mkdir(dir, os.ModeDir); err != nil {
			return g.rollback(name, err)
		}
	}

	var templates = []struct {
		templateName string
		fileName     string
	}{
		{
			templateName: "proto.tmpl",
			fileName:     fmt.Sprintf("%s/protos/%s.proto", lower, lower),
		},
		{
			templateName: "server.go.tmpl",
			fileName:     fmt.Sprintf("%s/server/server.go", lower),
		},
		{
			templateName: "notifier.go.tmpl",
			fileName:     fmt.Sprintf("%s/server/notifier.go", lower),
		},
		{
			templateName: "player.go.tmpl",
			fileName:     fmt.Sprintf("%s/server/player.go", lower),
		},
		{
			templateName: "main.go.tmpl",
			fileName:     fmt.Sprintf("%s/main.go", lower),
		},
	}

	var data = struct {
		Name      string
		NameLower string
	}{
		Name:      name,
		NameLower: lower,
	}

	for _, t := range templates {
		if err := g.writeFile(t.templateName, t.fileName, data); err != nil {
			return g.rollback(name, err)
		}
	}

	return nil
}

func (g *generator) writeFile(templateName, fileName string, data interface{}) error {
	text, err := g.loadTemplate(fmt.Sprintf("templates/%s", templateName))
	if err != nil {
		return err
	}

	t, err := template.New("deadpool").Parse(text)
	if err != nil {
		return err
	}

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	return t.Execute(f, data)
}

func (g *generator) loadTemplate(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (g *generator) rollback(name string, err error) error {
	if err != nil {
		os.RemoveAll(strings.ToLower(name))
	}
	return err
}
