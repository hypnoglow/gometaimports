package main

import (
	"bytes"
	ht "html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerator_Generate(t *testing.T) {
	tests := []struct {
		genFunc func() *Generator
		name    string
		host    string
		path    string
		want    string
	}{
		{
			genFunc: func() *Generator {
				gen := &Generator{
					config: &config{
						PackageHost: "go.mycorp.tld",
						Redirect: redirectOptions{
							Host: "https://docs.mycorp.tld",
							Path: "/{{ .Package }}",
						},
						Git: gitOptions{
							Host:       "ssh://git@mycorp.tld",
							PathPrefix: "group",
						},
						Http: httpOptions{
							Host:       "https://git.mycorp.tld",
							PathPrefix: "group",
						},
						MinDepth: 1,
						MaxDepth: 2,
						Rules:    nil,
					},
					template: ht.Must(ht.ParseFiles("templates/index.html.tmpl")),
				}
				return gen
			},
			name: "simple",
			host: "go.mycorp.tld",
			path: "/foo",
			want: `<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
    <meta name="go-import" content="go.mycorp.tld/foo git ssh://git@mycorp.tld/group/foo.git">
    <meta name="go-source" content="go.mycorp.tld/foo https://git.mycorp.tld/group/foo https://git.mycorp.tld/group/foo/tree/master{/dir} https://git.mycorp.tld/group/foo/blob/master{/dir}/{file}#L{line}">
    <meta http-equiv="refresh" content="0; url=https://docs.mycorp.tld/foo">
</head>
<body>
Nothing to see here; move along.
</body>
</html>
`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			err := test.genFunc().Generate(buf, test.host, test.path)
			assert.NoError(t, err)
			assert.Equal(t, test.want, buf.String())
		})
	}
}

func TestGenerator_deducePackage(t *testing.T) {
	testCases := map[string]struct {
		conf config
		path string
		pkg  string
		sub  string
	}{
		"when repo is longer than name": {
			conf: config{
				PackageHost: "",
				Redirect:    redirectOptions{},
				Git:         gitOptions{},
				Http:        httpOptions{},
				MinDepth:    1,
				MaxDepth:    3,
				Rules: []rule{
					{
						Name: "tools/*",
						Repo: "mygroup/tools/*",
					},
				},
			},
			path: "tools/foo",
			pkg:  "mygroup/tools/foo",
			sub:  "",
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			g := &Generator{
				config: &tc.conf,
			}

			pkg, sub, err := g.deducePackage(tc.path)
			assert.NoError(t, err)

			assert.Equal(t, tc.pkg, pkg)
			assert.Equal(t, tc.sub, sub)
		})
	}
}
