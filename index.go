package vanity

import (
	"html/template"
	"net/http"
)

// IndexTmpl is the html template used to render the index page.
//
// It's exported so users of this package could replace it to a different one if
// they so desire.
//
// The data used to execute the template is Config.
var IndexTmpl Templater = template.Must(template.New("index").Parse(`<!DOCTYPE html>
<html>
<head>
<title>Vanity go projects for {{.Prefix}}</title>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
</head>
<body>
{{if .Mappings}}
<ul>
{{range $mapping := .Mappings}}
{{if $mapping.HideInIndex -}}
{{else -}}
<li><a href="https://pkg.go.dev/{{$.Prefix}}{{$mapping.Path}}"><code>{{$.Prefix}}{{$mapping.Path}}</code></a>: (<a href="{{$mapping.URL}}">src</a>)
{{- if $mapping.Description -}}
&nbsp;{{$mapping.Description}}
{{- end -}}
</li>
{{- end}}
{{- end}}
</ul>
{{else}}
<p>Nothing to see here. Come back later.</p>
{{- end}}
</body>
</html>
`))

// IndexHandler creates a handler to render the index page.
//
// Note that by default this is already included in Handler so you don't need to
// use it directly. It's provided, along with Args.NoIndex, so that you could
// put the index page on a different path.
func IndexHandler(args Args) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		IndexTmpl.Execute(w, args.Config)
		return
	}
}
