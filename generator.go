package main

import (
	"bytes"
	"errors"
	"fmt"
	ht "html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	tt "text/template"

	"github.com/ghodss/yaml"
)

// NewGenerator returns a new template generator that generates HTML page
// with go-import meta tags.
func NewGenerator(configPath, templatePath string) (*Generator, error) {
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	conf, err := loadConfig(b)
	if err != nil {
		return nil, err
	}

	tpl, err := ht.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	return &Generator{
		config:   conf,
		template: tpl,
	}, nil
}

// Generator is a template generator that generates HTML pages
// with go-import meta tags.
type Generator struct {
	config   *config
	template *ht.Template
}

// Generate generates template based on path and writes it to dst.
func (g *Generator) Generate(dst io.Writer, path string) error {
	path = strings.Trim(path, "/")

	pathParts := strings.Split(path, "/")
	if len(pathParts) < g.config.MinDepth {
		return errors.New("path is below minimum depth")
	}

	// First of all, check against explicit rules.
	pkg := ""
	sub := ""

rules:
	for _, rule := range g.config.Rules {
		nameParts := strings.Split(rule.Name, "/")

	match:
		for i := range nameParts {
			if len(pathParts) < i {
				continue rules
			}

			if nameParts[i] == pathParts[i] {
				continue match
			}

			if nameParts[i] == "*" {
				continue match
			}

			continue rules
		}

		// The most specific (longest) match wins.
		if p := strings.Join(pathParts[:len(nameParts)], "/"); len(pkg) < len(p) {
			pkg = p
			sub = strings.Join(pathParts[len(nameParts):], "/")
		}
	}

	// If no rule matches path, then try to deduce automagically.
	if pkg == "" {
		pkg = strings.Join(pathParts[:g.config.MaxDepth], "/")
		sub = strings.Join(pathParts[g.config.MaxDepth:], "/")
	}

	data := map[string]interface{}{
		"ImportPrefix": filepath.Join(g.config.PackageHost, pkg),
		"Package":      pkg,
		"SubPackage":   sub,
	}

	g.makeGit(data, pkg)
	g.makeHTTP(data, pkg)
	if err := g.makeRedirectURL(data, g.config.RedirectURI); err != nil {
		return fmt.Errorf("execute template: %v", err)
	}

	if err := g.template.Execute(dst, data); err != nil {
		return fmt.Errorf("execute template: %v", err)
	}

	return nil
}

func (g *Generator) makeGit(data map[string]interface{}, pkg string) {
	gitURL := g.config.Git.Host
	if g.config.Http.PathPrefix != "" {
		gitURL += "/" + g.config.Git.PathPrefix
	}
	gitURL += "/" + pkg + ".git"

	data["Git"] = map[string]interface{}{
		"URL": gitURL,
	}
}

func (g *Generator) makeHTTP(data map[string]interface{}, pkg string) {
	httpURL := g.config.Http.Host
	if g.config.Http.PathPrefix != "" {
		httpURL += "/" + g.config.Http.PathPrefix
	}
	httpURL += "/" + pkg

	data["HTTP"] = map[string]interface{}{
		"URL": httpURL,
	}
}

func (g *Generator) makeRedirectURL(data map[string]interface{}, uri string) error {
	t, err := tt.New("redirect_url").Parse(uri)
	if err != nil {
		return err
	}

	buf := bytes.Buffer{}
	err = t.Execute(&buf, data)

	data["RedirectURI"] = buf.String()
	return err
}

func loadConfig(b []byte) (*config, error) {
	var conf config
	if err := yaml.Unmarshal(b, &conf); err != nil {
		return nil, err
	}

	// Defaults

	if conf.MinDepth == 0 {
		conf.MinDepth = 1
	}
	if conf.MaxDepth == 0 {
		conf.MaxDepth = 2
	}

	return &conf, nil
}

type config struct {
	PackageHost string `json:"packageHost"`
	RedirectURI string `json:"redirectURI"`

	Git  gitOptions  `json:"git"`
	Http httpOptions `json:"http"`

	MinDepth int `json:"minDepth"`
	MaxDepth int `json:"maxDepth"`

	Rules []rule `json:"rules"`
}

type gitOptions struct {
	Host       string `json:"host"`
	PathPrefix string `json:"pathPrefix"`
}

type httpOptions struct {
	Host       string `json:"host"`
	PathPrefix string `json:"pathPrefix"`
}

type rule struct {
	Name string `json:"name"`
	Repo string `json:"repo"`
}
