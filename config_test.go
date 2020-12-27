package vanity_test

import (
	"encoding/json"
	"testing"

	"go.yhsif.com/vanity"
)

func TestVCS(t *testing.T) {
	t.Run("unmarshal-json-empty", func(t *testing.T) {
		var vcs vanity.VCS
		if err := json.Unmarshal([]byte(`""`), &vcs); err != nil {
			t.Fatal(err)
		}
		if vcs.String() != vanity.DefaultVCS.String() {
			t.Errorf("Expected %q, got %q", vanity.DefaultVCS, vcs)
		}
	})

	for _, c := range []struct {
		label string
		vcs   vanity.VCS
		str   string
	}{
		{
			label: "zero",
			vcs:   "",
			str:   "git",
		},
		{
			label: "git",
			vcs:   "git",
			str:   "git",
		},
		{
			label: "mod",
			vcs:   "mod",
			str:   "mod",
		},
	} {
		t.Run(c.label, func(t *testing.T) {
			var output []byte

			t.Run("marshal-json", func(t *testing.T) {
				var err error
				output, err = json.Marshal(c.vcs)
				if err != nil {
					t.Error(err)
				}
			})
			if t.Failed() {
				t.FailNow()
			}

			t.Run("unmarshal-json-string", func(t *testing.T) {
				var str string
				if err := json.Unmarshal(output, &str); err != nil {
					t.Fatal(err)
				}
				if str != c.str {
					t.Errorf("Expected %q, got %q", c.str, str)
				}
			})

			t.Run("unmarshal-json-vcs", func(t *testing.T) {
				var vcs vanity.VCS
				if err := json.Unmarshal(output, &vcs); err != nil {
					t.Fatal(err)
				}
				if vcs.String() != c.vcs.String() {
					t.Errorf("Expected %q, got %q", c.vcs, vcs)
				}
			})
		})
	}
}
