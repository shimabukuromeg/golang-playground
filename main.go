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

	// contextの初期化
	// context.Background()関数の返り値からは、「キャンセルされない」「deadlineも持たない」「共有する値も何も持たない」状態のcontextが得られます。
	// いわば「context初期化のための関数」です。
	// https://zenn.dev/hsaki/books/golang-context/viewer/done#context%E3%81%AE%E5%88%9D%E6%9C%9F%E5%8C%96
	initialCtx := context.Background()

	// contextにキャンセル機能を追加
	// context.Background()から得たまっさらなcontextをcontext.WithCancel()関数に渡すことで、
	// 「Doneメソッドからキャンセル有無が判断できるcontext」と「第一返り値のコンテキストをキャンセルするための関数」を得ることができます。
	// https://zenn.dev/hsaki/books/golang-context/viewer/done#context%E3%81%AB%E3%82%AD%E3%83%A3%E3%83%B3%E3%82%BB%E3%83%AB%E6%A9%9F%E8%83%BD%E3%82%92%E8%BF%BD%E5%8A%A0
	ctx, cancel := context.WithCancel(initialCtx)

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
