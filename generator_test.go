package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
