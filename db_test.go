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
