package main

//Try Scan
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gorilla/mux"
)

type Item struct {
	Name        string `json:"Name"`
	Empty_slots int    `json:"empty_slots"`
	Free_bikes  int    `json:"free_bikes"`
}

type Status struct {
	Table       string `json:"table"`
	Recordcount int    `json:"recordCount"`
}

var status []Status

//Get All Table Items
func getTableItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	svc := dynamodb.New(sess)

	params := &dynamodb.ScanInput{
		TableName: aws.String("CityBike"),
	}

	result, err := svc.Scan(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	obj := []Item{}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &obj)
	if err != nil {
		fmt.Printf("failed to unmarshal Query result items, %v", err)
	}

	json.NewEncoder(w).Encode(obj)
}

func getTableInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func main() {
	os.Setenv("AWS_ACCESS_KEY_ID", "AKIA34XNLPJYLZGHNQGI")
	os.Setenv("AWS_SECRET_KEY", "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")

	//Init Router
	r := mux.NewRouter()

	status = append(status, Status{Table: "CityBike", Recordcount: 1350})

	//Route Handlers / Endpoints
	r.HandleFunc("/aadejo/all", getTableItems).Methods("GET")
	r.HandleFunc("/aadejo/status", getTableInfo).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
