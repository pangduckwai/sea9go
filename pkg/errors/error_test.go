package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	var err0 error
	var err1 error = errors.New("error 1")
	err := Append(err0, err1)
	if len(err.Errors) != 1 {
		t.Fatalf("TestError() expected 1 error but got %v", len(err.Errors))
	}
	if err.Fatal {
		t.Fatalf("TestError() expected non-fatal error but got fatal")
	}
	fmt.Printf("TestError() \"%v\"\n", err)
}

func TestFatal(t *testing.T) {
	var err0 error
	var err1 error = Fatal("error 1")
	err := Append(err0, err1)
	if len(err.Errors) != 1 {
		t.Fatalf("TestFatal() expected 1 error but got %v", len(err.Errors))
	}
	if !err.Fatal {
		t.Fatalf("TestFatal() expected fatal error but got non-fatal")
	}
	fmt.Printf("TestFatal() \"%v\"\n", err)
}

func TestChain(t *testing.T) {
	var err0 error
	var err1 error = Fatal("[CHAIN] error 1")
	err := Append(err0, err1)
	if len(err.Errors) != 1 {
		t.Fatalf("TestChain() expected 1 error but got %v", len(err.Errors))
	}
	if !err.Fatal {
		t.Fatalf("TestChain() expected fatal error but got non-fatal")
	}
	fmt.Printf("TestChain() \"%v\"\n", err)
}

func TestErrors(t *testing.T) {
	var err0 error
	var err1 error = errors.New("error 1")
	var err2 error = errors.New("error 2")
	err := Append(err0, err1, err2)
	if len(err.Errors) != 2 {
		t.Fatalf("TestErrors() expected 2 errors but got %v", len(err.Errors))
	}
	if err.Fatal {
		t.Fatalf("TestErrors() expected non-fatal error but got fatal")
	}
	fmt.Printf("TestErrors():\n%v\n", err)
}

func TestWithBase(t *testing.T) {
	err0 := Fatal("fatal error 0")
	var err1 error = errors.New("error 1")
	var err2 error = errors.New("error 2")
	err := Append(err0, err1, err2)
	if len(err.Errors) != 3 {
		t.Fatalf("TestWithBase() expected 3 errors but got %v", len(err.Errors))
	}
	if !err.Fatal {
		t.Fatalf("TestWithBase() expected fatal error but got non-fatal")
	}
	fmt.Printf("TestWithBase()\n%v\n", err)
}

func TestErr(t *testing.T) {
	var err0 error
	var err1 error = Fatal("error 1")
	var err2 error = errors.New("error 2")
	err := Append(err0, err1, err2)
	if len(err.Errors) != 2 {
		t.Fatalf("TestErr() expected 2 errors but got %v", len(err.Errors))
	}
	if !err.Fatal {
		t.Fatalf("TestErr() expected fatal error but got non-fatal")
	}
	fmt.Printf("TestErr()\n%v\n", err)
}

func TestNil(t *testing.T) {
	var err0 *Err
	var err1 error = NonFatal("error 1")
	var err2 error = errors.New("error 2")
	err := Append(err0, err1, err2)
	if len(err.Errors) != 2 {
		t.Fatalf("TestNil() expected 2 errors but got %v", len(err.Errors))
	}
	if err.Fatal {
		t.Fatalf("TestNil() expected non-fatal error but got fatal")
	}
	fmt.Printf("TestNil()\n%v\n", err)
}

func TestFatals(t *testing.T) {
	err0 := New(true)
	var err1 error = errors.New("error 1")
	var err2 error = errors.New("error 2")
	err := Append(err0, err1, err2)
	if len(err.Errors) != 2 {
		t.Fatalf("TestFatals() expected 2 errors but got %v", len(err.Errors))
	}
	if !err.Fatal {
		t.Fatalf("TestFatals() expected fatal error but got non-fatal")
	}
	fmt.Printf("TestFatals()\n%v\n", err)
}

func TestFatalss(t *testing.T) {
	err0 := New(true, "error 0", "error 3")
	var err1 error = errors.New("error 1")
	var err2 error = errors.New("error 2")
	err := Append(err0, err1, err2)
	if len(err.Errors) != 4 {
		t.Fatalf("TestFatalss() expected 4 errors but got %v", len(err.Errors))
	}
	if !err.Fatal {
		t.Fatalf("TestFatalss() expected fatal error but got non-fatal")
	}
	fmt.Printf("TestFatalss()\n%v\n", err)
}
