package vanity_test

import (
	"bytes"
	"testing"

	"go.yhsif.com/vanity"
)

func TestIndexTmpl(t *testing.T) {
	for _, c := range []struct {
		label string
		cfg   vanity.Config
	}{
		{
			label: "non-empty",
			cfg: vanity.Config{
				Prefix: "go.yhsif.com",
				Mappings: []vanity.Mapping{
					{
						Path:        "/vanity",
						URL:         "https://github.com/fishy/go-vanity",
						Description: "Go vanity URL handler",
					},
					{
						Path: "/url2epub",
						URL:  "https://github.com/fishy/url2epub",
					},
				},
			},
		},
		{
			label: "empty",
			cfg:   vanity.Config{},
		},
	} {
		t.Run(c.label, func(t *testing.T) {
			var buf bytes.Buffer
			err := vanity.IndexTmpl.Execute(&buf, c.cfg)
			t.Log("Output:")
			t.Log(buf.String())
			if err != nil {
				t.Error(err)
			}
		})
	}
}
