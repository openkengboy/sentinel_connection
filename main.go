package main

import (
	"fmt"
	"strings"
	"time"

	radix "github.com/mediocregopher/radix/v3"
)

func main() {
	fmt.Println("Hello, 世界")
	NewRedisSentinelConn()
	//NewRedisConn()
}

var snil *radix.Sentinel
var pool *radix.Pool

func NewRedisSentinelConn() error {
	redisMaster := "mymaster"
	redisSentinelList := "redis://:@ip:port"
	addr := strings.Split(redisSentinelList, ",")
	//snilClient, err := radix.NewSentinel(redisMaster, address)
	snilClient, err := radix.NewSentinel(redisMaster, addr, radix.SentinelPoolFunc(poolSentinelClientFunc))
	if err != nil {
		fmt.Printf("Cocnnect sentinel redis err: %v \n", err)
		return err
	}
	snil = snilClient
	fmt.Printf("Cocnnect sentinel redis success")
	return nil
}

func sentinelConnFunc(net, addr string) (radix.Conn, error) {
	redisIndex := 15

	pwd := "@2018"
	return radix.Dial(net, addr, radix.DialTimeout(5*time.Second), radix.DialSelectDB(redisIndex), radix.DialAuthPass(pwd))
}

func poolSentinelClientFunc(network, addr string) (radix.Client, error) {
	return radix.NewPool(network, addr, 10, radix.PoolConnFunc(sentinelConnFunc))

}

func NewRedisConn() error {
	redisHost := "redis://:password@ip:port"
	p, err := radix.NewPool("tcp", redisHost, 10, radix.PoolConnFunc(ConnFunc))
	if err != nil {
		fmt.Printf("Cocnnect redis err: %v \n", err)
		return err
	}
	pool = p
	fmt.Printf("Cocnnect redis success")
	return nil
}

func ConnFunc(network, addr string) (radix.Conn, error) {
	redisIndex := 15
	conn, err := radix.Dial(network, addr, radix.DialSelectDB(redisIndex))
	if err != nil {
		fmt.Printf("redis err: %#v", err)
	}
	return conn, err
}
