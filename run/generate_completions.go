package run

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"

	"github.com/aelsabbahy/dargs/data"
	"github.com/spf13/viper"
)

func GenerateCompletions(args []string) error {
	outDir := viper.GetString("completions-out")
	if outDir == "" {
		outDir = path.Join(viper.GetString("root"), ".dargs", "completions", viper.GetString("shell"))
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}
	cfg, err := AutoLoadConfig()
	if err != nil {
		return err
	}
	data, err := data.Asset(fmt.Sprintf("completion_%s.tmpl", viper.GetString("shell")))
	if err != nil {
		return err
	}
	t := template.New("test")
	tmpl, err := t.Parse(string(data))
	if err != nil {
		return err
	}
	tmpl.Option("missingkey=error")

	for _, cmd := range cfg.Commands {
		outName := path.Join(outDir, fmt.Sprintf("zzdargs_%s", cmd.Wrapper))
		if _, err := os.Stat(outName); !os.IsNotExist(err) && !viper.GetBool("force-completions") {
			return fmt.Errorf("File already exists: %s", outName)
		}
		if _, err := os.Stat(outDir); os.IsNotExist(err) {
			return fmt.Errorf("Directory does not exists: %s", outDir)
		}

		var tVars = struct {
			CmdFull string
			Cmd     string
		}{cmd.Name, path.Base(cmd.Name)}
		var doc bytes.Buffer
		err = tmpl.Execute(&doc, tVars)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(outName, doc.Bytes(), 0755)
		if err != nil {
			return err
		}
		log.Infof("Wrote %s", outName)
	}
	return nil
}
