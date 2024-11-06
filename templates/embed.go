package templates

import _ "embed"

//go:embed t.jsonc
var JsoncTmpl string

//go:embed t.yaml
var YamlTmpl string
