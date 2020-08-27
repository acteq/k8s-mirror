package main

import (
	"k8s.io/klog"
	"os"
	"strings"
	"fmt"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/parsers/hcl"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	flag "github.com/spf13/pflag"
)

var config = koanf.New(".")

func getConfigure() (mirror map[string]string, certFile string , keyFile string ) {

	config.Load(env.Provider("mirror_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "mirror_")), "_", ".", -1)
	}), nil)

	f := flag.NewFlagSet("config", flag.ContinueOnError)
	f.Usage = func() {
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}
	// Path to one or more config files to load into koanf along with some config params.
	f.StringSlice("conf", []string{"conf.hcl"}, "path to one or more .hcl config files")
	f.String( "tls-cert-file", "", 
		"File containing the default x509 Certificate for HTTPS. ")
	f.String("tls-private-key-file", "", 
		"File containing the default x509 private key matching --tls-cert-file.")

	f.Parse(os.Args[1:])

	// Load the config files provided in the commandline.
	cFiles, _ := f.GetStringSlice("conf")
	for _, c := range cFiles {
		if err := config.Load(file.Provider(c), hcl.Parser(true)); err != nil {
			klog.Fatalf("error loading config: %v", err)
		}
	}

	if err := config.Load(posflag.Provider(f, ".", config), nil); err != nil {
		klog.Fatalf("error loading config: %v", err)
	}

	mirror = config.StringMap("mirror")

	certFile = config.String("tls-cert-file")
	keyFile =  config.String("tls-private-key-file")
	
	return 
}