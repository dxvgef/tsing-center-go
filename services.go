package tsing_center_go

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/dxvgef/filter/v2"
)

type Service struct {
	ID          string
	LoadBalance string
	Meta        string
}

// 新建服务
func (self *Client) AddService(service Service) (int, error) {
	var (
		err    error
		req    *http.Request
		resp   *http.Response
		hc     http.Client
		status int
	)
	if err = filter.Batch(
		filter.String(service.ID, "服务ID").Require().Error(),
		filter.String(service.LoadBalance, "负载均衡算法").Require().EnumSliceString(",", []string{"SWRR", "WRR", "WR"}).Error(),
	); err != nil {
		return 400, err
	}

	values := url.Values{}
	values.Add("id", service.ID)
	values.Add("load_balance", service.LoadBalance)

	ctx, cancel := context.WithTimeout(context.Background(), self.timeout)
	defer cancel()
	req, err = http.NewRequestWithContext(ctx, "POST", self.addr+"/services/", strings.NewReader(values.Encode()))
	if err != nil {
		return 500, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("SECRET", self.secret)

	resp, err = hc.Do(req)
	if err != nil {
		return 500, err
	}

	_, status, err = parseBody(resp)

	return status, err
}

// 重写或新建服务
func (self *Client) SetService(service Service) (int, error) {
	var (
		err    error
		req    *http.Request
		resp   *http.Response
		hc     http.Client
		status int
	)
	if err = filter.Batch(
		filter.String(service.ID, "服务ID").Require().Error(),
		filter.String(service.LoadBalance, "负载均衡算法").Require().EnumSliceString(",", []string{"SWRR", "WRR", "WR"}).Error(),
	); err != nil {
		return 400, err
	}

	service.ID = base64.RawURLEncoding.EncodeToString([]byte(service.ID))

	values := url.Values{}
	values.Add("load_balance", service.LoadBalance)
	ctx, cancel := context.WithTimeout(context.Background(), self.timeout)
	defer cancel()

	req, err = http.NewRequestWithContext(ctx, "PUT", self.addr+"/services/"+service.ID, strings.NewReader(values.Encode()))
	if err != nil {
		return 500, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("SECRET", self.secret)

	resp, err = hc.Do(req)
	if err != nil {
		return 500, err
	}

	_, status, err = parseBody(resp)

	return status, err
}

// 移除服务
func (self *Client) RemoveService(id string) (int, error) {
	var (
		err    error
		req    *http.Request
		resp   *http.Response
		hc     http.Client
		status int
	)
	if id == "" {
		return 400, errors.New("服务ID参数无效")
	}
	id = base64.RawURLEncoding.EncodeToString([]byte(id))

	ctx, cancel := context.WithTimeout(context.Background(), self.timeout)
	defer cancel()
	req, err = http.NewRequestWithContext(ctx, "DELETE", self.addr+"/services/"+id, nil)
	if err != nil {
		return 500, err
	}

	req.Header.Set("SECRET", self.secret)
	resp, err = hc.Do(req)
	if err != nil {
		return 500, err
	}

	_, status, err = parseBody(resp)

	return status, err
}

// 发现服务(根据服务ID获取节点，返回ip, port)
func (self *Client) DiscoverService(serviceID string) (Node, int, error) {
	var (
		err    error
		req    *http.Request
		resp   *http.Response
		hc     http.Client
		body   []byte
		node   Node
		status int
	)
	if serviceID == "" {
		return node, 400, errors.New("服务ID参数无效")
	}
	serviceID = base64.RawURLEncoding.EncodeToString([]byte(serviceID))

	ctx, cancel := context.WithTimeout(context.Background(), self.timeout)
	defer cancel()
	req, err = http.NewRequestWithContext(ctx, "GET", self.addr+"/services/"+serviceID+"/select", nil)
	if err != nil {
		return node, 500, err
	}
	req.Header.Set("SECRET", self.secret)

	resp, err = hc.Do(req)
	if err != nil {
		return node, 500, err
	}
	if body, status, err = parseBody(resp); err != nil {
		return node, status, err
	}
	if status == 200 {
		if err = json.Unmarshal(body, &node); err != nil {
			return node, 500, err
		}
	}
	return node, status, nil
}
