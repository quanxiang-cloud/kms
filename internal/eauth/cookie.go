package eauth

import (
	"fmt"
	"kms/internal/enums"
	"kms/internal/httputil"
	"kms/internal/models"
	"kms/internal/models/redis"
	"kms/internal/rule"
	"kms/pkg/random"
	"kms/pkg/timestamp"
	"strconv"
	"time"
)

// Cookie auth by cookie
type Cookie struct {
	AuthBase
}

// CookieContent content
type CookieContent struct {
	httpClient *httputil.D            `json:"-"`
	Expiry     time.Duration          `json:"expiry"` // expiry date
	Header     map[string]string      `json:"header"`
	Body       map[string]interface{} `json:"body"`
	Query      map[string]string      `json:"query"`
	AuthURL    string                 `json:"authURL"` // sample: http://xxxxx.com/api/login
	Method     string                 `json:"method"`  // POST GET etc.
	Resp       []*CookieResp          `json:"resp"`
}

// CookieResp CookieResp
type CookieResp struct {
	RespField string `json:"respField"`
	RespType  string `json:"respType"`
	In        string `json:"in"`
}

// Init init
func (c *Cookie) Init(ak *models.AgencyKey) error {
	c.authContent = &CookieContent{
		Header: map[string]string{},
		Body:   map[string]interface{}{},
		Query:  map[string]string{},
		Resp:   []*CookieResp{},
	}
	return c.AuthBase.Init(ak)
}

// Name name
func (c *Cookie) Name() enums.Enum {
	return enums.AuthCookie
}

// Invoke invoke
func (c *Cookie) Invoke(d interface{}) ([]*AuthResp, error) {
	var tokens []*AuthResp
	content := c.authContent.(*CookieContent)

	// check cache
	rdb, err := getRedisCli()
	if err != nil {
		return nil, err
	}
	if err := rdb.Query(redis.CacheCookies, c.keyID, &tokens, time.Minute*content.Expiry); err == nil {
		return tokens, nil
	}

	content.httpClient = &httputil.D{
		Method: content.Method,
		URL:    content.AuthURL,
		Header: content.Header,
		Body:   content.Body,
		Query:  content.Query,
	}

	replace := []*httputil.R{{Old: ConfigKeyName, New: c.keyID}, {Old: ConfigSecretKeyName, New: c.secret}}
	_, err = content.httpClient.Do(nil, replace)
	if err != nil {
		return nil, err
	}

	tokens, err = content.getToken(content.Resp)
	if err != nil {
		return nil, err
	}

	err = rdb.Cache(redis.CacheCookies, c.keyID, tokens, time.Minute*content.Expiry)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

// Parse parse and verify config values for this auth type
func (c *CookieContent) Parse(vals rule.ConfigValueSet) error {
	for _, val := range vals {
		switch {
		case val.From != "":
			switch val.Type {
			case enums.BasicTypeKeyID.Val():
				val.Data = ConfigKeyName
			case enums.BasicTypeKeySecret.Val():
				val.Data = ConfigSecretKeyName
			case enums.BasicTypeFormatDate.Val():
				val.Data = timestamp.Format(time.Now(), val.Data)
			case enums.BasicTypeRandom.Val():
				l, err := strconv.Atoi(val.Data)
				if err != nil {
					return err
				}
				val.Data = random.String(l)
			}

			switch val.From {
			case enums.HTTPBody.Val():
				switch val.Type {
				case enums.BasicTypeNumber.Val():
					d, err := strconv.ParseFloat(val.Data, 64)
					if err != nil {
						return err
					}
					c.Body[val.Name] = d
				case enums.BasicTypeBoolean.Val():
					d, err := strconv.ParseBool(val.Data)
					if err != nil {
						return err
					}
					c.Body[val.Name] = d
				default:
					c.Body[val.Name] = val.Data
				}
			case enums.HTTPQuery.Val():
				c.Query[val.Name] = val.Data
			case enums.HTTPHeader.Val():
				c.Header[val.Name] = val.Data
			default:
				// do nothing
			}
		case val.In != "":
			c.Resp = append(c.Resp, &CookieResp{
				RespField: val.Name,
				RespType:  val.CaptureIn,
				In:        val.In,
			})
		case val.Type == enums.BasicTypeAuthURL.Val():
			c.AuthURL = val.Data
		case val.Type == enums.BasicTypeMethod.Val():
			c.Method = val.Data
		case val.Type == enums.BasicTypeExpire.Val():
			d, err := strconv.Atoi(val.Data)
			if err != nil {
				return err
			}
			c.Expiry = time.Duration(d)
		}
	}
	return nil
}

func (c *CookieContent) getToken(where []*CookieResp) ([]*AuthResp, error) {
	tokens := make([]*AuthResp, 0, len(where))
	for _, in := range where {
		var token string
		switch in.RespType {
		case enums.HTTPCookie.Val():
			token = c.httpClient.GetCookies()
			in.RespField = "Cookie"
		case enums.HTTPHeader.Val():
			token = c.httpClient.GetHeader(in.RespField)
		case enums.HTTPBody.Val():
			val, err := c.httpClient.GetJSONBody(in.RespField)
			if err != nil {
				return nil, err
			}
			if s, ok := val.(string); ok {
				token = s
			} else {
				return nil, fmt.Errorf("only string supported")
			}
		default:
			//
		}
		tokens = append(tokens, &AuthResp{
			Name:  in.RespField,
			Value: token,
			In:    in.In,
		})
	}
	return tokens, nil
}
