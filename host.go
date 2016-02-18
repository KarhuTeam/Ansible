package ansible

import (
	"os"
	"path"
	"text/template"
)

var hostsTemplate *template.Template

const (
	hostsFilename string = "hosts.ini"
)

func init() {

	var err error
	if hostsTemplate, err = template.New(hostsFilename).Parse(`[all]
{{ range . }}{{ .Name }} ansible_ssh_host={{ .Host }} ansible_ssh_port={{ .Port }} ansible_ssh_user={{ .User }}
{{ end }}`); err != nil {
		panic(err)
	}
}

type Host struct {
	Name string
	Host string
	Port int
	User string
}

type Hosts []*Host

func NewHost(name, ip string, port int, user string) *Host {
	return &Host{
		Name: name,
		Host: ip,
		Port: port,
		User: user,
	}
}

func NewHosts() Hosts {
	return make(Hosts, 0)
}

func (hs Hosts) AddHost(h *Host) Hosts {
	hs = append(hs, h)
	return hs
}

func (hs Hosts) write(rootPath string) error {

	// Write hosts
	file, err := os.OpenFile(path.Join(rootPath, hostsFilename), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := hostsTemplate.Execute(file, hs); err != nil {
		return err
	}

	return nil
}
