package proxy

import (
	"fmt"
	"log"
	"mudbot/botutil"
	"mudbot/parser"
	"net"
	"sync"
	"time"
)

type Worker struct {
	name string

	clientConn net.Conn
	mudConn    net.Conn

	clientToMudCopier *Copier
	mudToClientCopier *Copier

	stopping bool

	done chan struct{}

	mu sync.Mutex
	wg sync.WaitGroup
}

func (worker *Worker) Run(name string, clientConn net.Conn, mudConn net.Conn) {
	worker.wg.Add(1)
	worker.done = make(chan struct{})
	worker.name = name

	worker.clientToMudCopier = NewCopier(ACCUMULATION_POLICY_DONT, nil, botutil.NewLogger("cp_client"))
	worker.mudToClientCopier = NewCopier(ACCUMULATION_POLICY_DO, parser.NewParser(), botutil.NewLogger("cp_mud"))

	worker.clientConn = clientConn
	worker.mudConn = mudConn

	log.Printf("== Starting worker %s\n", worker.name)

	go worker.clientToMudCopier.Run(mudConn, clientConn, worker.done)
	go worker.mudToClientCopier.Run(clientConn, mudConn, worker.done)

	worker.waitForCopierClose()
	log.Printf("One of the copiers stopped, stopping worker\n")
	worker.stop()
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
func (worker *Worker) stop() {
	{
		worker.mu.Lock()
		if worker.stopping {
			return
		}

		worker.stopping = true
		fmt.Printf("%s worker stopping\n", time.Now())
		worker.mu.Unlock()
	}

	fmt.Printf("%s closing done channel\n", time.Now())
	close(worker.done)

	fmt.Printf("%s waiting for clienToMudCopier done\n", time.Now())
	<-worker.clientToMudCopier.done
	fmt.Printf("%s waiting for mudToClientCopier done\n", time.Now())
	<-worker.mudToClientCopier.done

	fmt.Printf("%s waiting on wait groups\n", time.Now())
	worker.clientToMudCopier.wg.Wait()
	worker.mudToClientCopier.wg.Wait()

	worker.clientConn.Close()
	worker.mudConn.Close()

	worker.wg.Done()
}

// func (worker *Worker) sendToClient(buf []byte, logger *log.Logger) {
// 	worker.sendToConn(buf, worker.clientConn, logger)
// }

// func (worker *Worker) sendToMud(buf []byte, logger *log.Logger) {
// 	worker.sendToConn(buf, worker.mudConn, logger)
// }
