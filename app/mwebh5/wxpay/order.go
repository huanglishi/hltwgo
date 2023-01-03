package wxpay

import (
	"context"
	"fmt"
	"huling/utils/results"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

func SubmitOrder(cgin *gin.Context) {
	order_id := cgin.DefaultQuery("order_id", "o")
	if order_id == "" {
		results.Failed(cgin, "请求参数（order_id）不能为空", nil)
	} else {
		orderdata, ordererr := DB().Table("client_product_order").Where("id", order_id).Fields("id,cuid,uid,product_id,title,price,out_trade_no,note").First()
		if ordererr != nil {
			results.Failed(cgin, "获取订单数据失败", ordererr)
		} else {
			paymentconfig, payconfrerr := DB().Table("client_system_paymentconfig").Where("cuid", orderdata["cuid"]).Fields("id,appId,mchID,mchAPIv3Key,mchCertificateSerialNumber,privatekey").First()
			if payconfrerr != nil {
				results.Failed(cgin, "获取支付配置失败", payconfrerr)
				return
			}
			if paymentconfig == nil {
				results.Failed(cgin, "获取支付配置未配置", orderdata)
				return
			}
			//配置参数
			mchID := paymentconfig["mchID"].(string)                                           // 商户号
			mchCertificateSerialNumber := paymentconfig["mchCertificateSerialNumber"].(string) // 商户证书序列号
			mchAPIv3Key := paymentconfig["mchAPIv3Key"].(string)                               // 商户APIv3密钥

			// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
			mchPrivateKey, err := utils.LoadPrivateKeyWithPath("/resource/staticfile/merchant/apiclient_key.pem")
			if err != nil {
				results.Failed(cgin, "加载本地商户私钥失败", err)
				return
			}

			ctx := context.Background()
			// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
			opts := []core.ClientOption{
				option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
			}
			client, err := core.NewClient(ctx, opts...)
			if err != nil {
				results.Failed(cgin, fmt.Sprintf("新微信支付客户端出错:%s", err), err)
				return
			}
			/************开始支付逻辑************************/
			svc := h5.H5ApiService{Client: client}
			// 得到prepay_id，以及调起支付所需的参数和签名
			resp, result, err := svc.Prepay(ctx,
				h5.PrepayRequest{
					Appid:       core.String(paymentconfig["appId"].(string)),           //公众号ID
					Mchid:       core.String(paymentconfig["mchID"].(string)),           //直连商户号
					Description: core.String(orderdata["title"].(string)),               // 商品描述
					OutTradeNo:  core.String(orderdata["out_trade_no"].(string)),        // 商户订单号
					Attach:      core.String(orderdata["note"].(string)),                // 附加数据
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
	}
}

// 支付成功回调
func WxPayNotify(cgin *gin.Context) {

}
