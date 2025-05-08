module mascapp

go 1.24.2

replace github.com/octoberswimmer/thunder => ../

require (
	github.com/octoberswimmer/masc v0.0.0-20250117215935-724533f95fd8
	github.com/octoberswimmer/thunder v0.0.0-00010101000000-000000000000
)

require golang.org/x/sync v0.6.0 // indirect
