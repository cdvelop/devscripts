module github.com/cdvelop/devscripts

go 1.25.2

require (
	github.com/cdvelop/badges v0.0.3
	github.com/cdvelop/gotest v0.0.1
	github.com/cdvelop/mdgo v0.0.9
)

replace (
	github.com/cdvelop/badges v0.0.2 => ../badges
	github.com/cdvelop/gotest v0.0.1 => ../gotest
	github.com/cdvelop/mdgo v0.0.1 => ../mdgo
)
