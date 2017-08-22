package run

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

// Run stuff..
func Run(args []string) error {
	cfg, err := AutoLoadConfig()
	if err != nil {
		return err
	}
	cn := args[0]
	var transformerNames []string
	transformersArg := viper.GetStringSlice("transformers")
	if len(transformersArg) != 0 {
		transformerNames = transformersArg
	} else {
		cmd := cfg.CommandLookup[cn]
		if cmd == nil {
			return fmt.Errorf("Could not find command %s in config and no '--transformers' provided", cn)
		}
		transformerNames = cmd.Transformers
	}
	transformers, err := cfg.getTransformers(transformerNames...)
	if err != nil {
		return err

	}
	cacheFile := viper.GetString("cacheFile")
	c, err := loadCache(cacheFile)
	if err != nil {
		return err
	}
	defer saveCache(c, cacheFile)

	args = args[1:]
	var nargs []string
	for _, f := range transformers {
		nargs = nil
		for i, _ := range args {
			m, err := f.CheckMatch(args, i)
			if err != nil {
				return err
			}
			if m == nil {
				nargs = append(nargs, args[i])
				continue
			}
			var v string
			key := f.Key(args[i])
			vi, ok := c.Get(key)
			if ok {
				log.Debugf("Cache hit")
				v = vi.(string)
			} else {
				v, err = f.Run(args, i, m)
				if err != nil {
					return err
				}
				if f.Cache > 0 {
					c.Set(key, v, time.Duration(f.Cache)*time.Second)
				}
			}
			for _, nv := range strings.Split(v, "\n") {
				nargs = append(nargs, nv)
			}

		}
		args = make([]string, len(nargs))
		copy(args, nargs)
		log.Debugf("Current args: %v", args)
	}

	// Insert command at the front
	args = append(args, "")
	copy(args[1:], args)
	args[0] = cn
	path := cn
	// expand to full path
	if !strings.HasPrefix(path, "/") {
		if path, err = exec.LookPath(path); err != nil {
			return err
		}
	}
	noop := viper.GetBool("noop")
	verbose := viper.GetBool("verbose")
	if noop || verbose {
		fmt.Println(strings.Join(args, " "))
		if noop {
			return nil
		}
	}

	if err := saveCache(c, cacheFile); err != nil {
		return err
	}

	if err := syscall.Exec(path, args, os.Environ()); err != nil {
		return err
	}
	return nil
}
