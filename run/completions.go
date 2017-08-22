package run

import (
	"fmt"
	"strings"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
)

var lastDuration = 10 * time.Second

func Completions(args []string) error {
	cfg, err := AutoLoadConfig()
	if err != nil {
		return err
	}
	cn := args[0]
	args = args[1:]
	// Hack to work with spaces
	args[0] = strings.Replace(args[0], `\ `, " ", 1)
	args[1] = strings.Replace(args[1], `\ `, " ", 1)
	var cmp_names []string
	cmd := cfg.CommandLookup[cn]
	completersArg := viper.GetStringSlice("completers")
	if len(completersArg) != 0 {
		cmp_names = completersArg
	} else {
		if cmd == nil {
			return fmt.Errorf("Could not find command %s in config and no '--completers' provided", cn)
		}
		cmp_names = cmd.Completers
	}
	completers, err := cfg.getCompleters(cmp_names...)
	if err != nil {
		return err

	}
	cacheFile := viper.GetString("cacheFile")
	var completions []string
	for _, f := range completers {
		m, err := f.CheckMatch(args, 1)
		if err != nil {
			return err
		}
		if m == nil {
			continue
		}
		c, err := loadCache(cacheFile)
		if err != nil {
			return err
		}
		defer saveCache(c, cacheFile)

		key := strings.Join([]string{f.Name, args[0], args[1]}, "-")
		var i interface{}
		i, ok := getLast(c, key, args[1])
		if ok {
			log.Debugf("Last cache hit")
		} else {
			key := f.Key(args[0], args[1])
			i, ok = c.Get(key)
			if ok {
				log.Debugf("Cache hit")
			} else {
				i, err = f.Run(args, 1, m)
				if err != nil {
					if _, ok := err.(*EmptyOutError); !ok {
						return err
					}
				}
				if f.Cache > 0 {
					c.Set(key, i, time.Duration(f.Cache)*time.Second)
				}

			}
		}
		cur := i.(string)
		// Remove completions that don't start with our prefix
		var curA []string
		for _, s := range strings.Split(cur, "\n") {
			//fmt.Println(va, key, cur)
			if strings.HasPrefix(s, args[1]) {
				curA = append(curA, s)
			}

		}
		cur = strings.Join(curA, "\n")
		if cur != "" {
			setLast(c, key, cur)
		}
		completions = append(completions, curA...)
		log.Debugf("Current completions: %v", completions)
	}

	if len(completions) == 0 {
		return nil
	}

	// Insert FzfComplete magic string
	if cmd.FzfComplete {
		completions = append(completions, "")
		copy(completions[1:], completions)
		completions[0] = "dargs_fzf\n"
	}
	// remove duplicates
	fmt.Println(strings.Join(completions, "\n"))
	return nil
}

func getLast(c *cache.Cache, key, cur string) (interface{}, bool) {
	k, ok := c.Get("last-key")
	if !ok {
		return nil, false
	}
	if !strings.HasPrefix(key, k.(string)) {
		return nil, false
	}
	v, ok := c.Get("last-" + k.(string))
	if !ok {
		return nil, false
	}
	return v, true
}
func setLast(c *cache.Cache, key, cur string) {
	c.Set("last-key", key, lastDuration)
	c.Set("last-"+key, cur, lastDuration)
}
