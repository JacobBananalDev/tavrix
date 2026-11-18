package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"go.yaml.in/yaml/v3"
)

const Filename = "tavrix.yaml"

const Starter = `version: "1"

tasks:
  hello:
    description: Print a greeting
    command: echo Hello from Tavrix!
`

var taskName = regexp.MustCompile(`^[A-Za-z0-9_][A-Za-z0-9_.-]*$`)

type text string

func (s *text) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.ScalarNode || node.ShortTag() != "!!str" {
		return fmt.Errorf("line %d: expected a string", node.Line)
	}
	*s = text(node.Value)
	return nil
}

type Task struct {
	Description text `yaml:"description"`
	Command     text `yaml:"command"`
}

type Config struct {
	Version text            `yaml:"version"`
	Tasks   map[string]Task `yaml:"tasks"`
}

func Load(path string) (Config, error) {
	f, err := os.Open(path) // #nosec G304 -- CLI loads the selected project's configuration file.
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	cfg, err := parse(f)
	if err != nil {
		return Config{}, fmt.Errorf("%s: %w", path, err)
	}
	return cfg, nil
}

func parse(r io.Reader) (Config, error) {
	var cfg Config
	dec := yaml.NewDecoder(r)
	dec.KnownFields(true)
	if err := dec.Decode(&cfg); err != nil {
		if errors.Is(err, io.EOF) {
			return Config{}, errors.New("configuration is empty")
		}
		return Config{}, err
	}
	var extra yaml.Node
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		if err != nil {
			return Config{}, err
		}
		return Config{}, errors.New("expected one YAML document")
	}
	switch cfg.Version {
	case "":
		return Config{}, errors.New(`missing configuration version; set version: "1"`)
	case "1":
	default:
		return Config{}, fmt.Errorf("unsupported configuration version %q (supported: \"1\")", cfg.Version)
	}
	if cfg.Tasks == nil {
		return Config{}, errors.New("tasks must be a mapping (use tasks: {} for no tasks)")
	}
	for _, name := range cfg.Names() {
		task := cfg.Tasks[name]
		if !taskName.MatchString(name) {
			return Config{}, fmt.Errorf("invalid task name %q: use letters, digits, underscores, dots or hyphens; start with a letter, digit or underscore", name)
		}
		if strings.TrimSpace(string(task.Command)) == "" {
			return Config{}, fmt.Errorf("task %q: command is required", name)
		}
		if strings.ContainsRune(string(task.Command), 0) {
			return Config{}, fmt.Errorf("task %q: command must not contain NUL", name)
		}
		if strings.ContainsFunc(string(task.Description), unicode.IsControl) {
			return Config{}, fmt.Errorf("task %q: description must be a single line without control characters", name)
		}
	}
	return cfg, nil
}

func (c Config) Names() []string {
	names := make([]string, 0, len(c.Tasks))
	for name := range c.Tasks {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
