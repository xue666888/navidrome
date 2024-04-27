package metadata

import (
	"strings"

	"github.com/navidrome/navidrome/conf"
	"github.com/navidrome/navidrome/log"
	"github.com/navidrome/navidrome/resources"
	"gopkg.in/yaml.v3"
)

type mapping struct {
	internalName string
	split        string
}

var mappings = map[string]mapping{}

// Load resources/metadata.yaml into memory
func Load() error {
	yml, err := resources.FS().Open("mappings.yaml")
	if err != nil {
		return err
	}
	defer yml.Close()
	data := map[string]struct {
		Aliases []string `yaml:"aliases"`
		Split   string   `yaml:"split"`
	}{}

	decoder := yaml.NewDecoder(yml)
	err = decoder.Decode(data)
	if err != nil {
		return err
	}
	for k, v := range data {
		internalName := strings.ToLower(k)
		for _, alias := range v.Aliases {
			alias = strings.ToLower(strings.TrimSpace(alias))
			if old, ok := mappings[alias]; ok {
				log.Warn("Duplicate alias in resources/mappings.yaml", "alias", alias, "internalName", internalName, "previousInternalName", old.internalName)
			}
			mappings[alias] = mapping{internalName, v.Split}
		}
	}
	return nil
}

func init() {
	conf.AddHook(func() {
		if err := Load(); err != nil {
			panic(err)
		}
	})
}
