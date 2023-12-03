package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/nfnt/resize"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

type uploadData struct {
	Image    string `json:"image"`
	Calories int    `json:"calories"`
	UserId   string `json:"userid"`
	Date     string // このフィールドはJSONから受け取らず、サーバー側で設定
}

func writeToSheet(data uploadData) error {
	ctx := context.Background()
	credsFilePath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credsFilePath == "" {
		log.Fatal("The GOOGLE_APPLICATION_CREDENTIALS environment variable is not set.")
	}

	creds, err := os.ReadFile(credsFilePath)
	if err != nil {
		return fmt.Errorf("Unable to read service account file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(creds, sheets.SpreadsheetsScope)
	if err != nil {
		return fmt.Errorf("Unable to parse service account file to config: %v", err)
	}

	client := config.Client(ctx)

	sheetsService, err := sheets.New(client)
	if err != nil {
		return fmt.Errorf("Unable to create sheets Service: %v", err)
	}

	spreadsheetId := os.Getenv("SPREADSHEET_ID")
	if spreadsheetId == "" {
		log.Fatal("The SPREADSHEET_ID environment variable must be set.")
	}

	rangeData := "Sheet1!A2:E"
	var values [][]interface{}
	row := []interface{}{
		data.UserId,
		data.Image,
		data.Calories,
		data.Date,
	}
	values = append(values, row)

	valueRange := &sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         values,
	}

	_, err = sheetsService.Spreadsheets.Values.Append(spreadsheetId, rangeData, valueRange).ValueInputOption("USER_ENTERED").Do()
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

	var data uploadData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if data.Calories <= 0 {
		http.Error(w, "Invalid calories value", http.StatusBadRequest)
		return
	}

	if data.UserId == "" {
		http.Error(w, "Invalid UserID", http.StatusBadRequest)
		return
	}

	// 現在の日付をYYYY-MM-DD形式で取得
	currentDate := time.Now().Format("2006-01-02")
	data.Date = currentDate // uploadData構造体に日付を設定

	// Resize Image
	resizedImage, err := resizeImageBase64(data.Image)
	if err != nil {
		http.Error(w, "Error resizing image", http.StatusInternalServerError)
		return
	}
	data.Image = resizedImage

	err = writeToSheet(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Data uploaded successfully")
}

func resizeImageBase64(base64Image string) (string, error) {
	imgData, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return "", fmt.Errorf("error decoding base64: %v", err)
	}

	resizedData, err := resizeImage(imgData, 800, 600)
	if err != nil {
		return "", fmt.Errorf("error resizing image: %v", err)
	}

	return base64.StdEncoding.EncodeToString(resizedData), nil
}

func resizeImage(imgData []byte, maxWidth, maxHeight uint) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, err
	}

	resizedImg := resize.Thumbnail(maxWidth, maxHeight, img, resize.Lanczos3)

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, resizedImg, nil); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func getDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Unsupported Method", http.StatusMethodNotAllowed)
		return
	}

	// クエリパラメータからユーザーIDを取得
	userID := r.URL.Query().Get("userid")
	if userID == "" {
		http.Error(w, "UserID is required", http.StatusBadRequest)
		return
	}

	data, err := readFromSheet(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 取得したデータをJSON形式で返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func readFromSheet(userID string) ([]uploadData, error) {
	ctx := context.Background()
	credsFilePath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credsFilePath == "" {
		return nil, fmt.Errorf("The GOOGLE_APPLICATION_CREDENTIALS environment variable is not set")
	}

	creds, err := os.ReadFile(credsFilePath)
	if err != nil {
		return nil, fmt.Errorf("Unable to read service account file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(creds, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse service account file to config: %v", err)
	}

	client := config.Client(ctx)

	sheetsService, err := sheets.New(client)
	if err != nil {
		return nil, fmt.Errorf("Unable to create sheets Service: %v", err)
	}

	spreadsheetId := os.Getenv("SPREADSHEET_ID")
	if spreadsheetId == "" {
		return nil, fmt.Errorf("The SPREADSHEET_ID environment variable must be set")
	}

	readRange := "Sheet1!A2:E" // スプレッドシートの読み取り範囲
	resp, err := sheetsService.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve data from sheet: %v", err)
	}

	var uploads []uploadData
	for _, row := range resp.Values {
		// スプレッドシートの各行を読み取り、uploadData構造体に変換
		if len(row) > 3 && row[0].(string) == userID {
			calories, _ := strconv.Atoi(row[2].(string)) // エラーハンドリングが必要
			upload := uploadData{
				UserId:   row[0].(string),
				Image:    row[1].(string),
				Calories: calories,
				Date:     row[3].(string),
			}
			uploads = append(uploads, upload)
		}
	}
	return uploads, nil
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/get", getDataHandler)

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
