package main

import (
	"bytes"
	"encoding/xml"
	"text/template"
  "fmt"
	"gorm.io/gorm"
)


const GET_CAPABILITES_XML_TEMPLATE string = 
`<?xml version="1.0" encoding="UTF-8"?>
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
        <ows:Get xlink:href="{{.BaseURL}}"/>
      </ows:HTTP></ows:DCP>
    </ows:Operation>
    <ows:Operation name="DescribeFeatureType">
      <ows:DCP><ows:HTTP>
        <ows:Get xlink:href="{{.BaseURL}}"/>
      </ows:HTTP></ows:DCP>
    </ows:Operation>
    <ows:Operation name="GetFeature">
      <ows:DCP><ows:HTTP>
        <ows:Get xlink:href="{{.BaseURL}}"/>
      </ows:HTTP></ows:DCP>
    </ows:Operation>
    
    <ows:Operation name="Transaction">
      <ows:DCP>
        <ows:HTTP>
          <ows:Post xlink:href="{{.BaseURL}}"/>
        </ows:HTTP>
      </ows:DCP>
    </ows:Operation>
  </ows:OperationsMetadata>

  <wfs:FeatureTypeList>
    {{.FeatureLayers}}
  </wfs:FeatureTypeList>

</wfs:WFS_Capabilities>
`



func GetCapabilities(db *gorm.DB) Result[string] {
	t, err := template.New("GetCapabilities").Parse(GET_CAPABILITES_XML_TEMPLATE);
	if err != nil {
		return Err[string](err);
	}
	
	
	layerMetadata, err := GetLayerMetadata(db);
	if err != nil {
		return Err[string](err);
	}

	featureXMLString := ""

	for _, metadata := range(layerMetadata) {
		ft := WFSFeatureType{};
		ft.Init(metadata.LayerName, metadata.LayerTitle);
		result := ft.Stringify()
		if result.Error != nil {
			return result;
		}
		featureXMLString += "\n" + result.Result;
	}

	templateMap := map[string]string {
		"BaseURL": BASE_URL,
		"FeatureLayers": featureXMLString,
	}
	var buf bytes.Buffer
    if err := t.Execute(&buf, templateMap); err != nil {
        return Err[string](err);
    }

	return Ok(buf.String());
}

const DESCRIBE_SCHEMA_TEMPLATE = `<?xml version="1.0" encoding="UTF-8"?>
<xsd:schema
    targetNamespace="http://example.com/gowfs"
    xmlns:xsd="http://www.w3.org/2001/XMLSchema"
    xmlns:gml="http://www.opengis.net/gml/3.2"
    xmlns:gowfs="http://example.com/gowfs"
    elementFormDefault="qualified">
    %s
</xsd:schema>
`

func DescribeFeatureType(db *gorm.DB, layerName string) (string, error) {
  
  var columns []Column;
  var err error;

  if columns, err = GetLayerSchema(db, layerName); err != nil {
    return "", err
  }

  var complexType DescribeFeatureType_ComplexType = 
    DescribeFeatureType_CreateColumnSchema(layerName, columns);
  var complexTypeBytes []byte;

  if complexTypeBytes, err = xml.MarshalIndent(&complexType, "", "  "); err != nil {
    return "", err
  }

  return fmt.Sprintf(DESCRIBE_SCHEMA_TEMPLATE, string(complexTypeBytes)), nil
}


func GetFeature(db *gorm.DB, layerName string, query interface{}) (string, error) {
  dbFeatures, err := DBLayer_GetAllFeatures(db, layerName)
  if err != nil {
    return "", err
  }
  var pointElements []GetFeature_WFSMember;
  for _, feature := range(dbFeatures) {
    var attrs []GetFeature_ColumnTag;

    // Load attr
    for k, v := range(feature.Attr) {
      attrs = append(attrs, GetFeature_ColumnTag{
        Tag: k,
        Value: v,
      })
    }

    pointCoords := feature.Geom.Coords();

    pointElements = append(pointElements, GetFeature_CreatePointMember(
        layerName,
        uint32(feature.Fid),
        pointCoords[0],
        pointCoords[1],
        attrs,
      ))
  }

  getFeature := GetFeature_CreateFeatureCollection(len(pointElements), pointElements);
  getFeatureBytes, err := xml.MarshalIndent(getFeature, "", "  ");
  if err != nil { 
    return "", err
  }
  return string(getFeatureBytes), nil
}