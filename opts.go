package main

import (
	"os"
	"github.com/jessevdk/go-flags"
	"fmt"
	"strings"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"github.com/kelseyhightower/envconfig"
)

const dockerFlowPath  = "docker-flow.yml"

type Opts struct {
	Host				string 		`short:"H" long:"host" description:"Docker host"`
	ComposePath 		string 		`short:"f" long:"compose-path" description:"Docker Compose configuration file" yaml:"compose_path" envconfig:"compose_path"`
	WebServer   		bool 		`short:"w" long:"web-server" description:"Whether a Web server should be started" yaml:"web_server" envconfig:"web_server"`
	BlueGreen			bool 		`short:"b" long:"blue-green" description:"Whether to perform blue-green desployment" yaml:"blue_green" envconfig:"blue_green"`
	Target				string 		`short:"t" long:"target" description:"Docker Compose target that will be deployed"`
	SideTargets			[]string 	`short:"T" long:"side-targets" description:"Side or auxiliary targets that will be deployed. Multiple values are allowed." yaml:"side_targets" envconfig:"side_targets"`
	SkipPullTarget		bool		`short:"P" long:"skip-pull-targets" description:"Whether to skip pulling target." yaml:"skip_pull_target" envconfig:"skip_pull_target"`
	PullSideTargets		bool		`short:"S" long:"pull-side-targets" description:"Whether side or auxiliary targets should be pulled." yaml:"pull_side_targets" envconfig:"pull_side_targets"`
	Project				string 		`short:"p" long:"project" description:"Docker Compose project. If not specified, current directory will be used instead."`
	Scale         		string 		`short:"s" long:"scale" description:"Number of instances that should be deployed. If value starts with the plug sign (+), the number of instances will be increased by the given number. If value starts with the minus sign (-), the number of instances will be decreased by the given number."`
	ConsulAddress 		string 		`short:"c" long:"consul-address" description:"The address of the consul server." yaml:"consul_address" envconfig:"consul_address"`
	ServiceName   		string
	CurrentColor  		string
	NextColor     		string
	CurrentTarget 		string
	NextTarget    		string
}

// TODO
func getArgs(opts *Opts) (err error) {
	if err := parseYml(opts); err != nil {
		return err
	}
	if err := parseEnvironmentVars(opts); err != nil {
		return err
	}
	if err := parseArgs(opts); err != nil {
		return err
	}
	if err := parseOpts(opts); err != nil {
		return err
	}
	fmt.Printf("%#v\n", opts)
	return nil
}

func parseYml(opts *Opts) error {
	data, err := ioutil.ReadFile(dockerFlowPath)
	if err != nil {
		return fmt.Errorf("Could not read the Docker Flow file %s\n%v", dockerFlowPath, err)
	}
	if err = yaml.Unmarshal([]byte(data), opts); err != nil {
		return fmt.Errorf("Could not parse the Docker Flow file %s\n%v", dockerFlowPath, err)
	}
	return nil
}

func parseArgs(opts *Opts) error {
	if _, err := flags.ParseArgs(opts, os.Args[1:]); err != nil {
		return fmt.Errorf("Could not parse command line arguments\n%v", err)
	}
	return nil
}


func parseEnvironmentVars(opts *Opts) error {
	if err := envconfig.Process("flow", opts); err != nil {
		return fmt.Errorf("Could not retrieve environment variables\n%v", err)
	}
	return nil
}

func parseOpts(opts *Opts) (err error) {
	if len(opts.Project) == 0 {
		dir, _ := os.Getwd()
		opts.Project = dir[strings.LastIndex(dir, "/") + 1:]
	}
	if len(opts.Target) == 0 {
		return fmt.Errorf("target argument is required")
	}
	if len(opts.ConsulAddress) == 0 {
		return fmt.Errorf("consul-address argument is required")
	}
	//	TODO: Verify that scale is a number (with or without +/-)
	opts.ServiceName = fmt.Sprintf("%s_%s", opts.Project, opts.Target)
	opts.CurrentColor, err = getConsulColor(opts.ConsulAddress, opts.ServiceName)
	if err != nil {
		return err
	}
	opts.NextColor, err = getConsulNextColor(opts.ConsulAddress, opts.ServiceName)
	if err != nil {
		return err
	}
	if opts.BlueGreen {
		opts.NextTarget = fmt.Sprintf("%s-%s", opts.Target, opts.NextColor)
		opts.CurrentTarget = fmt.Sprintf("%s-%s", opts.Target, opts.CurrentColor)
	} else {
		opts.NextTarget = opts.Target
		opts.CurrentTarget = opts.Target
	}
	return err
}
