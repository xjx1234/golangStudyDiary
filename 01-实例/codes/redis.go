/**
* @Author: XJX
* @Description: Redis操作示例
* @File: redis.go
* @Date: 2020/7/3 14:43
 */

package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func StringOp(client *redis.Client) {
	err := client.Set(ctx, "testKey", "xjx", 0).Err()
	fmt.Println(err)
	data, error := client.Get(ctx, "testKey").Result()
	if error == redis.Nil {
		fmt.Println("testKey does not exist")
	} else if error != nil {
		panic(error)
	} else {
		fmt.Println(data)
	}

	err2 := client.SetNX(ctx, "xjx", "zzzzz", 0).Err()
	fmt.Println(err2)
	data1, error1 := client.Get(ctx, "testKey").Result()
	if error1 == redis.Nil {
		fmt.Println("testKey does not exist")
	} else if error1 != nil {
		panic(error1)
	} else {
		fmt.Println(data1)
	}
}

func main() {

	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       1,
	})
	pong, err := client.Ping(ctx).Result()
	fmt.Println(pong, err)

	StringOp(client)

}
