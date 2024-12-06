package main

import (
	"context"
	"fmt"
	"log"
	"time"
	. "v5sdk_go/ws"
)

// 订阅私有频道
func wsPriv() {
	ep := "wss://wspap.okx.com:8443/ws/v5/private"

	// 填写您自己的APIKey信息
	apikey := "de6607bf-f39e-4781-83e3-b9b338319ae3"
	secretKey := "2A13AF232E1AE4F4C105AA25B3645C01"
	passphrase := "@Zxjchp1314520"

	// 创建ws客户端
	r, err := NewWsClient(ep)
	if err != nil {
		log.Println(err)
		return
	}

	// 设置连接超时
	r.SetDailTimeout(time.Second * 2)
	err = r.Start()
	if err != nil {
		log.Println(err)
		return
	}
	defer r.Stop()
	var res bool

	res, _, err = r.Login(apikey, secretKey, passphrase)
	if res {
		fmt.Println("登录成功！")
	} else {
		fmt.Println("登录失败！", err)
		return
	}

	// 订阅账户频道
	var args []map[string]string
	arg := make(map[string]string)
	arg["ccy"] = "BTC"
	args = append(args, arg)

	start := time.Now()
	res, _, err = r.PrivAccout(OP_SUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("订阅成功！耗时:", usedTime.String())
	} else {
		fmt.Println("订阅失败！", err)
	}

	time.Sleep(100 * time.Second)
	start = time.Now()
	res, _, err = r.PrivAccout(OP_UNSUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("取消订阅成功！", usedTime.String())
	} else {
		fmt.Println("取消订阅失败！", err)
	}

}

// websocket交易
func wsJrpc() {
	ep := "wss://wspap.okx.com:8443/ws/v5/private"

	// 填写您自己的APIKey信息
	// 填写您自己的APIKey信息
	apikey := "de6607bf-f39e-4781-83e3-b9b338319ae3"
	secretKey := "2A13AF232E1AE4F4C105AA25B3645C01"
	passphrase := "@Zxjchp1314520"

	var res bool
	var req_id string

	// 创建ws客户端
	r, err := NewWsClient(ep)
	if err != nil {
		log.Println(err)
		return
	}

	// 设置连接超时
	r.SetDailTimeout(time.Second * 2)
	err = r.Start()
	if err != nil {
		log.Println(err)
		return
	}

	defer r.Stop()

	res, _, err = r.Login(apikey, secretKey, passphrase)
	if res {
		fmt.Println("登录成功！")
	} else {
		fmt.Println("登录失败！", err)
		return
	}

	start := time.Now()
	param := map[string]interface{}{}
	param["instId"] = "BTC-USDT-SWAP"
	param["tdMode"] = "cross"
	param["side"] = "buy"
	param["ordType"] = "market"
	param["sz"] = "0.1"
	req_id = "00001"

	res, _, err = r.PlaceOrder(req_id, param)
	if res {
		usedTime := time.Since(start)
		fmt.Println("下单成功！", usedTime.String())
	} else {
		usedTime := time.Since(start)
		fmt.Println("下单失败！", usedTime.String(), err)
	}
}

func main() {

	// 私有订阅
	wsPriv()

	// websocket交易
	wsJrpc()


}