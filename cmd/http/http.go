package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pangduckwai/sea9go/pkg/http/client"
	"github.com/pangduckwai/sea9go/pkg/http/server"
	"github.com/pangduckwai/sea9go/pkg/logger"
)

const FRM_LOGF_SEC = "2006-01-02_15:04:05"
const PORT = 8999
const IDLE = 30

// Usage: ./http [delay] [verbose] [num-of-runs]
func main() {
	lo, le, lf := logger.Init()
	logsvr, errsvr := logger.AddPrefix(lo, "SERVER"), logger.AddPrefix(le, "SERVER")
	logcln, fatal := logger.AddPrefix(lo, "CLIENT"), logger.AddPrefix(lf, "CLIENT")

	var wg sync.WaitGroup
	var err error
	var run, tmp = 14, 0
	vbs := false
	dly := 0
	switch len(os.Args) {
	case 4:
		tmp, err = strconv.Atoi(os.Args[3])
		if err != nil {
			lf("%v", err)
		}
		if tmp < 0 {
			run = -tmp
		} else if tmp > 0 {
			run = tmp
		}
		fallthrough
	case 3:
		vbs, err = strconv.ParseBool(os.Args[2])
		if err != nil {
			lf("%v", err)
		}
		fallthrough
	case 2:
		dly, err = strconv.Atoi(os.Args[1])
		if err != nil {
			lf("%v", err)
		}
		switch {
		case dly <= 5:
			dly = 0
		case dly <= 10:
			dly = 6
		case dly <= 15:
			dly = 11
		case dly <= IDLE:
			dly = 16
		default:
			dly = IDLE + 1
		}
	}

	/////////////////
	// start server
	var hndr = map[string]func(http.ResponseWriter, *http.Request){
		"/test1/{env}/{ts}": func(w http.ResponseWriter, r *http.Request) {
			var err error
			if r.Method != "GET" {
				err = fmt.Errorf("unsupported method '%v'", r.Method)
				http.Error(w, err.Error(), http.StatusMethodNotAllowed)
				errsvr(" %v\n", err)
				return
			}

			auth, okay := r.Header["Authorization"]
			if !okay || !strings.HasPrefix(auth[0], "Bearer ") {
				err = fmt.Errorf("unauthorized access")
				http.Error(w, err.Error(), http.StatusUnauthorized)
				errsvr(" %v\n", err)
				return
			}

			ts := r.PathValue("ts")
			if ts == "" {
				err = fmt.Errorf("missing request path 'ts'")
				http.Error(w, err.Error(), http.StatusBadRequest)
				errsvr(" %v\n", err)
				return
			}

			env := r.PathValue("env")
			if env == "" {
				err = fmt.Errorf("missing request path 'env'")
				http.Error(w, err.Error(), http.StatusBadRequest)
				errsvr(" %v\n", err)
				return
			}

			_, err = fmt.Fprintf(w, "%v GET TEST1: ENV:%v TS:%v", time.Now().Format(FRM_LOGF_SEC), env, ts)
			if err != nil {
				err = fmt.Errorf("[RES] %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				errsvr(" %v\n", err)
				return
			}
		},
		"/test2/{delay}": func(w http.ResponseWriter, r *http.Request) {
			var err error
			if r.Method != "GET" {
				err = fmt.Errorf("unsupported method '%v'", r.Method)
				http.Error(w, err.Error(), http.StatusMethodNotAllowed)
				errsvr(" %v\n", err)
				return
			}

			auth, okay := r.Header["Authorization"]
			if !okay || !strings.HasPrefix(auth[0], "Bearer ") {
				err = fmt.Errorf("unauthorized access")
				http.Error(w, err.Error(), http.StatusUnauthorized)
				errsvr(" %v\n", err)
				return
			}

			delay := r.PathValue("delay")
			if delay == "" {
				err = fmt.Errorf("missing request path 'delay'")
				http.Error(w, err.Error(), http.StatusBadRequest)
				errsvr(" %v\n", err)
				return
			}

			_, err = fmt.Fprintf(w, "%v GET TEST2: DELAY:%v", time.Now().Format(FRM_LOGF_SEC), delay)
			if err != nil {
				err = fmt.Errorf("[RES] %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				errsvr(" %v\n", err)
				return
			}
		},
	}

	if dly > 0 {
		vbs = true
	}

	start, stop, stopped := server.Prepare(
		PORT, IDLE,
		10*time.Second, 30*time.Second,
		hndr,
		logsvr, errsvr,
		false, vbs,
	)
	go func() {
		start(nil)
	}()

	////////////////////
	// prepare clients
	client, err := client.Client(30 * time.Second)
	if err != nil {
		fatal("[CLNT]%v", err)
	}
	if dly > 0 {
		time.Sleep(time.Duration(dly) * time.Second) // delay 6 to 16 seconds
	}

	//////////////
	// thread #1
	wg.Add(1)
	go func(client *http.Client) {
		defer wg.Done()
		for range run {
			req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%v/test1/TEST/%v", PORT, time.Now().Unix()), nil)
			if err != nil {
				fatal("[REQ][1][PREP]%v", err)
			}
			req.Header.Add("Authorization", "Bearer 1234567001")

			res, err := client.Do(req)
			if err != nil {
				fatal("[RES][1][EXEC]%v", err)
			}
			defer res.Body.Close()
			buf, err := io.ReadAll(res.Body)
			if err != nil {
				fatal("[RES][1][READ]%v", err)
			}

			logcln("[RES][1] %v\n", string(buf))
			time.Sleep(350 * time.Millisecond)
		}
	}(client)

	//////////////
	// thread #2
	wg.Add(1)
	go func(client *http.Client) {
		defer wg.Done()
		for i := range run {
			req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%v/test2/%v", PORT, 900+i), nil)
			if err != nil {
				fatal("[REQ][2]%v", err)
			}
			req.Header.Add("Authorization", "Bearer 1234567002")

			res, err := client.Do(req)
			if err != nil {
				fatal("[RES][2][EXEC]%v", err)
			}
			defer res.Body.Close()
			buf, err := io.ReadAll(res.Body)
			if err != nil {
				fatal("[RES][2][READ]%v", err)
			}

			logcln("[RES][2] %v\n", string(buf))
			time.Sleep(350 * time.Millisecond)
		}
	}(client)

	wg.Wait()
	go func() {
		<-stopped
		logsvr(" received `stopped` signal\n")
	}()
	stop()
	time.Sleep(1 * time.Second)
}
