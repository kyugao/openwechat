package openwechat

import (
	"errors"
	"sync"
)

/**
微信开放平台需要的基本参数.
这里尽可能多的添加了参数.各种 接口的需求不一样..各取所需.
*/
type Secret struct {
	AppId     string
	AppSecret string
	AppKey    string
	ApiCert   []byte
}

/**
秘钥校验
*/
func (s *Secret) valid() error {
	if len(s.AppId) == 0 {
		return errors.New("app_id 不能为空")
	}

	if len(s.AppSecret) == 0 {
		return errors.New("app_secret 不能为空")
	}

	return nil
}

/**
存储秘钥列表
*/
var secretLst sync.Map

/**
调用者进行注册
以 app_id 为 key
*/
func RegisterSecret(s ...Secret) error {
	if len(s) == 0 {
		return errors.New("配置参数不能为空!")
	}

	for _, v := range s {
		if err := v.valid(); err != nil {
			return err
		}
		secretLst.Store(v.AppId, v)
	}

	return nil
}

/**
获取秘钥
*/
func getSecret(app_id string) Secret {
	if v, ok := secretLst.Load(app_id); !ok || v == nil {
		return Secret{}
	} else {
		return v.(Secret)
	}
}


