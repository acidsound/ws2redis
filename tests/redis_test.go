package tests

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"strings"
	"testing"
)

func TestConnect(t *testing.T) {
	msg := "REQCON 127.0.0.1:6380 1234qwer"
	t.Log("init")
	cmds := strings.Split(msg, " ")
	client := redis.NewClient(&redis.Options{
		Addr:     cmds[1],
		Password: cmds[2],
		DB:       0,
	})
	t.Log(client)
}

func TestPing(t *testing.T) {
	msg := "REQCON 127.0.0.1:6380 1234qwer"
	t.Log("init")
	cmds := strings.Split(msg, " ")
	client := redis.NewClient(&redis.Options{
		Addr:     cmds[1],
		Password: cmds[2],
		DB:       0,
	})
	t.Log(client)
	res, err := client.Ping().Result()
	t.Log(res, err)
	if res != "PONG" {
		panic(err)
	}
}

func TestCustomCmd(t *testing.T) {
	msg := "REQCON 127.0.0.1:6380 1234qwer"
	t.Log("init")
	cmds := strings.Split(msg, " ")
	client := redis.NewClient(&redis.Options{
		Addr:     cmds[1],
		Password: cmds[2],
		DB:       0,
	})
	t.Log(client)
	res, err := client.Do("PING").Result()
	t.Log(res, err)
	if res != "PONG" {
		panic(err)
	}
}

func TestEval(t *testing.T) {
	msg := "REQCON 127.0.0.1:6380 1234qwer"
	t.Log("init")
	cmds := strings.Split(msg, " ")
	client := redis.NewClient(&redis.Options{
		Addr:     cmds[1],
		Password: cmds[2],
		DB:       0,
	})
	t.Log(client)
	res, err := client.Do("EVAL", "12+34", "0").Result()
	t.Log(res, err)

	if fmt.Sprintf("%v", res) != "46" {
		panic(err)
	}
}

func TestArrayCmd(t *testing.T) {
	msg := "REQCON 127.0.0.1:6380 1234qwer"
	t.Log("init")
	cmds := strings.Split(msg, " ")
	client := redis.NewClient(&redis.Options{
		Addr:     cmds[1],
		Password: cmds[2],
		DB:       0,
	})
	t.Log(client)
	c := strings.Split("EVAL 12+34 0", " ")
	d := make([]interface{}, len(c))
	for i, v := range c {
		d[i] = v
	}
	res, err := client.Do(d...).Result()
	t.Log(res, err)

	if fmt.Sprintf("%v", res) != "46" {
		panic(err)
	}
}
