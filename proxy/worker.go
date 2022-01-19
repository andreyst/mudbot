package proxy

import (
	"fmt"
	"mudbot/bot"
	"mudbot/botutil"
	"net"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Worker struct {
	bot *bot.Bot

	clientConn net.Conn
	mudConn    net.Conn

	clientToMudCopier *Copier
	mudToClientCopier *Copier

	stopping bool

	done chan struct{}

	mu sync.Mutex
	wg sync.WaitGroup

	logger *zap.SugaredLogger
}

func NewWorker(clientConn net.Conn, mudConn net.Conn) *Worker {
	bot := bot.NewBot()

	w := Worker{
		bot: bot,

		logger: botutil.NewLogger("worker"),

		done: make(chan struct{}),

		clientToMudCopier: NewCopier(ACCUMULATION_POLICY_DONT, nil, botutil.NewLogger("cp_client")),
		mudToClientCopier: NewCopier(ACCUMULATION_POLICY_DO, bot, botutil.NewLogger("cp_mud")),

		clientConn: clientConn,
		mudConn:    mudConn,
	}

	return &w
}

func (w *Worker) Run() {
	w.logger.Info("Starting worker")

	w.wg.Add(1)

	go func() {
		for {
			select {
			case <-w.done:
				return
			default:
			}
			fmt.Printf("Bot: %+v\n", w.bot)
			time.Sleep(time.Duration(2) * time.Second)
		}
	}()

	go w.clientToMudCopier.Run(w.mudConn, w.clientConn, w.done)
	go w.mudToClientCopier.Run(w.clientConn, w.mudConn, w.done)

	w.waitForCopierClose()
	w.logger.Info("One of the copiers stopped, stopping worker")
	w.stop()
}

func (worker *Worker) waitForCopierClose() {
	for {
		select {
		case <-worker.clientToMudCopier.done:
		case <-worker.mudToClientCopier.done:
			return
		}
	}
}
func (w *Worker) stop() {
	{
		w.mu.Lock()
		if w.stopping {
			return
		}

		w.stopping = true
		w.mu.Unlock()
	}

	close(w.done)

	<-w.clientToMudCopier.done
	<-w.mudToClientCopier.done

	w.clientToMudCopier.wg.Wait()
	w.mudToClientCopier.wg.Wait()

	w.clientConn.Close()
	w.mudConn.Close()

	w.wg.Done()
}

// func (worker *Worker) sendToClient(buf []byte, logger *log.Logger) {
// 	worker.sendToConn(buf, worker.clientConn, logger)
// }

// func (worker *Worker) sendToMud(buf []byte, logger *log.Logger) {
// 	worker.sendToConn(buf, worker.mudConn, logger)
// }
