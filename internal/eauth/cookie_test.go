package eauth

import (
	"encoding/json"
	"fmt"
	"io"
	"kms/internal/models"
	"kms/internal/rule"
	"kms/pkg/dbcli"
	"kms/pkg/misc/config"
	"net/http"
	"testing"
	"time"

	"github.com/quanxiang-cloud/cabin/logger"
)

func _TestCookie(t *testing.T) {
	config, err := config.NewConfig("./testdata/config.yml")
	if err != nil {
		panic(err)
	}

	log := logger.New(&config.Logger)

	dbcli.InitDB(config, log)

	// cookieServer()

	type TestingCase struct {
		sample *models.AgencyKey
	}

	cases := []TestingCase{
		{
			sample: &models.AgencyKey{
				KeyID:     "qxp",
				KeySecret: "SfRbeO5pM21jNdbzIiU_ZwwKwEvhJTfM",
				// https://home.yunify.com/distributor.action
				AuthContent: `
				[{
					"name": "format",
					"type": "number",
					"data": "1",
					"from": "body"
				},
				{
					"name": "useragent",
					"type": "string",
					"data": "ApiClient",
					"from": "body"
				},
				{
					"name": "rid",
					"type": "random",
					"data": "8",
					"from": "body"
				},
				{
					"name": "parameters",
					"type": "string",
					"data": "[\"5fac019554b458\", \"${KEY_ID}\", \"${SECRET_KEY}\", 2052]",
					"from": "body"
				},
				{
					"name": "timestamp",
					"type": "date",
					"data": "yyyy-MM-ddTHH:mm:ss.Szz",
					"from": "body"
				},
				{
					"name": "v",
					"type": "string",
					"data": "1.0",
					"from": "body"
				},
				{
					"type": "authurl",
					"data": "https://kingdee.yunify.com/K3Cloud/Kingdee.BOS.WebApi.ServicesStub.AuthService.ValidateUser.common.kdsvc"
				},
				{
					"type": "method",
					"data": "POST"
				},
				{
					"type": 	 "cookie",
					"captureIn": "cookie",
					"in": 		 "header"
				},
				{
					"type": "expire",
					"data": "01"
				}]`,
				// TODO: type=cookie captureIn=header
				AuthType: "cookie",
			},
		},
	}

	for _, cs := range cases {
		if _, err := rule.ParseAuthContent(cs.sample.AuthContent, true); err != nil {
			panic(err)
		}

		auth, err := factory.Create(cs.sample)
		if err != nil {
			t.Fatal(err)
		}
		tokens, err := auth.Invoke(nil)
		if err != nil {
			t.Fatal(err)
		}

		for _, v := range tokens {
			fmt.Printf("%+v\n", v)
		}
	}
}

func cookieServer() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Printf("method: %s\n", r.Method)

		b, _ := io.ReadAll(r.Body)
		type Req struct {
			Const     string `json:"serviceName"`
			KeyID     string `json:"userName"`
			KeySecret string `json:"password"`
		}

		req := &Req{}
		json.Unmarshal(b, req)
		req.Const = r.URL.Query().Get("serviceName")
		req.KeyID = r.URL.Query().Get("userName")
		req.KeySecret = r.URL.Query().Get("password")
		fmt.Printf("req: %+v\n", req)

		type Resp struct {
			Token string `json:"binding"`
		}
		// if req.KeyID == "UjjNR4B_qO6925q7afx8dQ" && req.KeySecret == "fEdVqD9KojiSBfYMfs61gUmmdxLoM9KcqFB9I8dTHKw" {
		// 	b, _ := json.Marshal(Resp{
		// 		Token: "this is a token",
		// 	})
		// 	rw.Write(b)
		// }
		notb, _ := json.Marshal(Resp{
			Token: "token",
		})
		http.SetCookie(rw, &http.Cookie{
			Name:  "cookiename",
			Value: "mycookie",
		})
		rw.Write(notb)
	})
	go http.ListenAndServe(":8080", nil)
	<-time.After(1 * time.Second)
}
