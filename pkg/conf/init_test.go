package conf

import (
	"context"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
	"time"
)

var remoteCfg = `
app:
  name: cxx
  mode: dev

server:
  port: 8888

database:
  type: postgresql
  host: localhost
  port: 5432
`

var remoteDevCfg = `
database:
  host: postgres
  port: 5432
`

func initEtcd() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 1 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	_, err = cli.Put(context.Background(), "/config/app", remoteCfg)
	if err != nil {
		panic(err)
	}
	_, err = cli.Put(context.Background(), "/config/app-dev", remoteDevCfg)
	if err != nil {
		panic(err)
	}
}

func TestInitRemote(t *testing.T) {
	asserts := assert.New(t)

	asserts.NotPanics(func() {
		initEtcd()
		Init()
	})

	asserts.Equal("cxx", AppConf.Name)
	asserts.Equal("dev", AppConf.Mode)
	asserts.Equal("8888", ServerConf.Port)
	asserts.Equal("postgresql", DatabaseConf.Type)
	// dev配置覆盖了remote配置
	asserts.Equal("postgres", DatabaseConf.Host)
}

func TestInitLocal(t *testing.T) {

	asserts := assert.New(t)
	asserts.NotPanics(func() {
		Init("local.yaml")
	})
	// local-dev配置覆盖了local配置
	asserts.Equal("local-dev", AppConf.Name)
	asserts.Equal("dev", AppConf.Mode)
	asserts.Equal("error", LogConf.Level)
}
