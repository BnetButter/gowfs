package main;

import (
	"encoding/xml"
	"strings"
	"strconv"
)

type Transaction struct {
	XMLName xml.Name `xml:"Transaction"`
	Inserts []Insert `xml:"Insert"`
}

type Insert struct {
	Layer DynamicElement `xml:",any"`
}

type DynamicElement struct {
	XMLName xml.Name
	Fields  map[string]string
	Geom    InsertGeometry `xml:"geom"`
}

type InsertGeometry struct {
	XMLName xml.Name   `xml:"geom"`
	Point   GmlPoint   `xml:"Point"`
}

type GmlPoint struct {
	XMLName     xml.Name     `xml:"Point"`
	SrsName     string       `xml:"srsName,attr"`
	Coordinates GmlCoords    `xml:"coordinates"`
}

type GmlCoords struct {
	XMLName xml.Name `xml:"coordinates"`
	Ts      string   `xml:"ts,attr"`
	Cs      string   `xml:"cs,attr"`
	Value   string   `xml:",chardata"`
}


func (d *DynamicElement) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	d.XMLName = start.Name
	d.Fields = make(map[string]string)
	
	for {
		tok, err := dec.Token()
		if err != nil {
			return err
		}

		switch elem := tok.(type) {
		case xml.StartElement:
			if elem.Name.Local == "geom" {
				if err := dec.DecodeElement(&d.Geom, &elem); err != nil {
					return err
				}
			} else {
				var content string
				if err := dec.DecodeElement(&content, &elem); err != nil {
					return err
				}
				d.Fields[elem.Name.Local] = strings.TrimSpace(content)
			}

		case xml.EndElement:
			if elem.Name.Local == d.XMLName.Local {
				return nil
			}
		}
	}
}


type InsertRequestParams struct {
	LayerName string
	Coordinates [2]float64
	Fields map[string]string
}

func XMLLayer_ParseInsertionRequest(insertionXML string) ([]InsertRequestParams, error) {
	var tx Transaction;
	if err := xml.Unmarshal([]byte(insertionXML), &tx); err != nil {
		return []InsertRequestParams{}, err;
	}

	var insertionRequests []InsertRequestParams;

	for _, ins := range tx.Inserts {
		
		parts := strings.Split(ins.Layer.Geom.Point.Coordinates.Value, ",")
		X, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return []InsertRequestParams{},err
		}

		Y, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return []InsertRequestParams{}, err
		}



		insertionRequest := InsertRequestParams{
			LayerName: ins.Layer.XMLName.Local,
			Fields: ins.Layer.Fields,
			Coordinates: [2]float64{
				X, 
				Y,
			},
		}
		insertionRequests = append(insertionRequests, insertionRequest)
	}
	return insertionRequests, nil
}


func main() {
	const INSERTION_XML = `<Transaction xmlns="http://www.opengis.net/wfs" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:gml="http://www.opengis.net/gml" xsi:schemaLocation="http://example.com/gowfs http://localhost:8000/ows?access_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOjF9.n0_TCv6aixmt4LzFzEf18kB7ivf2SPN2SEaAdHuoAYU&amp;SERVICE=WFS&amp;REQUEST=DescribeFeatureType&amp;VERSION=1.0.0&amp;TYPENAME=parks_ny" version="1.0.0" service="WFS">
   <Insert xmlns="http://www.opengis.net/wfs">
       <parks_ny xmlns="http://example.com/gowfs">
           <park_name xmlns="http://example.com/gowfs">f</park_name>
           <size_acres xmlns="http://example.com/gowfs">d</size_acres>
           <geom xmlns="http://example.com/gowfs">
               <gml:Point>
                   <gml:coordinates ts=" " cs=",">-74.21244970386565853,40.89557153692449987</gml:coordinates>
               </gml:Point>
           </geom>
       </parks_ny>
   </Insert>
   <Insert xmlns="http://www.opengis.net/wfs">
       <parks_mo xmlns="http://example.com/gowfs">
           <park_mo xmlns="http://example.com/gowfs">12</park_mo>
           <size_mo xmlns="http://example.com/gowfs">12</size_mo>
           <geom xmlns="http://example.com/gowfs">
               <gml:Point>
                   <gml:coordinates ts=" " cs=",">-74.24520207755853107,40.83474570006630699</gml:coordinates>
               </gml:Point>
           </geom>
       </parks_mo>
   </Insert>
</Transaction>`

	InsertionRequests, err := XMLLayer_ParseInsertionRequest(INSERTION_XML);
	if err != nil {
		panic("Failed to parse")
	}
	I0 := InsertionRequests[0]
	I1 := InsertionRequests[1]

	if I0.LayerName != "parks_ny" && I1.LayerName != "parks_mo" {
		panic("Wrong Layer")
	}

	conditions := I0.Fields["park_name"] == "f" && I0.Fields["size_acres"] == "d"
	if ! conditions {
		panic("Bad condition")
	}

	if I1.Coordinates[0] != -74.24520207755853107 && I1.Coordinates[1] != 40.83474570006630699 {
		panic("Bad coords")
	}

}