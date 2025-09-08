package main

import (
	"encoding/xml"
	"fmt"
)

type WFSFeatureType struct {
	XMLName xml.Name `xml:"wfs:FeatureType"`
	Name string `xml:"wfs:Name"`
	Title string `xml:"wfs:Title"`
	DefaultCRS string `xml:"wfs:DefaultCRS"`	
}

func (feature *WFSFeatureType) Init(name string, title string) {
	*feature = WFSFeatureType{
		Name: name,
		Title: title,
		DefaultCRS: "https://www.opengis.net/def/crs/EPSG/0/4326",
	}
}

func (feature *WFSFeatureType) Stringify() Result[string] {
	output, err := xml.MarshalIndent(*feature, "", " ")
	return Result[string]{
		string(output),
		err,
	}
}


// DescribeFeatureType

type DescribeFeatureType_XMLColumnElement struct {
	XMLName xml.Name `xml:"xsd:element"`
	Name string `xml:"name,attr"` // ColumnName
	Type string `xml:"type,attr"`
	MinOccurs uint `xml:"minOccurs,attr"`
	MaxOccurs uint `xml:"maxOccurs,attr"`
	Nillable bool `xml:"nillable,attr"`
}



func (xmlElement *DescribeFeatureType_XMLColumnElement) GeomColumn() {
	*xmlElement = DescribeFeatureType_XMLColumnElement{
		Name: "geom",
		Type: "gml:PointPropertyType",
		MinOccurs: 0,
		MaxOccurs: 1,
		Nillable: true,
	}
}

type DescribeFeatureType_Sequence struct {
	XMLName xml.Name `xml:"xsd:sequence"`
	Elements []DescribeFeatureType_XMLColumnElement
}

type DescribeFeatureType_Extension struct {
	XMLName xml.Name `xml:"xsd:extension"`
	Sequence DescribeFeatureType_Sequence
	Base string `xml:"base,attr"`
}

type DescribeFeatureType_ComplexContent struct {
	XMLName xml.Name `xml:"xsd:complexContent"`
	Extension DescribeFeatureType_Extension
}

type DescribeFeatureType_ComplexType struct {
	XMLName xml.Name `xml:"xsd:complexType"`
	ComplexContent DescribeFeatureType_ComplexContent
	Name string `xml:"name,attr"`
}

type DescribeFeatureType_Element struct {
	XMLName xml.Name `xml:"xsd:element"`
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
	SubstitutionGroup string `xml:"substitutionGroup,attr"`
}

func DescribeFeatureType_Element_create(layerName string) DescribeFeatureType_Element {
	with_ns := fmt.Sprintf("%s:%sType", "gowfs", layerName);
	return DescribeFeatureType_Element{
		Name:layerName,
		Type:with_ns,
		SubstitutionGroup: "gml:AbstractFeature",
	}
}

func convertSQLTypeToXSD(sqlType string) string {
	return "xsd:string"
}

func convertSqlBool(sqlbool string) bool {
	return true;
}


func DescribeFeatureType_CreateColumnSchema(layerName string, columns []Column) DescribeFeatureType_ComplexType {
	elements := make([]DescribeFeatureType_XMLColumnElement, 1, len(columns) + 1)
	
	geom_col := DescribeFeatureType_XMLColumnElement{};
	geom_col.GeomColumn() // init geom
	elements[0] = geom_col;

	for _, col := range(columns) {
		elements = append(elements, DescribeFeatureType_XMLColumnElement{
			Name: col.ColumnName,
			Type: convertSQLTypeToXSD(col.DataType),
			Nillable: convertSqlBool(col.IsNullable),
			MinOccurs: 0,
			MaxOccurs: 1,
		})
	}

	return DescribeFeatureType_ComplexType{
		Name: fmt.Sprintf("%sType", layerName),
		ComplexContent: DescribeFeatureType_ComplexContent{
			Extension: DescribeFeatureType_Extension{
				Base: "gml:AbstractFeatureType",
				Sequence: DescribeFeatureType_Sequence{
					Elements: elements,
				},
			},
		},
	}
} 

type GetFeature_ColumnTag struct {
	Tag string
	Value string
}

func (tag GetFeature_ColumnTag) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	// Just emit the tag name with the prefix as part of Local (Go won't manage prefixes otherwise)
	startElem := xml.StartElement{
		Name: xml.Name{
			Local: "gowfs:" + tag.Tag,
		},
	}

	// Write start, content, and end
	if err := e.EncodeToken(startElem); err != nil {
		return err
	}
	if err := e.EncodeToken(xml.CharData([]byte(tag.Value))); err != nil {
		return err
	}
	return e.EncodeToken(startElem.End())
}


type GMLPos struct {
	X float64
	Y float64
}

func (p GMLPos) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	start := xml.StartElement{
		Name: xml.Name{Local: "gml:pos"}, // fixed tag name
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	
	valueString := fmt.Sprintf("%f %f", p.Y, p.X);

	if err := e.EncodeToken(xml.CharData([]byte(valueString))); err != nil {
		return err
	}
	return e.EncodeToken(start.End())
}


type GMLPoint struct {
	XMLName xml.Name `xml:"gml:Point"`
	SrsName string `xml:"srsName,attr"`
	Pos GMLPos;
}


type GMLGeomColumn struct {
	XMLName string `xml:"gowfs:geom"`
	Point GMLPoint
}

type GMLGeomElement struct {
	LayerName string
	Geom GMLGeomColumn
	GmlID string // Will be rendered as gml:id
	Columns []GetFeature_ColumnTag
}



func (g GMLGeomElement) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	// Split prefix and local name from Tag
	

	// Compose the element name
	start := xml.StartElement{
		Name: xml.Name{
			Local: "gowfs:" + g.LayerName,
		},
		Attr: []xml.Attr{
			{
				Name:  xml.Name{Local: "gml:id"},
				Value: g.GmlID,
			},
		},
	}

	// Write start tag
	if err := e.EncodeToken(start); err != nil {
		return err
	}

	// Encode the Point element
	if err := e.Encode(g.Geom); err != nil {
		return err
	}
		// Encode the Point element
	if err := e.Encode(g.Columns); err != nil {
		return err
	}

	// Write end tag
	return e.EncodeToken(start.End())
}

type GetFeature_WFSMember struct {
	XMLName xml.Name `xml:"wfs:member"`
	Member GMLGeomElement
}

func GetFeature_CreatePointMember(layerName string, layerId uint32, X float64, Y float64, columns []GetFeature_ColumnTag) GetFeature_WFSMember {
	return GetFeature_WFSMember{
			Member:GMLGeomElement{
			LayerName: layerName,
			Geom: GMLGeomColumn{
				Point: GMLPoint {
					SrsName: "http://www.opengis.net/def/crs/EPSG/0/4326",
					Pos: GMLPos{
						X,
						Y,
					},
				},
			},
			Columns: columns,
			GmlID: fmt.Sprintf("fid.%d", layerId),
		},
	}
}

type GetFeature_FeatureCollection struct {
	XMLName xml.Name `xml:"wfs:FeatureCollection"`
	WFSNS string `xml:"xmlns:wfs,attr"`
	GMLNS string `xml:"xmlns:gml,attr"`
	GOWFS string `xml:"xmlns:gowfs,attr"`
	NumberMatched int `xml:"numberMatched,attr"`
	NumberReturned int `xml:"numberReturned,attr"`
	Members []GetFeature_WFSMember
}

func GetFeature_CreateFeatureCollection (numMatched int, members []GetFeature_WFSMember) GetFeature_FeatureCollection {
	return GetFeature_FeatureCollection{
		WFSNS: "http://www.opengis.net/wfs/2.0",
		GMLNS: "http://www.opengis.net/gml/3.2",
		GOWFS: "http://example.com/gowfs",
		NumberMatched: numMatched,
		NumberReturned: len(members),
		Members:members,
	}
}

