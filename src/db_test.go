package main

import (
	"log"
	"reflect"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCreateTableString(t *testing.T) {
	create_layer_table_params := CreateLayerTable {
		LayerName: "user_location",
		Columns: []ColumnType { 
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

	actual := createTableStatement(create_layer_table_params).Unwrap()
	expected := `CREATE TABLE user_location (fid SERIAL PRIMARY KEY, geom GEOMETRY(Point, 4326), name TEXT, address TEXT);`
	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual);
	}
}


func TestCreateTable(t *testing.T) {
	create_layer_table_params := CreateLayerTable {
		LayerName: "user_location",
		LayerTitle: "User Location",
		Columns: []ColumnType { 
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

	db := Ensure(gorm.Open(postgres.Open(CONNECTION_STRING), &gorm.Config{}));
	db.AutoMigrate(& LayerMetadata{});

	sqlDB := Ensure(db.DB());
	defer sqlDB.Close()
	
	layername := CreateLayer(db, create_layer_table_params).Unwrap();

	if layername != "user_location" {
		t.Errorf("create layer failed to return proper tablename")
	}

	if err := DeleteLayer(db, "user_location"); err != nil {
		t.Errorf("%s", err.Error())
	}
}

func TestGetTableName(t *testing.T) {
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

	db, err := gorm.Open(postgres.Open(CONNECTION_STRING), &gorm.Config{});
	if err != nil {
		t.Errorf("%s", err.Error());
	}
	db.AutoMigrate(& LayerMetadata{});

	sqlDB, err := db.DB();
	if err != nil {
		log.Fatal(err);
	}

	defer sqlDB.Close()

	CreateLayer(db, tableName1).Unwrap();
	defer DeleteLayer(db, "layer_1");
	
	CreateLayer(db, tableName2).Unwrap();
	defer DeleteLayer(db, "layer_2")

	
	CreateLayer(db, tableName3).Unwrap();
	defer DeleteLayer(db, "layer_3")


	expected_value := []string{ "layer_1", "layer_2", "layer_3" }
	actual_value, err := GetLayerNames(db);
	if err != nil {
		t.Errorf("failed to return layer names");
	}
	if !reflect.DeepEqual(expected_value, actual_value) {
		t.Errorf("expected error names does not match actual layer names");
	}
}


func TestGetLayerMetadata(t *testing.T) {
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

	db, err := gorm.Open(postgres.Open(CONNECTION_STRING), &gorm.Config{});
	if err != nil {
		t.Errorf("%s", err.Error());
	}
	db.AutoMigrate(& LayerMetadata{});

	sqlDB, err := db.DB();
	if err != nil {
		log.Fatal(err);
	}

	defer sqlDB.Close()

	CreateLayer(db, tableName1).Unwrap();
	defer DeleteLayer(db, "layer_1");
	
	CreateLayer(db, tableName2).Unwrap();
	defer DeleteLayer(db, "layer_2")

	
	CreateLayer(db, tableName3).Unwrap();
	defer DeleteLayer(db, "layer_3")


	actualLayerMetadata := Ensure(GetLayerMetadata(db));
	expectedLayerMetadata := []LayerMetadata{
		{
			LayerName: "layer_1",
			LayerTitle: "Layer 1",
		},
		{
			LayerName: "layer_2",
			LayerTitle: "Layer 2",
		},
		{
			LayerName: "layer_3",
			LayerTitle: "Layer 3",
		},
	}

	// We can't ensure that ID starts counting at 0 so we can't use deep reflect
	for i := 0; i < 3; i++ {
		actual := actualLayerMetadata[i]
		expected := expectedLayerMetadata[i]
		layerEQ := actual.LayerName == expected.LayerName
		titleEQ := actual.LayerTitle == expected.LayerTitle
		if ! layerEQ && titleEQ {
			t.Errorf("%+v != %+v", actual, expected)
		}
	}
}

func TestGetLayerSchema(t *testing.T) {
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
	db, err := gorm.Open(postgres.Open(CONNECTION_STRING), &gorm.Config{});
	if err != nil {
		t.Errorf("%s", err.Error());
	}
	db.AutoMigrate(& LayerMetadata{});

	sqlDB, err := db.DB();
	if err != nil {
		log.Fatal(err);
	}

	defer sqlDB.Close()
	CreateLayer(db, tableName1).Unwrap();

	defer DeleteLayer(db, "layer_1");

	actual_columns, err := GetLayerSchema(db, "layer_1");
	if err != nil {
		t.Errorf("%s", err.Error());
	}	

	var expected_columns = []Column{
		{
			ColumnName: "name",
			DataType: "text",
			IsNullable: "YES",
		},
		{
			ColumnName: "address",
			DataType: "text",
			IsNullable: "YES",
		},
	};

	if !reflect.DeepEqual(expected_columns, actual_columns) {
		t.Errorf("expected columns don't match actual columns");
	}
}

func TestGetAllFeatures(t *testing.T) {
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
		(ST_GeomFromText('POINT(12.123 42.789)', 4326), 'Kevin',   '4510 Laclede Ave'),
		(ST_GeomFromText('POINT(14.128 12.238)', 4326), 'Lawrance','2145 Tenbrink'),
		(ST_GeomFromText('POINT(28.479 94.177)', 4326), 'Peter',   '23 Olde Hamlet Dr.');
	`

	if db.Exec(sqlInsert).Error != nil {
		t.Errorf("Failed to insert")
	}
	
	features := Ensure(DBLayer_GetAllFeatures(db, "layer_1"))
	
	c1 := features[0].Attr["name"] == "Kevin"
	c2 := features[1].Attr["name"] == "Lawrance"
	c3 := features[2].Attr["name"] == "Peter"

	if ! (c1 && c2 && c3) {
		t.Errorf("Features don't match")
	}

	expectedCoord := [][]float64{
		{
			12.123,
			42.789,
		},
		{
			14.128,
			12.238,
		},
		{
			28.479,
			94.177,
		},
	}
	for i, row := range(features) {
		point := row.Geom.Coords();
		expectedPoint := expectedCoord[i];
		if ! (point[0] == expectedPoint[0] && point[1] == expectedPoint[1]) {
			t.Errorf("Points don't match")
		}
	}

}



func TestGetLayerMetadataByUser(t *testing.T) {
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

	db, err := gorm.Open(postgres.Open(CONNECTION_STRING), &gorm.Config{});
	if err != nil {
		t.Errorf("%s", err.Error());
	}
	db.AutoMigrate(& LayerMetadata{});

	sqlDB, err := db.DB();
	if err != nil {
		log.Fatal(err);
	}

	defer sqlDB.Close()

	CreateLayerByUser(db, tableName1, 1).Unwrap();
	defer DeleteLayer(db, "layer_1");
	
	CreateLayerByUser(db, tableName2, 2).Unwrap();
	defer DeleteLayer(db, "layer_2")

	
	CreateLayerByUser(db, tableName3, 1).Unwrap();
	defer DeleteLayer(db, "layer_3")


	actualLayerMetadata := Ensure(GetLayerMetadataByUser(db, 1));
	expectedLayerMetadata := []LayerMetadata{
		{
			LayerName: "layer_1",
			LayerTitle: "Layer 1",
		},
		{
			LayerName: "layer_3",
			LayerTitle: "Layer 3",
		},
	}

	if len(actualLayerMetadata) != 2 {
		t.Errorf("should only see 2 layers for User 1")
	}

	// We can't ensure that ID starts counting at 0 so we can't use deep reflect
	for i := 0; i < 2; i++ {
		actual := actualLayerMetadata[i]
		expected := expectedLayerMetadata[i]
		layerEQ := actual.LayerName == expected.LayerName
		titleEQ := actual.LayerTitle == expected.LayerTitle
		if ! layerEQ && titleEQ {
			t.Errorf("%+v != %+v", actual, expected)
		}
	}

	expectedUser2Metadata := LayerMetadata{
			LayerName: "layer_2",
	}

	actualUser2LayerMetadata := Ensure(GetLayerMetadataByUser(db, 2))[0];

	if expectedUser2Metadata.LayerName != actualUser2LayerMetadata.LayerName {
		t.Errorf("User 2 metadata does not match!, %s != %s\n", expectedUser2Metadata.LayerName, actualUser2LayerMetadata.LayerName)
	}

}