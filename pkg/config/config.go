package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"good.site/pkg/utils/convert"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Build struct {
		Public   string `yaml:"public" json:"public"`
		Nonsense string `yaml:"nonsense" json:"nonsense"`
		Markdown string `yaml:"markdown" json:"markdown"`
	} `yaml:"build" json:"build"`
	Site struct {
		Name        string `yaml:"name" json:"name"`
		Description string `yaml:"description" json:"description"`
		Url         string `yaml:"url" json:"url"`
		Author      string `yaml:"author" json:"author"`
		Language    string `yaml:"language" json:"language"`
		Copyright   string `yaml:"copyright" json:"copyright"`
		Icons       struct {
			Sm string `yaml:"sm" json:"sm"`
			Md string `yaml:"md" json:"md"`
			Lg string `yaml:"lg" json:"lg"`
		} `yaml:"icons" json:"icons"`
		Logo2x1 string `yaml:"logo2x1" json:"logo2x1"`
	} `yaml:"site" json:"site"`
}

var cfg Config

func GetConfig() *Config {
	return &cfg
}

func init() {
	readYaml(&cfg)
	readJson(&cfg)
	readEnv(&cfg)
}

func processError(err error) {
	slog.Debug(fmt.Sprintln(err))
}

func readYaml(cfg *Config) {
	f, err := os.Open("config.yml")
	if err != nil {
		processError(err)
		return
	}

	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		processError(err)
		return
	}
}

func readJson(cfg *Config) {
	f, err := os.Open("config.json")
	if err != nil {
		processError(err)
		return
	}

	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		processError(err)
		return
	}
}

func readEnv(cfg *Config) {
	err := godotenv.Load()

	if err != nil {
		processError(err)
		return
	}

	cfg.Build.Public = convert.DefaultTo(os.Getenv("PUBLIC_PATH"), cfg.Build.Public).(string)
	cfg.Build.Nonsense = convert.DefaultTo(os.Getenv("NONSENSE_PATH"), cfg.Build.Nonsense).(string)
	cfg.Build.Markdown = convert.DefaultTo(os.Getenv("MARKDOWN_PATH"), cfg.Build.Markdown).(string)
	cfg.Site.Name = convert.DefaultTo(os.Getenv("SITE_NAME"), cfg.Site.Name).(string)
	cfg.Site.Description = convert.DefaultTo(os.Getenv("SITE_DESCRIPTION"), cfg.Site.Description).(string)
	cfg.Site.Url = convert.DefaultTo(os.Getenv("SITE_URL"), cfg.Site.Url).(string)
	cfg.Site.Author = convert.DefaultTo(os.Getenv("SITE_AUTHOR"), cfg.Site.Author).(string)
	cfg.Site.Language = convert.DefaultTo(os.Getenv("SITE_LANGUAGE"), cfg.Site.Language).(string)
	cfg.Site.Copyright = convert.DefaultTo(os.Getenv("SITE_COPYRIGHT"), cfg.Site.Copyright).(string)
	cfg.Site.Icons.Sm = convert.DefaultTo(os.Getenv("SITE_ICON_SM"), cfg.Site.Icons.Sm).(string)
	cfg.Site.Icons.Md = convert.DefaultTo(os.Getenv("SITE_ICON_MD"), cfg.Site.Icons.Md).(string)
	cfg.Site.Icons.Lg = convert.DefaultTo(os.Getenv("SITE_ICON_LG"), cfg.Site.Icons.Lg).(string)
	cfg.Site.Logo2x1 = convert.DefaultTo(os.Getenv("SITE_LOGO_2X1"), cfg.Site.Logo2x1).(string)
}
