package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/golang-jwt/jwt/v5"
	"path/filepath"
	"io"
	"strings"

)

func ParseSub(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(JWT_SECRET), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	// Since you always store `sub` as an int, just assert directly
	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, fmt.Errorf("sub claim not an int")
	}

	return int(sub), nil
}


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

	query := r.URL.Query()
	version := getAny(query, "version", "VERSION") // must be 2.0.0
	service := getAny(query, "service", "SERVICE") // Must be WFS
	request := getAny(query, "request", "REQUEST") // Could be GetFeature, GetCapabilities, DescribeFeatureType
	accessToken := getAny(query, "access_token", "ACCESS_TOKEN")

	if accessToken == "" {
		authHeader := r.Header.Get("Authorization")
		const prefix = "Bearer "
    	if !strings.HasPrefix(authHeader, prefix) {
        	http.Error(w, "invalid auth header", http.StatusUnauthorized)
        	return
    	}
    	accessToken = strings.TrimPrefix(authHeader, prefix)
	}


	userId, err := ParseSub(accessToken);
	if err != nil {
		http.Error(w, "Failed to extract JWT", http.StatusInternalServerError)
		return
	}


	if r.Method == "POST" && service == "WFS" {
		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		insertionRequests, err := XMLLayer_ParseInsertionRequest(string(body))
		if err != nil {
			http.Error(w, "failed to insert", http.StatusInternalServerError);
			return
		}

		fids, err := DBLayer_InsertLayer(db, insertionRequests); 
		if err != nil {
			http.Error(w, "failed to write to db", http.StatusInternalServerError);
			return
		}
		
		responseStr, err := XMLLayer_CreateInsertResponse(fids);
		if err != nil {
			http.Error(w, "Failed to create insert response", http.StatusInternalServerError);
		}

		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(responseStr))


		return
	}

	if err := validateInitialQuery(version, service, request); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(OWSError { Error: err.Error() })
		return
	}

	if request == "GetCapabilities" {
		xmlstring, err := GetCapabilities(db, userId).Maybe()
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
		outputFormat := getAny(query, "outputformat", "OUTPUTFORMAT", "outputFormat")
		if outputFormat == "application/json" || outputFormat == "json" {
			jsonstr, err := GetFeatureGeoJSON(db, layerName, nil)
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError);
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(jsonstr))
		} else if outputFormat == "" {
			xmlstring, err := GetFeature(db, layerName, nil);
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError);
			}
			w.Header().Set("Content-Type", "application/xml");
			w.Write([]byte(xmlstring));
		}
		
	}

}

func main() {

	fs := http.FileServer(http.Dir("public"))

	db := Ensure(gorm.Open(postgres.Open(CONNECTION_STRING), &gorm.Config{}));
	db.AutoMigrate(& LayerMetadata{});
	sqlDB := Ensure(db.DB());
	defer sqlDB.Close()

	// Handle "/" by serving index.html explicitly
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join("public", "index.html"))
			return
		}
		fs.ServeHTTP(w, r)
	})
	
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