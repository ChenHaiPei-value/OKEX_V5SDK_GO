package main

import (
	"fmt"
	"log"
	"time"
	"sync"
	. "okex_v5sdk_go/ws"
	. "okex_v5sdk_go/ws/wImpl"
	"encoding/json"
	"io/ioutil"
)


// 定义与JSON结构相对应的Go结构体
type MjsonConfig struct {
	APIKey       string `json:"apiKey"`
	SecretKey    string `json:"secretKey"`
	Passphrase   string `json:"passphrase"`
	MEndPoint   string `json:"mEndPoint"`
	DelayTime    int    `json:"delay_time"`
	WaitTime     int    `json:"whait_time"`
	WaitTimes    int    `json:"whait_times"`
	FollowAccounts []FollowAccount `json:"followAccounts"`
	Telegram     TelegramConfig  `json:"telegram"`
}

type FollowAccount struct {
	APIKey      string `json:"apiKey"`
	SecretKey   string `json:"secretKey"`
	Passphrase  string `json:"passphrase"`
	FollowRatio int    `json:"followRatio"`
}

type TelegramConfig struct {
	ChatIDInfo string `json:"chatID_Info"`
	TokenInfo  string `json:"token_Info"`
	ChatIDError string `json:"chatID_Error"`
	TokenError  string `json:"token_Error"`
}

// 读取和解析JSON文件的函数
func LoadConfig(filename string) (MjsonConfig, error) {
	var config MjsonConfig
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
	return config, err
}

// 保存配置到JSON文件的函数
func SaveConfig(filename string, config MjsonConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

// 添加新的followAccount的函数
func AddFollowAccount(config *MjsonConfig, account FollowAccount) {
	config.FollowAccounts = append(config.FollowAccounts, account)
}

// 删除指定apiKey的followAccount的函数
func RemoveFollowAccount(config *MjsonConfig, apiKey string) {
	for i, account := range config.FollowAccounts {
		if account.APIKey == apiKey {
			config.FollowAccounts = append(config.FollowAccounts[:i], config.FollowAccounts[i+1:]...)
			break
		}
	}
}
/* 
// 示例：更新WebSocket实例的逻辑（伪代码，具体实现依赖于WebSocket库）
func UpdateWebSocketInstances(config jsonConfig) {
	// 遍历所有followAccounts并更新其WebSocket实例
	for _, account := range config.FollowAccounts {
		// 假设CreateWebSocketInstance是一个创建并返回WebSocket实例的函数
		// wsInstance := CreateWebSocketInstance(account)
		// 这里可以添加逻辑来管理这些WebSocket实例，例如保存在一个map中
	}
}
*/
/*  
func main() {
	// 加载配置
	config, err := LoadConfig("config.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	// 添加新的followAccount
	newAccount := FollowAccount{
		APIKey:      "new-api-key",
		SecretKey:   "new-secret-key",
		Passphrase:  "new-passphrase",
		FollowRatio: 10,
	}
	AddFollowAccount(&config, newAccount)

	// 删除一个followAccount
	RemoveFollowAccount(&config, "909cf8b1-265c-419c-8a70-9a71c80eca90")

	// 保存更新后的配置
	err = SaveConfig("config.json", config)
	if err != nil {
		fmt.Println("Error saving config:", err)
		return
	}

	// 更新WebSocket实例
	UpdateWebSocketInstances(config)
}
*/


var (
    signalClients = make(map[string]*WsClient)
    followClient  *WsClient
    mutex         sync.Mutex
)

func monitorSignalAccounts() {
    for apiKey, client := range signalClients {
        go func(apiKey string, client *WsClient) {
            for msg := range client.Messages {
                switch msg.Type {
                case "orders":
                    handleOrder(apiKey, msg.Info)
                case "balance_and_position_update":
                    handleBalanceAndPositionUpdate(apiKey, msg.Info)
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
        fmt.Println("Failed to place order in follow account: %v", err)
    } else {
        fmt.Println("Order placed successfully in follow account")
    }
}

// 登录和订阅多个信号账户
func con_login_sub_s(config *MjsonConfig) {
	for _, account := range config.FollowAccounts {
		if r, err := NewWsClient(config.MEndPoint); err == nil {
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

			
			res, _, err = r.PrivBalAndPos(OP_SUBSCRIBE, argsa)
			if res {
				
				fmt.Println("	订阅账户余额和持仓频道成功！耗时:")
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
func con_login_sub_f(config *MjsonConfig) {
	if r, err := NewWsClient(config.MEndPoint); err == nil {
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
		res, _, err = r.Login(config.APIKey, config.SecretKey, config.Passphrase)
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

		res, _, err = r.PrivBalAndPos(OP_SUBSCRIBE, argsa)
		if res {
			fmt.Println("跟单订阅账户余额和持仓频道成功！耗时:")
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
func loadWsClients(config *MjsonConfig) error {
    con_login_sub_s(config)
	con_login_sub_f(config)
    return nil
}

// 根据配置更新WebSocket实例
/* 
func updateWsClients(newConfig *jsonConfig) {
    // 遍历新的配置，添加新的实例或更新现有实例
    for _, newAccount := range newConfig.FollowAccounts {
        if client, exists := signalClients[newAccount.APIKey]; exists {
            // 更新现有实例的配置（如果需要）
            //client.UpdateConfig(newAccount)
        } else {
            // 添加新的实例
            if newClient, err := NewWsClient(newAccount.EndPoint); err == nil {
                signalClients[newAccount.APIKey] = newClient
                go newClient.Start()
                //newClient.Subscribe("order-book")
                //newClient.Subscribe("balance")
                // ... 订阅其他需要的频道
            }
        }
    }

    // 遍历现有实例，删除不再需要的实例
    for APIKey, client := range signalClients {
        found := false
        for _, account := range newConfig.FollowAccounts {
            if account.APIKey == APIKey {
                found = true
                break
            }
        }
        if !found {
            client.Stop()
            delete(signalClients, APIKey)
        }
    }
}

// 监控配置文件的变化
func watchConfigChanges(filePath string, onChange func(*jsonConfig)) {
    // ... 实现文件监控逻辑
	updateWsClients(newConfig)
}
*/

func main() {

	// 加载配置
	var config MjsonConfig
	config, err := LoadConfig("ws/config.json")
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


    // 阻塞主goroutine，防止程序退出
    select {}
	
}


