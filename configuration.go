package main

import (
	"fmt"
	toml "github.com/pelletier/go-toml/v2"
)

type MyConfig struct {
	Version int
	Name    string
	Tags    []string
}

func config() {
	doc := `
	version = 2
	name = "go-toml"
	tags = ["go", "toml"]
	`

	var cfg MyConfig
	err := toml.Unmarshal([]byte(doc), &cfg)
	if err != nil {
		panic(err)
	}
	fmt.Println("version:", cfg.Version)
	fmt.Println("name:", cfg.Name)
	fmt.Println("tags:", cfg.Tags)

}
