package descriptor

import (
	"encoding/xml"
	"strings"
)

type PropertyRef struct {
	Name string `xml:"Name,attr"`
}

type Key struct {
	PropertyRefs []PropertyRef `xml:"PropertyRef"`
}

type NavigationProperty struct {
	Name           string `xml:"Name,attr"`
	Type           string `xml:"Type,attr"`
	Nullable       bool   `xml:"Nullable,attr"`
	ContainsTarget bool   `xml:"ContainsTarget,attr"`
}

type Property struct {
	Name        string `xml:"Name,attr"`
	Type        string `xml:"Type,attr"`
	Nullable    bool   `xml:"Nullable,attr"`
	MaxLength   int32  `xml:"MaxLength,attr"`
	FixedLength bool   `xml:"FixedLength,attr"`
	Unicode     bool   `xml:"Unicode,attr"`
}

type EntityType struct {
	XMLName              xml.Name             `xml:"EntityType"`
	Name                 string               `xml:"Name,attr"`
	BaseType             string               `xml:"BaseType,attr"`
	Key                  Key                  `xml:"Key"`
	Properties           []Property           `xml:"Property"`
	NavigationProperties []NavigationProperty `xml:"NavigationProperty"`
}

type ComplexType struct {
	Name       string     `xml:"Name,attr"`
	Properties []Property `xml:"Property"`
}

type EnumType struct {
	Name    string   `xml:"Name,attr"`
	Members []Member `xml:"Member"`
}

type Member struct {
	Name  string `xml:"Name,attr"`
	Value string `xml:"Value,attr"`
}

type End struct {
	Multiplicity string `xml:"Multiplicity,attr"`
	Type         string `xml:"Type,attr"`
	Role         string `xml:"Role,attr"`
	EntitySet    string `xml:"EntitySet,attr"`
}

type Dependent struct {
	Role        string `xml:"Role,attr"`
	PropertyRef PropertyRef
}

type Principal struct {
	Role        string `xml:"Role,attr"`
	PropertyRef PropertyRef
}

type ReferentialConstraint struct {
	Principal Principal
	Dependent Dependent
}

type Association struct {
	Name                  string `xml:"Name,attr"`
	Ends                  []End  `xml:"End"`
	ReferentialConstraint ReferentialConstraint
}

type AssociationSet struct {
	Name        string `xml:"Name,attr"`
	Association string `xml:"Association,attr"`
	Ends        []End  `xml:"End"`
}

type EntitySet struct {
	Name       string `xml:"Name,attr"`
	EntityType string `xml:"EntityType,attr"`
}

type EntityContainer struct {
	Name                     string           `xml:"Name,attr"`
	IsDefaultEntityContainer bool             `xml:"IsDefaultEntityContainer,attr"`
	LazyLoadingEnabled       bool             `xml:"LazyLoadingEnabled,attr"`
	EntitySets               []EntitySet      `xml:"EntitySet"`
	AssociationSets          []AssociationSet `xml:"AssociationSet"`
}

type Schema struct {
	Namespace       string        `xml:"Namespace,attr"`
	EntityTypes     []EntityType  `xml:"EntityType"`
	EnumTypes       []EnumType    `xml:"EnumType"`
	ComplexTypes    []ComplexType `xml:"ComplexType"`
	EntityContainer EntityContainer
}

type DataServices struct {
	DataServiceVersion    string   `xml:"DataServiceVersion,attr"`
	MaxDataServiceVersion float32  `xml:"MaxDataServiceVersion,attr"`
	Schemas               []Schema `xml:"Schema"`
}

type Edmx struct {
	XMLName      xml.Name `xml:"Edmx"`
	DataServices DataServices
}

func (p *Property) ConvertTypes() string {
	var typesMap = map[string]string{
		"Edm.Binary":                   "[]byte",
		"Edm.Boolean":                  "bool",
		"Edm.Byte":                     "byte",
		"Edm.DateTime":                 "*carbon.DateTime",
		"Edm.Date":                     "*carbon.Date",
		"Edm.Decimal":                  "float32",
		"Edm.Double":                   "float64",
		"Edm.Single":                   "float32",
		"Edm.Guid":                     "uuid.NullUUID",
		"Edm.Int16":                    "int16",
		"Edm.Int32":                    "int32",
		"Edm.Int64":                    "int64",
		"Edm.SByte":                    "",
		"Edm.String":                   "string",
		"Edm.Time":                     "*carbon.Time",
		"Edm.DateTimeOffset":           "time.Time",
		"Edm.Geography":                "",
		"Edm.GeographyPoint":           "",
		"Edm.GeographyLineString":      "",
		"Edm.GeographyPolygon":         "",
		"Edm.GeographyMultiPoint":      "",
		"Edm.GeographyMultiLineString": "",
		"Edm.GeographyMultiPolygon":    "",
		"Edm.GeographyCollection":      "",
		"Edm.Geometry":                 "",
		"Edm.GeometryPoint":            "",
		"Edm.GeometryLineString":       "",
		"Edm.GeometryPolygon":          "",
		"Edm.GeometryMultiPoint":       "",
		"Edm.GeometryMultiLineString":  "",
		"Edm.GeometryMultiPolygon":     "",
		"Edm.GeometryCollection":       "",
		"Edm.Stream":                   "",
		"Edm.Untyped":                  "any",
	}
	if strings.HasPrefix(p.Type, "Collection(") {
		wrappedType := p.Type[strings.Index(p.Type, "(")+1 : strings.Index(p.Type, ")")]
		if resultType, ok := typesMap[wrappedType]; ok {
			return "[]" + resultType
		}
	}
	return typesMap[p.Type]
}
