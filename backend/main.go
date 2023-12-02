package main

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

// JSONリクエストを解析するための構造体
type ImageRequest struct {
    Image string `json:"image"`
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
    var imgReq ImageRequest

    // JSONリクエストボディを解析
    if err := json.NewDecoder(r.Body).Decode(&imgReq); err != nil {
        http.Error(w, "Error decoding JSON", http.StatusBadRequest)
        return
    }

    // Base64データをデコード
    imgData, err := base64.StdEncoding.DecodeString(imgReq.Image)
    if err != nil {
        http.Error(w, "Error decoding base64", http.StatusBadRequest)
        return
    }

    // デコードされたデータをファイルに書き込む
    filePath := "./success.jpg" // 適切なファイルパスを指定
    if err := ioutil.WriteFile(filePath, imgData, 0644); err != nil {
        http.Error(w, "Error saving the file", http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "File uploaded successfully")
}

func main() {
    http.HandleFunc("/upload", uploadFile)

    fmt.Println("Server started on :8080")
    http.ListenAndServe(":8080", nil)
}
