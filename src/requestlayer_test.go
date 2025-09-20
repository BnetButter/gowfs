package main

import (
	"testing"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func TestGetCapabilitiesAsString(t *testing.T) {
	db := Ensure(gorm.Open(postgres.Open(CONNECTION_STRING), &gorm.Config{}));
	db.AutoMigrate(& LayerMetadata{});
	sqlDB := Ensure(db.DB());
	defer sqlDB.Close()

	tableName1 := CreateLayerTable{
		LayerName: "layer_1",
		LayerTitle: "Layer 1",
		Columns: []ColumnType{ 
			{ 
				Name: "name",
				Dtype: "TEXT",
			},
			{
				Name: "address",
				Dtype: "TEXT",
			},
		},
	}

	tableName2 := CreateLayerTable{
		LayerName: "layer_2",
		LayerTitle: "Layer 2",
		Columns: []ColumnType{ 
			{ 
				Name: "foo",
				Dtype: "TEXT",
			},
			{
				Name: "bar",
				Dtype: "TEXT",
			},
		},
	}

	tableName3 := CreateLayerTable{
		LayerName: "layer_3",
		LayerTitle: "Layer 3",
		Columns: []ColumnType{ 
			{ 
				Name: "spam",
				Dtype: "TEXT",
			},
			{
				Name: "egg",
				Dtype: "TEXT",
			},
		},
	}


	
	CreateLayer(db, tableName1).Unwrap();
	defer DeleteLayer(db, "layer_1");
	
	CreateLayer(db, tableName2).Unwrap();
	defer DeleteLayer(db, "layer_2")

	
	CreateLayer(db, tableName3).Unwrap();
	defer DeleteLayer(db, "layer_3")

    expectedString := `<?xml version="1.0" encoding="UTF-8"?>
<wfs:WFS_Capabilities
    version="2.0.0"
    xmlns:wfs="http://www.opengis.net/wfs/2.0"
    xmlns:fes="http://www.opengis.net/fes/2.0"
    xmlns:gml="http://www.opengis.net/gml/3.2"
    xmlns:ows="http://www.opengis.net/ows/1.1"
    xmlns:xlink="http://www.w3.org/1999/xlink"
    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xsi:schemaLocation="
        http://www.opengis.net/wfs/2.0 http://schemas.opengis.net/wfs/2.0/wfs.xsd
        http://www.opengis.net/ows/1.1 http://schemas.opengis.net/ows/1.1.0/owsAll.xsd">

  <ows:ServiceIdentification>
    <ows:Title>Minimal Mock WFS (2.0.0)</ows:Title>
    <ows:ServiceType>WFS</ows:ServiceType>
    <ows:ServiceTypeVersion>2.0.0</ows:ServiceTypeVersion>
  </ows:ServiceIdentification>

  <ows:OperationsMetadata>
    <ows:Operation name="GetCapabilities">
      <ows:DCP><ows:HTTP>
        <ows:Get xlink:href="http://localhost:8000/ows"/>
      </ows:HTTP></ows:DCP>
    </ows:Operation>
    <ows:Operation name="DescribeFeatureType">
      <ows:DCP><ows:HTTP>
        <ows:Get xlink:href="http://localhost:8000/ows"/>
      </ows:HTTP></ows:DCP>
    </ows:Operation>
    <ows:Operation name="GetFeature">
      <ows:DCP><ows:HTTP>
        <ows:Get xlink:href="http://localhost:8000/ows"/>
      </ows:HTTP></ows:DCP>
    </ows:Operation>
    
    <ows:Operation name="Transaction">
      <ows:DCP>
        <ows:HTTP>
          <ows:Post xlink:href="http://localhost:8000/ows"/>
        </ows:HTTP>
      </ows:DCP>
    </ows:Operation>
  </ows:OperationsMetadata>

  <wfs:FeatureTypeList>
    
<wfs:FeatureType>
 <wfs:Name>layer_1</wfs:Name>
 <wfs:Title>Layer 1</wfs:Title>
 <wfs:DefaultCRS>https://www.opengis.net/def/crs/EPSG/0/4326</wfs:DefaultCRS>
</wfs:FeatureType>
<wfs:FeatureType>
 <wfs:Name>layer_2</wfs:Name>
 <wfs:Title>Layer 2</wfs:Title>
 <wfs:DefaultCRS>https://www.opengis.net/def/crs/EPSG/0/4326</wfs:DefaultCRS>
</wfs:FeatureType>
<wfs:FeatureType>
 <wfs:Name>layer_3</wfs:Name>
 <wfs:Title>Layer 3</wfs:Title>
 <wfs:DefaultCRS>https://www.opengis.net/def/crs/EPSG/0/4326</wfs:DefaultCRS>
</wfs:FeatureType>
  </wfs:FeatureTypeList>

</wfs:WFS_Capabilities>
`
	xmlString := GetCapabilities(db).Unwrap()
    if xmlString != expectedString {
        t.Errorf("XML strings don't match")
    }

}

func TestDescribeFeatureType(t *testing.T) {
	db := Ensure(gorm.Open(postgres.Open(CONNECTION_STRING), &gorm.Config{}));
	db.AutoMigrate(& LayerMetadata{});
	sqlDB := Ensure(db.DB());
	defer sqlDB.Close()

	tableName1 := CreateLayerTable{
		LayerName: "layer_1",
		LayerTitle: "Layer 1",
		Columns: []ColumnType{ 
			{ 
				Name: "name",
				Dtype: "TEXT",
			},
			{
				Name: "address",
				Dtype: "TEXT",
			},
		},
	}

	tableName2 := CreateLayerTable{
		LayerName: "layer_2",
		LayerTitle: "Layer 2",
		Columns: []ColumnType{ 
			{ 
				Name: "foo",
				Dtype: "TEXT",
			},
			{
				Name: "bar",
				Dtype: "TEXT",
			},
		},
	}

	tableName3 := CreateLayerTable{
		LayerName: "layer_3",
		LayerTitle: "Layer 3",
		Columns: []ColumnType{ 
			{ 
				Name: "spam",
				Dtype: "TEXT",
			},
			{
				Name: "egg",
				Dtype: "TEXT",
			},
		},
	}


	
	CreateLayer(db, tableName1).Unwrap();
	defer DeleteLayer(db, "layer_1");
	
	CreateLayer(db, tableName2).Unwrap();
	defer DeleteLayer(db, "layer_2")

	
	CreateLayer(db, tableName3).Unwrap();
	defer DeleteLayer(db, "layer_3")

	expectedOutput := `<?xml version="1.0" encoding="UTF-8"?>
<xsd:schema
    targetNamespace="http://example.com/gowfs"
    xmlns:xsd="http://www.w3.org/2001/XMLSchema"
    xmlns:gml="http://www.opengis.net/gml/3.2"
    xmlns:gowfs="http://example.com/gowfs"
    elementFormDefault="qualified">
    <xsd:complexType name="layer_1Type">
  <xsd:complexContent>
    <xsd:extension base="gml:AbstractFeatureType">
      <xsd:sequence>
        <xsd:element name="geom" type="gml:PointPropertyType" minOccurs="0" maxOccurs="1" nillable="true"></xsd:element>
        <xsd:element name="name" type="xsd:string" minOccurs="0" maxOccurs="1" nillable="true"></xsd:element>
        <xsd:element name="address" type="xsd:string" minOccurs="0" maxOccurs="1" nillable="true"></xsd:element>
      </xsd:sequence>
    </xsd:extension>
  </xsd:complexContent>
</xsd:complexType>
</xsd:schema>
`
	describeFeatureTypeString := Ensure(DescribeFeatureType(db, "layer_1"));
	if describeFeatureTypeString != expectedOutput {
		t.Errorf("string not eq: \n%s\n%s", describeFeatureTypeString, expectedOutput)
	}
}

func TestGetFeature(t *testing.T) {
	tableName1 := CreateLayerTable{
		LayerName: "layer_1",
		LayerTitle: "Layer 1",
		Columns: []ColumnType{ 
			{ 
				Name: "name",
				Dtype: "TEXT",
			},
			{
				Name: "address",
				Dtype: "TEXT",
			},
		},
	}
	db := Ensure(gorm.Open(postgres.Open(CONNECTION_STRING), &gorm.Config{}))
	sqlDb := Ensure(db.DB());
	defer sqlDb.Close()
	CreateLayer(db, tableName1).Unwrap();
	defer DeleteLayer(db, "layer_1");

	sqlInsert := `
		INSERT INTO layer_1 (geom, name, address) VALUES
		(ST_GeomFromText('POINT(12.123 42.789)', 4326), 'Kevin',   '1234 Somewhere Ave'),
		(ST_GeomFromText('POINT(14.128 12.238)', 4326), 'Lawrance','8291 Foobar Drive'),
		(ST_GeomFromText('POINT(28.479 94.177)', 4326), 'Peter',   '192 Hello World Blvd');
	`

	if db.Exec(sqlInsert).Error != nil {
		t.Errorf("Failed to insert")
	}

	outputString := Ensure(GetFeature(db, "layer_1", nil));
	actualString := `<wfs:FeatureCollection xmlns:wfs="http://www.opengis.net/wfs/2.0" xmlns:gml="http://www.opengis.net/gml/3.2" xmlns:gowfs="http://example.com/gowfs" numberMatched="3" numberReturned="3">
  <wfs:member>
    <gowfs:layer_1 gml:id="fid.1">
      <gowfs:geom>
        <gml:Point srsName="http://www.opengis.net/def/crs/EPSG/0/4326">
          <gml:pos>42.789000 12.123000</gml:pos>
        </gml:Point>
      </gowfs:geom>
      <gowfs:name>Kevin</gowfs:name>
      <gowfs:address>1234 Somewhere Ave</gowfs:address>
    </gowfs:layer_1>
  </wfs:member>
  <wfs:member>
    <gowfs:layer_1 gml:id="fid.2">
      <gowfs:geom>
        <gml:Point srsName="http://www.opengis.net/def/crs/EPSG/0/4326">
          <gml:pos>12.238000 14.128000</gml:pos>
        </gml:Point>
      </gowfs:geom>
      <gowfs:name>Lawrance</gowfs:name>
      <gowfs:address>8291 Foobar Drive</gowfs:address>
    </gowfs:layer_1>
  </wfs:member>
  <wfs:member>
    <gowfs:layer_1 gml:id="fid.3">
      <gowfs:geom>
        <gml:Point srsName="http://www.opengis.net/def/crs/EPSG/0/4326">
          <gml:pos>94.177000 28.479000</gml:pos>
        </gml:Point>
      </gowfs:geom>
      <gowfs:name>Peter</gowfs:name>
      <gowfs:address>192 Hello World Blvd</gowfs:address>
    </gowfs:layer_1>
  </wfs:member>
</wfs:FeatureCollection>`
	if outputString != actualString {
		t.Errorf("Maybe the attribute orders are different? string does not eq \n%s\n%s\n", actualString, outputString)
	}
}



func TestGetFeatureGeoJSON(t *testing.T) {
	tableName1 := CreateLayerTable{
		LayerName: "layer_1",
		LayerTitle: "Layer 1",
		Columns: []ColumnType{ 
			{ 
				Name: "name",
				Dtype: "TEXT",
			},
			{
				Name: "address",
				Dtype: "TEXT",
			},
		},
	}
	db := Ensure(gorm.Open(postgres.Open(CONNECTION_STRING), &gorm.Config{}))
	sqlDb := Ensure(db.DB());
	defer sqlDb.Close()
	CreateLayer(db, tableName1).Unwrap();
	defer DeleteLayer(db, "layer_1");

	sqlInsert := `
		INSERT INTO layer_1 (geom, name, address) VALUES
		(ST_GeomFromText('POINT(12.123 42.789)', 4326), 'Kevin',   '1234 Somewhere Ave'),
		(ST_GeomFromText('POINT(14.128 12.238)', 4326), 'Lawrance','8291 Foobar Drive'),
		(ST_GeomFromText('POINT(28.479 94.177)', 4326), 'Peter',   '192 Hello World Blvd');
	`

	if db.Exec(sqlInsert).Error != nil {
		t.Errorf("Failed to insert")
	}

	outputString := Ensure(GetFeatureGeoJSON(db, "layer_1", nil));
	actualString := `{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "id": "1",
      "geometry": {
        "type": "Point",
        "coordinates": [
          12.123,
          42.789
        ]
      },
      "properties": {
        "address": "1234 Somewhere Ave",
        "name": "Kevin"
      }
    },
    {
      "type": "Feature",
      "id": "2",
      "geometry": {
        "type": "Point",
        "coordinates": [
          14.128,
          12.238
        ]
      },
      "properties": {
        "address": "8291 Foobar Drive",
        "name": "Lawrance"
      }
    },
    {
      "type": "Feature",
      "id": "3",
      "geometry": {
        "type": "Point",
        "coordinates": [
          28.479,
          94.177
        ]
      },
      "properties": {
        "address": "192 Hello World Blvd",
        "name": "Peter"
      }
    }
  ]
}`
	if outputString != actualString {
		t.Errorf("Maybe the attribute orders are different? string does not eq \n%s\n%s\n", actualString, outputString)
	}
}


