package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/logrusorgru/aurora/v3"
	"github.com/signedsecurity/sigsubfind3r/pkg/runner"
)

type options struct {
	sourcesList bool
	noColor     bool
	silent      bool
}

var (
	co options
	so runner.Options
	au aurora.Aurora
)

func banner() {
	fmt.Fprintln(os.Stderr, aurora.BrightBlue(`
     _                 _      __ _           _ _____
 ___(_) __ _ ___ _   _| |__  / _(_)_ __   __| |___ / _ __
/ __| |/ _`+"`"+` / __| | | | '_ \| |_| | '_ \ / _`+"`"+` | |_ \| '__|
\__ \ | (_| \__ \ |_| | |_) |  _| | | | | (_| |___) | |
|___/_|\__, |___/\__,_|_.__/|_| |_|_| |_|\__,_|____/|_| V1.0.0
       |___/
`).Bold())
}

func init() {
	flag.StringVar(&so.Domain, "d", "", "")
	flag.StringVar(&so.Domain, "domain", "", "")
	flag.StringVar(&so.SourcesExclude, "eS", "", "")
	flag.StringVar(&so.SourcesExclude, "exclude-sources", "", "")
	flag.BoolVar(&co.sourcesList, "lS", false, "")
	flag.BoolVar(&co.sourcesList, "list-sources", false, "")
	flag.BoolVar(&co.noColor, "nC", false, "")
	flag.BoolVar(&co.noColor, "no-color", false, "")
	flag.BoolVar(&co.silent, "s", false, "")
	flag.BoolVar(&co.silent, "silent", false, "")
	flag.StringVar(&so.SourcesUse, "uS", "", "")
	flag.StringVar(&so.SourcesUse, "use-sources", "", "")

	flag.Usage = func() {
		banner()

		h := "USAGE:\n"
		h += "  sigsubfind3r [OPTIONS]\n"

		h += "\nOPTIONS:\n"
		h += "  -d,  --domain            domain to find subdomains for\n"
		h += "  -eS, --exclude-sources   comma(,) separated list of sources to exclude\n"
		h += "  -lS, --list-sources      list all the sources available\n"
		h += "  -nC, --no-color          no color mode: Don't use colors in output\n"
		h += "  -s,  --silent            silent mode: Output subdomains only\n"
		h += "  -uS, --use-sources       comma(,) separated list of sources to use\n"

		fmt.Println(h)
	}

	flag.Parse()

	au = aurora.NewAurora(!co.noColor)
}

func main() {
	options, err := runner.ParseOptions(&so)
	if err != nil {
		log.Fatalln(err)
	}

	if !co.silent {
		banner()
	}

	if co.sourcesList {
		fmt.Println("[", au.BrightBlue("INF"), "] current list of the available", au.Underline(strconv.Itoa(len(options.YAMLConfig.Sources))+" sources").Bold())
		fmt.Println("[", au.BrightBlue("INF"), "] sources marked with an * needs key or token")
		fmt.Println("")

		keys := options.YAMLConfig.GetKeys()
		needsKey := make(map[string]interface{})
		keysElem := reflect.ValueOf(&keys).Elem()

		for i := 0; i < keysElem.NumField(); i++ {
			needsKey[strings.ToLower(keysElem.Type().Field(i).Name)] = keysElem.Field(i).Interface()
		}

		for _, source := range options.YAMLConfig.Sources {
			if _, ok := needsKey[source]; ok {
				fmt.Println(">", source, "*")
			} else {
				fmt.Println(">", source)
			}
		}

		os.Exit(0)
	}

	if !co.silent {
		fmt.Println("[", au.BrightBlue("INF"), "]", au.Underline(so.Domain).Bold(), "subdomain enumeration.")
		fmt.Println("")
	}

	runner := runner.New(options)

	subdomains, err := runner.Run()
	if err != nil {
		log.Fatalln(err)
	}

	for n := range subdomains {
		if co.silent {
			fmt.Println(n.Value)
		} else {
			fmt.Printf("[%s] %s\n", au.BrightBlue(n.Source), n.Value)
		}
	}
}
