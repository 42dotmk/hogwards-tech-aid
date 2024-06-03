package handlers

import (
	"net/http"
	"github.com/angelofallars/htmx-go"
)

const html = `
<html>
<body>
	<h1>Hello World</h1>
</body>
</html>
`

func TestHandler(w http.ResponseWriter, r *http.Request) {
	if htmx.IsHTMX(r) {
		w.Write([]byte("Hello World HTMX"))
	} else {
		w.Write([]byte(html))
	}
}
