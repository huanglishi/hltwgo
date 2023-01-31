package wxpay

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

// 4 基础支付的回调通知
func Paynotify(cgin *gin.Context) {
	idstr := cgin.Param("id")
	string_slice := strings.Split(idstr, "/")
	cuid, _ := DB().Table("client_product_order").Where("id", string_slice[1]).Value("cuid")
	paymentconfig, _ := DB().Table("client_system_paymentconfig").Where("cuid", cuid).Fields("id,appId,mchID,mchAPIv3Key,mchCertificateSerialNumber,privatekey").First()
	//配置参数
	mchID := paymentconfig["mchID"].(string)                                           // 商户号
	mchCertificateSerialNumber := paymentconfig["mchCertificateSerialNumber"].(string) // 商户证书序列号
	mchAPIv3Key := paymentconfig["mchAPIv3Key"].(string)
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	pemPath := fmt.Sprintf("./%s", paymentconfig["privatekey"])
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(pemPath)
	if err != nil {
		log.Fatal("load merchant private key error")
		cgin.JSON(500, gin.H{
			"code":    "FAIL",
			"message": "失败",
		})
		return
	}
	//支付回调
	fmt.Println("-----支付回调---------")
	/************开始支付逻辑************************/
	ctx := context.Background()
	// 1. 使用 `RegisterDownloaderWithPrivateKey` 注册下载器
	derr := downloader.MgrInstance().RegisterDownloaderWithPrivateKey(ctx, mchPrivateKey, mchCertificateSerialNumber, mchID, mchAPIv3Key)
	if derr != nil {
		fmt.Println("-----注册下载器失败---------")
		fmt.Println(derr)
		cgin.JSON(500, gin.H{
			"code":    "FAIL",
			"message": "失败",
		})
	} else {
		// 2. 获取商户号对应的微信支付平台证书访问器
		certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(mchID)
		// 3. 使用证书访问器初始化 `notify.Handler`
		handler := notify.NewNotifyHandler(mchAPIv3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
		transaction := new(payments.Transaction)
		_, err := handler.ParseNotifyRequest(context.Background(), cgin.Request, transaction)
		// 如果验签未通过，或者解密失败
		if err != nil {
			fmt.Println("如果验签未通过，或者解密失败")
			fmt.Println(err)
			cgin.JSON(500, gin.H{
				"code":    "FAIL",
				"message": "失败",
			})
		} else { //支付成功
			DB().Table("logs").Data(map[string]interface{}{
				"path":       "处理通知内容",
				"createtime": time.Now().Unix(),
				"param":      JSONMarshalToString(transaction),
			}).Insert()
			//更新支付单号
			paytime := time.Now().Unix()
			DB().Table("client_product_order").Where("out_trade_no", transaction.OutTradeNo).Data(map[string]interface{}{"transaction_id": transaction.TransactionId, "total_fee": transaction.Amount.Total, "status": 1, "paytime": paytime, "time_end": paytime}).Update()
			cgin.JSON(200, nil) //支付成功-通知应答
		}
	}

}

// 4 基础支付的回调通知
func Paynotify1(cgin *gin.Context) {
	idstr := cgin.Param("id")
	string_slice := strings.Split(idstr, "/")
	cgin.JSON(200, gin.H{
		"code":    1,
		"message": "获取参数",
		"result":  string_slice[1],
		"time":    time.Now().Unix(),
	})
}
