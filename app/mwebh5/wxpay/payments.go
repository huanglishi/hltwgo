package wxpay

import (
	"context"
	"fmt"
	"huling/utils/results"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/certificates"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

// 1微信支付平台证书
func Getwechatpay(cgin *gin.Context) {
	const (
		mchID                      string = "190000****"                               // 商户号
		mchCertificateSerialNumber string = "3775B6A45ACD588826D15E583A95F5DD********" // 商户证书序列号
		mchAPIv3Key                string = "2ab9****************************"         // 商户APIv3密钥
	)
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

	// 发送请求，以下载微信支付平台证书为例
	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay5_1.shtml
	svc := certificates.CertificatesApiService{Client: client}
	resp, result, err := svc.DownloadCertificates(ctx)
	log.Printf("status=%d resp=%s", result.Response.StatusCode, resp)
	// 通过私钥的文件路径内容加载私钥
	results.Success(cgin, "获取页面数据", string(""), nil)
}

// 2 JSAPI下单 为例
func Getjsapi(cgin *gin.Context) {
	const (
		mchID                      string = "190000****"                               // 商户号
		mchCertificateSerialNumber string = "3775B6A45ACD588826D15E583A95F5DD********" // 商户证书序列号
		mchAPIv3Key                string = "2ab9****************************"         // 商户APIv3密钥
	)

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
	svc := jsapi.JsapiApiService{Client: client}
	// 得到prepay_id，以及调起支付所需的参数和签名
	resp, result, err := svc.PrepayWithRequestPayment(ctx,
		jsapi.PrepayRequest{
			Appid:       core.String("wxd678efh567hg6787"),
			Mchid:       core.String("1900009191"),
			Description: core.String("Image形象店-深圳腾大-QQ公仔"),
			OutTradeNo:  core.String("1217752501201407033233368018"),
			Attach:      core.String("自定义数据说明"),
			NotifyUrl:   core.String("https://www.weixin.qq.com/wxpay/pay.php"),
			Amount: &jsapi.Amount{
				Total: core.Int64(100),
			},
			Payer: &jsapi.Payer{
				Openid: core.String("oUpF8uMuAJO_M2pxb1Q9zNjWeS6o"),
			},
		},
	)
	if err == nil {
		log.Println(resp)
	} else {
		log.Println(err)
	}
	// 通过私钥的文件路径内容加载私钥
	results.Success(cgin, "获取页面数据", result, nil)
}

// 3获取h5支付跳转链接
func Geth5url(cgin *gin.Context) {
	const (
		mchID                      string = "190000****"                               // 商户号
		mchCertificateSerialNumber string = "3775B6A45ACD588826D15E583A95F5DD********" // 商户证书序列号
		mchAPIv3Key                string = "2ab9****************************"         // 商户APIv3密钥
	)
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
	results.Success(cgin, "获取页面数据", result, resp)
}

// 4 基础支付的回调通知
func Paynotify(cgin *gin.Context) {
	const (
		mchID                      string = "190000****"                               // 商户号
		mchCertificateSerialNumber string = "3775B6A45ACD588826D15E583A95F5DD********" // 商户证书序列号
		mchAPIv3Key                string = "2ab9****************************"         // 商户APIv3密钥
	)

	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath("/resource/staticfile/merchant/apiclient_key.pem")
	if err != nil {
		log.Fatal("load merchant private key error")
	}
	//支付回调
	fmt.Println("-----支付回调---------")
	/************开始支付逻辑************************/
	ctx := context.Background()
	// 1. 使用 `RegisterDownloaderWithPrivateKey` 注册下载器
	derr := downloader.MgrInstance().RegisterDownloaderWithPrivateKey(ctx, mchPrivateKey, mchCertificateSerialNumber, mchID, mchAPIv3Key)
	// 2. 获取商户号对应的微信支付平台证书访问器
	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(mchID)
	// 3. 使用证书访问器初始化 `notify.Handler`
	handler := notify.NewNotifyHandler(mchAPIv3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))

	transaction := new(payments.Transaction)
	notifyReq, err := handler.ParseNotifyRequest(context.Background(), cgin.Request, transaction)
	// 如果验签未通过，或者解密失败
	if err != nil {
		fmt.Println(err)
		return
	}
	// 处理通知内容
	fmt.Println(notifyReq.Summary)
	fmt.Println(transaction.TransactionId)
	results.Success(cgin, "获取页面数据", handler, derr)
}
