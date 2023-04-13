package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/router"
)

func main() {
	err := realMain()
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain() error {
	// config values
	const (
		defaultPort   = ":8080"
		defaultDBPath = ".sqlite3/todo.db"
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// set time zone
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	// set up sqlite3
	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		return err
	}
	defer todoDB.Close()

	// NOTE: 新しいエンドポイントの登録はrouter.NewRouterの内部で行うようにする
	mux := router.NewRouter(todoDB)

	// * HTTPサーバーを立ち上げる
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	// シグナルを受け取るチャネルを用意する
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)

	// ウェイトグループを用意する
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// Graceful Shutdown を行うためのゴルーチン
	go func() {
		defer wg.Done()

		// シグナルを受け取るまで待つ
		<-sigCh
		log.Println("Shutting down...")

		// コンテキストを作成する
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Graceful Shutdown を実行する
		if err := server.Shutdown(ctx); err != nil {
			log.Println("Graceful Shutdown:", err)
		}
	}()

	log.Println("Starting server...")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	// ウェイトグループを待つ
	wg.Wait()

	return nil
}
