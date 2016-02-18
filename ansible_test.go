package ansible

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestAnsible(t *testing.T) {

	roleTesting := NewRole("testing").AddTask(Task{
		"name":    "first task",
		"command": "/bin/echo 'Hello World!'",
	}).AddTask(Task{

		"name":    "second task",
		"command": "/bin/echo 'Foo Bar'",
	})

	roleDebug := NewRole("debug").AddTask(Task{
		"name":    "debug first task",
		"command": "/bin/echo 'debug Hello World!'",
	}).AddTask(Task{
		"name":    "debug second task",
		"command": "/bin/echo 'debug Foo Bar'",
	}).AddTemplate(NewTemplate("binary.service", []byte(`hello world
		123 345
		i'm a template {{ value }}`)))

	playbook := NewPlaybook(roleTesting, roleDebug).SetVar("hello", "world").SetVar("bool", true).SetVar("int", 42)

	hosts := NewHosts().AddHost(NewHost("localhost", "127.0.0.1", 22, "root")).AddHost(NewHost("staging", "1.2.3.4", 22, "admin"))

	a, err := NewAnsible(playbook, hosts)
	if err != nil {
		t.Fatal(err)
	}
	defer a.Clean()

	if err := a.Write(); err != nil {
		t.Fatal(err)
	}

	// Show output
	if err := filepath.Walk(a.Workdir, func(path string, info os.FileInfo, err error) error {
		fmt.Println(path)
		if !info.IsDir() {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}

			fmt.Println(string(data))
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}

}
