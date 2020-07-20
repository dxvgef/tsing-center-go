package tsing_center_go

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// 客户端实例
type Client struct {
	addr          string        // 服务端地址
	secret        string        // 服务端api调用密码
	ttl           time.Duration // 生命周期(秒)
	touchInterval time.Duration // 自动触活的间隔时间(秒)
	timeout       time.Duration // 操作超时时间(秒)
}

// 实例配置
type Config struct {
	Addr          string // 服务端地址
	Secret        string // 服务端api调用密码
	TTL           int    // 生命周期(秒)
	TouchInterval int    // 自动触活的间隔时间(秒)
	Timeout       int    // 操作超时时间(秒)
}

// 节点信息
type Node struct {
	IP   string `json:"ip"`
	Port uint16 `json:"port"`
}

// 接收自动触活时的错误的处理器
type AutoTouchErrorHandler func(err error)

// 新建客户端实例
func New(config Config) (*Client, error) {
	if config.TouchInterval < config.TTL-2 {
		return nil, errors.New("touchInterval参数至少要比ttl参数小2秒")
	}
	var cli Client
	cli.addr = config.Addr
	cli.secret = config.Secret
	cli.ttl = time.Duration(config.TTL) * time.Second
	cli.touchInterval = time.Duration(config.TouchInterval) * time.Second
	cli.timeout = time.Duration(config.Timeout) * time.Second
	// todo 这里要对addr做一些格式检查
	return &cli, nil
}

// 解析响应
func parseResp(resp *http.Response) (map[string]string, error) {
	respData := make(map[string]string)

	switch resp.StatusCode {
	case 204:
		return respData, nil
	case 401:
		return respData, errors.New(http.StatusText(http.StatusUnauthorized))
	case 500:
		return respData, errors.New(http.StatusText(http.StatusInternalServerError))
	case 501:
		return respData, errors.New(http.StatusText(http.StatusNotImplemented))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return respData, err
	}

	if resp.StatusCode == 200 {
		var node Node
		if err = json.Unmarshal(body, &node); err != nil {
			return respData, err
		}
		respData["ip"] = node.IP
		respData["port"] = strconv.FormatUint(uint64(node.Port), 10)
		return respData, nil
	}

	if resp.StatusCode == 400 {
		return respData, errors.New(respData["error"])
	}

	return respData, nil
}
