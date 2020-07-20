package tsing_center_go

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// 新建节点
func (self *Client) AddNode(serviceID, ip string, port uint16, weight uint) error {
	var (
		err  error
		req  *http.Request
		resp *http.Response
		hc   http.Client
	)
	if serviceID == "" {
		return errors.New("服务ID参数不能为空")
	}
	serviceID = base64.RawURLEncoding.EncodeToString([]byte(serviceID))
	if port == 0 {
		return errors.New("端口无效")
	}

	values := url.Values{}
	values.Add("service_id", serviceID)
	values.Add("ip", ip)
	values.Add("port", strconv.FormatUint(uint64(port), 10))
	values.Add("weight", strconv.FormatUint(uint64(weight), 10))
	values.Add("expires", strconv.FormatInt(time.Now().Add(self.ttl).Unix(), 10))

	ctx, cancel := context.WithTimeout(context.Background(), self.timeout)
	defer cancel()
	req, err = http.NewRequestWithContext(ctx, "POST", self.addr+"/nodes/", strings.NewReader(values.Encode()))
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

// 重写或新建节点
func (self *Client) SetNode(serviceID, node string, weight uint) error {
	var (
		err  error
		req  *http.Request
		resp *http.Response
		hc   http.Client
	)
	if serviceID == "" {
		return errors.New("服务ID参数不能为空")
	}
	serviceID = base64.RawURLEncoding.EncodeToString([]byte(serviceID))
	if node == "" {
		return errors.New("节点参数不能为空")
	}
	node = base64.RawURLEncoding.EncodeToString([]byte(node))

	values := url.Values{}
	values.Add("weight", strconv.FormatUint(uint64(weight), 10))
	values.Add("expires", strconv.FormatInt(time.Now().Add(self.ttl).Unix(), 10))

	ctx, cancel := context.WithTimeout(context.Background(), self.timeout)
	defer cancel()
	req, err = http.NewRequestWithContext(ctx, "PUT", self.addr+"/nodes/"+serviceID+"/"+node, strings.NewReader(values.Encode()))
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

// 删除节点
func (self *Client) RemoveNode(serviceID, node string) error {
	var (
		err  error
		req  *http.Request
		resp *http.Response
		hc   http.Client
	)
	if serviceID == "" {
		return errors.New("服务ID参数不能为空")
	}
	serviceID = base64.RawURLEncoding.EncodeToString([]byte(serviceID))
	if node == "" {
		return errors.New("节点参数不能为空")
	}
	node = base64.RawURLEncoding.EncodeToString([]byte(node))

	ctx, cancel := context.WithTimeout(context.Background(), self.timeout)
	defer cancel()
	req, err = http.NewRequestWithContext(ctx, "DELETE", self.addr+"/nodes/"+serviceID+"/"+node, nil)
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

// node触活
func (self *Client) TouchNode(serviceID, node string) error {
	var (
		err     error
		req     *http.Request
		resp    *http.Response
		expires int64
		hc      http.Client
	)
	if serviceID == "" {
		return errors.New("服务ID参数不能为空")
	}
	serviceID = base64.RawURLEncoding.EncodeToString([]byte(serviceID))
	if node == "" {
		return errors.New("节点参数不能为空")
	}
	node = base64.RawURLEncoding.EncodeToString([]byte(node))

	expires = time.Now().Add(self.ttl).Unix()

	values := url.Values{}
	values.Add("expires", strconv.FormatInt(expires, 10))
	ctx, cancel := context.WithTimeout(context.Background(), self.timeout)
	defer cancel()
	req, err = http.NewRequestWithContext(ctx, "PATCH", self.addr+"/nodes/"+serviceID+"/"+node+"/expires", strings.NewReader(values.Encode()))
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

// node自动触活
func (self *Client) AutoTouchNode(serviceID, node string, errorHandler AutoTouchErrorHandler) {
	for {
		time.Sleep(self.touchInterval)
		if err := self.TouchNode(serviceID, node); err != nil {
			errorHandler(err)
		}
	}
}
