// Package server is a wrapper of `http.Server` with start/stop/timeout handling.
package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const SECOND_1 = 1 * time.Second

// Prepare prepare a http server and return the functions to start/stop it.
// * Parameters:
//   - 'idle': the idle time (in seconds) before the server is automatically stopped, < 0 means no timeout
//   - 'handler': its map keys are the route of each http handler
//
// * Returns:
//   - 'start': the function to start the server, it accepts a callback function which is called when the server is stopped
//   - 'stop': the function to stop the server
//   - 'stopped': the channel to signal when the server is stopped
func Prepare(
	name string,
	port, idle int,
	readTimeout, writeTimeout time.Duration,
	handler map[string]func(http.ResponseWriter, *http.Request),
	log, logerr func(string, ...any),
	verbose bool,
) (
	start func(func()) bool,
	stop func(),
	stopped chan bool,
) {
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		ReadTimeout:    readTimeout * time.Second,
		WriteTimeout:   writeTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	running := true
	idleCnt := 0
	incr := 1
	if idle < 0 {
		idle = 1
		incr = 0
	}

	stopped = make(chan bool)
	cctrl := make(chan bool)

	for route, fn := range handler {
		if verbose {
			log(" Adding handler \"%v\"\n", route)
		}
		http.HandleFunc(
			route,
			func(w http.ResponseWriter, r *http.Request) {
				cctrl <- true // reset timeout counter
				fn(w, r)
			},
		)
	}

	start = func(stopping func()) bool {
		var wg sync.WaitGroup
		log(" \"%v\" (r:%v|w:%v) listening on port %d...\n", name, readTimeout, writeTimeout, port)

		wg.Add(1)
		go func() {
			for flag := range cctrl {
				if flag {
					idleCnt = 0 // timeout counter resetted
					if verbose && incr > 0 {
						log(" server timeout reset to %3vs\n", idle)
					}
				} else {
					running = false
					if verbose {
						log(" server stopping...\n")
					}
				}
			}
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			for ; idleCnt < idle && running; idleCnt += incr {
				if (idleCnt%5 == 0) && verbose && incr > 0 {
					log(" server timeout in %3vs\n", idle-idleCnt)
				}
				time.Sleep(SECOND_1)
			}
			if verbose && running {
				log(" server shutting down...\n")
			}
			server.Shutdown(context.Background())
			close(cctrl)
			if stopping != nil {
				stopping()
			}
			close(stopped)
			wg.Done()
		}()

		err := server.ListenAndServe()
		wg.Wait()
		if err != http.ErrServerClosed {
			logerr("%v\n", err)
			return false
		} else {
			log(" %v\n", err)
			return true
		}
	}

	stop = func() {
		cctrl <- false
	}

	return
}
