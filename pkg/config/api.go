package config

import (
	"fmt"
	"log"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

func InitFromJson(cfgAsJson string) error {
	// Load JSON config.
	if err := k.Load(file.Provider("mock/mock.json"), json.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	// Load YAML config and merge into the previously loaded config (because we can).
	k.Load(file.Provider("mock/mock.yml"), yaml.Parser())

	k.All()

	fmt.Println("parent's name is = ", k.String("parent1.name"))
	fmt.Println("parent's ID is = ", k.Int("parent1.id"))
	return nil
}
