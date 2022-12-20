package wxpay

import (
	"context"
	"huling/utils/results"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

func SubmitOrder(cgin *gin.Context) {

	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath("/resource/staticfile/merchant/apiclient_key.pem")
	if err != nil {
		log.Fatal("load merchant private key error")
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		log.Fatalf("new wechat pay client err:%s", err)
	}
	/************开始支付逻辑************************/
	svc := h5.H5ApiService{Client: client}
	// 得到prepay_id，以及调起支付所需的参数和签名
	resp, result, err := svc.Prepay(ctx,
		h5.PrepayRequest{
			Appid:       core.String("wxd678efh567hg6787"),                      //公众号ID
			Mchid:       core.String("1900009191"),                              //直连商户号
			Description: core.String("Image形象店-深圳腾大-QQ公仔"),                      // 商品描述
			OutTradeNo:  core.String("1217752501201407033233368018"),            // 商户订单号
			Attach:      core.String("自定义数据说明"),                                 // 附加数据
			NotifyUrl:   core.String("https://www.weixin.qq.com/wxpay/pay.php"), // 有效性：1. HTTPS；2. 不允许携带查询串。
			Amount: &h5.Amount{
				Total:    core.Int64(1),
				Currency: core.String("CNY"),
			},
			SceneInfo: &h5.SceneInfo{
				PayerClientIp: core.String("127.0.0.1"),
				// 商户端设备号
				DeviceId: core.String("12346785"),
				H5Info: &h5.H5Info{
					Type: core.String("iOSAndroidWap"),
				},
			},
		},
	)
	if err == nil {
		log.Println(resp)
	} else {
		log.Println(err)
	}
	results.Success(cgin, "支付统一订单", result, resp)
}
