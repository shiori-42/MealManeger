package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
    // リクエストがPOSTかどうか確認
    if r.Method != "POST" {
        http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }

    // ファイルと他のフォームデータを解析
    r.ParseMultipartForm(10 << 20) // 10 MBのファイルサイズ制限

    // フォームから "image" というキーでファイルを取得
    file, _, err := r.FormFile("image")
    if err != nil {
        http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
        fmt.Println("Error Retrieving the File")
        fmt.Println(err)
        return
    }
    defer file.Close()

    // ファイルをサーバ上に保存
    dst, err := os.Create("/tmp/uploaded_image.jpg") // 適切なファイルパスを指定
    if err != nil {
        http.Error(w, "Error saving the file", http.StatusInternalServerError)
        fmt.Println("Error Saving the File")
        fmt.Println(err)
        return
    }
    defer dst.Close()

    // ファイルをコピー
    _, err = io.Copy(dst, file)
    if err != nil {
        http.Error(w, "Error saving the file", http.StatusInternalServerError)
        fmt.Println("Error Copying the File")
        fmt.Println(err)
        return
    }

    fmt.Fprintf(w, "File uploaded successfully")
}

func main() {
    http.HandleFunc("/upload", uploadFile)

    fmt.Println("Server started on :8080")
    http.ListenAndServe(":8080", nil)
}
