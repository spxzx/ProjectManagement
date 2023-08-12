package common

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(e *gin.Engine, addr, srvName string, stop func()) {
	srv := &http.Server{
		Addr:    addr,
		Handler: e,
	}

	// 优雅启停
	go func() {
		log.Printf("%s is running on %s \n", srvName, addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}() // 就这样写的话程序会直接停止，所以需要阻塞

	quit := make(chan os.Signal)
	// SIGINT  用户发送INTR（程序终止信号）字符（通常是Ctrl+C）触发 kill -2
	// SIGTERM 程序结束（可被捕获、阻塞或忽略）
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 没有收到信号就阻塞在这里，收到就将其立即读出去进行后面的步骤
	<-quit
	log.Printf("shutting down project %s... \n", srvName)

	/*
		它主要的用处如果用一句话来说，是在于控制goroutine的生命周期。
		当一个计算任务被goroutine承接了之后，由于某种原因（超时，或者强制退出）
		我们希望中止这个goroutine的计算任务，那么就用得到这个Context了。
		1. context包的WithTimeout()函数接受一个 Context 和超时时间作为参数，返回其子Context和取消函数cancel
		2. 新创建协程中传入子Context做参数，且需监控子Context的Done通道，若收到消息，则退出
		3. 需要新协程结束时，在外面调用 cancel 函数，即会往子Context的Done通道发送消息
		4. 若不调用cancel函数，到了原先创建Context时的超时时间，它也会自动调用cancel()函数，即会往子Context的Done通道发送消息
	*/
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if stop != nil {
		stop()
	}
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("%s shutdown, cause by : %s \n", srvName, err)
	}
	select {
	case <-ctx.Done():
		log.Println("timeout shutdown...")
	}
	log.Printf("%s stop success... \n", srvName)
}
