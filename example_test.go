package tsing_center_go

import (
	"testing"
)

func TestAdd(t *testing.T) {
	cli, err := New(Config{
		Addr:          "http://127.0.0.1:20080",
		Secret:        "123456",
		TTL:           10,
		TouchInterval: 8,
		Timeout:       5,
	})
	if err != nil {
		t.Error(err)
		return
	}
	ip, status, e := cli.GetIP()
	if e != nil {
		t.Error(e, status)
		return
	}
	t.Log(ip)

	// 添加服务
	if status, err = cli.AddService(Service{
		ID:          "demo",
		LoadBalance: "SWRR",
	}); err != nil {
		t.Error(err, status)
		return
	}

	// 添加节点
	if status, err = cli.AddNode("demo", Node{
		IP:     ip,
		Port:   80,
		Weight: 1,
		TTL:    0,
	}); err != nil {
		t.Error(err, status)
		return
	}

}

func TestSet(t *testing.T) {
	cli, err := New(Config{
		Addr:          "http://127.0.0.1:20080",
		Secret:        "123456",
		TTL:           10,
		TouchInterval: 8,
		Timeout:       5,
	})
	if err != nil {
		t.Error(err)
		return
	}
	ip, status, e := cli.GetIP()
	if e != nil {
		t.Error(e, status)
		return
	}
	t.Log(ip)

	// 重写或添加服务
	if status, err = cli.SetService(Service{
		ID:          "demo",
		LoadBalance: "SWRR",
		Meta:        `{"SECRET":"123456"}`,
	}); err != nil {
		t.Error(err, status)
		return
	}

	// 重写或添加节点
	if status, err = cli.SetNode("demo", Node{
		IP:     ip,
		Port:   80,
		Weight: 1,
	}); err != nil {
		t.Error(err, status)
		return
	}
	// 重写或添加节点
	if status, err = cli.SetNode("demo", Node{
		IP:     ip,
		Port:   90,
		Weight: 1,
		TTL:    10,
		Meta:   `{"OS":"Mac OS"}`,
	}); err != nil {
		t.Error(err, status)
		return
	}

}

func TestDiscovery(t *testing.T) {
	cli, err := New(Config{
		Addr:          "http://127.0.0.1:20080",
		Secret:        "123456",
		TTL:           10,
		TouchInterval: 8,
		Timeout:       5,
	})
	if err != nil {
		t.Error(err)
		return
	}
	ip, status, e := cli.GetIP()
	if e != nil {
		t.Error(e, status)
		return
	}
	t.Log(ip)

	// 获取服务的节点
	var node Node
	if node, status, err = cli.DiscoverService("demo"); err != nil {
		t.Error(err, status)
		return
	}
	t.Log(node)
}
