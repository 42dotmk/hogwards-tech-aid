This is a Go web APP that I made to learn some basics with go.

This project uses [HTMX](https://github.com/bigskysoftware/htmx) to send AJAX requests, this expects a response in html format.
Due to this, all the CRUD based operations retrieve small html partials that are inserted in to the DOM.
This increases user experience because no page reload is needed.

Technologies used in this project and why:

BACKEND:
  - Go: The main programming language, the http server is running here.
    - [Go Gin](https://github.com/gin-gonic/gin): Web Framework 
    - [GORM](https://github.com/go-gorm/gorm): ORM for interacting with the database

FRONTEND:
  - [HTMX](https://github.com/bigskysoftware/htmx): Used for adding reactivity without the need of refreshing the page. Acomplished sending and receiving AJAX request.
- [TailwindCSS](https://github.com/tailwindlabs/tailwindcss): CSS framework for rapid UI development.

## How to start developing further
1. Clone the repo
1. Install the live reload library with `go install https://github.com/codegangsta/gin`
2. `go mod tidy` to install the dependencies
3. `gin` to start the server, on every file change the server will rebuild



