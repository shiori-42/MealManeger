import base64
import requests
import json

# 画像ファイルをBase64エンコードする
with open("hello.jpg", "rb") as image_file:
    encoded_string = base64.b64encode(image_file.read()).decode('utf-8')

# エンコードされた画像データをJSONオブジェクトに格納
data = json.dumps({"image": encoded_string})

# APIエンドポイント
url = "http://localhost:8080/upload"

# POSTリクエストを送信
response = requests.post(url, data=data, headers={"Content-Type": "application/json"})

# レスポンスを表示
print(response.text)
