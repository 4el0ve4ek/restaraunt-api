module orders

go 1.20

//replace github.com/4el0ve4ek/restaraunt-api/library => ../library

require (
	github.com/4el0ve4ek/restaraunt-api/library v0.0.0-20230528205726-f0baec0f52fd
	github.com/go-chi/chi/v5 v5.0.8
	github.com/pkg/errors v0.9.1
)

require (
	github.com/lib/pq v1.10.9 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
