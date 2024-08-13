package main

import "flag"

//go:generate go run main.go -path https://api.timetta.com/odata/$metadata
func main() {
	var path, outDir string
	flag.StringVar(&path, "path", "", "path to odata metadata edmx file")
	flag.StringVar(&outDir, "out", "", "output directory")
	flag.Parse()
	if path == "" {
		panic("path arg is required")
	}
	g := New(path, outDir)
	g.Generate()
}
