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
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
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
				var newCmds []string
				if cmds[0] == "SET" {
					newCmds = append(newCmds, cmds[0])
					newCmds = append(newCmds, cmds[1])
					newCmds = append(newCmds, strings.Join(cmds[2:], " "))
				} else {
					newCmds = cmds
				}
				res, err := client.Do(convStr2Interface(newCmds)...).Result()
				session.Send(fmt.Sprintf("%v", res))
				if err != nil && res != nil {
					session.Send(res.(string))
				}
			}
			if cmds[0] == "REQCON" {
				if len(cmds) != 3 {
					log.Println("wrong params")
					break
				}
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
