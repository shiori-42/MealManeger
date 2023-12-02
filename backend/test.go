package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"net/http"

	"github.com/nfnt/resize"
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

    // 初期サイズ設定
    maxWidth, maxHeight := uint(800), uint(600)

    var resizedData []byte
    var resizedBase64 string

    // 画像を縮小し、Base64エンコードされたサイズが50,000文字以下になるまでループ
    for {
        // 画像のサイズを調整
        resizedData, err = resizeImage(imgData, maxWidth, maxHeight)
        if err != nil {
            http.Error(w, "Error resizing image", http.StatusInternalServerError)
            return
        }

        // 調整された画像データを再度Base64エンコード
        resizedBase64 = base64.StdEncoding.EncodeToString(resizedData)

        // サイズが50,000文字以下かチェック
        if len(resizedBase64) <= 50000 {
            break
        }

        // サイズをさらに小さくする
        maxWidth, maxHeight = maxWidth*9/10, maxHeight*9/10
    }

    // ファイルに書き込む
    filePath := "./success.jpg"
    if err := ioutil.WriteFile(filePath, resizedData, 0644); err != nil {
        http.Error(w, "Error saving the file", http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "File uploaded successfully")
}

func resizeImage(imgData []byte, maxWidth, maxHeight uint) ([]byte, error) {
	// デコードされた画像データから画像を作成
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, err
	}

	// 画像のサイズを調整
	resizedImg := resize.Thumbnail(maxWidth, maxHeight, img, resize.Lanczos3)

	// 新しいバッファを作成して、そこにJPEG形式で画像をエンコード
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, resizedImg, nil); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func main() {
	http.HandleFunc("/upload", uploadFile)

	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
