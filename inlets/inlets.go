package inlets

import (
	"net/http"
	"path/filepath"
)

// TODO: watch inlets.d folder and reload automatically (#61)

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path"
)

type Inlet interface {
	Name() string
	SupportedMethods() []string
	Handler(http.Handler) http.Handler
}

func LoadInlets(dir string) []Inlet {
	configPaths, err := filepath.Glob(path.Join(dir, "*.yaml"))
	if err != nil {
		log.Printf("Warning: Failed to load inlets from config: %v\n", err)
		return []Inlet{}
	}

	names := make(map[string]bool)
	inlets := make([]Inlet, 0, len(configPaths))

	for _, c := range configPaths {
		f, err := os.Open(c)
		if err != nil {
			log.Printf("Warning: Failed to load inlet '%s' from config: %v\n", c, err)
			continue
		}
		defer f.Close()

		var inletConfig InletConfig
		if err := yaml.NewDecoder(f).Decode(&inletConfig); err != nil {
			log.Printf("Warning: Faield to parse inlet config from '%s': %v\n", c, err)
			continue
		}

		if _, ok := names[inletConfig.Name]; ok {
			log.Printf("Warning: Ignoring inlet definition from '%s', because name was found twice\n")
			continue
		}
		names[inletConfig.Name] = true

		inlet, err := NewConfigInlet(&inletConfig)
		if err != nil {
			log.Printf("Warning: %v\n", err)
			continue
		}
		inlets = append(inlets, inlet)
	}

	return inlets
}
