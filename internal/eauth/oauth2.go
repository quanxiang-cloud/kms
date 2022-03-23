package eauth

// import (
// 	"context"
// 	"encoding/json"
// 	"io"
// 	"kms/internal/enums"
// 	"kms/internal/models"
// 	"kms/internal/models/redis"
// 	code2 "kms/pkg/misc/code"
// 	"kms/pkg/misc/config"
// 	"net/http"
// 	"strings"
// 	"time"

// 	"git.internal.yunify.com/qxp/misc/error2"
// 	"git.internal.yunify.com/qxp/misc/redis2"
// 	"git.internal.yunify.com/qxp/misc/resp"
// 	"github.com/gin-gonic/gin"
// 	"golang.org/x/oauth2"
// )

// var oauth2Redis models.Cache

// // OAuth2 oauth2
// type OAuth2 struct {
// 	AuthBase
// }

// // OAuth2Content auth content of oauth2
// type OAuth2Content struct {
// 	emptyAuthContentParser
// 	Type      string   `json:"type"`      // authorization-code, password,
// 	AuthURL   string   `json:"authURL"`   //
// 	TokenURL  string   `json:"tokenURL"`  //
// 	Scopes    []string `json:"scopes"`    // user define
// 	AuthStyle string   `json:"authStyle"` //
// }

// //Name name
// func (o *OAuth2) Name() enums.Enum {
// 	return enums.AuthOAuth2
// }

// // Init init
// func (o *OAuth2) Init(ak *models.AgencyKey) error {
// 	o.authContent = &OAuth2Content{}

// 	if oauth2Redis == nil {
// 		redisConfig, err := redis2.NewClient(config.Conf.Redis)
// 		if err != nil {
// 			return err
// 		}
// 		oauth2Redis = redis.NewOAuth2RedisClient(redisConfig)
// 	}

// 	return o.AuthBase.Init(ak)
// }

// // Invoke invoke
// func (o *OAuth2) Invoke(d interface{}) (string, error) {
// 	content := o.authContent.(OAuth2Content)

// 	oauth2c, err := oauth2Redis.QueryConfig(o.agencyKey.ID)
// 	if err != nil {
// 		oauth2c = &oauth2.Config{
// 			ClientID:     o.keyID,
// 			ClientSecret: o.secret,
// 			Scopes:       oauth2c.Scopes,
// 			RedirectURL:  config.Conf.OAuth2Host,
// 			Endpoint: oauth2.Endpoint{
// 				AuthURL:  content.AuthURL,
// 				TokenURL: content.TokenURL,
// 			},
// 		}
// 	}

// 	token, err := oauth2Redis.QueryToken(o.agencyKey.ID)
// 	// if token existed and expired, refresh
// 	if err == nil && time.Until(token.Expiry) < 0 {
// 		if token, err = oauth2c.TokenSource(context.Background(), token).Token(); err != nil {
// 			return "", err
// 		}
// 		oauth2Redis.CacheToken(o.agencyKey.ID, token, time.Until(token.Expiry.Add(time.Hour*24)))
// 		return token.AccessToken, nil
// 	}

// 	url := oauth2c.AuthCodeURL(o.agencyKey.ID)
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return "", err
// 	}
// 	b, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return "", err
// 	}
// 	token = &oauth2.Token{}
// 	json.Unmarshal(b, token)
// 	oauth2Redis.CacheToken(o.agencyKey.ID, token, time.Until(token.Expiry.Add(time.Hour*24)))

// 	return token.AccessToken, nil
// }

// // OAuth2Token oauth2 redirect to get token
// // 		   ** External API **
// // Redirect this API by oauth2 server
// // To get token when oauth invoked
// func OAuth2Token(c *gin.Context) {
// 	state, code, clientID := getParams(c)
// 	if code == "" {
// 		resp.Format(nil, error2.NewError(code2.ErrInvalidOAuth2Code)).Context(c)
// 		return
// 	}
// 	if clientID == "" {
// 		resp.Format(nil, error2.NewError(code2.ErrInvalidOAuth2ClientID)).Context(c)
// 		return
// 	}

// 	oauth2c, err := oauth2Redis.QueryConfig(state)
// 	if err != nil {
// 		resp.Format(nil, err).Context(c)
// 	}

// 	token, err := oauth2c.Exchange(context.Background(), code)
// 	resp.Format(token, err).Context(c)
// }

// func getParams(c *gin.Context) (state, code, clientID string) {
// 	switch c.Request.Method {
// 	case "GET":
// 		state = c.Query("state")
// 		code = c.Query("code")
// 		clientID = c.Query("client_id")
// 	case "POST":
// 		contentType := strings.ToLower(c.ContentType())
// 		if contentType == "application/json" {
// 			codeState := &struct {
// 				State    string `json:"state"`
// 				Code     string `json:"code"`
// 				ClientID string `json:"client_id"`
// 			}{}
// 			c.ShouldBind(codeState)
// 			state = codeState.State
// 			code = codeState.Code
// 			clientID = codeState.ClientID
// 		}
// 		if contentType == "multipart/form-date" {
// 			state = c.PostForm("state")
// 			code = c.PostForm("code")
// 			clientID = c.PostForm("client_id")
// 		}
// 	}
// 	return
// }
