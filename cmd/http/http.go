package main

import (
	"crypto/sha256"
	"encoding/hex"
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

func hasher() (
	hash func([]byte) (string, error),
	hashString func(string) (string, error),
	hashInt64 func(int64) (string, error),
) {
	hsh := sha256.New()

	hash = func(v []byte) (h string, err error) {
		hsh.Reset()
		_, err = hsh.Write(v)
		if err != nil {
			return
		}
		h = hex.EncodeToString(hsh.Sum(nil))
		return
	}
	hashString = func(v string) (string, error) {
		return hash([]byte(v))
	}
	hashInt64 = func(v int64) (string, error) {
		return hashString(strconv.FormatInt(v, 10))
	}

	return
}

func getHndlr(hsh func(string) (string, error), log func(string, ...any)) map[string]func(http.ResponseWriter, *http.Request) {
	return map[string]func(http.ResponseWriter, *http.Request){
		"/test1/{env}/{ts}": func(w http.ResponseWriter, r *http.Request) {
			var err error
			if r.Method != "GET" {
				err = fmt.Errorf("unsupported method '%v'", r.Method)
				http.Error(w, err.Error(), http.StatusMethodNotAllowed)
				log(" %v\n", err)
				return
			}

			auth, okay := r.Header["Authorization"]
			if !okay || !strings.HasPrefix(auth[0], "Bearer ") {
				err = fmt.Errorf("unauthorized access")
				http.Error(w, err.Error(), http.StatusUnauthorized)
				log(" %v\n", err)
				return
			}

			ts := r.PathValue("ts")
			if ts == "" {
				err = fmt.Errorf("missing request path 'ts'")
				http.Error(w, err.Error(), http.StatusBadRequest)
				log(" %v\n", err)
				return
			}
			tk, err := hsh(ts)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				log(" %v\n", err)
				return
			}
			if tk != auth[0][7:] {
				err = fmt.Errorf("authorization failed")
				http.Error(w, err.Error(), http.StatusUnauthorized)
				log(" %v\n", err)
				return
			}

			env := r.PathValue("env")
			if env == "" {
				err = fmt.Errorf("missing request path 'env'")
				http.Error(w, err.Error(), http.StatusBadRequest)
				log(" %v\n", err)
				return
			}

			time.Sleep(100 * time.Millisecond)
			_, err = fmt.Fprintf(w, "%v GET TEST1: TSTMP:%v TK:%v", time.Now().Format(FRM_LOGF_SEC), ts, tk)
			if err != nil {
				err = fmt.Errorf("[RES] %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log(" %v\n", err)
				return
			}
		},
		"/test2/{delay}": func(w http.ResponseWriter, r *http.Request) {
			var err error
			if r.Method != "GET" {
				err = fmt.Errorf("unsupported method '%v'", r.Method)
				http.Error(w, err.Error(), http.StatusMethodNotAllowed)
				log(" %v\n", err)
				return
			}

			auth, okay := r.Header["Authorization"]
			if !okay || !strings.HasPrefix(auth[0], "Bearer ") {
				err = fmt.Errorf("unauthorized access")
				http.Error(w, err.Error(), http.StatusUnauthorized)
				log(" %v\n", err)
				return
			}

			delay := r.PathValue("delay")
			if delay == "" {
				err = fmt.Errorf("missing request path 'delay'")
				http.Error(w, err.Error(), http.StatusBadRequest)
				log(" %v\n", err)
				return
			}
			tk, err := hsh(delay)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				log(" %v\n", err)
				return
			}
			if tk != auth[0][7:] {
				err = fmt.Errorf("authorization failed")
				http.Error(w, err.Error(), http.StatusUnauthorized)
				log(" %v\n", err)
				return
			}

			time.Sleep(100 * time.Millisecond)
			_, err = fmt.Fprintf(w, "%v GET TEST2: DELAY:%013v TK:%v", time.Now().Format(FRM_LOGF_SEC), delay, tk)
			if err != nil {
				err = fmt.Errorf("[RES] %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log(" %v\n", err)
				return
			}
		},
	}
}

func args() (dly, run, idle, dgt int, wait bool, err error) {
	var tmp int
	dly, run, idle, dgt = 0, 7, -1, 2
	switch len(os.Args) {
	case 4:
		tmp, err = strconv.Atoi(os.Args[3])
		if err != nil {
			return
		}
		switch {
		case tmp < 0:
			idle = -1
		case tmp < 15:
			idle = 15
		case tmp > 30:
			idle = 30
		default:
			idle = tmp
		}
		fallthrough
	case 3:
		tmp, err = strconv.Atoi(os.Args[2])
		if err != nil {
			return
		}
		run = tmp
		dgt = logger.DigitCount(uint64(run))
		fallthrough
	case 2:
		tmp, err = strconv.Atoi(os.Args[1])
		if err != nil {
			return
		}
		switch {
		case tmp <= 5:
			dly = 0
		case tmp <= 10:
			dly = 6
		case tmp <= 15:
			dly = 11
		case tmp <= 30:
			dly = 16
		default:
			dly = 31
		}
	}
	if dly > 0 {
		idle = 15
	}
	if idle > 0 {
		wait = true
	}
	return
}

// Usage: ./http [start-delay] [num-of-runs] [server-idle]
func main() {
	var wg sync.WaitGroup
	var wait bool
	lo, le, lf := logger.Init()

	dly, run, idle, dgt, wait, err := args()
	if err != nil {
		lf("error parsing arguments: %v", err)
	}

	logtst, logsvr, errsvr := logger.AddPrefix(lo, "TESTER"), logger.AddPrefix(lo, "SERVER"), logger.AddPrefix(le, "SERVER")
	fatal := logger.AddPrefix(lf, "CLIENT")
	logcln, _ := logger.AddLabels(logger.AddPrefix(lo, "CLIENT"), 2, 1, 2, dgt)
	errcln, _ := logger.AddLabels(logger.AddPrefix(le, "CLIENT"), 2, 1, 2, dgt)
	fatalc, _ := logger.AddLabels(fatal, 2, 1, 2, dgt)

	_, hashString, hashInt64 := hasher()

	/////////////////
	// start server
	logtst(" Starting http package test: delay: %v; run: %v; idle: %v\n", dly, run, idle)
	start, stop, stopped := server.Prepare(
		"Test server",
		PORT, idle,
		10*time.Second, 30*time.Second,
		getHndlr(hashString, errsvr),
		logsvr, errsvr,
		true,
	)
	go func() {
		start(nil)
	}()

	if !wait {
		go func() {
			<-stopped
			logtst(" received `stopped` signal\n")
		}()
	}

	////////////////////
	// prepare clients
	client, err := client.Client(30 * time.Second)
	if err != nil {
		fatal("[NEW] %v", err)
	}
	if dly > 0 {
		time.Sleep(time.Duration(dly) * time.Second) // delay 6 to 16 seconds
	}

	////////////////////
	// client thread #1
	wg.Add(1)
	go func(client *http.Client) {
		var err error
		var ts int64
		var tk string
		var req *http.Request
		var res *http.Response
		var buf []byte
		defer wg.Done()
		for i := range run {
			ts = time.Now().UnixMilli()
			tk, err = hashInt64(ts)

			req, err = http.NewRequest("GET", fmt.Sprintf("http://localhost:%v/test1/SHA2/%v", PORT, ts), nil)
			if err != nil {
				fatalc("[REQ][PREP]%v", 1, i, err)
			}
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", tk))

			res, err = client.Do(req)
			if err != nil {
				fatalc("[RES][EXEC]%v", 1, i, err)
			}
			defer res.Body.Close()
			buf, err = io.ReadAll(res.Body)
			if err != nil {
				fatalc("[RES][READ]%v", 1, i, err)
			}

			if len(buf) > 54 {
				if string(buf[54:]) == tk {
					logcln("[RES] OKAY \"%v\"\n", 1, i, string(buf[:50]))
				} else {
					errcln("[RES] FAILED %v\n", 1, i, string(buf))
				}
			} else {
				errcln("[RES] FAILED: %v\n", 1, i, string(buf))
			}
		}
	}(client)

	////////////////////
	// client thread #2
	wg.Add(1)
	go func(client *http.Client) {
		var err error
		var cnt int64
		var tk string
		var req *http.Request
		var res *http.Response
		var buf []byte
		defer wg.Done()
		for i := range run {
			cnt = int64(i) + 900
			tk, err = hashInt64(cnt)

			req, err = http.NewRequest("GET", fmt.Sprintf("http://localhost:%v/test2/%v", PORT, cnt), nil)
			if err != nil {
				fatalc("[REQ]%v", 2, i, err)
			}
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", tk))

			res, err = client.Do(req)
			if err != nil {
				fatalc("[RES][EXEC]%v", 2, i, err)
			}
			defer res.Body.Close()
			buf, err = io.ReadAll(res.Body)
			if err != nil {
				fatalc("[RES][READ]%v", 2, i, err)
			}

			if len(buf) > 54 {
				if string(buf[54:]) == tk {
					logcln("[RES] OKAY \"%v\"\n", 2, i, string(buf[:50]))
				} else {
					errcln("[RES] FAILED %v\n", 2, i, string(buf))
				}
			} else {
				errcln("[RES] FAILED: %v\n", 2, i, string(buf))
			}
		}
	}(client)

	wg.Wait()

	if !wait {
		logtst(" stopping the server in 3 seconds...\n")
		time.Sleep(3 * time.Second)
		stop()
	} else {
		<-stopped
		logtst(" server stopped\n")
	}
	time.Sleep(1 * time.Second)
}
