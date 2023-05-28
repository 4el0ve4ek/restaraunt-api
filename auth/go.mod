module auth

go 1.20

require (
	github.com/4el0ve4ek/restaraunt-api/library v0.0.0-20230528113506-d3e9b12cded8
	github.com/go-chi/chi/v5 v5.0.8
	github.com/lib/pq v1.10.9
	github.com/pkg/errors v0.9.1
	golang.org/x/crypto v0.9.0
)

require gopkg.in/yaml.v3 v3.0.1 // indirect

//replace github.com/4el0ve4ek/restaraunt-api/library => ../library
