package ansible

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
)

// Roles directories names
const (
	tasksDir     string = "tasks"
	handlersDir         = "handlers"
	templatesDir        = "templates"
	filesDir            = "files"
)

// Roles subdir list, main purpose is to create when on write()
var roleSubDirs = []string{tasksDir, handlersDir, templatesDir}

// Role definition
type Role struct {
	Name      string
	tasks     Tasks
	handlers  Tasks
	templates Templates
	files     Files
}

type Roles []*Role

// Create a new role
func NewRole(name string) *Role {

	return &Role{
		Name: name,
	}
}

// join a task to a role
func (r *Role) AddTask(t Task) *Role {

	if _, ok := t["name"]; !ok {
		log.Panic("missing 'name' in task")
	}
	r.tasks = append(r.tasks, t)
	return r
}

// join a handler to a role
func (r *Role) AddHandler(t Task) *Role {

	if _, ok := t["name"]; !ok {
		log.Panic("missing 'name' in task")
	}
	r.handlers = append(r.handlers, t)
	return r
}

func (r *Role) AddTemplate(t *Template) *Role {

	r.templates = append(r.templates, t)
	return r
}

func (r *Role) AddFile(f *File) *Role {

	r.files = append(r.files, f)
	return r
}

// Write role directories and files to path
// The path MUST have the /roles directory
// Ex: /my/path/roles
func (r *Role) write(rolesDir string) error {

	rolesPath := path.Join(rolesDir, r.Name)

	// Create role subdirectories
	for _, dir := range roleSubDirs {

		if err := os.MkdirAll(path.Join(rolesPath, dir), 0755); err != nil {
			return err
		}
	}

	// Marshal tasks
	if data, err := yaml.Marshal(r.tasks); err != nil {
		return err
	} else {

		// Write tasks
		if err := ioutil.WriteFile(path.Join(rolesPath, tasksDir, "main.yml"), data, 0644); err != nil {
			return err
		}
	}

	// Marshal handlers
	if data, err := yaml.Marshal(r.handlers); err != nil {
		return err
	} else {

		// Write tasks
		if err := ioutil.WriteFile(path.Join(rolesPath, handlersDir, "main.yml"), data, 0644); err != nil {
			return err
		}
	}

	// Write tempalte
	for _, t := range r.templates {

		filepath := path.Join(rolesPath, templatesDir, t.Path)
		if err := os.MkdirAll(path.Dir(filepath), 0755); err != nil {
			return err
		}

		if err := ioutil.WriteFile(filepath, t.Data, 0644); err != nil {
			return err
		}
	}

	// Write files
	for _, f := range r.files {

		filepath := path.Join(rolesPath, filesDir, f.Path)
		if err := os.MkdirAll(path.Dir(filepath), 0755); err != nil {
			return err
		}

		if err := ioutil.WriteFile(filepath, f.Data, 0644); err != nil {
			return err
		}
	}

	return nil
}
