package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/exec"
	"regexp"
)

type Resolved struct {
	Object  Object
	Version int
}

type Object struct {
	Pins []Pin
}

type Pin struct {
	Package       string
	RepositoryURL string `json:"repositoryURL"`
	State         State
}

type State struct {
	Branch   *string
	Revision string
	Version  string
}

func main() {
	var (
		urlReRaw      string
		reverseRegexp bool
	)
	flag.StringVar(&urlReRaw, "re", "", "specify regexp for urls")
	flag.BoolVar(&reverseRegexp, "rev", false, "reverses regexp")

	flag.Parse()

	var urlRe *regexp.Regexp
	if urlReRaw != "" {
		var err error
		if urlRe, err = regexp.Compile(urlReRaw); err != nil {
			log.Fatalf("failed to parse regexp: %s", err)
		}
	}

	args := flag.Args()
	if len(args) < 1 {
		log.Fatalf("filename should be specified")
	}
	resolvedPath := args[0]

	f, err := os.Open(resolvedPath)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}

	var res Resolved
	if err = json.NewDecoder(f).Decode(&res); err != nil {
		f.Close()
		log.Printf("parse json: %s", err)
		os.Exit(1)
	}
	f.Close()

	if res.Version != 1 {
		log.Printf("version %d is not supported", res.Version)
		os.Exit(1)
	}

	if len(res.Object.Pins) == 0 {
		log.Print("no pins available")
		return
	}

	var urls []string
	for _, p := range res.Object.Pins {
		if urlRe == nil || urlRe.MatchString(p.RepositoryURL) != reverseRegexp {
			urls = append(urls, p.RepositoryURL)
		}
	}

	if len(urls) == 0 {
		log.Print("no pins matched specified regexp")
		return
	}

	wd, _ := os.Getwd()
	log.Printf("fetching %d repos to %s...", len(urls), wd)
	for i, u := range urls {
		if err = cloneRepo(context.Background(), u); err != nil {
			log.Printf("[%d/%d] FAILED – %s, reason: %s", i+1, len(urls), u, err)
		} else {
			log.Printf("[%d/%d] OK – %s", i+1, len(urls), u)
		}
	}
}

// TODO(ArtAndreev): timeout.
func cloneRepo(ctx context.Context, u string) error {
	cmd := exec.CommandContext(ctx, "git", "clone", u)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
