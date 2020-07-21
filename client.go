package tsing_center_go

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
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
	TTL           uint   // 生命周期(秒)
	TouchInterval uint   // 自动触活的间隔时间(秒)
	Timeout       uint   // 操作超时时间(秒)
}

type ErrorResponse struct {
	Text string `json:"error"`
}

// 接收自动触活时的错误的处理器(状态码，错误文本)
type AutoTouchErrorHandler func(int, error)

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

// 获取当前节点的IP地址
func (self *Client) GetIP() (string, int, error) {
	var (
		err    error
		req    *http.Request
		resp   *http.Response
		hc     http.Client
		body   []byte
		status int
	)
	ctx, cancel := context.WithTimeout(context.Background(), self.timeout)
	defer cancel()
	req, err = http.NewRequestWithContext(ctx, "GET", self.addr+"/ip", nil)
	if err != nil {
		return "", 500, err
	}
	req.Header.Set("SECRET", self.secret)

	resp, err = hc.Do(req)
	if err != nil {
		return "", 500, err
	}

	body, status, err = parseBody(resp)
	if err != nil {
		return "", status, err
	}
	return string(body), 200, nil
}

// 解析响应body
func parseBody(resp *http.Response) ([]byte, int, error) {
	switch resp.StatusCode {
	case 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, 500, err
		}
		return body, resp.StatusCode, nil
	case 204:
		return nil, resp.StatusCode, nil
	case 400:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, 500, err
		}
		var e ErrorResponse
		if err = json.Unmarshal(body, &e); err != nil {
			return nil, 500, err
		}
		return nil, resp.StatusCode, errors.New(e.Text)
	case 401:
		return nil, resp.StatusCode, errors.New(http.StatusText(http.StatusUnauthorized))
	case 404:
		return nil, resp.StatusCode, errors.New(http.StatusText(http.StatusNotFound))
	case 500:
		return nil, resp.StatusCode, errors.New(http.StatusText(http.StatusInternalServerError))
	case 501:
		return nil, resp.StatusCode, errors.New(http.StatusText(http.StatusNotImplemented))
	default:
		return nil, resp.StatusCode, errors.New("意外的响应")
	}
}
