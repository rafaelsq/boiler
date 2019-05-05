package website

import (
	"fmt"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
  <head>
    <title>WebSite</title>
  </head>
  <body>
    <div id="app"></div>
    <script type="module" src="/static/app.js"></script>
  </body>
</html>`)
}
