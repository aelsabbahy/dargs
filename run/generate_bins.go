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

func GenerateBins(args []string) error {
	outDir := viper.GetString("bin-out")
	if outDir == "" {
		return fmt.Errorf("Please specifiy output directory using the -o flag")
	}
	cfg, err := AutoLoadConfig()
	if err != nil {
		return err
	}
	data, err := data.Asset("bin.tmpl")
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
		outName := path.Join(outDir, cmd.Wrapper)
		if _, err := os.Stat(outName); !os.IsNotExist(err) && !viper.GetBool("force-bins") {
			return fmt.Errorf("File already exists: %s", outName)
		}
		if _, err := os.Stat(outDir); os.IsNotExist(err) {
			return fmt.Errorf("Directory does not exists: %s", outDir)
		}

		var tVars = struct {
			Shell   string
			CmdFull string
		}{cfg.Shell, cmd.Name}
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
