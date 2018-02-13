package transpiler

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/alecthomas/template"
	yaml "gopkg.in/yaml.v2"
)

const loaderTplText = `// Code generated by goautocfg.
// source: {{ .Source }}
// DO NOT EDIT!

package {{ .Package }}

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/spf13/cast"
)

type AutoCfg struct {
{{range $item := .Items}}
	_{{$item.Key}} {{$item.Type}}
{{end}}
}

var cfg *AutoCfg

func init() {
	cfg = New()
}

func New() *AutoCfg {
	a := new(AutoCfg)
	return a
}

{{range $item := .Items}}
func {{$item.KEY}}() {{$item.Type}} {
	return cfg._{{$item.Key}}
}
func (cfg *AutoCfg) {{$item.KEY}}() {{$item.Type}} {
	return cfg._{{$item.Key}}
}
{{end}}

func Load(name string, paths ...string) {
	cfg.Load(name, paths...)
}
func (cfg *AutoCfg) Load(name string, paths ...string) {
	vpr := viper.New()
	vpr.SetConfigName(name)
	for _, path := range paths {
		vpr.AddConfigPath(path)
	}
	err := vpr.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
{{range .Items}}
	if v := vpr.Get("{{.KEY}}"); v != nil {
		cfg._{{.Key}} = cast.{{.Function}}(v)
	} else {
		{{if .Optional}}cfg._{{.Key}} = cast.{{.Function}}("{{.Default}}") {{end}}
		{{if eq .Optional false}}panic(fmt.Errorf("Required config missing: '{{.KEY}}'")) {{end}}
	}
{{end}}
}
`

type loader struct {
	Items           []item
	Source, Package string
}
type item struct {
	Optional                 bool
	Default                  string
	KEY, Key, Type, Function string
}

type schema map[string]string

var tokenRegex = regexp.MustCompile(`(?:(?:\\.)|(?:[^,]))*`)
var modifierRegex = regexp.MustCompile(`@(\w+)\((.*)\)$`)

func Transpile(sourceFile string, sourceContent []byte, targetPackage string, targetWriter io.Writer) error {
	schemaENV := make(schema)
	if err := yaml.Unmarshal(sourceContent, &schemaENV); err != nil {
		return err
	}

	loaderTpl, err := template.New("setup").Parse(loaderTplText)
	if err != nil {
		return err
	}

	items, err := getTplItems(schemaENV)
	if err != nil {
		return err
	}

	err = loaderTpl.Execute(targetWriter, &loader{
		Items:   items,
		Source:  sourceFile,
		Package: targetPackage,
	})
	if err != nil {
		return err
	}
	return nil
}

func getTplItems(schemaENV schema) ([]item, error) {
	items := make([]item, 0, len(schemaENV))
	for envKey, modifiersS := range schemaENV {
		var i item
		i.KEY = strings.ToUpper(envKey)
		i.Key = strings.ToLower(envKey)

		modifiers := tokenRegex.FindAllStringSubmatch(modifiersS, -1)
		for _, modifier := range modifiers {
			err := updateModifier(&i, modifier[0])
			if err != nil {
				return nil, err
			}
		}

		items = append(items, i)
	}
	return items, nil
}

func updateModifier(i *item, token string) error {
	switch token {
	case "string":
		i.Type, i.Function = "string", "ToString"
	case "integer":
		i.Type, i.Function = "int", "ToInt"
	case "bool":
		i.Type, i.Function = "bool", "ToBool"
	case "float":
		i.Type, i.Function = "float64", "ToFloat64"
	case "strings":
		i.Type, i.Function = "[]string", "ToStringSlice"
	default:
		if modifier := modifierRegex.FindStringSubmatch(token); modifier != nil {
			switch modifier[1] {
			case "optional":
				i.Optional = true
				i.Default = modifier[2]
			default:
				return invalidToken(token)
			}
		} else {
			return invalidToken(token)
		}
	}
	return nil
}

func invalidToken(token string) error {
	return fmt.Errorf("parsing error: Invalid token '%s'. "+
		"Valid tokens: string integer bool float strings @optional", token)
}