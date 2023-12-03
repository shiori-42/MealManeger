import base64
import requests
import json
import time

# 画像ファイルをBase64エンコードする
with open("hello.jpg", "rb") as image_file:
    encoded_string = base64.b64encode(image_file.read()).decode('utf-8')

# # エンコードされた画像データをJSONオブジェクトに格納
# # data = json.dumps({"image": encoded_string})

data = json.dumps({"image": encoded_string, "calories": 123,
                  "userid": "userid"})

# # APIエンドポイント
url = "http://localhost:8080"

# # POSTリクエストを送信
response = requests.post(url+"/upload", data=data, headers={
                         "Content-Type": "application/json"})


# response=requests.get(url+"/get?userid=userid")
# レスポンスを表示
# print(len(response.json()))
print(response.text)