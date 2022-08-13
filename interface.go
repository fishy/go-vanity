package vanity

import (
	html "html/template"
	"io"
	text "text/template"
)

// Templater is the minimal interface shared between go's text and html
// templates.
type Templater interface {
	Execute(io.Writer, any) error
}

var (
	_ Templater = (*html.Template)(nil)
	_ Templater = (*text.Template)(nil)
)
