package selfmadecsrverifier

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	//"io"

	"github.com/go-kit/kit/log"
)

// New creates a executablecsrverifier.ExecutableCSRVerifier.
func New(logger log.Logger) (*selfmadeCSRVerifier, error) {
	return &selfmadeCSRVerifier{logger: logger}, nil
}

// selfmadeCSRVerifier implements a csrverifier.CSRVerifier.
// It send Httprequest for nodeserver in order to verify csr
// If the command exit code is 0, the CSR is considered valid.
// In any other cases, the CSR is considered invalid.
type selfmadeCSRVerifier struct {
	logger log.Logger
}

type Values struct {
	Cn     string `json: "Cn"`
	Secret string `json: "Secret"`
}

func (v *selfmadeCSRVerifier) Verify(data []byte, ChallengePassword string, CSR *x509.CertificateRequest) (bool, error) {
	CN := CSR.Subject.String()
	challenge := ChallengePassword
	// build sample structure
	values := new(Values)
	values.Cn = CN
	values.Secret = challenge
	// encode json
	values_json, _ := json.Marshal(values)
	// タイムアウトを30秒に指定してClient構造体を生成
	cli := &http.Client{Timeout: time.Duration(30) * time.Second}
	// 生成したURLを元にRequest構造体を生成
	URL := os.Getenv("URL")
	req, _ := http.NewRequest("POST", URL, bytes.NewBuffer(values_json))
	// リクエストにヘッダ情報を追加
	token := os.Getenv("JWT_TOKEN")
	req.Header.Add("authorization", "Bearer"+" "+token)
	req.Header.Set("Content-Type", "application/json")
	// POSTリクエスト発行
	rsp, err := cli.Do(req)
	if err != nil {
		fmt.Print("debug2:POSTリクエスト発行\n")
		fmt.Println(err)
		fmt.Print("\n")
		fmt.Print("unsupported protocol schemeの時は環境変数が設定されているかをチェック\n")
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
	fmt.Print("unsupported protocol schemeの時は環境変数が設定されているかをチェック\n")
	fmt.Print("\n")
	return false, errors.New(string(body))
}
