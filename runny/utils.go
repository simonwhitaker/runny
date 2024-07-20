package runny

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func commandStringToSingleLine(command string) string {
	command = strings.TrimSpace(command)
	lines := strings.Split(command, "\n")
	trimmedLines := []string{}
	for _, line := range lines {
		trimmedLines = append(trimmedLines, strings.TrimSpace(line))
	}
	return strings.Join(trimmedLines, "; ")
}

func readConfig() (Config, error) {
	// Read .runny.yaml from the current directory
	var conf Config
	yamlFile, err := os.ReadFile(".runny.yaml")
	if err != nil {
		return conf, err
	}

	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		return conf, err
	}
	return conf, nil
}
