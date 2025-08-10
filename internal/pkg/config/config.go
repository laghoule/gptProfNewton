package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	openaiEnvKey = "OPENAI_API_KEY"
)

type Config struct {
	Student Student `yaml:"eleve"`
	OpenAI  OpenAI  `yaml:"openai"`
}

type Student struct {
	Name    string `yaml:"nom"`
	Grade   int    `yaml:"niveau"`
	Details string `yaml:"details,omitempty"`
}

type OpenAI struct {
	APIKey   string `yaml:"clef_api"`
	Model    string `yaml:"modele,omitempty"`
}

func New(cfgFile string) (*Config, error) {
	c := &Config{}
	if err := c.Load(cfgFile); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) Load(cfgFile string) error {
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		return fmt.Errorf("échec de la lecture du fichier de configuration: %w", err)
	}

	if err := yaml.Unmarshal(data, c); err != nil {
		return fmt.Errorf("échec de l'analyse du fichier de configuration: %w", err)
	}

	if err := c.validate(); err != nil {
		return fmt.Errorf("échec de la validation du fichier de configuration: %w", err)
	}

	return nil
}

func (c *Config) validate() error {
	if c.Student.Name == "" {
		return fmt.Errorf("nom est requis")
	}

	if c.Student.Grade < 1 && c.Student.Grade > 12 {
		return fmt.Errorf("niveau doit être entre 1 et 12")
	}

	openaiKey, _ := os.LookupEnv(openaiEnvKey)
	if c.OpenAI.APIKey == "" && openaiKey == "" {
		return fmt.Errorf("la variable d'environnement %s ou clef_api est requise", openaiEnvKey)
	}

	return nil
}
