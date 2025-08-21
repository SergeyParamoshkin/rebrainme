package app

import "fmt"

type App struct {
	Name    string
	Version string
	commit  string
}

func NewApp() {
	fmt.Println("aaaapp")
}
