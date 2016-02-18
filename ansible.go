package ansible

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"text/template"
)

var configTemplate *template.Template

const (
	configFilename string = "ansible.cfg"
)

func init() {

	var err error
	if configTemplate, err = template.New(configFilename).Parse(`[defaults]

jinja2_extensions = jinja2.ext.loopcontrols
inventory      = {{ .HostFile }}

# uncomment this to disable SSH key host checking
host_key_checking = False

# SSH timeout
timeout = 10

# default module name for /usr/bin/ansible
module_name = setup

# if set to a persistent type (not 'memory', for example 'redis') fact values
# from previous runs in Ansible will be stored.  This may be useful when
# wanting to use, for example, IP information from one group of servers
# without having to talk to them in the same playbook run to get their
# current IP information.
fact_caching = memory


# retry files
retry_files_enabled = False

[ssh_connection]
ssh_args = -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -i {{ .SshKeyPath }}`); err != nil {
		panic(err)
	}
}

type AnsibleConfig struct {
	SshKeyPath string
	HostFile   string
}

func NewDefaultConfig() *AnsibleConfig {

	return &AnsibleConfig{
		SshKeyPath: "~/.ssh/id_rsa",
		HostFile:   hostsFilename,
	}
}

type Ansible struct {
	playbook *Playbook
	hosts    Hosts
	config   *AnsibleConfig
	Workdir  string
}

func NewAnsible(playbook *Playbook, hosts Hosts) (*Ansible, error) {

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}

	playbook.SetVar("ansible_workdir", dir)

	return &Ansible{
		playbook: playbook,
		hosts:    hosts,
		Workdir:  dir,
		config:   NewDefaultConfig(),
	}, nil
}

func (a *Ansible) UseConfig(cfg *AnsibleConfig) {
	a.config = cfg
}

func (a *Ansible) UseKey(path string) {
	a.config.SshKeyPath = path
}

func (a *Ansible) Write() error {

	if err := a.playbook.write(a.Workdir); err != nil {
		return err
	}

	if err := a.hosts.write(a.Workdir); err != nil {
		return err
	}

	// Write ansible config
	file, err := os.OpenFile(path.Join(a.Workdir, configFilename), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := configTemplate.Execute(file, a.config); err != nil {
		return err
	}

	return nil
}

func (a *Ansible) Clean() error {
	return os.RemoveAll(a.Workdir)
}

func (a *Ansible) Run() ([]byte, error) {

	command := fmt.Sprintf("ansible-playbook -i %s %s", hostsFilename, playbookFilename)
	cmd := exec.Command("sh", "-c", fmt.Sprintf("cd %s && %s", a.Workdir, command))

	return cmd.CombinedOutput()
}
