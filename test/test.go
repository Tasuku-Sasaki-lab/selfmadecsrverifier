package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

type Values struct {
	Cn     string `json: "Cn"`
	Secret string `json: "Secret"`
}

func testVerrify() (bool, error) {
	CN := "CN=scepclient,OU=MDM,O=scep-client,C=US"
	challenge := "pass"

	// build sample structure
	values := new(Values)
	values.Cn = CN
	values.Secret = challenge

	// encode json
	values_json, _ := json.Marshal(values)
	fmt.Printf("[+] %s\n", string(values_json))

	// タイムアウトを30秒に指定してClient構造体を生成
	cli := &http.Client{Timeout: time.Duration(30) * time.Second}

	// 生成したURLを元にRequest構造体を生成
	URL := os.Getenv("URL")
	//fmt.Print(URL)
	//URL := " localhost:5000/scep"
	//fmt.Print(bytes.NewBuffer(values_json))
	req, _ := http.NewRequest("POST", URL, bytes.NewBuffer(values_json))

	// リクエストにヘッダ情報を追加
	token := os.Getenv("JWT_TOKEN")
	req.Header.Add("authorization", "Bearer"+" "+token)
	req.Header.Set("Content-Type", "application/json")

	// リクエストヘッダの内容を出力
	header, _ := httputil.DumpRequestOut(req, true)
	fmt.Print("debug1:リクエストヘッダの内容を出力\n")
	fmt.Println(string(header))
	fmt.Print("\n")

	// POSTリクエスト発行
	rsp, err := cli.Do(req)
	if err != nil {
		fmt.Print("debug2:POSTリクエスト発行\n")
		fmt.Println(err)
		fmt.Print("\n")
		return false, err
	}
	// 関数を抜ける際に必ずresponseをcloseするようにdeferでcloseを呼ぶ
	defer rsp.Body.Close()

	// レスポンスを取得し出力
	if rsp.StatusCode == 200 {
		return true, nil
	}
	body, _ := ioutil.ReadAll(rsp.Body)
	fmt.Print("debug3:エラー！レスポンスを取得し出力\n")
	fmt.Println(string(body))
	fmt.Print("\n")
	return false, errors.New(string(body))
}
func main() {
	out, err := testVerrify()
	fmt.Print("debug4\n")
	fmt.Print(out, err)
	fmt.Print("\n")
}
