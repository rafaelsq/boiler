package website

import (
	"fmt"
	"net/http"
)

// Handle handle an http request
func Handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
  <head>
    <title>WebSite</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" href="/static/bulma.min.css">
  </head>
  <body>
    <div id="app"></div>
    <script type="module" src="/static/app.js"></script>
  </body>
</html>`)
}
