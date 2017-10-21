package benchmarker

import (
	"fmt"
	"io/ioutil"
	"math/rand"

	"github.com/go-yaml/yaml"
	"github.com/ryanuber/go-glob"
)

type ymlParams struct {
	Action  string
	Enctype string
	Method  string
	Data    []struct {
		Types  map[string]string
		Values []map[string][]string
	}
}

type FormSetter struct {
	params []ymlParams
}

func NewFormSetter(configPath string) (*FormSetter, error) {
	f := &FormSetter{}
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(buf, &f.params)
	return f, err
}

func (f *FormSetter) Set(form *HtmlForm) error {
	for _, val := range f.params {
		if !glob.Glob(val.Action, form.Action) || val.Method != form.Method || val.Enctype != form.EncType {
			continue
		}
		for _, chunk := range val.Data {
			for name, values := range chunk.Values[rand.Intn(len(chunk.Values))] {
				for _, v := range values {
					if !form.ExistKey(name) {
						fmt.Printf("not exist key: %s\n", name)
						continue
					}
					switch chunk.Types[name] {
					case "file":
						form.AddParam(name, FileParam(v))
					default:
						form.AddParam(name, TextParam(v))
					}
				}
			}
		}
	}
	return nil
}
