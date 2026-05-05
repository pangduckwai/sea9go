package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestBase(t *testing.T) {
	var err error = errors.New("error!")
	if !IsFatal(err) {
		t.Fatalf("TestBase() expected fatal error but got non-fatal")
	}
	fmt.Printf("TestBase() \"%v\"\n", err)
}

func TestError(t *testing.T) {
	var err0 error
	var err1 error = errors.New("error 1")
	err := Append(err0, err1)
	if Count(err) != 1 {
		t.Fatalf("TestError() expected 1 error but got %v", Count(err))
	}
	if !IsFatal(err) {
		t.Fatalf("TestError() expected fatal error but got non-fatal")
	}
	fmt.Printf("TestError() \"%v\"\n", err)
}

func TestNonFatal(t *testing.T) {
	var err0 error
	var err1 error = NonFatal("error 1")
	err := Append(err0, err1)
	if Count(err) != 1 {
		t.Fatalf("TestNonFatal() expected 1 error but got %v", Count(err))
	}
	if IsFatal(err) {
		t.Fatalf("TestNonFatal() expected non-fatal error but got fatal")
	}
	fmt.Printf("TestNonFatal() \"%v\"\n", err)
}

func TestChain(t *testing.T) {
	var err0 error
	var err1 error = Fatal("[CHAIN] error 1")
	err := Append(err0, err1)
	if Count(err) != 1 {
		t.Fatalf("TestChain() expected 1 error but got %v", Count(err))
	}
	if !IsFatal(err) {
		t.Fatalf("TestChain() expected fatal error but got non-fatal")
	}
	fmt.Printf("TestChain() \"%v\"\n\n", err)
}

func TestErrors(t *testing.T) {
	var err0 error
	var err1 error = errors.New("error 1")
	var err2 error = errors.New("error 2")
	err := Append(err0, err1, err2)
	if Count(err) != 2 {
		t.Fatalf("TestErrors() expected 2 errors but got %v", Count(err))
	}
	if !IsFatal(err) {
		t.Fatalf("TestErrors() expected fatal error but got non-fatal")
	}
	fmt.Printf("TestErrors():\n%v\n\n", err)
}

func TestWithBase(t *testing.T) {
	err0 := NonFatal("error 0")
	var err1 error = errors.New("error 1")
	var err2 error = errors.New("error 2")
	err := Append(err0, err1, err2)
	if Count(err) != 3 {
		t.Fatalf("TestWithBase() expected 3 errors but got %v", Count(err))
	}
	if IsFatal(err) {
		t.Fatalf("TestWithBase() expected non-fatal error but got fatal")
	}
	fmt.Printf("TestWithBase()\n%v\n\n", err)
}

func TestErr(t *testing.T) {
	var err0 error
	var err1 error = NonFatal("error 1")
	var err2 error = errors.New("error 2")
	err := Append(err0, err1, err2)
	if Count(err) != 2 {
		t.Fatalf("TestErr() expected 2 errors but got %v", Count(err))
	}
	if IsFatal(err) {
		t.Fatalf("TestErr() expected non-fatal error but got fatal")
	}
	fmt.Printf("TestErr()\n%v\n\n", err)
}

func TestNil(t *testing.T) {
	var err0 *Err
	var err1 error = NonFatal("error 1")
	var err2 error = errors.New("error 2")
	err := Append(err0, err1, err2)
	if Count(err) != 2 {
		t.Fatalf("TestNil() expected 2 errors but got %v", Count(err))
	}
	if IsFatal(err) {
		t.Fatalf("TestNil() expected non-fatal error but got fatal")
	}
	fmt.Printf("TestNil()\n%v\n\n", err)
}

func TestNils(t *testing.T) {
	var err0 error = NonFatal("error 0")
	var err1 *Err
	var err2 error = errors.New("error 2")
	err := Append(err0, err1, err2)
	if Count(err) != 2 {
		t.Fatalf("TestNils() expected 2 errors but got %v", Count(err))
	}
	if IsFatal(err) {
		t.Fatalf("TestNils() expected non-fatal error but got fatal")
	}
	fmt.Printf("TestNils()\n%v\n\n", err)
}

func TestNilss(t *testing.T) {
	var err0 error = NonFatal("error 0")
	var err1 error
	var err2 error = errors.New("error 2")
	err := Append(err0, err1, err2)
	if Count(err) != 2 {
		t.Fatalf("TestNilss() expected 2 errors but got %v", Count(err))
	}
	if IsFatal(err) {
		t.Fatalf("TestNilss() expected non-fatal error but got fatal")
	}
	fmt.Printf("TestNilss()\n%v\n\n", err)
}

func TestPlains(t *testing.T) {
	var err0 *Err
	var err1 error = errors.New("error 1")
	var err2 error = errors.New("error 2")
	err := Append(err0, err1, err2)
	if Count(err) != 2 {
		t.Fatalf("TestPlains() expected 2 errors but got %v", Count(err))
	}
	if !IsFatal(err) {
		t.Fatalf("TestPlains() expected fatal error but got non-fatal")
	}
	fmt.Printf("TestPlains()\n%v\n\n", err)
}

func TestNonFatals(t *testing.T) {
	err0 := New(false)
	var err1 error = errors.New("error 1")
	var err2 error = errors.New("error 2")
	err := Append(err0, err1, err2)
	if Count(err) != 2 {
		t.Fatalf("TestNonFatals() expected 2 errors but got %v", Count(err))
	}
	if IsFatal(err) {
		t.Fatalf("TestNonFatals() expected non-fatal error but got fatal")
	}
	fmt.Printf("TestNonFatals()\n%v\n\n", err)
}

func TestNonFatalss(t *testing.T) {
	err0 := New(false, "error 0", "error 3")
	var err1 error = errors.New("error 1")
	var err2 error = errors.New("error 2")
	err := Append(err0, err1, err2)
	if Count(err) != 4 {
		t.Fatalf("TestNonFatalss() expected 4 errors but got %v", Count(err))
	}
	if IsFatal(err) {
		t.Fatalf("TestNonFatalss() expected non-fatal error but got fatal")
	}
	fmt.Printf("TestNonFatalss()\n%v\n\n", err)
}

func TestErrs(t *testing.T) {
	var err0 error = errors.New("error 0")
	err1 := New(false, "error 1", "error 3")
	var err2 error = errors.New("error 2")
	err := Append(err0, err1, err2)
	if Count(err) != 4 {
		t.Fatalf("TestErrs() expected 4 errors but got %v", Count(err))
	}
	if IsFatal(err) {
		t.Fatalf("TestErrs() expected non-fatal error but got fatal")
	}
	fmt.Printf("TestErrs()\n%v\n\n", err)
}

func TestPackage(t *testing.T) {
	err0 := New(false, "yo!")
	var err1 error = errors.New("error 1")
	var err2 error = NonFatal("error 2")
	err := Append(err0, err1)
	err.errors = append(err.errors, err2)
	if Count(err) != 3 {
		t.Fatalf("TestErrs() expected 3 errors but got %v", Count(err))
	}
	if IsFatal(err) {
		t.Fatalf("TestErrs() expected non-fatal error but got fatal")
	}
	fmt.Printf("TestErrs()\n%v\n\n", err)
}
