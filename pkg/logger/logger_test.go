package logger

// Usage:
// $ cd .../logger
// $ go test -bench . -benchmem

import (
	"fmt"
	"math"
	"testing"

	"gopkg.in/yaml.v3"
)

var vals = []uint64{345, 90000, 21, 567890, 1, 1232345, math.MaxUint64, 34958769857, 5678}

func TestDigits(t *testing.T) {
	var ctrl, dcnt int
	for _, val := range vals {
		ctrl = int(math.Log10(float64(val))) + 1
		dcnt = DigitCount(val)
		if dcnt != ctrl {
			t.Fatalf("TestDigits() %20v -> %2v != %v\n", val, dcnt, ctrl)
		}
		fmt.Printf("TestDigits() %20v -> %2v == %v\n", val, dcnt, ctrl)
	}
}

func BenchmarkDigits(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		DigitCount(vals[i%9])
		i++
	}
}

func BenchmarkLog10(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		_ = int(math.Log10(float64(vals[i%9]))) + 1
		i++
	}
}

func TestLogger(t *testing.T) {
	lbls := []string{"this", "is", "a", "test!!!"}
	msgs := []string{"How", "are", "you", "today", "?"}
	log, err, _ := Init()

	log1 := AddPrefix(AddPrefix(log, "TestLogger() - "), lbls...)
	for _, m := range msgs {
		log1(" ==> %v\n", m)
	}

	loge, _ := AddLabels(AddPrefix(err, "TestLogger() - "), 1, 7, 0, 5, 2, 2)
	for i, l := range lbls {
		for j, m := range msgs {
			loge(" ==> IDX:%v\n", l, m, i, j)
		}
	}
}

func getInput() []byte {
	return []byte(
		"header:\n" +
			"  version: 1.1.3\n" +
			"  type: runner\n" +
			"content:\n" +
			"  logLevel: info\n" +
			"  logFile: default\n" +
			"  concurrent: 10\n" +
			"  relogin: true\n" +
			"  retry: 8\n" +
			"  backoff: 2\n" +
			"  wait: 5\n" +
			"  list:\n" +
			"  - id: 1\n" +
			"    name: Search\n" +
			"    path: dosrch.yaml\n" +
			"    body:\n" +
			"      query:\n" +
			"        status:\n" +
			"          $eq: Valid\n" +
			"  - id: 2\n" +
			"    name: Another test\n" +
			"    path: here.json\n" +
			"  - id: 3\n" +
			"    name: Last test\n" +
			"    path: there.json\n" +
			"    body:\n" +
			"      query:\n" +
			"        value:\n" +
			"        - How are you today\n" +
			"        - I'm fine thank you\n" +
			"        - Good to know bye\n" +
			"  envVar:\n" +
			"    LOG_LEVEL: logLevel\n" +
			"    THREAD_NUM: concurrent\n" +
			"    TLS_SERVER: serverCert\n" +
			"  envList:\n" +
			"  - retries:\n" +
			"    - retry1\n" +
			"    - retry2\n" +
			"    - retry3\n" +
			"  - relogin\n" +
			"  - - tester1\n" +
			"    - tester2\n" +
			"    - tester3\n" +
			"  - verifier",
	)
}

// func TestShow(t *testing.T) {
// 	buf := getInput()
// 	fmt.Printf("%s\n", buf)
// }

func TestLogging(t *testing.T) {
	inp := make(map[string]interface{})
	err := yaml.Unmarshal(getInput(), &inp)
	if err != nil {
		t.Fatal(err)
	}
	str, err := Format("", inp)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(str)
}

func TestLogging2(t *testing.T) {
	inp := make(map[string]interface{})
	inp["this"] = []string{"THIS", "IS", "REALLY"}
	inp["is"] = []string{"IS", "REALLY", "A"}
	inp["a"] = []string{"REALLY", "A", "TEST"}
	inp["test"] = []string{"A", "TEST", "!!!"}
	str, err := Format("", inp)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(str)
}

func TestLogging3(t *testing.T) {
	inp := make(map[string]interface{})
	inp["this"] = []interface{}{"THIS", "IS", "REALLY"}
	inp["is"] = []interface{}{"IS", "REALLY", "A"}
	inp["a"] = []interface{}{"REALLY", "A", "TEST"}
	inp["test"] = []interface{}{"A", "TEST", "!!!"}
	str, err := Format("", inp)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(str)
}

func TestMoss(t *testing.T) {
	inp := make(map[string][]string)
	inp["this"] = []string{"THIS", "IS", "REALLY"}
	inp["is"] = []string{"IS", "REALLY", "A"}
	inp["a"] = []string{"REALLY", "A", "TEST"}
	inp["test"] = []string{"A", "TEST", "!!!"}
	fmt.Println(FormatMoss("\t", inp))
}
