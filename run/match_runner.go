package run

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Runner interface {
	Run(string, map[string]string) (string, error)
	Key() string
}

type MatchRunner struct {
	Match     string `yaml:"match"`
	PrevMatch string `yaml:"prev-match"`
	Name      string `yaml:"name"`
	Cache     int    `yaml:"cache"`
	re        *regexp.Regexp
	prevre    *regexp.Regexp
	loaded    bool
	r         Runner
}

func (m *MatchRunner) Key(a ...string) string {
	a = append(a, []string{m.Match, m.PrevMatch, m.r.Key()}...)
	data := fmt.Sprint(strings.Join(a, "-"))
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

func (m *MatchRunner) setup() error {
	if m.loaded {
		return nil
	}
	re, err := regexp.Compile(m.Match)
	if err != nil {
		return err
	}
	m.re = re
	re, err = regexp.Compile(m.PrevMatch)
	if err != nil {
		return err
	}
	m.prevre = re
	m.loaded = true
	return nil

}

func (m *MatchRunner) Run(args []string, i int, match map[string]string) (string, error) {
	cur := args[i]
	return m.r.Run(cur, match)
}

func (m *MatchRunner) CheckMatch(args []string, i int) (map[string]string, error) {
	log := log.WithField("prefix", m.Name)
	if err := m.setup(); err != nil {
		return nil, err
	}
	cur := args[i]
	if m.PrevMatch != "" {
		var prev string
		if i == 0 {
			prev = ""
		} else {
			prev = args[i-1]
		}
		if !m.prevre.MatchString(prev) {
			log.Debugf("miss for prev matcher on arg '%s'", cur)
			return nil, nil
		}
	}
	match := m.re.FindStringSubmatch(cur)
	if match == nil {
		log.Debugf("miss for arg '%s'", cur)
		return nil, nil
	}
	result := make(map[string]string)
	for i, name := range m.re.SubexpNames() {
		if name != "" {
			result[name] = match[i]
		}
		result[strconv.Itoa(i)] = match[i]
	}
	log.Debugf("hit for arg '%s'", cur)
	return result, nil
}

func (m *MatchRunner) UnmarshalYAML(unmarshal func(v interface{}) error) error {
	var tmp struct {
		Match     string `yaml:"match"`
		PrevMatch string `yaml:"prev-match"`
		Name      string `yaml:"name"`
		Cache     int    `yaml:"cache"`
	}
	if err := unmarshal(&tmp); err != nil {
		return err
	}
	m.Match = tmp.Match
	m.PrevMatch = tmp.PrevMatch
	m.Name = tmp.Name
	m.Cache = tmp.Cache
	var runner *CommandRunner
	if err := unmarshal(&runner); err != nil {
		return err
	}
	m.r = runner
	return nil
}
