package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

var DB *gorm.DB

type Guest struct {
	gorm.Model
	FirstName       string `json:"firstName"`
	LastName        string `json:"LastName"`
	Email           string `json:"email"`
	CreateBatchSize int
}

func create(w http.ResponseWriter, req *http.Request) {
	// Read from the request
	var guest Guest
	decoder := json.NewDecoder(req.Body) // read the request
	err := decoder.Decode(&guest)
	if err != nil {
		fmt.Printf("decoding body failed with error %s", err)

	}

	//TODO: the guest variable now contains name, email, etc. > create a database entry for the guest variable via
	DB.Create(&guest)
	result := DB.Create(&guest)
	if result.Error != nil {
		fmt.Printf("Error has confirmed")
	}

	//gorm and the CREATE method
	//be aware: the DB var is a globally defined variable (see line 14), so you can use it here like that: DB.Create(...)

	//alternative storage: file
	/*

		f, err := os.OpenFile("UserList.csv", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) // methode to open and edit it
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()                                                                           //closing the file and save it
		_, err = f.WriteString(guest.FirstName + "" + guest.LastName + "  " + guest.Email + "\n") //writing the File
		if err != nil {
			log.Fatal(err)
		}
		//data := guest.Name + " " + guest.Email
		//storing the guest info in the Database

	*/

	w.WriteHeader(http.StatusOK)
	//TODO: send a response message back "registration succeeded"
}

func registration(w http.ResponseWriter) {

}

/*
func list(w http.ResponseWriter, req *http.Request) {
	var guests []Guest
	db.Find(guests)
	data, err := json.Marshal(guests)
	if err != nil {
		fmt.Println("crazy error happened", err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

*/

func main() {
	//setup a database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PORT"))
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to Database: %s", err)

	}

	DB.AutoMigrate(&Guest{})

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/registration", create).Methods(http.MethodPost)
	//router.HandleFunc("/guest", list).Methods(http.MethodGet)
	err = http.ListenAndServe(":8090", router)
	if err != nil {
		return
	}
}
