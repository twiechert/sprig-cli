package main

import (
	"flag"
	"github.com/Masterminds/sprig"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

var (
	datafileFlag = flag.String("data", "", "Datafile")
	tmplFlag     = flag.String("tmpl", "", "Template")

	tmpl []byte
	data map[interface{}]interface{}
)

func init() {

	flag.Parse()

	stat, _ := os.Stdin.Stat()
	if (stat.Mode()&os.ModeNamedPipe == 0) && *tmplFlag == "" {
		panic("No template ")
	}

	var (
		source *os.File
		err    error
	)

	if *tmplFlag == "-" {
		source = os.Stdin
	} else {
		source, err = os.Open(*tmplFlag)
		if err != nil {
			panic(err)
		}
	}
	defer source.Close()

	tmpl, err = ioutil.ReadAll(source)

	if *datafileFlag != "" {

		dataFiles := strings.Split(*datafileFlag, ",")

		for _, dataFile := range dataFiles {
			var dataInner map[interface{}]interface{}

			dataBytes, err := ioutil.ReadFile(dataFile)
			if err != nil {
				panic(err)
			}

			err = yaml.Unmarshal(dataBytes, &dataInner)
			if err != nil {
				panic(err)
			}
			mergo.Merge(&data, dataInner)
		}

	}

}

func main() {
	t := template.Must(template.New(*tmplFlag).Funcs(sprig.TxtFuncMap()).Parse(string(tmpl)))
	if err := t.Execute(os.Stdout, data); err != nil {
		panic(err)
	}
}
