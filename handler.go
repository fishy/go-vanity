package vanity

import (
	"html/template"
	"net/http"
	"strings"
)

// MetaTemplateLine is the html->head->meta line that must exist in vanity
// page html template.
const MetaTemplateLine = `<meta name="go-import" content="{{.Prefix}}{{.Path}} {{.VCS}} {{.URL}}">`

// PageTmpl is the html template used to render the vanity pages.
//
// It's exported so users of this package could replace it to a different one if
// they so desire. But when replacing it, MetaTemplateLine (or equivalent)
// must exists in the template, or the vanity page might be invalid.
//
// The data used to execute the template is defined in PageTmplData.
var PageTmpl Templater = template.Must(template.New("vanity").Parse(`<!DOCTYPE html>
<html>
<head>
<title>Vanity page for {{.Prefix}}{{.Reqpath}}</title>
` + MetaTemplateLine + `
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta http-equiv="refresh" content="0; url=https://pkg.go.dev/{{.Prefix}}{{.Reqpath}}">
</head>
<body><p>
Nothing to see here, <a href="https://pkg.go.dev/{{.Prefix}}{{.Reqpath}}">move along</a>.
</p></body>
</html>
`))

// PageTmplData is the data used to execute PageTmpl.
type PageTmplData struct {
	// From Config.Prefix.
	Prefix string

	// From Config.Mappings
	Path string
	URL  string
	VCS  VCS

	// The path from the request,
	// could be different from Path but always have Path as prefix.
	Reqpath string
}

// Args defines the args used by Handler function.
type Args struct {
	Config Config

	// If set to true, do not render and handle the index page.
	NoIndex bool
}

// Handler creates an HTTP handler that handles go vanity URL requests.
func Handler(args Args) http.HandlerFunc {
	for i := range args.Config.Mappings {
		if !strings.HasPrefix(args.Config.Mappings[i].Path, "/") {
			// Common mistake, just append the trailing "/" for them.
			args.Config.Mappings[i].Path = "/" + args.Config.Mappings[i].Path
		}
	}
	indexHandler := IndexHandler(args)

	return func(w http.ResponseWriter, r *http.Request) {
		if !args.NoIndex && r.URL.Path == "/" {
			indexHandler(w, r)
			return
		}

		for _, m := range args.Config.Mappings {
			if r.URL.Path != m.Path && !strings.HasPrefix(r.URL.Path, m.Path+"/") {
				continue
			}
			data := PageTmplData{
				Prefix:  args.Config.Prefix,
				Path:    m.Path,
				URL:     m.URL,
				VCS:     m.VCS,
				Reqpath: r.URL.Path,
			}
			PageTmpl.Execute(w, data)
			return
		}
		// No mapping found, return 404.
		http.NotFound(w, r)
	}
}
