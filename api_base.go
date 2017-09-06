package openwechat

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"sort"

	"github.com/astaxie/beego"
	"github.com/liteck/logs"
	"github.com/liteck/tools"
	"github.com/liteck/tools/httplib"
)

type reqInterface interface {
	valid() error
}

type responseInterface interface{}

type Response struct {
	ErrCode string `json:"errcode,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`
}

type wechatApi struct {
	//公众账号ID
	app_id string
	//保留秘钥参数.留作签名和验证签名使用.
	secret Secret
	//请求参数
	req reqInterface
}

func (w *wechatApi) SetAppId(app_id string) error {
	w.app_id = app_id
	if len(w.app_id) == 0 {
		return errors.New("app_id can not be nil")
	}

	if s := getSecret(app_id); len(s.AppKey) == 0 && len(s.AppSecret) == 0 {
		return errors.New("secret can not be nil")
	} else {
		w.secret = s
	}
	return nil
}

/**
设置请求内容..
*/
func (w *wechatApi) SetReqContent(v reqInterface) error {
	if err := v.valid(); err != nil {
		return err
	}
	w.req = v
	return nil
}

func (w *wechatApi) toMap() map[string]interface{} {
	var data = make(map[string]interface{})

	t := reflect.TypeOf(w.req)
	v := reflect.ValueOf(w.req)

	for i := 0; i < t.NumField(); i++ {
		key := t.Field(i).Name
		value := v.Field(i).Interface()
		tag := t.Field(i).Tag.Get("json")
		if tag != "" {
			if strings.Contains(tag, ",") {
				ps := strings.Split(tag, ",")
				key = ps[0]
			} else {
				key = tag
			}
		}
		data[key] = value
	}

	return data
}

func (w *wechatApi) toXml(input map[string]interface{}) (xml string) {
	xml = "<xml>"
	for k, v := range input {
		value := fmt.Sprintf("%v", input[k])
		if value != "" {
			switch v.(type) {
			case string:
				if len(value) > 0 {
					xml += fmt.Sprintf("<%s>%s</%s>", k, v, k)
				}
			case int:
				xml += fmt.Sprintf("<%s>%d</%s>", k, v, k)
			}
		}
	}
	xml += "</xml>"
	return
}

func (w *wechatApi) marshalXML(input string) map[string]interface{} {
	m := map[string]interface{}{}
	var t xml.Token
	var err error
	inputReader := strings.NewReader(input)

	decoder := xml.NewDecoder(inputReader)
	var key string
	for t, err = decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		// 处理元素开始（标签）
		case xml.StartElement:
			key = token.Name.Local
			// 处理元素结束（标签）
		case xml.EndElement:
			key = "root" //保证此时不做处理
			// 处理字符数据（这里就是元素的文本）
		case xml.CharData:
			content := string([]byte(token))
			if key != "root" {
				m[key] = content
			}
		default:
		}
	}
	if err == nil || (err != nil && err.Error() == "EOF") {
		return m
	} else {
		return nil
	}
}

/**
做签名
*/
func (w *wechatApi) doSign(m map[string]interface{}) string {
	//对key进行升序排序.
	sorted_keys := make([]string, 0)
	for k, _ := range m {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)
	//对key=value的键值对用&连接起来，略过空值
	var signStrings string
	for _, k := range sorted_keys {
		value := fmt.Sprintf("%v", m[k])
		if value != "" {
			signStrings = signStrings + k + "=" + value + "&"
		}
	}
	//在键值对的最后加上key=API_KEY
	signStrings = signStrings + "key=" + w.secret.AppKey
	logs.Debug(fmt.Sprintf("==[sign tobe]==[%s]", signStrings))

	//5.进行MD5签名并且将所有字符转为大写.
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signStrings))
	cipherStr := md5Ctx.Sum(nil)
	upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))
	logs.Debug(fmt.Sprintf("==[sign end]==[%s]", upperSign))
	return upperSign
}

//微信验证签名
func (w *wechatApi) verifySign(resp_map map[string]interface{}) (pass bool) {
	sorted_keys := make([]string, 0)
	for k, _ := range resp_map {
		if k != "sign" {
			sorted_keys = append(sorted_keys, k)
		}
	}
	sort.Strings(sorted_keys)
	//STEP2, 对key=value的键值对用&连接起来，略过空值
	var signStrings string
	for _, k := range sorted_keys {
		value := fmt.Sprintf("%v", resp_map[k])
		if value != "" {
			signStrings = signStrings + k + "=" + value + "&"
		}
	}
	//STEP3, 在键值对的最后加上key=API_KEY
	signStrings = signStrings + "key=" + w.secret.AppKey
	beego.Debug(signStrings)

	//STEP4, 进行MD5签名并且将所有字符转为大写.
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signStrings))
	cipherStr := md5Ctx.Sum(nil)
	upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))
	if upperSign == resp_map["sign"] {
		pass = true
	}
	return
}

func (w *wechatApi) request(link string) ([]byte, error) {
	//struct 转 map
	m := w.toMap()
	m["nonce_str"] = tools.RandomNumeric(32)
	m["appid"] = w.app_id
	m["sign"] = w.doSign(m)

	http_request := httplib.Post(link)
	logs.Debug("==[request_params]==", m)
	str_requst := w.toXml(m)
	logs.Debug("==[request_params]==", str_requst)
	http_request.Body(str_requst)
	if v, err := http_request.Bytes(); err != nil {
		return nil, err
	} else {
		return v, nil
	}
}
