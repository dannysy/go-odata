package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dannysy/go-odata/go-odata-gen/descriptor"
)

type Generator struct {
	entityBuf *bytes.Buffer
	enumBuf   *bytes.Buffer
	nameBuf   *bytes.Buffer

	edmxData       []byte
	enums          map[string]struct{}
	complexTypes   map[string]struct{}
	entitySetTypes map[string]struct{}
	edmxURL        string
	outDir         string
}

func New(path string, outDir string) *Generator {
	g := new(Generator)
	g.entityBuf = new(bytes.Buffer)
	g.enumBuf = new(bytes.Buffer)
	g.nameBuf = new(bytes.Buffer)
	g.enums = make(map[string]struct{})
	g.complexTypes = make(map[string]struct{})
	g.entitySetTypes = make(map[string]struct{})
	g.outDir = outDir
	if strings.HasPrefix(path, "http") {
		resp, err := http.Get(path)
		if err != nil {
			panic(err)
		}
		g.edmxData, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil
		}
	} else {
		var err error
		g.edmxData, err = os.ReadFile(path)
		if err != nil {
			panic(err)
		}
	}

	g.edmxURL = path
	return g
}

func (g *Generator) Generate() {
	// Unmarshal xml
	var edmx descriptor.Edmx
	err := xml.Unmarshal(g.edmxData, &edmx)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	err = os.Mkdir(g.outDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	for _, schema := range edmx.DataServices.Schemas {
		g.reset()

		g.parseEntitySets(schema)
		g.printHeader(schema)
		g.generateEnumTypes(schema)
		g.generateComplexTypes(schema)
		g.generateEntityTypes(schema)

		g.reformat()

		g.writeToFiles(schema)
	}
}

func (g *Generator) writeToFiles(schema descriptor.Schema) {
	packageDir := g.outDir + "/" + strings.ToLower(schema.Namespace)
	err := os.Mkdir(packageDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	file, err := os.Create(packageDir + "/entity.od.go")
	if err != nil {
		log.Fatal(err)
	}
	_, _ = g.entityBuf.WriteTo(file)
	file, err = os.Create(packageDir + "/name.od.go")
	if err != nil {
		log.Fatal(err)
	}
	_, _ = g.nameBuf.WriteTo(file)
	file, err = os.Create(packageDir + "/enum.od.go")
	if err != nil {
		log.Fatal(err)
	}
	_, _ = g.enumBuf.WriteTo(file)
}

func (g *Generator) printHeader(schema descriptor.Schema) {
	generateHeader(g.entityBuf, g.edmxURL)
	generateHeader(g.enumBuf, g.edmxURL)
	generateHeader(g.nameBuf, g.edmxURL)
	generatePackage(g.entityBuf, schema.Namespace)
	generatePackage(g.nameBuf, schema.Namespace)
	generatePackage(g.enumBuf, schema.Namespace)
	generateImports(g.entityBuf)
}

func (g *Generator) parseEntitySets(schema descriptor.Schema) {
	for _, entitySet := range schema.EntityContainer.EntitySets {
		g.entitySetTypes[ToUpperFirstChar(entitySet.EntityType)] = struct{}{}
	}
	for _, entityType := range schema.EntityTypes {
		g.entitySetTypes[ToUpperFirstChar(schema.Namespace)+"."+entityType.Name] = struct{}{}
	}
}

func (g *Generator) generateEnumTypes(schema descriptor.Schema) {
	for _, enumType := range schema.EnumTypes {
		enumName := enumType.Name
		enumTypeName := ToUpperFirstChar(enumName)
		g.enums[schema.Namespace+"."+enumName] = struct{}{}
		P(g.enumBuf, "type ", enumTypeName, " string")
		P(g.enumBuf, "const (")
		for _, member := range enumType.Members {
			P(g.enumBuf, enumTypeName, "_", member.Name, " ", enumTypeName, " = \"", member.Name, "\"")
		}
		P(g.enumBuf, ")")
		P(g.enumBuf)
	}
}

func (g *Generator) generateComplexTypes(schema descriptor.Schema) {
	for _, complexType := range schema.ComplexTypes {
		g.complexTypes[schema.Namespace+"."+complexType.Name] = struct{}{}
	}
	for _, complexType := range schema.ComplexTypes {
		P(g.entityBuf, "type ", ToUpperFirstChar(complexType.Name), " struct {")
		for _, property := range complexType.Properties {
			P(
				g.entityBuf,
				ToUpperFirstChar(property.Name),
				" ",
				g.getPropertyType(property),
				"`json:\"", property.Name, ",omitempty\"`",
			)
		}
		P(g.entityBuf, "}")
		P(g.entityBuf)
	}
}

func (g *Generator) generateEntityTypes(schema descriptor.Schema) {
	entityNamesDeclarationBuf := new(bytes.Buffer)
	P(entityNamesDeclarationBuf, "var EntityNames = struct {")
	entityNamesAssignmentBuf := new(bytes.Buffer)
	entityPropertyNamesBuf := new(bytes.Buffer)
	entityNavNamesBuf := new(bytes.Buffer)
	for _, entityType := range schema.EntityTypes {
		entityPropertyNamesDeclarationBuf := new(bytes.Buffer)
		entityPropertyNamesAssignmentBuf := new(bytes.Buffer)
		entityNavNamesDeclarationBuf := new(bytes.Buffer)
		entityNavNamesAssignmentBuf := new(bytes.Buffer)
		// entity nav property names
		if len(entityType.NavigationProperties) > 0 {
			P(entityNavNamesDeclarationBuf, "var ", ToUpperFirstChar(entityType.Name), "RefNames = struct {")
		}
		//
		// entity property names
		if len(entityType.Properties) > 0 {
			P(entityPropertyNamesDeclarationBuf, "var ", ToUpperFirstChar(entityType.Name), "PropertyNames = struct {")
		}
		//
		// entity names
		P(entityNamesDeclarationBuf, ToUpperFirstChar(entityType.Name), " string")
		P(entityNamesAssignmentBuf, ToUpperFirstChar(entityType.Name), ": \"", entityType.Name, "s\",")
		//
		// entity structs
		P(g.entityBuf, "type ", ToUpperFirstChar(entityType.Name), " struct {")
		for _, property := range entityType.Properties {
			// entity property names
			P(entityPropertyNamesDeclarationBuf, ToUpperFirstChar(property.Name), " string")
			P(entityPropertyNamesAssignmentBuf, ToUpperFirstChar(property.Name), ": \"", property.Name, "\",")
			//
			// entity structs
			P(
				g.entityBuf,
				ToUpperFirstChar(property.Name),
				" ",
				g.getPropertyType(property),
				"`json:\"", property.Name, ",omitempty\"`",
			)
			//
		}
		for _, navProp := range entityType.NavigationProperties {
			// entity nav property names
			P(entityNavNamesDeclarationBuf, ToUpperFirstChar(navProp.Name), " string")
			P(entityNavNamesAssignmentBuf, ToUpperFirstChar(navProp.Name), ": \"", navProp.Name, "\",")
			//
			// entity structs
			P(
				g.entityBuf,
				ToUpperFirstChar(navProp.Name),
				" ",
				g.getNavigationPropertyType(navProp),
				"`json:\"", navProp.Name, ",omitempty\"`",
			)
		}
		// entity nav property names
		if len(entityType.NavigationProperties) > 0 {
			P(entityNavNamesDeclarationBuf, "}{")
			P(entityNavNamesAssignmentBuf, "}")
			P(entityNavNamesBuf, entityNavNamesDeclarationBuf.String(), entityNavNamesAssignmentBuf.String())
		}
		//
		// entity property names
		if len(entityType.Properties) > 0 {
			P(entityPropertyNamesDeclarationBuf, "}{")
			P(entityPropertyNamesAssignmentBuf, "}")
			P(entityPropertyNamesBuf, entityPropertyNamesDeclarationBuf.String(),
				entityPropertyNamesAssignmentBuf.String())
		}
		// entity structs
		P(g.entityBuf, "}")
		P(g.entityBuf)
		//
	}
	// entity names
	P(entityNamesDeclarationBuf, "}{")
	P(entityNamesAssignmentBuf, "}")
	P(g.nameBuf, entityNamesDeclarationBuf.String())
	P(g.nameBuf, entityNamesAssignmentBuf.String())
	//
	// entity property names
	P(g.nameBuf, entityPropertyNamesBuf.String())
	P(g.nameBuf)
	//
	// entity nav property names
	P(g.nameBuf, entityNavNamesBuf.String())
	P(g.nameBuf)
}

func (g *Generator) getNavigationPropertyType(property descriptor.NavigationProperty) string {
	// Making every OneToOne reference to be a pointer
	prefix := "*"
	propertyType := property.Type
	if strings.HasPrefix(propertyType, "Collection(") {
		prefix = "[]"
		// Should cut wrapper
		propertyType = strings.TrimSuffix(
			strings.TrimPrefix(property.Type, "Collection("),
			")",
		)
	}
	if _, ok := g.entitySetTypes[propertyType]; ok {
		// Should cut schema namespace
		return prefix + propertyType[strings.LastIndex(propertyType, ".")+1:]
	}
	Fail("unknown navigation property type", propertyType)
	return ""
}

func (g *Generator) getPropertyType(property descriptor.Property) string {
	propType := property.ConvertTypes()
	if propType != "" {
		return propType
	}

	if strings.HasPrefix(property.Type, "Collection(") {
		// Should cut wrapper and
		// schema namespace
		// example: Collection(WP.CustomType) -> CustomType
		unwrapped := strings.TrimSuffix(
			strings.TrimPrefix(property.Type, "Collection("),
			")",
		)
		return "[]" + unwrapped[strings.LastIndex(unwrapped, ".")+1:]
	}
	if _, ok := g.enums[property.Type]; ok {
		// Should cut schema namespace
		// example: WP.CustomType -> CustomType
		return property.Type[strings.LastIndex(property.Type, ".")+1:]
	}
	if _, ok := g.complexTypes[property.Type]; ok {
		// Should cut schema namespace
		// example: WP.CustomType -> CustomType
		return property.Type[strings.LastIndex(property.Type, ".")+1:]
	}
	Fail("unknown type", property.Type)
	return ""
}

func (g *Generator) reset() {
	g.entityBuf.Reset()
	g.enumBuf.Reset()
	g.nameBuf.Reset()
}

func (g *Generator) reformat() {
	reformat(g.entityBuf)
	reformat(g.enumBuf)
	reformat(g.nameBuf)
}

// reformat generated code.
func reformat(buf *bytes.Buffer) {
	fset := token.NewFileSet()
	ast, err := parser.ParseFile(fset, "", buf, parser.ParseComments)
	if err != nil {
		Fail("bad Go source code was generated:", err.Error())
	}
	buf.Reset()
	err = (&printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}).Fprint(buf, fset, ast)
	if err != nil {
		Fail("generated Go source code could not be reformatted:", err.Error())
	}
}

func generateHeader(buf *bytes.Buffer, edmxUrl string) {
	P(buf, "// Code generated by go-odata-gen. DO NOT EDIT.")
	P(buf, "// source: ", edmxUrl)
	P(buf)
}

func generatePackage(buf *bytes.Buffer, pkg string) {
	P(buf, "package ", strings.ToLower(pkg))
	P(buf)
}

// Generate the imports
func generateImports(buf *bytes.Buffer) {
	P(buf)

	P(buf, "import (")
	P(buf, "\"time\"")
	P(buf, "\"github.com/golang-module/carbon/v2\"")
	P(buf, "\"github.com/google/uuid\"")
	P(buf, ")")

	P(buf, "// Reference imports to suppress errors if they are not otherwise used.")
	P(buf, "var _ = ", "time", ".RubyDate")
	P(buf, "var _ = ", "carbon", ".Version")
	P(buf, "var _ = ", "uuid", ".Reserved")
	P(buf)
}

func ToUpperFirstChar(in string) string {
	return strings.ToUpper(in[0:1]) + in[1:]
}

// P prints the arguments to the generated output.  It handles strings and int32s, plus
// handling indirections because they may be *string, etc.
func P(buf *bytes.Buffer, str ...interface{}) {
	for _, v := range str {
		switch s := v.(type) {
		case string:
			buf.WriteString(s)
		case *string:
			buf.WriteString(*s)
		case bool:
			buf.WriteString(fmt.Sprintf("%t", s))
		case *bool:
			buf.WriteString(fmt.Sprintf("%t", *s))
		case int:
			buf.WriteString(fmt.Sprintf("%d", s))
		case *int32:
			buf.WriteString(fmt.Sprintf("%d", *s))
		case *int64:
			buf.WriteString(fmt.Sprintf("%d", *s))
		case float64:
			buf.WriteString(fmt.Sprintf("%g", s))
		case *float64:
			buf.WriteString(fmt.Sprintf("%g", *s))
		default:
			Fail(fmt.Sprintf("unknown type in printer: %T", v))
		}
	}
	buf.WriteByte('\n')
}

// Fail reports a problem and exits the program.
func Fail(msgs ...string) {
	s := strings.Join(msgs, " ")
	log.Print("go-odata-gen: error:", s)
	os.Exit(1)
}
