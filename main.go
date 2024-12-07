package main

import (
	"fmt"
	"log"
	"time"
	"os"
	"sync"
	. "v5sdk_go/ws"
)

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
		fmt.Println("websocket交易登录成功！")
	} else {
		fmt.Println("websocket交易登录失败！", err)
		return
	}

	start := time.Now()
	param := map[string]interface{}{}
	param["instId"] = "BTC-USDT-SWAP"
	param["tdMode"] = "cross"
	param["side"] = "sell"
	param["posSide"] = "short"
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

var (
    signalClients = make(map[string]*WsClient)
    followClient  *WsClient
    mutex         sync.Mutex
)

func monitorSignalAccounts() {
    for apiKey, client := range signalClients {
        go func(apiKey string, client *WsClient) {
            for msg := range client.Messages() {
                switch msg.Type {
                case "orders":
                    handleOrder(apiKey, msg.Info)
                case "balance_and_position_update":
                    handleBalanceAndPositionUpdate(apiKey, msg.Data)
                // 处理其他类型的更新
                }
            }
        }(apiKey, client)
    }
}

func handleOrder(apiKey string, orderData interface{}) {
    // 解析订单数据，并根据需要在跟单账户上下单或撤单
    // ...
    //placeOrderInFollowAccount(orderData)
}

func handleBalanceAndPositionUpdate(apiKey string, balanceAndPositionData interface{}) {
    // 解析持仓数据，并根据需要在跟单账户上调整持仓
    // ...
}

func placeOrderInFollowAccount(orderData interface{}) {
    // 根据订单数据在跟单账户上下单
    // ...
    param := map[string]interface{}{
        // 设置订单参数
    }
    _, _, err := followClient.PlaceOrder("00001", param)
    if err != nil {
        log.Errorf("Failed to place order in follow account: %v", err)
    } else {
        log.Info("Order placed successfully in follow account")
    }
}

// 登录和订阅多个信号账户
func con_login_sub_s(config *jsonConfig) {
	for _, account := range config.Accounts {
		if r, err := NewWsClient(config.EndPoint); err == nil {
			signalClients[account.APIKey] = r
			// 启动客户端并订阅必要的频道
			// 设置连接超时
			r.SetDailTimeout(time.Second * 2)
			err = r.Start()
			if err != nil {
				log.Println(err)
				return
			}
			defer r.Stop()
			var res bool
			var err error

			res, _, err = r.Login(account.APIKey, account.SecretKey, account.Passphrase)
			if res {
				fmt.Println("登录成功！")
			} else {
				fmt.Println("登录失败！", err)
				return
			}
			// 添加自定义消息钩子
			/* 
			r.AddBookMsgHook(func(ts time.Time, data MsgData) error {
				// 添加你的方法
				fmt.Println("这是自定义AddBookMsgHook")
				fmt.Println("当前数据是", data)
				return nil
			})*/

			var args []map[string]string
			arg := make(map[string]string)
			arg["channel"] = "orders"
			arg["instType"] = "SWAP"
			args = append(args, arg)

			start := time.Now()
			res, _, err = r.PrivBookOrder(OP_SUBSCRIBE, args)
			if res {
				usedTime := time.Since(start)
				fmt.Println("订阅订单频道成功！耗时:", usedTime.String())
			} else {
				fmt.Println("订阅订单频道失败！", err)
			}

			// 订阅账户余额和持仓频道
			var argsa []map[string]string
			arga := make(map[string]string)
			arga["channel"] = "balance_and_position"
			argsa = append(argsa, arga)

			start := time.Now()
			res, _, err = r.PrivBalAndPos(OP_SUBSCRIBE, argsa)
			if res {
				usedTime := time.Since(start)
				fmt.Println("	订阅账户余额和持仓频道成功！耗时:", usedTime.String())
			} else {
				fmt.Println("	订阅账户余额和持仓频道失败！", err)
			}


		} else {
			// 处理错误
			log.Println(err)
			return
		}
	}
}
// 跟单登录和订阅
func con_login_sub_f(config *jsonConfig) {
	if r, err := NewWsClient(config.EndPoint); err == nil {
		followClient = r
		// 启动客户端并订阅必要的频道
		// 设置连接超时
		r.SetDailTimeout(time.Second * 2)
		err = r.Start()
		if err != nil {
			log.Println(err)
			return
		}
		defer r.Stop()
		var res bool
		var err error
		res, _, err = r.Login(account.APIKey, account.SecretKey, account.Passphrase)
		if res {
			fmt.Println("跟单登录成功！")
		} else {
			fmt.Println("跟单登录失败！", err)
			return
		}
		// 添加自定义消息钩子
		r.AddBookMsgHook(func(ts time.Time, data MsgData) error {
			// 添加你的方法
			fmt.Println("这是自定义AddBookMsgHook")
			fmt.Println("当前数据是", data)
			return nil
		})

		// 订阅订单频道
		var args []map[string]string
		arg := make(map[string]string)
		arg["channel"] = "orders"
		arg["instType"] = "SWAP"
		args = append(args, arg)

		start := time.Now()
		res, _, err = r.PrivBookOrder(OP_SUBSCRIBE, args)
		if res {
			usedTime := time.Since(start)
			fmt.Println("跟单订阅订单频道成功！耗时:", usedTime.String())
		} else {
			fmt.Println("跟单订阅订单频道失败！", err)
		}

		// 订阅账户余额和持仓频道
		var argsa []map[string]string
		arga := make(map[string]string)
		arga["channel"] = "balance_and_position"
		argsa = append(argsa, arga)

		start := time.Now()
		res, _, err = r.PrivBalAndPos(OP_SUBSCRIBE, argsa)
		if res {
			usedTime := time.Since(start)
			fmt.Println("跟单订阅账户余额和持仓频道成功！耗时:", usedTime.String())
		} else {
			fmt.Println("跟单订阅账户余额和持仓频道失败！", err)
		}


	} else {
		// 处理错误
		log.Println(err)
		return
	}

}

// 根据配置加载WebSocket实例
func loadWsClients(config *jsonConfig) error {
    con_login_sub_s(config)
	con_login_sub_f(config)
    return nil
}

// 根据配置更新WebSocket实例
func updateWsClients(newConfig *jsonConfig) {
    // 遍历新的配置，添加新的实例或更新现有实例
    for _, newAccount := range newConfig.Accounts {
        if client, exists := wsClients[newAccount.ID]; exists {
            // 更新现有实例的配置（如果需要）
            client.UpdateConfig(newAccount)
        } else {
            // 添加新的实例
            if newClient, err := NewWsClient(newAccount.EndPoint); err == nil {
                wsClients[newAccount.ID] = newClient
                go newClient.Start()
                newClient.Subscribe("order-book")
                newClient.Subscribe("balance")
                // ... 订阅其他需要的频道
            }
        }
    }

    // 遍历现有实例，删除不再需要的实例
    for id, client := range wsClients {
        found := false
        for _, account := range newConfig.Accounts {
            if account.ID == id {
                found = true
                break
            }
        }
        if !found {
            client.Stop()
            delete(wsClients, id)
        }
    }
}

// 监控配置文件的变化
func watchConfigChanges(filePath string, onChange func(*jsonConfig)) {
    // ... 实现文件监控逻辑
}

func main() {

	// 加载配置
	var config jsonConfig
	config, err := LoadConfig("config.json")
	if err != nil {
		log.Println("Error loading config:", err)
		return
	}
	// 加载WebSocket实例
	if err := loadWsClients(&config); err != nil {
        log.Fatalf("Failed to load WebSocket clients: %v", err)
    }

	 // 监控信号账户的更新
	 monitorSignalAccounts()

    // 监控配置文件的变化，并在变化时更新WebSocket实例
    watchConfigChanges("config.json", func(newConfig *jsonConfig) {
        updateWsClients(newConfig)
    })

    // 阻塞主goroutine，防止程序退出
    select {}
	
}


