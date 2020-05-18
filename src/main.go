package main

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/igm/sockjs-go/v3/sockjs"
	"log"
	"net/http"
	"strings"
)

func main() {
	option := sockjs.DefaultOptions
	option.Origin = "*"
	option.CheckOrigin = func(req *http.Request) bool {
		return true
	}

	http.Handle("/ws/", sockjs.NewHandler("/ws", option, wsHandler))
	log.Println("sockjs server initiated")
	log.Fatal(http.ListenAndServe(":8090", nil))
}

func convStr2Interface(c []string) []interface{} {
	d := make([]interface{}, len(c))
	for i, v := range c {
		d[i] = v
	}
	return d
}

func wsHandler(session sockjs.Session) {
	log.Println("session connected", session.ID())
	var client *redis.Client
	isRedisConnected := false
	for {
		if msg, err := session.Recv(); err == nil {
			cmds := strings.Split(msg, " ")
			if len(cmds) < 1 {
				continue
			}
			if isRedisConnected {
				res, err := client.Do(convStr2Interface(cmds)...).Result()
				session.Send(fmt.Sprintf("%v", res))
				if err != nil && res != nil {
					session.Send(res.(string))
				}
			}
			if cmds[0] == "REQCON" {
				client = redis.NewClient(&redis.Options{
					Addr:     cmds[1],
					Password: cmds[2],
					DB:       0,
				})
				isRedisConnected = true
				log.Println("redis connected")
				session.Send("redis connected")
			}
			continue
		}
		break
	}
}
