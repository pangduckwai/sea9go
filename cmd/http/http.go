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
const IDLE = 30

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

// Usage: ./http [delay] [verbose] [num-of-runs]
func main() {
	var wg sync.WaitGroup
	var err error
	var run, dgt, tmp = 14, 2, 0
	vbs := false
	dly := 0
	lo, le, lf := logger.Init()

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
		dgt = logger.DigitCount(uint64(run))
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

	logsvr, errsvr := logger.AddPrefix(lo, "SERVER"), logger.AddPrefix(le, "SERVER")
	fatal := logger.AddPrefix(lf, "CLIENT")
	logcln, _ := logger.AddLabels(logger.AddPrefix(lo, "CLIENT"), 2, 1, 2, dgt)
	errcln, _ := logger.AddLabels(logger.AddPrefix(le, "CLIENT"), 2, 1, 2, dgt)
	fatalc, _ := logger.AddLabels(fatal, 2, 1, 2, dgt)

	_, hashString, hashInt64 := hasher()

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
			tk, err := hashString(ts)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				errsvr(" %v\n", err)
				return
			}
			if tk != auth[0][7:] {
				err = fmt.Errorf("authorization failed")
				http.Error(w, err.Error(), http.StatusUnauthorized)
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

			time.Sleep(100 * time.Millisecond)
			_, err = fmt.Fprintf(w, "%v GET TEST1: TSTMP:%v TK:%v", time.Now().Format(FRM_LOGF_SEC), ts, tk)
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
			tk, err := hashString(delay)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				errsvr(" %v\n", err)
				return
			}
			if tk != auth[0][7:] {
				err = fmt.Errorf("authorization failed")
				http.Error(w, err.Error(), http.StatusUnauthorized)
				errsvr(" %v\n", err)
				return
			}

			time.Sleep(100 * time.Millisecond)
			_, err = fmt.Fprintf(w, "%v GET TEST2: DELAY:%013v TK:%v", time.Now().Format(FRM_LOGF_SEC), delay, tk)
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
		"Test server",
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
		fatal("[NEW] %v", err)
	}
	if dly > 0 {
		time.Sleep(time.Duration(dly) * time.Second) // delay 6 to 16 seconds
	}

	//////////////
	// thread #1
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

			if string(buf[54:]) == tk {
				logcln("[RES] OKAY \"%v\"\n", 1, i, string(buf[:50]))
			} else {
				errcln("[RES] FAILED %v\n", 1, i, string(buf))
			}
		}
	}(client)

	//////////////
	// thread #2
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

			if string(buf[54:]) == tk {
				logcln("[RES] OKAY \"%v\"\n", 2, i, string(buf[:50]))
			} else {
				errcln("[RES] FAILED %v\n", 2, i, string(buf))
			}
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
