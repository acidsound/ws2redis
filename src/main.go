package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v7"
	"github.com/igm/sockjs-go/v3/sockjs"
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
	if msg, err := session.Recv(); err == nil {
		cmds := strings.Split(msg, " ")
		if len(cmds) != 2 {
			log.Println("wrong params")
			return
		}
		client = redis.NewClient(&redis.Options{
			Addr:     cmds[0],
			Password: cmds[1],
			DB:       0,
		})
		log.Println("redis connected")
		session.Send("$connected")
		go ws2redis(client, session)
	}
}

func ws2redis(client *redis.Client, session sockjs.Session) {
	for {
		if msg, err := session.Recv(); err == nil {
			cmds := strings.Split(msg, " ")
			if len(cmds) < 1 {
				continue
			}
			var newCmds []string
			switch cmds[1] {
			case "SET", "PUBLISH":
				newCmds = append(newCmds, cmds[1])
				newCmds = append(newCmds, cmds[2])
				newCmds = append(newCmds, strings.Join(cmds[3:], " "))
			default:
				newCmds = cmds[1:]
			}

			res, err := client.Do(convStr2Interface(newCmds)...).Result()
			session.Send(fmt.Sprintf("%v %v", cmds[0], res))
			if err != nil && res != nil {
				session.Send(res.(string))
			}

		}
	}
}
