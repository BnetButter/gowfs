package main

import (
	"testing"
	"encoding/xml"
)

func TestCreateFeature(t *testing.T) {
	featuretype := WFSFeatureType{};
	featuretype.Init("Foobar", "Foobarbaz");
	
	as_string := featuretype.Stringify().Unwrap();
	expected_string :=
`<wfs:FeatureType>
 <wfs:Name>Foobar</wfs:Name>
 <wfs:Title>Foobarbaz</wfs:Title>
 <wfs:DefaultCRS>https://www.opengis.net/def/crs/EPSG/0/4326</wfs:DefaultCRS>
</wfs:FeatureType>`
	if as_string != expected_string {
		t.Errorf("%s != %s", expected_string, as_string)
	}
}


func TestCreateDescribeFeatureType_Element(t *testing.T) {
	expectedString := `<xsd:element name="mylayer" type="gowfs:mylayerType" substitutionGroup="gml:AbstractFeature"></xsd:element>`
	xmlElement := DescribeFeatureType_Element_create("mylayer")
	actualString := string(Ensure(xml.MarshalIndent(&xmlElement, "", "")));

	if expectedString != actualString {
		t.Errorf("expected != actual: %s", actualString)
	}
}

func TestCreateDescribeFeatureType_ComplexType(t *testing.T) {
	var expectedString string =
`<xsd:complexType name="mylayerType">
    <xsd:complexContent>
        <xsd:extension base="gml:AbstractFeatureType">
            <xsd:sequence>
                <xsd:element name="geom" type="gml:PointPropertyType" minOccurs="0" maxOccurs="1" nillable="true"></xsd:element>
                <xsd:element name="foo" type="xsd:string" minOccurs="0" maxOccurs="1" nillable="true"></xsd:element>
                <xsd:element name="bar" type="xsd:string" minOccurs="0" maxOccurs="1" nillable="true"></xsd:element>
                <xsd:element name="baz" type="xsd:string" minOccurs="0" maxOccurs="1" nillable="true"></xsd:element>
            </xsd:sequence>
        </xsd:extension>
    </xsd:complexContent>
</xsd:complexType>`

	columns := []Column{
		{
			ColumnName: "foo",
			DataType: "TEXT",
			IsNullable: "TRUE",
		},
		{
			ColumnName: "bar",
			DataType: "TEXT",
			IsNullable: "TRUE",
		},
		{
			ColumnName: "baz",
			DataType: "TEXT",
			IsNullable: "TRUE",
		},
	}

	xmlElement := DescribeFeatureType_CreateColumnSchema("mylayer", columns);
	actualString := string(Ensure(xml.MarshalIndent(&xmlElement, "", "    ")));
	if actualString != expectedString {
		t.Errorf("actual != expected: \n%s\n%s", actualString, expectedString)
	}

}

func TestMarshalColumnTag(t *testing.T) {
	columnTag := GetFeature_ColumnTag{
		Tag: "gowfs:foo",
		Value: "bar",
	}

	expected := `<gowfs:foo>bar</gowfs:foo>`
	columnString := string(Ensure(xml.Marshal(columnTag)))
	if expected != columnString {
		t.Errorf("\n%s != %s", expected, columnString);
	}
}

func TestGMLPos(t *testing.T) {
	pos := GMLPos {
		12.2123,
		64.1233,
	}

	geomString := string(Ensure(xml.Marshal(pos)));
	expecteString := `<gml:pos>64.123300 12.212300</gml:pos>`
	if geomString != expecteString {
		t.Errorf("%s", geomString);
	}
}


func TestGMLGeometryCreateMember(t *testing.T) {
	pointElement := GetFeature_CreatePointMember("mylayer", 1, 12.2123, 64.1233, []GetFeature_ColumnTag{
		{
			Tag:"filename",
			Value: "foo.txt",
		},
		{
			Tag: "name",
			Value: "Kevin",
		},
	});
	gmlString := string(Ensure(xml.Marshal(pointElement)))
	expected := `<wfs:member><gowfs:mylayer gml:id="fid.1"><gowfs:geom><gml:Point srsName="http://www.opengis.net/def/crs/EPSG/0/4326"><gml:pos>64.123300 12.212300</gml:pos></gml:Point></gowfs:geom><gowfs:filename>foo.txt</gowfs:filename><gowfs:name>Kevin</gowfs:name></gowfs:mylayer></wfs:member>`
	if expected != gmlString {
		t.Errorf("%s != %s", expected, gmlString);
	}
}


func TestXMLGetFeature(t *testing.T) {

	filenames := []string{"foo.txt", "bar.txt", "hello.c"}
	names := []string{"kevin", "ian", "lawrance"}
	coords := [][]float64{
		{12.321, 54.4358},
		{89.123, 64.1235},
		{0.1234, 89.1278},
	}

	var pointElements []GetFeature_WFSMember;
	for i := 0; i < 3; i++ {
		coord := coords[i];

		filename := filenames[i]
		name := names[i]

		columns := []GetFeature_ColumnTag{
			{
				Tag:"filename",
				Value: filename,
			},
			{
				Tag: "name",
				Value: name,
			},
		}
		pointElements = append(pointElements, GetFeature_CreatePointMember("mylayer", uint32(i), coord[0], coord[1], columns));
	}

	describeFeature := GetFeature_CreateFeatureCollection(3, pointElements);
	describeFeatureString := string(Ensure(xml.Marshal(describeFeature)))
	expectedString := `<wfs:FeatureCollection xmlns:wfs="http://www.opengis.net/wfs/2.0" xmlns:gml="http://www.opengis.net/gml/3.2" xmlns:gowfs="http://example.com/gowfs" numberMatched="3" numberReturned="3"><wfs:member><gowfs:mylayer gml:id="fid.0"><gowfs:geom><gml:Point srsName="http://www.opengis.net/def/crs/EPSG/0/4326"><gml:pos>54.435800 12.321000</gml:pos></gml:Point></gowfs:geom><gowfs:filename>foo.txt</gowfs:filename><gowfs:name>kevin</gowfs:name></gowfs:mylayer></wfs:member><wfs:member><gowfs:mylayer gml:id="fid.1"><gowfs:geom><gml:Point srsName="http://www.opengis.net/def/crs/EPSG/0/4326"><gml:pos>64.123500 89.123000</gml:pos></gml:Point></gowfs:geom><gowfs:filename>bar.txt</gowfs:filename><gowfs:name>ian</gowfs:name></gowfs:mylayer></wfs:member><wfs:member><gowfs:mylayer gml:id="fid.2"><gowfs:geom><gml:Point srsName="http://www.opengis.net/def/crs/EPSG/0/4326"><gml:pos>89.127800 0.123400</gml:pos></gml:Point></gowfs:geom><gowfs:filename>hello.c</gowfs:filename><gowfs:name>lawrance</gowfs:name></gowfs:mylayer></wfs:member></wfs:FeatureCollection>`
	if describeFeatureString != expectedString {
		t.Errorf("%s != %s", expectedString, describeFeatureString);
	}
}