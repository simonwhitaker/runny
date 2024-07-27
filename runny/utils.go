package runny

import (
	"fmt"
	"os"
	"strings"

	"github.com/dominikbraun/graph"
	"gopkg.in/yaml.v3"
)

func commandStringToSingleLine(command string, maxlength int) string {
	command = strings.TrimSpace(command)
	lines := strings.Split(command, "\n")
	trimmedLines := []string{}
	for _, line := range lines {
		trimmedLines = append(trimmedLines, strings.TrimSpace(line))
	}
	result := strings.Join(trimmedLines, "; ")
	if maxlength > 0 && len(result) > maxlength {
		result = result[:maxlength-1] + "…"
	}
	return result
}

func readConfig(path string) (Config, error) {
	// Read .runny.yaml from the current directory
	var conf Config
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return conf, err
	}

	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal") {
			return conf, fmt.Errorf("invalid runny config file: %s", path)
		}
		return conf, err
	}

	return conf, conf.validate()
}

func (c *Config) validate() error {
	hash := func(c CommandName) string {
		return string(c)
	}
	g := graph.New(hash, graph.Directed(), graph.PreventCycles())

	for cmdName := range c.Commands {
		err := g.AddVertex(cmdName)
		if err != nil {
			return fmt.Errorf("error declaring command %s: %v", cmdName, err)
		}
	}

	for cmdName, cmd := range c.Commands {
		for _, needsName := range cmd.Needs {
			err := g.AddEdge(hash(cmdName), hash(needsName))
			if err != nil {
				return fmt.Errorf("error declaring %s as dependency of %s: %v", needsName, cmdName, err)
			}
		}
	}
	return nil
}
