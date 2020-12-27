package vanity

import (
	"encoding"
)

// Config defines the configuration used by the http handler.
type Config struct {
	// The vanity import prefix shared by all paths, without URL schema and
	// trailing "/".
	//
	// Example:
	//
	//     go.yhsif.com
	//
	// Or:
	//
	//     yhsif.com/go
	Prefix string `yaml:"prefix" json:"prefix"`

	// The mappings of the repositories.
	Mappings []Mapping `yaml:"mappings" json:"mappings"`
}

// Mapping defines a mapping from a vanity path to an actual repository.
type Mapping struct {
	// The path of the vanity URL, with leading "/"/
	//
	// Example:
	//
	//     /vanity
	Path string `yaml:"path" json:"path"`

	// The full URL of the actual repository.
	//
	// Example:
	//
	//     https://github.com/fishy/go-vanity
	URL string `yaml:"url" json:"url"`

	// The VCS of the repository. Default to DefaultVCS ("git").
	VCS VCS `yaml:"vcs" json:"vcs"`
}

// The default VCS to be used.
const (
	DefaultVCS VCS = "git"
)

// VCS defines the vcs ("git", "mod", "hg", etc.) used in Mapping.
//
// It treats zero value (empty string) as DefaultVCS,
// and behave the same as strings otherwise.
type VCS string

var (
	_ encoding.TextMarshaler   = VCS("")
	_ encoding.TextUnmarshaler = (*VCS)(nil)
)

func (v VCS) String() string {
	if v == "" {
		return string(DefaultVCS)
	}
	return string(v)
}

// MarshalText implements encoding.TextMarshaler.
func (v VCS) MarshalText() (text []byte, err error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (v *VCS) UnmarshalText(text []byte) error {
	if len(text) <= 0 {
		*v = DefaultVCS
	} else {
		*v = VCS(text)
	}
	return nil
}
