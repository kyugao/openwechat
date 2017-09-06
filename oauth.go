/*
** ===============================================
** USER NAME: garlic(QQ:3173413)
** FILE NAME: api_auth.go
** DATE TIME: 2017-07-21 09:09:23
** ===============================================
 */

package openwechat

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/liteck/logs"
	"github.com/liteck/tools/httplib"
)

/**
授权参考文档
https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=4_4
微信这里授权方式比较多样.针对各种场景下.微信做了独立的区分...
1. 网页授权 开放平台用户授权  https://open.weixin.qq.com/connect
2.
**/

//开放平台的网页授权模式.获取授权链接
func OpenWebAuth(app_id, scope, redirect_uri string) string {
	uri := "https://open.weixin.qq.com/connect/oauth2/authorize"
	uri += "?appid=" + app_id
	uri += "&redirect_uri=" + redirect_uri
	uri += "&response_type=code"
	uri += "&scope=" + scope
	uri += "&state=" + app_id
	uri += "#wechat_redirect"
	return uri
}

//通过授权回调之后的 code 换取 access_token
type Api_wechat_sns_oauth2_access_token struct {
	wechatApi
}

func (a *Api_wechat_sns_oauth2_access_token) apiUrl() string {
	return "https://api.weixin.qq.com/sns/oauth2/access_token"
}

func (a *Api_wechat_sns_oauth2_access_token) Run(resp *Resp_api_wechat_sns_oauth2_access_token) error {
	m := a.toMap()
	m["appid"] = a.app_id
	m["secret"] = a.secret.AppSecret

	http_request := httplib.Get(a.apiUrl())
	tmp_string := ""
	for k, v := range m {
		value := fmt.Sprintf("%v", v)
		if value != "" {
			http_request.Param(k, value)
			tmp_string = tmp_string + k + "=" + value + "\t"
		}
	}
	logs.Debug(fmt.Sprintf("==[request params]==[%s,%s]", a.apiUrl(), tmp_string))
	if v, err := http_request.Bytes(); err != nil {
		return err
	} else if err := json.Unmarshal(v, resp); err != nil {
		return err
	} else {
		logs.Debug(string(v))
	}

	return nil
}

type Req_api_wechat_sns_oauth2_access_token struct {
	Appid     string `json:"appid"`
	Secret    string `json:"secret"`
	Code      string `json:"code"`
	GrantType string `json:"grant_type"`
}

func (p Req_api_wechat_sns_oauth2_access_token) valid() error {
	if len(p.GrantType) == 0 {
		return errors.New("grant_type can not be nil")
	}
	if len(p.Code) == 0 {
		return errors.New("code can not be nil")
	}

	return nil
}

type Resp_api_wechat_sns_oauth2_access_token struct {
	ErrCode      float64 `json:"errcode,omitempty"`
	ErrMsg       string  `json:"errmsg,omitempty"`
	AccessToken  string  `json:"access_token,omitempty"`
	ExpiresIn    float64 `json:"expires_in,omitempty"`
	RefreshToken string  `json:"refresh_token,omitempty"`
	OpenId       string  `json:"openid,omitempty"`
	Scope        string  `json:"scope,omitempty"`
	UnoinId      string  `json:"unionid,omitempty"`
}

/**
通过授权回调之后的 access_token 换取 userinfo
接口说明
此接口用于获取用户个人信息。开发者可通过OpenID来获取用户基本信息。
特别需要注意的是，如果开发者拥有多个移动应用、网站应用和公众帐号，
可通过获取用户基本信息中的unionid来区分用户的唯一性，
因为只要是同一个微信开放平台帐号下的移动应用、网站应用和公众帐号，用户的unionid是唯一的。
换句话说，同一用户，对同一个微信开放平台下的不同应用，unionid是相同的。
请注意，在用户修改微信头像后，旧的微信头像URL将会失效，
因此开发者应该自己在获取用户信息后，将头像图片保存下来，避免微信头像URL失效后的异常情况。
**/
type Api_wechat_sns_userinfo struct {
	wechatApi
}

func (o *Api_wechat_sns_userinfo) apiUrl() string {
	return "https://api.weixin.qq.com/sns/userinfo"
}

func (a *Api_wechat_sns_userinfo) Run(resp *Resp_api_wechat_sns_userinfo) error {
	m := a.toMap()
	http_request := httplib.Get(a.apiUrl())
	tmp_string := ""
	for k, v := range m {
		value := fmt.Sprintf("%v", v)
		if value != "" {
			http_request.Param(k, value)
			tmp_string = tmp_string + k + "=" + value + "\t"
		}
	}
	logs.Debug(fmt.Sprintf("==[request params]==[%s,%s]", a.apiUrl(), tmp_string))
	if v, err := http_request.Bytes(); err != nil {
		return err
	} else if err := json.Unmarshal(v, resp); err != nil {
		return err
	} else {
		logs.Debug(string(v))
	}

	return nil
}

type Req_api_wechat_sns_userinfo struct {
	AccessToken string `json:"access_token"`
	OpenId      string `json:"openid"`
	Lang        string `json:"lang"`
}

func (p Req_api_wechat_sns_userinfo) valid() error {
	if len(p.AccessToken) == 0 {
		return errors.New("access_token can not be nil")
	}
	if len(p.OpenId) == 0 {
		return errors.New("openid can not be nil")
	}

	return nil
}

type Resp_api_wechat_sns_userinfo struct {
	ErrCode    float64 `json:"errcode,omitempty"`
	ErrMsg     string  `json:"errmsg,omitempty"`
	OpenId     string  `json:"openid,omitempty"`
	NickName   string  `json:"nickname,omitempty"`
	Sex        string  `json:"sex,omitempty"`
	Province   string  `json:"province,omitempty"`
	City       string  `json:"city,omitempty"`
	Country    string  `json:"country,omitempty"`
	HeadimgUrl string  `json:"headimgurl,omitempty"`
	Privilege  string  `json:"privilege,omitempty"`
	UnoinId    string  `json:"unionid,omitempty"`
}
