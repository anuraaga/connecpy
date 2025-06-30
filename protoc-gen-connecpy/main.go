package main

import (
	"flag"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/i2y/connecpy/protoc-gen-connecpy/generator"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	flag.Parse()

	var flags flag.FlagSet

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(plugin.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generator.GenerateFile(gen, f)
		}
		return nil
	})
}
