package openwechat

import (
	"errors"

	"encoding/json"
)

/**
留着其它接口用
else if m := a.marshalXML(string(v)); m == nil {
		return errors.New("response xml data parse error!")
	} else if d, err := json.Marshal(&m); err != nil {
		return err
	} else if err := json.Unmarshal(d, resp); err != nil {
		return err
	} else {
		logs.Debug(resp)
	}



*/
//https://api.mch.weixin.qq.com/pay/downloadbill
//https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=9_6
//商户可以通过该接口下载历史交易清单
//微信在次日9点启动生成前一天的对账单，建议商户10点后再获取；
//对账单接口只能下载三个月以内的账单。
type Api_wechat_pay_downloadbill struct {
	wechatApi
}

func (a *Api_wechat_pay_downloadbill) apiUrl() string {
	return "https://api.mch.weixin.qq.com/pay/downloadbill"
}

func (a *Api_wechat_pay_downloadbill) Run(resp *Resp_api_wechat_pay_downloadbill) error {
	//准备请求
	result := []byte{}
	if v, err := a.request(a.apiUrl()); err != nil {
		//接口错误.直接返回.
		return err
	} else if m := a.marshalXML(string(v)); m != nil {
		//把xml 数据转 map
		//如果 m 不是 nil.. 表示转换成功...说明返回了 xml 数据..
		//应该就能表明肯定是失败的..
		if d, err := json.Marshal(&m); err != nil {
			return err
		} else if err := json.Unmarshal(d, resp); err != nil {
			return err
		}
	} else {
		result = v
	}
	//那么这里 的数据是对账单数据...直接返回
	//当这2个没值的时候..才返回数据.
	if len(resp.ReturnCode) == 0 && len(resp.ErrCode) == 0 {
		resp.Data = result
	}
	return nil
}

//对账单参数
type Req_api_wechat_pay_downloadbill struct {
	//商户号
	MchId string `json:"mch_id"`
	//随机字符串
	NonceStr string `json:"nonce_str"`
	//子商户公众账号ID
	SubAppId string `json:"sub_appid"`
	//子商户号
	SubMchId string `json:"sub_mch_id"`
	//传了这个.就只会下载这里的对账单
	DeviceInfo string `json:"device_info"`
	//下载对账单的日期，格式：20140603
	BillDate string `json:"bill_date"`
	//ALL，返回当日所有订单信息，默认值
	//SUCCESS，返回当日成功支付的订单
	//REFUND，返回当日退款订单
	//RECHARGE_REFUND，返回当日充值退款订单（相比其他对账单多一栏“返还手续费”）
	BillType string `json:"bill_type"`
	//非必传参数，固定值：GZIP，返回格式为.gzip的压缩包账单。不传则默认为数据流形式。
	TarType string `json:"tar_type"`
}

func (p Req_api_wechat_pay_downloadbill) valid() error {
	if len(p.MchId) == 0 {
		return errors.New(" mch_id can not be nil")
	}
	if len(p.TarType) > 0 && p.TarType != "GZIP" {
		return errors.New("grant_type must be GZIP")
	}
	if len(p.BillDate) == 0 {
		return errors.New("bill_date can not be nil")
	}
	return nil
}

type Resp_api_wechat_pay_downloadbill struct {
	Response
	Data []byte //文件内容
}
