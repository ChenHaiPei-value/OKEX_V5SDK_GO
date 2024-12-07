package ws

import (
	"encoding/json"

	"io/ioutil"
)

// 定义与JSON结构相对应的Go结构体
type jsonConfig struct {
	APIKey       string `json:"apiKey"`
	SecretKey    string `json:"secretKey"`
	Passphrase   string `json:"passphrase"`
	EndPoint   string `json:"EndPoint"`
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
func LoadConfig(filename string) (jsonConfig, error) {
	var config jsonConfig
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
	return config, err
}

// 保存配置到JSON文件的函数
func SaveConfig(filename string, config jsonConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

// 添加新的followAccount的函数
func AddFollowAccount(config *jsonConfig, account FollowAccount) {
	config.FollowAccounts = append(config.FollowAccounts, account)
}

// 删除指定apiKey的followAccount的函数
func RemoveFollowAccount(config *jsonConfig, apiKey string) {
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
