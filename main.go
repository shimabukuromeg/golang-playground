package main

import (
	"context"
	"fmt"
	"sync"
)

var wg sync.WaitGroup

// キャンセルされるまでnumをひたすら送信し続けるチャネルを生成
func generator(ctx context.Context, num int) <-chan int {
	// チャネル作成
	outChannel := make(chan int)

	// goroutine
	go func() {
		defer wg.Done()

	LOOP:
		for {
			select {
			case <-ctx.Done():
				break LOOP
			case outChannel <- num: // キャンセルされてなければnumを送信
			}
		}
		close(outChannel)
		fmt.Println("generator closed")
	}()

	return outChannel
}

func main() {

	// チェネルを作成する
	// doneChannel := make(chan string)

	ctx, cancel := context.WithCancel(context.Background())

	gen := generator(ctx, 1)

	wg.Add(1)

	for i := 0; i < 5; i++ {
		a := <-gen
		fmt.Println(a)
	}

	// close(doneChannel) // 5回genを使ったら、doneチャネルをcloseしてキャンセルを実行
	cancel()

	wg.Wait()
}
