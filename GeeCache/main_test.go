package geecache

/*
$ curl http://localhost:9999/_geecache/scores/Tom
630

$ curl http://localhost:9999/_geecache/scores/kkk
kkk not exist
*/

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

var db2 = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestName1(t *testing.T) {
	NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) { // 缓存汇总没有数据 获取数据存入缓存的回调
			log.Println("[SlowDB] search key", key)
			if v, ok := db2[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := NewHTTPPool(addr)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
