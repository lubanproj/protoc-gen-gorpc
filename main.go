package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/lubanproj/protoc-gen-gorpc/gorpc"
	gengo "google.golang.org/protobuf/cmd/protoc-gen-go/internal_gengo"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	var (
		flags        flag.FlagSet
		plugins      = flags.String("plugins", "", "list of plugins to enable (supported values: grpc)")
		importPrefix = flags.String("import_prefix", "", "prefix to prepend to import paths")
	)
	importRewriteFunc := func(importPath protogen.GoImportPath) protogen.GoImportPath {
		switch importPath {
		case "context", "fmt", "math":
			return importPath
		}
		if *importPrefix != "" {
			return protogen.GoImportPath(*importPrefix) + importPath
		}
		return importPath
	}
	protogen.Options{
		ParamFunc:         flags.Set,
		ImportRewriteFunc: importRewriteFunc,
	}.Run(func(gen *protogen.Plugin) error {
		isGorpc := false
		for _, plugin := range strings.Split(*plugins, ",") {
			switch plugin {
			case "gorpc":
				isGorpc = true
			case "":
			default:
				return fmt.Errorf("protoc-gen-go: unknown plugin %q", plugin)
			}
		}
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			g := gengo.GenerateFile(gen, f)
			if isGorpc {
				gorpc.GenerateFileContent(gen, f, g)
			}
		}
		gen.SupportedFeatures = gengo.SupportedFeatures
		return nil
	})
}
