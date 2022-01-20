package proxy

import (
	"mudbot/botutil"
	"net"
	"sync"

	"go.uber.org/zap"
)

type Worker struct {
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

func NewWorker(clientConn net.Conn, mudConn net.Conn, botParseCallback AccumulatorCallback) *Worker {
	w := Worker{
		logger: botutil.NewLogger("worker"),

		done: make(chan struct{}),

		clientToMudCopier: NewCopier(ACCUMULATION_POLICY_DONT, nil, botutil.NewLogger("cp_client")),
		mudToClientCopier: NewCopier(ACCUMULATION_POLICY_DO, botParseCallback, botutil.NewLogger("cp_mud")),

		clientConn: clientConn,
		mudConn:    mudConn,
	}

	return &w
}

func (w *Worker) Run() {
	w.logger.Info("Starting worker")

	w.wg.Add(1)

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

func (w *Worker) sendToClient(buf []byte) {
	w.writeToConn(buf, w.clientConn)
}

func (w *Worker) sendToMud(buf []byte) {
	w.writeToConn(buf, w.mudConn)
}

func (w *Worker) writeToConn(buf []byte, conn net.Conn) {
	n, writeErr := conn.Write(buf)
	if writeErr != nil {
		w.logger.Fatalf("Error while writing to conn: %v", writeErr)
	}
	w.logger.Debugf("Wrote string of len %v to conn", n)
}
