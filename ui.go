package postdb

import (
	"net/http"
	"text/template"
)

const t = `
<!DOCTYPE html>
<html>
<head>
  <meta charSet="utf-8" />
  <meta name="viewport" content="user-scalable=no, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0, minimal-ui">
  <title>{{ .Namespace }} - PostDB</title>
  <link rel="stylesheet" href="//cdn.jsdelivr.net/npm/postdb/dist/app.css" />
  <link rel="shortcut icon" href="//cdn.jsdelivr.net/npm/postdb/dist/favicon.png" />
  <script src="//cdn.jsdelivr.net/npm/postdb/dist/app.js"></script>
</head>
<body>
  <main></main>
  <script>
    (function(app) {
		var fetcher = app.fetcher
		fetcher.setBaseUrl('{{ .apiUrl }}')
		fetcher.setHeader('Authorization', '{{ .auth }}')
		fetcher.setHeader('X-Namespace', '{{ .namespace }}')
        app.render('main')
    })(window.PostDBApp)
  </script>
</body>
</html>
`

type UIConfig struct {
	APIUrl    string
	Namespace string
}

type UIMux struct {
	config *UIConfig
	tpl    *template.Template
}

func UI(config UIConfig) *UIMux {
	tpl, err := template.New("-").Parse(t)
	if err != nil {
		panic(err)
	}

	return &UIMux{
		config: &config,
		tpl:    tpl,
	}
}

func (mux *UIMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.tpl.Execute(w, map[string]interface{}{
		"apiUrl":    mux.config.APIUrl,
		"namespace": mux.config.Namespace,
		"auth":      r.Header.Get("Authorization"),
	})
}
