package ansible

import (
	"os"
	"path"
	"text/template"
)

var playbookTemplate *template.Template

const (
	playbookFilename string = "playbook.yml"
)

func init() {

	var err error
	if playbookTemplate, err = template.New(playbookFilename).Parse(`---

- hosts: {{ .Hosts }}
  sudo: {{ .Sudo }}
  {{ if len .Vars }}vars:{{ range $key, $value := .Vars}}
    {{ $key }}: {{ $value }}{{ end }}
  {{ end }}
  roles:{{ range .Roles }}
  - {{ .Name }}{{ end }}
`); err != nil {
		panic(err)
	}
}

// Playbook definition
type Playbook struct {
	Hosts string
	Sudo  string
	Roles Roles
	Vars  PlaybookVars
}

type PlaybookVars map[string]interface{}

func (pv PlaybookVars) Set(key string, value interface{}) {
	pv[key] = value
}

const (
	rolesDir string = "roles"
)

var playbookSubDirs = []string{rolesDir}

// Create a new Playbook with roles
func NewPlaybook(roles ...*Role) *Playbook {

	return &Playbook{
		Hosts: "all",
		Sudo:  "yes",
		Roles: roles,
		Vars:  make(PlaybookVars),
	}
}

func (p *Playbook) AddRole(r *Role) *Playbook {

	p.Roles = append(p.Roles, r)
	return p
}

// Write Playbook and dependents roles into rootPath
func (p *Playbook) write(rootPath string) error {

	// Create playbook subdirectories
	for _, dir := range playbookSubDirs {

		if err := os.MkdirAll(path.Join(rootPath, dir), 0755); err != nil {
			return err
		}
	}

	// Write all roles
	for _, role := range p.Roles {

		if err := role.write(path.Join(rootPath, rolesDir)); err != nil {
			return err
		}
	}

	// Write playbook himself
	file, err := os.OpenFile(path.Join(rootPath, playbookFilename), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := playbookTemplate.Execute(file, p); err != nil {
		return err
	}

	return nil
}

func (p *Playbook) SetVar(key string, value interface{}) *Playbook {
	p.Vars.Set(key, value)
	return p
}
