package main

import (
	"fmt"
	"gorm.io/gorm"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"encoding/hex"
)

type ColumnType struct {
	Name string `json:"name"`
	Dtype string `json:"dtype"`
}

type CreateLayerTable struct {
	LayerName string `json:"layername"`
	LayerTitle string `json:"layertitle"`
	Columns []ColumnType
}


func sql_string_is_clean(input string) (error) {
	// not implemented. I'm on a deadline
	return nil;
}

/* 
We need to dynamically generate the create SQL query because Geospatial layer schema are defined by the user. 
So we can't have a fixed schema.
*/
func createTableStatement(params CreateLayerTable) Result[string] {
	// This function needs sanitization	
	
	if err := sql_string_is_clean(params.LayerName); err != nil {
		return Err[string](fmt.Errorf("bad LayerName String"))
	}
	// We're only gonna deal with Points for this implementation, but ostensibly we can also add Lines, Polygon, etc.
	sql_statement := fmt.Sprintf(`CREATE TABLE %s (fid SERIAL PRIMARY KEY, geom GEOMETRY(Point, 4326)`, params.LayerName)

	for  _, columntype := range params.Columns {
		if err := sql_string_is_clean(columntype.Name); err != nil {
			return Err[string](fmt.Errorf("bad column name"));
		}
		if err := sql_string_is_clean(columntype.Dtype); err != nil {
			return Err[string](fmt.Errorf("bad column dtype"));
		}
		sql_statement += fmt.Sprintf(`, %s %s`, columntype.Name, columntype.Dtype)
	}

	sql_statement += `);`

	return Ok(sql_statement);
}


// Since I'm only implementing POINT, I don't have to keep track of EPSG and Geom Type
type LayerMetadata struct {
	Id uint `gorm:"primaryKey"`
	LayerName string `gorm:"unique"`
	LayerTitle string `gorm:"type:text"`
	// GeomType string 
	// EPSG uint
}

func CreateLayer(db *gorm.DB, params CreateLayerTable) Result[string] {
	
	sql_statement, err := createTableStatement(params).Maybe()
	if err != nil {
		return Err[string](err)
	}
	
	if err := db.Exec(sql_statement).Error; err != nil {
		return Err[string](err)
	}

	if err := db.Create(& LayerMetadata{ LayerName: params.LayerName, LayerTitle: params.LayerTitle }).Error; err != nil {
		return Err[string](err)
	}

	return Ok(params.LayerName)
}

func DeleteLayer(db *gorm.DB, layername string) error {
	err := db.Where("layer_name = ?", layername).
		Delete(&LayerMetadata{}).
		Error;
	if err != nil {
		fmt.Println(err)
		return err;
	}

	return db.Exec(fmt.Sprintf("DROP TABLE %s;", layername)).Error;
}

func GetLayerMetadata(db *gorm.DB) ([]LayerMetadata, error){
	var layers []LayerMetadata;
	err := db.Find(&layers).Error;
	return layers, err;
}

func GetLayerNames(db *gorm.DB) ([]string, error) {
	var layerNames []string;
	if err := db.Model(&LayerMetadata{}).Pluck("layer_name", &layerNames).Error; err != nil {
		return []string{}, err
	}
	return layerNames, nil
}


type Column struct {
	ColumnName string
	DataType string
	IsNullable string
}

func GetLayerSchema(db *gorm.DB, layerName string) ([]Column, error) {
	// We're gonna exclude the geometry column because PSQL returns it as 'USER-DEFINED' instead of GEOMETRY(Point, 4326)

	var cols []Column;
	if err := db.Raw(
		`SELECT column_name, data_type, is_nullable 
		FROM information_schema.columns
			WHERE table_schema = 'public'
			AND table_name = ?
			AND column_name != 'geom'
			AND column_name != 'fid'
		ORDER BY ordinal_position;`, layerName,
	).Scan(&cols).Error; err != nil {
		return []Column{}, err
	}

	return cols, nil;
}

func GetLayerGeometry(_ *gorm.DB) (string, uint, error) {
	// Leave room for expansion layter
	return "Point", 4326, nil
}

type FeatureRow struct {
	Fid int32
	Geom geom.Point
	Attr map[string]string;
}

func DBLayer_GetAllFeatures(db *gorm.DB, layername string) ([]FeatureRow, error) {
	var rows []map[string]interface{};
	err := db.Table(layername).Find(&rows).Error
	if err != nil {
		return []FeatureRow{}, err;
	}

	features := make([]FeatureRow, 0, len(rows));
	for _, row := range(rows) {

		hexStr := row["geom"].(string)
		pointBytes, err := hex.DecodeString(hexStr)
		if err != nil {
			return []FeatureRow{}, err
		}

		point, err := ewkb.Unmarshal(pointBytes)
		if err != nil {
			return []FeatureRow{}, err
		}	

		f := FeatureRow{
			Fid: row["fid"].(int32),
			Geom: *(point.(*geom.Point)),
			Attr: make(map[string]string),

		}
		for k, v := range(row) {
			if k != "fid" && k != "geom" {
				f.Attr[k] = v.(string)
			}
		}
		features = append(features, f)
	}

	return features, nil
}
