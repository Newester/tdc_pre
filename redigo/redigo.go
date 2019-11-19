package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/gomodule/redigo/redis"
)

func main() {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("connect redis error :", err)
		return
	}
	fmt.Println("connect success ...")
	defer conn.Close()
	_, err = conn.Do("SET", "name", "sepro")
	if err != nil {
		fmt.Println("redis set error:", err)
	}
	name, err := redis.String(conn.Do("GET", "name"))
	if err != nil {
		fmt.Println("redis get error:", err)
	} else {
		fmt.Printf("Got name: %s \n", name)
	}
	_, err = conn.Do("MSET", "name", "sepro", "age", 24)
	if err != nil {
		fmt.Println("redis mset error:", err)
	}
	res, err := redis.Strings(conn.Do("MGET", "name", "age"))
	if err != nil {
		fmt.Println("redis get error:", err)
	} else {
		resType := reflect.TypeOf(res)
		fmt.Printf("res type : %s \n", resType)
		fmt.Printf("MGET name: %s \n", res)
		fmt.Println(len(res))
	}

	_, err = conn.Do("LPUSH", "list1", "ele1", "ele2", "ele3", "ele4")
	if err != nil {
		fmt.Println("redis mset error:", err)
	}
	//res, err := redis.String(conn.Do("LPOP", "list1"))//获取栈顶元素
	//res, err := redis.String(conn.Do("LINDEX", "list1", 3)) //获取指定位置的元素
	res, err = redis.Strings(conn.Do("LRANGE", "list1", 0, 3)) //获取指定下标范围的元素
	if err != nil {
		fmt.Println("redis POP error:", err)
	} else {
		resType := reflect.TypeOf(res)
		fmt.Printf("res type : %s \n", resType)
		fmt.Printf("res  : %s \n", res)
	}

	_, err = conn.Do("HSET", "user", "name", "sepro", "age", 24)
	if err != nil {
		fmt.Println("redis mset error:", err)
	}
	var result, errInt = redis.Int64(conn.Do("HGET", "user", "age"))
	if errInt != nil {
		fmt.Println("redis HGET error:", errInt)
	} else {
		resType := reflect.TypeOf(result)
		fmt.Printf("res type : %s \n", resType)
		fmt.Printf("res  : %d \n", result)
	}

	conn.Send("HSET", "user", "name", "sepro", "age", "24")
	conn.Send("HSET", "user", "sex", "male")
	conn.Send("HGET", "user", "age")
	conn.Flush()

	res1, err := conn.Receive()
	fmt.Printf("Receive res1:%v \n", res1)
	res2, err := conn.Receive()
	fmt.Printf("Receive res2:%v\n", res2)
	res3, err := conn.Receive()
	fmt.Printf("Receive res3:%s\n", res3)

	go subs()
	go push("this is sepro")
	time.Sleep(time.Second * 3)

	conn.Send("MULTI")
	conn.Send("INCR", "foo")
	conn.Send("INCR", "bar")
	r, err := conn.Do("EXEC")
	fmt.Println(r)

	conn2 := Pool.Get()
	res3, err = conn2.Do("HSET", "user", "name", "sepro")
	fmt.Println(res3, err)
	res1, err = redis.String(conn.Do("HGET", "user", "name"))
	fmt.Printf("res:%s,error:%v", res1, err)
}

// Pool ...
var Pool redis.Pool

func init() { //init 用于初始化一些参数，先于main执行
	Pool = redis.Pool{
		MaxIdle:     16,
		MaxActive:   32,
		IdleTimeout: 120,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1:6379")
		},
	}
}

func subs() { //订阅者
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("connect redis error :", err)
		return
	}
	defer conn.Close()
	psc := redis.PubSubConn{Conn: conn}
	psc.Subscribe("channel1") //订阅channel1频道
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
		case redis.Subscription:
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			fmt.Println(v)
			return
		}
	}

}

func push(message string) { //发布者
	conn, _ := redis.Dial("tcp", "127.0.0.1:6379")
	_, err1 := conn.Do("PUBLISH", "channel1", message)
	if err1 != nil {
		fmt.Println("pub err: ", err1)
		return
	}

}
