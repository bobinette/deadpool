package main

import (
	"io/ioutil"
	"log"
)

type Lister interface {
	List() []string
}

type lister struct{}

func NewLister() Lister {
	return &lister{}
}

func (l *lister) List() []string {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	var games []string
	for _, file := range files {
		if file.Mode().IsDir() && l.isGameDir(file.Name()) {
			games = append(games, file.Name())
		}
	}

	return games
}

func (l *lister) isGameDir(dir string) bool {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.Name() == "server" && file.Mode().IsDir() {
			return true
		}
	}
	return false
}
