package generator

import (
	"bytes"
	"errors"
	"fmt"
	"path"
	"slices"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/protobuf/compiler/protogen"
)

func GenerateFile(gen *protogen.Plugin, file *protogen.File) (*protogen.GeneratedFile, error) {
	name := *file.Proto.Name

	fileNameWithoutSuffix := strings.TrimSuffix(name, path.Ext(name))
	moduleName := strings.Join(strings.Split(fileNameWithoutSuffix, "/"), ".")

	vars := ConnecpyTemplateVariables{
		FileName:   name,
		ModuleName: moduleName,
		Imports:    importStatements(file),
	}

	packageName := *file.Proto.Package
	for _, svc := range file.Services {
		connecpySvc := &ConnecpyService{
			Name:    string(svc.Desc.Name()),
			Package: packageName,
		}

		for _, method := range svc.Methods {
			// idempotencyLevel := method.Desc.Options.GetIdempotencyLevel()
			noSideEffects := false // idempotencyLevel == descriptor.MethodOptions_NO_SIDE_EFFECTS
			connecpyMethod := &ConnecpyMethod{
				Package:               packageName,
				ServiceName:           connecpySvc.Name,
				Name:                  string(method.Desc.Name()),
				InputType:             symbolName(method.Input),
				InputTypeForProtocol:  symbolName(method.Input),
				OutputType:            symbolName(method.Output),
				OutputTypeForProtocol: symbolName(method.Output),
				NoSideEffects:         noSideEffects,
			}

			connecpySvc.Methods = append(connecpySvc.Methods, connecpyMethod)
		}
		vars.Services = append(vars.Services, connecpySvc)
	}

	var buf = &bytes.Buffer{}
	err := ConnecpyTemplate.Execute(buf, vars)
	if err != nil {
		return nil, err
	}

	g := gen.NewGeneratedFile(file.GeneratedFilenamePrefix+"_connecpy.py", file.GoImportPath)
	g.P(buf.String())
	return g, nil
}

// https://github.com/grpc/grpc/blob/0dd1b2cad21d89984f9a1b3c6249d649381eeb65/src/compiler/python_generator_helpers.h#L67
func moduleName(filename string) string {
	fn, ok := strings.CutSuffix(filename, ".protodevel")
	if !ok {
		fn, _ = strings.CutSuffix(filename, ".proto")
	}
	fn = strings.ReplaceAll(fn, "-", "_")
	fn = strings.ReplaceAll(fn, "/", ".")
	return fn
}

// https://github.com/grpc/grpc/blob/0dd1b2cad21d89984f9a1b3c6249d649381eeb65/src/compiler/python_generator_helpers.h#L80
func moduleAlias(filename string) string {
	mn := moduleName(filename)
	mn = strings.ReplaceAll(mn, "_", "__")
	mn = strings.ReplaceAll(mn, ".", "_dot_")
	return mn
}

func symbolName(msg *protogen.Message) string {
	msg.Desc.Parent() // Ensure the parent is set
	packageName := string(msg.Desc.Parent().FullName())
	name := string(msg.Desc.Name())
	return fmt.Sprintf("%s.%s", moduleAlias(packageName), name)
}

func findFileDescriptor(files []*descriptor.FileDescriptorProto, name string) (*descriptor.FileDescriptorProto, error) {
	//Assumption: Number of files will not be large enough to justify making a map
	for _, f := range files {
		if f.GetName() == name {
			return f, nil
		}
	}
	return nil, errors.New("could not find descriptor")
}

func importStatements(file *protogen.File) []ImportStatement {
	mods := map[string]string{}
	for _, svc := range file.Services {
		for _, method := range svc.Methods {
			method.Input.Desc.Parent()
			inPkg := string(method.Input.Desc.ParentFile().Package())
			mods[moduleName(inPkg)] = moduleAlias(inPkg)
			outPkg := string(method.Output.Desc.ParentFile().Package())
			mods[moduleName(outPkg)] = moduleAlias(outPkg)
		}
	}

	imports := make([]ImportStatement, 0, len(mods))
	for mod, alias := range mods {
		imports = append(imports, ImportStatement{
			Name:  mod,
			Alias: alias,
		})
	}
	slices.SortFunc(imports, func(a, b ImportStatement) int {
		return strings.Compare(a.Name, b.Name)
	})
	return imports
}
