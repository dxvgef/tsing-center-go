package tsing_center_go

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

// 新建服务
func (self *Client) AddService(id, loadBalance string) error {
	var (
		err  error
		req  *http.Request
		resp *http.Response
		hc   http.Client
	)
	if id == "" {
		return errors.New("服务ID参数不能为空")
	}
	id = base64.RawURLEncoding.EncodeToString([]byte(id))
	if loadBalance != "SWRR" && loadBalance != "WRR" && loadBalance != "WR" {
		return errors.New("不支持的负载算法")
	}

	values := url.Values{}
	values.Add("id", id)
	values.Add("load_balance", loadBalance)

	ctx, cancel := context.WithTimeout(context.Background(), self.timeout)
	defer cancel()
	req, err = http.NewRequestWithContext(ctx, "POST", self.addr+"/services/", strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("SECRET", self.secret)

	resp, err = hc.Do(req)
	if err != nil {
		return err
	}

	_, err = parseResp(resp)

	return err
}

// 重写或新建服务
func (self *Client) SetService(id, loadBalance string) error {
	var (
		err  error
		req  *http.Request
		resp *http.Response
		hc   http.Client
	)
	if id == "" {
		return errors.New("服务ID参数不能为空")
	}
	id = base64.RawURLEncoding.EncodeToString([]byte(id))
	if loadBalance != "SWRR" && loadBalance != "WRR" && loadBalance != "WR" {
		return errors.New("不支持的负载算法")
	}
	values := url.Values{}
	values.Add("load_balance", loadBalance)
	ctx, cancel := context.WithTimeout(context.Background(), self.timeout)
	defer cancel()

	req, err = http.NewRequestWithContext(ctx, "PUT", self.addr+"/services/"+id, strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("SECRET", self.secret)

	resp, err = hc.Do(req)
	if err != nil {
		return err
	}

	_, err = parseResp(resp)

	return err
}

// 移除服务
func (self *Client) RemoveService(id, loadBalance string) error {
	var (
		err  error
		req  *http.Request
		resp *http.Response
		hc   http.Client
	)
	if id == "" {
		return errors.New("服务ID参数不能为空")
	}
	id = base64.RawURLEncoding.EncodeToString([]byte(id))
	ctx, cancel := context.WithTimeout(context.Background(), self.timeout)
	defer cancel()

	req, err = http.NewRequestWithContext(ctx, "DELETE", self.addr+"/services/"+id, nil)
	if err != nil {
		return err
	}
	req.Header.Set("SECRET", self.secret)

	resp, err = hc.Do(req)
	if err != nil {
		return err
	}

	_, err = parseResp(resp)

	return err
}

// 发现服务(根据服务ID获取节点，返回ip, port)
func (self *Client) DiscoverService(serviceID string) (string, string, error) {
	var (
		err      error
		req      *http.Request
		resp     *http.Response
		hc       http.Client
		respData map[string]string
	)
	if serviceID == "" {
		return "", "", errors.New("服务ID参数不能为空")
	}
	serviceID = base64.RawURLEncoding.EncodeToString([]byte(serviceID))

	ctx, cancel := context.WithTimeout(context.Background(), self.timeout)
	defer cancel()
	req, err = http.NewRequestWithContext(ctx, "GET", self.addr+"/services/"+serviceID+"/select", nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("SECRET", self.secret)

	resp, err = hc.Do(req)
	if err != nil {
		return "", "", err
	}
	if respData, err = parseResp(resp); err != nil {
		return "", "", err
	}

	return respData["ip"], respData["port"], nil
}
