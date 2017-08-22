package run

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Imports           []string
	Transformers      []*MatchRunner
	Completers        []*MatchRunner
	Commands          []*Command
	CommandLookup     map[string]*Command
	TransformerLookup map[string]*MatchRunner
	CompleterLookup   map[string]*MatchRunner
	Shell             string
}

type Command struct {
	Name         string
	Wrapper      string
	Transformers []string
	FzfComplete  bool `yaml:"fzf-complete"`
	Completers   []string
}

func AutoLoadConfig() (*Config, error) {
	// Find home directory.
	root := viper.GetString("root")
	return LoadConfig(path.Join(root, ".dargs.yml"))
}

func mergeConfig(c, nc Config) Config {
	c.Commands = append(c.Commands, nc.Commands...)
	c.Transformers = append(c.Transformers, nc.Transformers...)
	c.Completers = append(c.Completers, nc.Completers...)
	return c
}

func LoadConfig(f string) (*Config, error) {
	var config Config
	files, err := getConfig(f)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		file, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(file, &config)
		if err != nil {
			return nil, err
		}
		for _, i := range config.Imports {
			nc, err := LoadConfig(i)
			if err != nil {
				return nil, err
			}
			config = mergeConfig(config, *nc)
		}
		config.CommandLookup = make(map[string]*Command)
		for _, c := range config.Commands {
			config.CommandLookup[c.Name] = c
			config.CommandLookup[path.Base(c.Name)] = c
		}
		config.TransformerLookup = make(map[string]*MatchRunner)
		for _, f := range config.Transformers {
			config.TransformerLookup[f.Name] = f
		}
		config.CompleterLookup = make(map[string]*MatchRunner)
		for _, f := range config.Completers {
			config.CompleterLookup[f.Name] = f
		}
		if config.Shell != "" {
			p, err := exec.LookPath(config.Shell)
			if err != nil {
				return nil, err
			}
			config.Shell = p
		} else {
			if p, err := exec.LookPath("bash"); err == nil {
				config.Shell = p
			} else if p, err := exec.LookPath("sh"); err == nil {
				config.Shell = p
			} else {
				config.Shell = os.Getenv("SHELL")
			}
		}
	}
	return &config, nil
}

func (c *Config) getCompleters(a ...string) ([]*MatchRunner, error) {
	var ret []*MatchRunner
	for _, n := range a {
		cm := c.CompleterLookup[n]
		if cm == nil {
			return nil, fmt.Errorf("Could not find completer %s", n)
		}
		ret = append(ret, cm)
	}
	return ret, nil
}

func (c *Config) getTransformers(a ...string) ([]*MatchRunner, error) {
	var ret []*MatchRunner
	for _, n := range a {
		cm := c.TransformerLookup[n]
		if cm == nil {
			return nil, fmt.Errorf("Could not find transformer %s", n)
		}
		ret = append(ret, cm)
	}
	return ret, nil
}

func getConfig(f string) ([]string, error) {
	if !strings.HasPrefix(f, "http") {
		return filepath.Glob(f)
	}
	dlRoot := viper.GetString("downloadRoot")
	if err := os.MkdirAll(dlRoot, 0755); err != nil {
		return nil, err
	}
	dlPath := path.Join(dlRoot, strings.Replace(f, "/", "_", -1))
	if _, err := os.Stat(dlPath); !os.IsNotExist(err) {
		// File already downloaded
		return []string{dlPath}, nil
	}
	if err := downloadFile(dlPath, f); err != nil {
		return nil, err
	}
	return []string{dlPath}, nil
}

func downloadFile(filepath string, url string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 {
		return fmt.Errorf("Bad URL response status: %d: %s", resp.StatusCode, url)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
