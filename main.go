package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

// キャンセルされるまでnumをひたすら送信し続けるチャネルを生成
func generator(doneChannel chan string, num int) <-chan int {
	// チャネル作成
	outChannel := make(chan int)

	// goroutine
	go func() {
		defer wg.Done()

	LOOP:
		for {
			select {
			case <-doneChannel: // doneチャネルがcloseされたらbreakが実行される
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
	doneChannel := make(chan string)

	gen := generator(doneChannel, 1)

	wg.Add(1)

	for i := 0; i < 5; i++ {
		a := <-gen
		fmt.Println(a)
	}

	close(doneChannel) // 5回genを使ったら、doneチャネルをcloseしてキャンセルを実行

	wg.Wait()
}
