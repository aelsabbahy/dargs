package run

import (
	"io/ioutil"
	"os"
	"time"

	yaml "gopkg.in/yaml.v2"

	cache "github.com/patrickmn/go-cache"
)

func loadCache(f string) (*cache.Cache, error) {
	if _, err := os.Stat(f); os.IsNotExist(err) {
		return cache.New(1*time.Nanosecond, 10*time.Minute), nil
	}
	d, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	var out map[string]cache.Item
	err = yaml.Unmarshal(d, &out)
	c := cache.NewFrom(5*time.Minute, 10*time.Minute, out)
	return c, nil
}

func saveCache(c *cache.Cache, f string) error {
	d, err := yaml.Marshal(c.Items())
	if err != nil {
		return err
	}
	return ioutil.WriteFile(f, d, 0644)
}
