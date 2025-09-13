package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type OWSResponse struct {
	Message string `json:"message"`
}

type OWSError struct {
	Error string `json:"error"`
}

func getAny(q url.Values, keys ...string) string {
    for _, k := range keys {
        if v := q.Get(k); v != "" {
            return v
        }
    }
    return ""
}

func validateInitialQuery(version string, service string, request string) (error) {
	if service != "WFS" {
		return fmt.Errorf("service must be WFS")
	}
	if request != "DescribeFeatureType" && request != "GetCapabilities" && request != "GetFeature" {
		return fmt.Errorf("invalid request")
	}
	return nil
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}


func owsHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	version := getAny(query, "version", "VERSION") // must be 2.0.0
	service := getAny(query, "service", "SERVICE") // Must be WFS
	request := getAny(query, "request", "REQUEST") // Could be GetFeature, GetCapabilities, DescribeFeatureType


	if err := validateInitialQuery(version, service, request); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(OWSError { Error: err.Error() })
		return
	}

	if request == "GetCapabilities" {
		xmlstring, err := GetCapabilities(db).Maybe()
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(xmlstring))
		return 
	} else if request == "DescribeFeatureType" {
		layerName := getAny(query, "typeName", "TYPENAME", "typeNames", "TYPENAMES", "typename")
		var xmlstring string;
		var err error;
		if xmlstring, err = DescribeFeatureType(db, layerName); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError);
		}
		w.Header().Set("Content-Type","application/xml")
		w.Write([]byte(xmlstring));
	} else if request == "GetFeature" {
		layerName := getAny(query, "typeName", "TYPENAME", "typeNames", "TYPENAMES", "typename")
		xmlstring, err := GetFeature(db, layerName, nil);
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError);
		}
		w.Header().Set("Content-Type", "application/xml");
		w.Write([]byte(xmlstring));
	}

}

func main() {
	db := Ensure(gorm.Open(postgres.Open(CONNECTION_STRING), &gorm.Config{}));
	db.AutoMigrate(& LayerMetadata{});
	sqlDB := Ensure(db.DB());
	defer sqlDB.Close()


	
	http.HandleFunc("/ows", func (w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s - %s\n", r.Method, r.URL.String());
		owsHandler(w, r, db);
	})
	http.HandleFunc("/healthcheck", healthCheckHandler)

	log.Println("listening on 8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}