package wxpay

import (
	"context"
	"fmt"
	"huling/utils/results"
	"log"
	"net"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/certificates"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

// 下单
func SubmitOrder(cgin *gin.Context) {
	order_id := cgin.DefaultQuery("order_id", "")
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
				results.Failed(cgin, "支付配置未配置", orderdata)
				return
			}
			//配置参数
			mchID := paymentconfig["mchID"].(string)                                           // 商户号
			mchCertificateSerialNumber := paymentconfig["mchCertificateSerialNumber"].(string) // 商户证书序列号
			mchAPIv3Key := paymentconfig["mchAPIv3Key"].(string)                               // 商户APIv3密钥

			// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
			pemPath := fmt.Sprintf("./%s", paymentconfig["privatekey"])
			mchPrivateKey, err := utils.LoadPrivateKeyWithPath(pemPath)
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
					Description: core.String("Image形象店-深圳腾大-QQ公仔"),                      // 商品描述
					OutTradeNo:  core.String("202377525012033332"),                      // 商户订单号
					Attach:      core.String("自定义数据说明"),                                 // 附加数据
					NotifyUrl:   core.String("https://www.weixin.qq.com/wxpay/pay.php"), // 有效性：1. HTTPS；2. 不允许携带查询串。
					Amount: &h5.Amount{
						Total:    core.Int64(1),
						Currency: core.String("CNY"),
					},
					SceneInfo: &h5.SceneInfo{
						PayerClientIp: core.String(GetOutboundIP()),
						H5Info: &h5.H5Info{
							Type: core.String("Wap"),
						},
					},
				},
				// h5.PrepayRequest{
				// 	Appid:       core.String(paymentconfig["appId"].(string)),           //公众号ID
				// 	Mchid:       core.String(paymentconfig["mchID"].(string)),           //直连商户号
				// 	Description: core.String(orderdata["title"].(string)),               // 商品描述
				// 	OutTradeNo:  core.String(orderdata["out_trade_no"].(string)),        // 商户订单号
				// 	Attach:      core.String(orderdata["note"].(string)),                // 附加数据
				// 	NotifyUrl:   core.String("https://www.weixin.qq.com/wxpay/pay.php"), // 有效性：1. HTTPS；2. 不允许携带查询串。
				// 	Amount: &h5.Amount{
				// 		Total:    core.Int64(1),
				// 		Currency: core.String("CNY"),
				// 	},
				// },
			)
			if err == nil {
				results.Success(cgin, "支付统一订单", resp, result.Response.StatusCode)
			} else {
				results.Success(cgin, fmt.Sprintf("h5统一订单失败,微信返回状态码：%d", result.Response.StatusCode), err, resp)
				// results.Success(cgin, fmt.Sprintf("h5统一订单失败,微信返回错误码：%d", result.Response.StatusCode), resp, result.Response.StatusCode)
			}
		}
	}
}

// JSAPI支付
func WxJsapiPay(cgin *gin.Context) {
	order_id := cgin.DefaultQuery("order_id", "")
	wxopenid := cgin.DefaultQuery("openid", "")
	if order_id == "" {
		results.Failed(cgin, "请求参数（order_id）不能为空", nil)
	} else if wxopenid == "" {
		results.Failed(cgin, "请求参数（openid）不能为空", nil)
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
				results.Failed(cgin, "支付配置未配置", orderdata)
				return
			}
			//配置参数
			mchID := paymentconfig["mchID"].(string)                                           // 商户号
			mchCertificateSerialNumber := paymentconfig["mchCertificateSerialNumber"].(string) // 商户证书序列号
			mchAPIv3Key := paymentconfig["mchAPIv3Key"].(string)                               // 商户APIv3密钥

			// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
			pemPath := fmt.Sprintf("./%s", paymentconfig["privatekey"])
			mchPrivateKey, err := utils.LoadPrivateKeyWithPath(pemPath)
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
			svc := jsapi.JsapiApiService{Client: client}
			// 得到prepay_id，以及调起支付所需的参数和签名
			//计算价格
			price_fl, _ := strconv.ParseFloat(orderdata["price"].(string), 64)
			price_int := int64(price_fl * 100)
			// pay_out_order := GenerateCode()
			resp, result, err := svc.PrepayWithRequestPayment(ctx,
				jsapi.PrepayRequest{
					Appid:       core.String(paymentconfig["appId"].(string)),
					Mchid:       core.String(paymentconfig["mchID"].(string)),
					Description: core.String(orderdata["title"].(string)),
					OutTradeNo:  core.String(orderdata["out_trade_no"].(string)),
					Attach:      core.String(orderdata["note"].(string)),
					NotifyUrl:   core.String("https://tuwen.hulingyun.cn/mwebh5/wxpay/paynotify/"),
					Amount: &jsapi.Amount{
						Total: core.Int64(price_int),
					},
					Payer: &jsapi.Payer{
						Openid: core.String(wxopenid), //微信公众号用户openid
					},
				},
			)
			if err == nil {
				//更新支付单号
				DB().Table("client_product_order").Where("id", order_id).Data(map[string]interface{}{"prepay_id": resp.PrepayId}).Update()
				results.Success(cgin, "jsAPi支付统一订单", resp, result.Response.StatusCode)
			} else {
				results.Success(cgin, fmt.Sprintf("jsAPi支付统一订单失败,微信返回状态码：%d", result.Response.StatusCode), err, resp)
			}
		}
	}
}

// 下载微信支付平台证书
func WxDownCert(cgin *gin.Context) {
	order_id := cgin.DefaultQuery("order_id", "")
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
				results.Failed(cgin, "支付配置未配置", orderdata)
				return
			}
			//配置参数
			mchID := paymentconfig["mchID"].(string)                                           // 商户号
			mchCertificateSerialNumber := paymentconfig["mchCertificateSerialNumber"].(string) // 商户证书序列号
			mchAPIv3Key := paymentconfig["mchAPIv3Key"].(string)                               // 商户APIv3密钥

			// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
			pemPath := fmt.Sprintf("./%s", paymentconfig["privatekey"])
			mchPrivateKey, err := utils.LoadPrivateKeyWithPath(pemPath)
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
			svc := certificates.CertificatesApiService{Client: client}
			resp, result, err := svc.DownloadCertificates(ctx)
			// log.Printf("status=%d resp=%s", result.Response.StatusCode, resp)
			if err == nil {
				results.Success(cgin, "测试下载商户证书成功", resp, result.Response.StatusCode)
			} else {
				results.Failed(cgin, "测试下载商户证书失败", err)
			}
		}
	}
}

// 查询订单
func WxFindOrder(cgin *gin.Context) {
	order_id := cgin.DefaultQuery("order_id", "")
	if order_id == "" {
		results.Failed(cgin, "请求参数（order_id）不能为空", nil)
	} else {
		orderdata, ordererr := DB().Table("client_product_order").Where("id", order_id).Fields("id,cuid,uid,product_id,title,price,out_trade_no,note,transaction_id").First()
		if ordererr != nil {
			results.Failed(cgin, "获取订单数据失败", ordererr)
		} else {
			paymentconfig, payconfrerr := DB().Table("client_system_paymentconfig").Where("cuid", orderdata["cuid"]).Fields("id,appId,mchID,mchAPIv3Key,mchCertificateSerialNumber,privatekey").First()
			if payconfrerr != nil {
				results.Failed(cgin, "获取支付配置失败", payconfrerr)
				return
			}
			if paymentconfig == nil {
				results.Failed(cgin, "支付配置未配置", orderdata)
				return
			}
			//配置参数
			mchID := paymentconfig["mchID"].(string)                                           // 商户号
			mchCertificateSerialNumber := paymentconfig["mchCertificateSerialNumber"].(string) // 商户证书序列号
			mchAPIv3Key := paymentconfig["mchAPIv3Key"].(string)                               // 商户APIv3密钥

			// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
			pemPath := fmt.Sprintf("./%s", paymentconfig["privatekey"])
			mchPrivateKey, err := utils.LoadPrivateKeyWithPath(pemPath)
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
			svc := jsapi.JsapiApiService{Client: client}

			resp, result, err := svc.QueryOrderById(ctx,
				jsapi.QueryOrderByIdRequest{
					TransactionId: core.String(orderdata["transaction_id"].(string)),
					Mchid:         core.String(paymentconfig["mchID"].(string)),
				},
			)
			if err == nil {
				results.Success(cgin, "查询订单", result, resp)
			} else {
				results.Failed(cgin, "查询订单失败", err)
			}
		}
	}
}

// 获取当前服务器IP
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	return localAddr.IP.String()
}
func GetLocalIP() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, addr := range addrs {
		ipAddr, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if ipAddr.IP.IsLoopback() {
			continue
		}
		if !ipAddr.IP.IsGlobalUnicast() {
			continue
		}
		return ipAddr.IP.String(), nil
	}
	return
}
