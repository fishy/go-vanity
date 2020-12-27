package vanity

import (
	"html/template"
	"net/http"
	"strings"
)

var vanityTmpl = template.Must(template.New("vanity").Parse(`<!DOCTYPE html>
<html>
<head>
<title>Vanity page for {{.Prefix}}{{.Path}}{{.Subpath}}</title>">
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="{{.Prefix}}{{.Path}} {{.VCS}} {{.URL}}">
<meta http-equiv="refresh" content="0; url=https://pkg.go.dev/{{.Prefix}}{{.Path}}{{.Subpath}}">
</head>
<body><p>
Nothing to see here, <a href="https://pkg.go.dev/{{.Prefix}}{{.Path}}{{.Subpath}}">move along</a>.
</p></body>
</html>
`))

type vanityData struct {
	Prefix  string
	Path    string
	URL     string
	VCS     VCS
	Subpath string
}

// Args defines the args used by Handler function.
type Args struct {
	Config Config

	// TODO: Add index page handling.
}

// Handler creates an HTTP handler that handles go vanity URL requests.
func Handler(args Args) http.HandlerFunc {
	for i := range args.Config.Mappings {
		if !strings.HasPrefix(args.Config.Mappings[i].Path, "/") {
			// Common mistake, just append the trailing "/" for them.
			args.Config.Mappings[i].Path = "/" + args.Config.Mappings[i].Path
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		for _, m := range args.Config.Mappings {
			if !strings.HasPrefix(r.URL.Path, m.Path) {
				continue
			}
			// We also don't want to handle /foobar requests with /foo mapping.
			if r.URL.Path != m.Path && !strings.HasPrefix(r.URL.Path, m.Path+"/") {
				continue
			}
			data := vanityData{
				Prefix:  args.Config.Prefix,
				Path:    m.Path,
				URL:     m.URL,
				VCS:     m.VCS,
				Subpath: strings.TrimPrefix(r.URL.Path, m.Path),
			}
			vanityTmpl.Execute(w, data)
			return
		}
		// No mapping found, return 404.
		// TODO: Add index page handling.
		http.NotFound(w, r)
	}
}
