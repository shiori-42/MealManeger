package main

import (
	// "encoding/json"
	// "os"
	"fmt"
	"log"
	"net/http"
	// "golang.org/x/api/sheets/v4"
)

type uploadData struct {
	ImageBase64 string `json:"imageBase64"`
	Calories    int    `json:"calories"`
	UserId      string `json:"userid"`
	Date        string `json:"date"`
	MealType    string `json:"mealType"`
	Comment     string `json:"comment"`
}

// func getUserDataHandler(w http.ResponseWriter,r *http.Request){
// 	userId:=r.URL.Query().Get("userid")
// 	if userId==""{
// 		http.Error(w,"User ID is required",http.StatusBadRequest)
// 		return
// 	}
// 	ctx:=context.Background()
// 	creds,err:=google.FindDegaultCredentials(ctx,sheets.SpreadsheetsScope)
// 	if err!=nil{
// 		log.Fatalf("Unable to create sheets service:%v",err)
// 	}

//var resonse

// }

func main() {
	http.HandleFunc("/upload", uploadHandler)
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
