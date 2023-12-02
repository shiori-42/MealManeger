package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

type uploadData struct {
	ImageURL string `json:"imageUrl"`
	Calories int    `json:"calories"`
	UserId   string `json:"userid"`
	Date     string `json:"date"`
	MealType string `json:"mealType"`
	Comment  string `json:"comment"`
}

func writeToSheet(data uploadData) error {
	// Google Sheets APIクライアントの初期化
	ctx := context.Background()

	//サービスアカウントキーをファイルから読み込む
	credsFilePath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credsFilePath == "" {
		log.Fatal("The GOOGLE_APPLICATION_CREDENTIALS environment variable is not set.")
	}

	creds, err := os.ReadFile(credsFilePath)
	if err != nil {
		return fmt.Errorf("Unable to read service account file:%v", err)
	}

	//JWT認証情報を取得
	config, err := google.JWTConfigFromJSON(creds, sheets.SpreadsheetsScope)
	if err != nil {
		return fmt.Errorf("Unable to parse service account file to config:%v", err)
	}

	//認証済みクライアントを作成
	client := config.Client(ctx)

	//Google Sheetsサービスを作成
	sheetsService, err := sheets.New(client)
	if err != nil {
		return fmt.Errorf("Unable to create sheets Service:%v", err)
	}

	spreadsheetId := os.Getenv("SPREADSHEET_ID")
	if spreadsheetId == "" {
		log.Fatal("The SPREADSHEET_ID environment variables must be set.")
	}
	rangeData := "Sheet1!A2:F2"
	var values [][]interface{} // 書き込むデータの配列を初期化
	row := []interface{}{
		data.UserId,
		data.ImageURL,
		data.Calories,
		data.Date,
		data.MealType,
		data.Comment,
	}
	values = append(values, row)

	//スプレッドシートへの追加
	valueRange := &sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         values,
	}

	_, err = (sheetsService.Spreadsheets.Values.Append(spreadsheetId, rangeData, valueRange).ValueInputOption("USER_ENTERED").Do())
	if err != nil {
		return err
	}
	return nil
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Unsupported Method", http.StatusMethodNotAllowed)
		return
	}

	//リクエストボディ解析
	var data uploadData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// if _,err:=url.ParseRequestURI(data.ImageURL);
	// err!=nil{
	// 	http.Error(w,"Invalid ImageURL",http.StatusBadRequest)
	// 	return
	// }

	if data.Calories <= 0 {
		http.Error(w, "Invalid calories value", http.StatusBadRequest)
		return
	}

	if data.UserId == "" {
		http.Error(w, "Invalid UserID", http.StatusBadRequest)
		return
	}

	err = writeToSheet(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}