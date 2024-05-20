package main

import (
	"bufio"
	"errors"
	yaml "gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"regexp"
)

// CommentConfig is the configuration to sort all comments.
type CommentConfig struct {
	// Local is the root Go module name. All subpackages of this module
	// will be separated from the external packages.
	Local string `yaml:"local"`

	// Prefixes is a list of relative Go packages from the root package.
	// All comments with these prefixes will be separated from each other.
	Prefixes []string `yaml:"prefixes"`

	// Signature is the personal signature affixed at the end of the comment in order to allow users and
	// other iterations to find the autogenerated comments and to be able to evolve
	// them with future updates of the script.
	// If empty, the automatically added prefix is "auto".
	Signature string `yaml:"signature"`
	// Allows you to know if you update the tagged comments each time the script is executed.
	UpdateComments bool `yaml:"update-comments"`
	// do we use OPENAI to generate function comments
	OpenAIActive bool   `yaml:"openai-active"`
	OpenAIAPIKey string `yaml:"openai-api_key"`
	OpenAIURL    string `yaml:"openai-url"`
}

// Merge merges the given CommentConfig with this configure and return
// a new CommentConfig with the merged result.
// - Local attribute value is overriden.
// - Prefixes attribute values are appended.
func (cfg *CommentConfig) Merge(newCfg *CommentConfig) *CommentConfig {
	local := cfg.Local
	if newCfg.Local != "" {
		local = newCfg.Local
	}

	return &CommentConfig{
		Local:          local,
		OpenAIActive:   false,
		OpenAIAPIKey:   "",
		OpenAIURL:      "https://api.openai.com/v1/engines/davinci-codex/completions",
		Prefixes:       append(cfg.Prefixes, newCfg.Prefixes...),
		Signature:      "AutoComBOT",
		UpdateComments: false,
	}
}

// CommentConfigCache is a cache to contains the configuration for all processed files.
type CommentConfigCache struct {
	rootConfig CommentConfig
	configs    map[string]*CommentConfig
}

// NewCommentConfigCache instantiates a new cache to store the configuration for all processed files.
func NewCommentConfigCache(local string, prefixes []string) *CommentConfigCache {
	return &CommentConfigCache{
		rootConfig: CommentConfig{
			Local:    local,
			Prefixes: prefixes,
		},
		configs: make(map[string]*CommentConfig),
	}
}

// Get returns the configuration for the given processed file.
// Keep all intermediate configurations in the cache.
func (cache *CommentConfigCache) Get(filename string) (*CommentConfig, error) {
	absFilepath, _ := filepath.Abs(filename)
	dirpath := filepath.Dir(absFilepath)

	cfg, err := cache.get(dirpath)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cache *CommentConfigCache) get(dirpath string) (*CommentConfig, error) {
	cfg, ok := cache.configs[dirpath]
	if ok {
		if cfg == nil {
			panic("cfg should not be nil")
		}
		return cfg, nil
	}

	var parentCfg *CommentConfig

	goModFilepath := filepath.Join(dirpath, "go.mod")
	modname, goModFileExists, err := getModuleNameFromGoModFile(goModFilepath)
	if err != nil {
		return nil, err
	}

	if !goModFileExists && dirpath != "." && dirpath != "/" {
		parentCfg, err = cache.get(filepath.Dir(dirpath))
		if err != nil {
			return nil, err
		}
	} else {
		parentCfg = &cache.rootConfig
		if parentCfg.Local == "" {
			parentCfg.Local = modname
		}
	}

	localCfg, err := readConfigFile(filepath.Join(dirpath, ".gocomments"))
	if err != nil {
		return nil, err
	}

	cfg = parentCfg
	if localCfg != nil {
		cfg = cfg.Merge(localCfg)
	}

	if cfg == nil {
		panic("cfg should not be nil")
	}

	cache.configs[dirpath] = cfg

	return cfg, nil
}

func readConfigFile(filename string) (*CommentConfig, error) {
	f, err := os.Open(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}

		return nil, err
	}

	defer func() {
		_ = f.Close()
	}()

	var cfg CommentConfig
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

var goModModuleRegexp = regexp.MustCompile(`^module\s+(\S+)$`)

func getModuleNameFromGoModFile(goModFilepath string) (string, bool, error) {
	var modname string

	f, err := os.Open(goModFilepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", false, nil
		}

		return "", false, err
	}

	defer func() {
		_ = f.Close()
	}()

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		if m := goModModuleRegexp.FindStringSubmatch(line); m != nil {
			modname = m[1]
			break
		}
	}

	if err := s.Err(); err != nil {
		return "", true, err
	}

	return modname, true, nil
}
