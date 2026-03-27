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
// - map keys of the argument 'handler' are the route of each http handler
func Prepare(
	port, idle int,
	readTimeout, writeTimeout time.Duration,
	handler map[string]func(http.ResponseWriter, *http.Request),
	log, logerr func(string, ...any),
	debug, verbose bool,
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
	idlecnt := 0

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
		log(" HTTP server (r:%v|w:%v) listening on port %d...\n", readTimeout, writeTimeout, port)

		wg.Add(1)
		go func() {
			for flag := range cctrl {
				if flag {
					idlecnt = 0 // timeout counter resetted
					if verbose {
						log(" server timeout reset to %3vs\n", idle)
					}
				} else {
					idlecnt = idle
					if verbose {
						log(" server stopping...\n")
					}
				}
			}
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			for ; idlecnt < idle; idlecnt++ {
				if (idlecnt%5 == 0) && (debug || verbose) {
					log(" server timeout in %3vs\n", idle-idlecnt)
				}
				time.Sleep(SECOND_1)
			}
			if verbose {
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
