package tsing_center_go

import (
	"context"
	"encoding/base64"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dxvgef/filter/v2"
)

type Node struct {
	IP     string
	Port   uint16
	Weight uint
	TTL    uint
	Meta   string
}

// 新建节点
func (self *Client) AddNode(serviceID string, node Node) (int, error) {
	var (
		err    error
		req    *http.Request
		resp   *http.Response
		hc     http.Client
		status int
	)
	if err = filter.Batch(
		filter.String(serviceID, "服务ID").Require().Error(),
		filter.String(node.IP, "节点IP").Require().IsIP().Error(),
		filter.String(strconv.FormatUint(uint64(node.Port), 10), "节点端口").Require().IsDigit().MinInteger(1).MaxInteger(math.MaxUint16).Error(),
		filter.String(node.Meta, "节点元信息").IsJSON().Error(),
	); err != nil {
		return 400, err
	}

	values := url.Values{}
	values.Add("service_id", serviceID)
	values.Add("ip", node.IP)
	values.Add("port", strconv.FormatUint(uint64(node.Port), 10))
	values.Add("weight", strconv.FormatUint(uint64(node.Weight), 10))
	values.Add("ttl", strconv.FormatUint(uint64(node.TTL), 10))
	values.Add("meta", node.Meta)

	ctx, cancel := context.WithTimeout(context.Background(), self.timeout)
	defer cancel()
	req, err = http.NewRequestWithContext(ctx, "POST", self.addr+"/nodes/", strings.NewReader(values.Encode()))
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

// 重写或新建节点
func (self *Client) SetNode(serviceID string, node Node) (int, error) {
	var (
		err     error
		req     *http.Request
		resp    *http.Response
		hc      http.Client
		status  int
		nodeStr string
	)
	if err = filter.Batch(
		filter.String(serviceID, "服务ID").Require().Error(),
		filter.String(node.IP, "节点IP").Require().IsIP().Error(),
		filter.String(strconv.FormatUint(uint64(node.Port), 10), "节点端口").Require().IsDigit().MinInteger(1).MaxInteger(math.MaxUint16).Error(),
		filter.String(node.Meta, "节点元信息").IsJSON().Error(),
	); err != nil {
		return 400, err
	}
	serviceID = base64.RawURLEncoding.EncodeToString([]byte(serviceID))
	nodeStr = base64.RawURLEncoding.EncodeToString([]byte(node.IP + ":" + strconv.FormatUint(uint64(node.Port), 10)))

	values := url.Values{}
	values.Add("weight", strconv.FormatUint(uint64(node.Weight), 10))
	values.Add("ttl", strconv.FormatUint(uint64(node.TTL), 10))
	values.Add("meta", node.Meta)

	ctx, cancel := context.WithTimeout(context.Background(), self.timeout)
	defer cancel()
	log.Println(self.addr + "/nodes/" + serviceID + "/" + nodeStr)
	req, err = http.NewRequestWithContext(ctx, "PUT", self.addr+"/nodes/"+serviceID+"/"+nodeStr, strings.NewReader(values.Encode()))
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

// 删除节点
func (self *Client) RemoveNode(serviceID, ip string, port uint16) (int, error) {
	var (
		err     error
		req     *http.Request
		resp    *http.Response
		hc      http.Client
		status  int
		nodeStr string
	)

	if err = filter.Batch(
		filter.String(serviceID, "服务ID").Require().Error(),
		filter.String(ip, "节点IP").Require().IsIP().Error(),
		filter.String(strconv.FormatUint(uint64(port), 10), "节点端口").Require().IsDigit().MinInteger(1).MaxInteger(math.MaxUint16).Error(),
	); err != nil {
		return 400, err
	}
	serviceID = base64.RawURLEncoding.EncodeToString([]byte(serviceID))
	nodeStr = base64.RawURLEncoding.EncodeToString([]byte(ip + ":" + strconv.FormatUint(uint64(port), 10)))

	ctx, cancel := context.WithTimeout(context.Background(), self.timeout)
	defer cancel()
	req, err = http.NewRequestWithContext(ctx, "DELETE", self.addr+"/nodes/"+serviceID+"/"+nodeStr, nil)
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

// node触活
func (self *Client) TouchNode(serviceID, ip string, port uint16) (int, error) {
	var (
		err     error
		req     *http.Request
		resp    *http.Response
		hc      http.Client
		status  int
		nodeStr string
	)
	if err = filter.Batch(
		filter.String(serviceID, "服务ID").Require().Error(),
		filter.String(ip, "节点IP").Require().IsIP().Error(),
		filter.String(strconv.FormatUint(uint64(port), 10), "节点端口").Require().IsDigit().MinInteger(1).MaxInteger(math.MaxUint16).Error(),
	); err != nil {
		return 400, err
	}
	serviceID = base64.RawURLEncoding.EncodeToString([]byte(serviceID))
	nodeStr = base64.RawURLEncoding.EncodeToString([]byte(ip + ":" + strconv.FormatUint(uint64(port), 10)))

	ctx, cancel := context.WithTimeout(context.Background(), self.timeout)
	defer cancel()
	req, err = http.NewRequestWithContext(ctx, "POST", self.addr+"/nodes/"+serviceID+"/"+nodeStr, nil)
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

// node自动触活
func (self *Client) AutoTouchNode(serviceID, ip string, port uint16, errorHandler AutoTouchErrorHandler) {
	for {
		time.Sleep(self.touchInterval)
		if status, err := self.TouchNode(serviceID, ip, port); err != nil {
			errorHandler(status, err)
		}
	}
}
