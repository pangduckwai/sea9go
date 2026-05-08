package errs

import (
	"errors"
	"fmt"
	"testing"
)

const ERR0 = "error 0"
const ERR1 = "error 1"
const ERR2 = "error 2"
const ERR3 = "error 3"

func TestBuiltInType(t *testing.T) {
	var err error = errors.New(ERR0)
	if !IsFatal(err) {
		t.Fatalf("TestBuiltInType() expected fatal error but got non-fatal")
	}
	if err.Error() != ERR0 {
		t.Fatalf("TestBuiltInType() expected \"%v\" but got \"%v\"", ERR0, err.Error())
	}
}

func TestNonFatalType(t *testing.T) {
	var err error = NonFatal(ERR1)
	if IsFatal(err) {
		t.Fatalf("TestNonFatalType() expected non-fatal error but got fatal")
	}
	if err.Error()[12:] != ERR1 {
		t.Fatalf("TestNonFatalType() expected \"%v\" but got \"%v\"", ERR1, err.Error())
	}
}

func TestFatalType(t *testing.T) {
	var err error = Fatal(ERR2)
	if !IsFatal(err) {
		t.Fatalf("TestFatalType() expected fatal error but got non-fatal")
	}
	if err.Error() != ERR2 {
		t.Fatalf("TestFatalType() expected \"%v\" but got \"%v\"", ERR2, err.Error())
	}
}

func TestWrapBuiltInType(t *testing.T) {
	var err error = Wrap(errors.New(ERR0), "[CHAIN]", "OF", "[LABELS")
	if !IsFatal(err) {
		t.Fatalf("TestWrapBuiltInType() expected fatal error but got non-fatal")
	}
	fmt.Printf("TestWrapBuiltInType() - %v\n", err)
}

func TestLabelChain(t *testing.T) {
	var err error = Wrap(NonFatal(ERR3), "[CHAIN]", "OF", "[LABELS")
	if Count(err) != 1 {
		t.Fatalf("TestLabelChain() expected 1 error but got %v", Count(err))
	}
	if IsFatal(err) {
		t.Fatalf("TestLabelChain() expected non-fatal error but got fatal")
	}
	xpt := fmt.Sprintf("[non-fatal][CHAIN][OF][LABELS] %v", ERR3)
	if err.Error() != xpt {
		t.Fatalf("TestLabelChain() expected \"%v\" but got \"%v\"", xpt, err.Error())
	}
	fmt.Printf("TestLabelChain() - %v\n\n", err)
}

func TestLabelChain2(t *testing.T) {
	var errx *Err
	for i := range 5 {
		errx = Append(errx, NonFatalf("error %v", i))
	}
	err := Wrap(errx, "THIS", "IS", "A", "LIST")
	if Count(err) != 5 {
		t.Fatalf("TestLabelChain2() expected 5 errors but got %v", Count(err))
	}
	if IsFatal(err) {
		t.Fatalf("TestLabelChain2() expected non-fatal error but got fatal")
	}
	fmt.Printf("TestLabelChain2()\n%v\n\n", err)
}

func TestAppendBuiltInBase(t *testing.T) {
	var err error = errors.New(ERR0)
	var errs []error = []error{
		errors.New(ERR1), nil, errors.New(ERR3), errors.New(""), nil,
	}
	cnt := 0
	for _, e := range errs {
		err = Append(err, e)
		cnt++
	}
	if cnt != 5 {
		t.Fatalf("TestAppendBuiltInBase() expected to run 5 times but actually %v", cnt)
	}
	if Count(err) != 3 {
		t.Fatalf("TestAppendBuiltInBase() expected 3 errors but got %v", Count(err))
	}
	if !IsFatal(err) {
		t.Fatalf("TestAppendBuiltInBase() expected fatal error but got non-fatal")
	}
}

func TestAppendNilBase1(t *testing.T) {
	var err error
	var errs []error = []error{
		errors.New(ERR1), nil, errors.New(ERR3), errors.New(""), nil,
	}
	cnt := 0
	for _, e := range errs {
		err = Append(err, e)
		cnt++
	}
	if cnt != 5 {
		t.Fatalf("TestAppendNilBase1() expected to run 5 times but actually %v", cnt)
	}
	if Count(err) != 2 {
		t.Fatalf("TestAppendNilBase1() expected 2 errors but got %v", Count(err))
	}
	if !IsFatal(err) {
		t.Fatalf("TestAppendNilBase1() expected fatal error but got non-fatal")
	}
}

func TestAppendNilBase2(t *testing.T) {
	var err *Err
	var errs []error = []error{
		nil, errors.New(ERR1), nil, errors.New(ERR3), errors.New(""),
	}
	cnt := 0
	for _, e := range errs {
		err = Append(err, e)
		cnt++
	}
	if cnt != 5 {
		t.Fatalf("TestAppendNilBase2() expected to run 5 times but actually %v", cnt)
	}
	if Count(err) != 2 {
		t.Fatalf("TestAppendNilBase2() expected 2 errors but got %v", Count(err))
	}
	if !IsFatal(err) {
		t.Fatalf("TestAppendNilBase2() expected fatal error but got non-fatal")
	}
}

func TestAppendNonFatalsBase(t *testing.T) {
	var err error = NonFatal(ERR0)
	var errs []error = []error{
		errors.New(ERR1), nil, errors.New(ERR3), errors.New(""), nil,
	}
	cnt := 0
	for _, e := range errs {
		err = Append(err, e)
		cnt++
	}
	if cnt != 5 {
		t.Fatalf("TestAppendNonFatalsBase() expected to run 5 times but actually %v", cnt)
	}
	if Count(err) != 3 {
		t.Fatalf("TestAppendNonFatalsBase() expected 3 errors but got %v", Count(err))
	}
	if IsFatal(err) {
		t.Fatalf("TestAppendNonFatalsBase() expected non-fatal error but got fatal")
	}
	fmt.Printf("TestAppendNonFatalsBase()\n%v\n\n", err)
}

func TestAppendNonFatals1(t *testing.T) {
	var err error = NonFatal(ERR0)
	var errs []error = []error{
		NonFatal(ERR1), nil, NonFatal(ERR3), New(false), nil,
	}
	cnt := 0
	for _, e := range errs {
		err = Append(err, e)
		cnt++
	}
	if cnt != 5 {
		t.Fatalf("TestAppendNonFatals1() expected to run 5 times but actually %v", cnt)
	}
	if Count(err) != 3 {
		t.Fatalf("TestAppendNonFatals1() expected 3 errors but got %v", Count(err))
	}
	if IsFatal(err) {
		t.Fatalf("TestAppendNonFatals1() expected non-fatal error but got fatal")
	}
}

func TestAppendNonFatals2(t *testing.T) {
	var err error = NonFatal(ERR0)
	var errs []error = []error{
		NonFatal(ERR1), nil, NonFatal(ERR3), NonFatal(""), nil,
	}
	cnt := 0
	for _, e := range errs {
		err = Append(err, e)
		cnt++
	}
	if cnt != 5 {
		t.Fatalf("TestAppendNonFatals2() expected to run 5 times but actually %v", cnt)
	}
	if Count(err) != 3 {
		t.Fatalf("TestAppendNonFatals2() expected 3 errors but got %v", Count(err))
	}
	if IsFatal(err) {
		t.Fatalf("TestAppendNonFatals2() expected non-fatal error but got fatal")
	}
}

func TestAppendNonFatals3(t *testing.T) {
	var err error = New(false)
	var errs []error = []error{
		NonFatal(ERR1), nil, NonFatal(ERR3), NonFatal(""), nil,
	}
	cnt := 0
	for _, e := range errs {
		err = Append(err, e)
		cnt++
	}
	if cnt != 5 {
		t.Fatalf("TestAppendNonFatals3() expected to run 5 times but actually %v", cnt)
	}
	if Count(err) != 2 {
		t.Fatalf("TestAppendNonFatals3() expected 2 errors but got %v", Count(err))
	}
	if IsFatal(err) {
		t.Fatalf("TestAppendNonFatals3() expected non-fatal error but got fatal")
	}
}

func TestAppendNonFatals4(t *testing.T) {
	var err error = New(false, ERR0, ERR2)
	var errs []error = []error{
		errors.New(ERR1), nil, errors.New(ERR3), errors.New(""), nil,
	}
	cnt := 0
	for _, e := range errs {
		err = Append(err, e)
		cnt++
	}
	if cnt != 5 {
		t.Fatalf("TestAppendNonFatals4() expected to run 5 times but actually %v", cnt)
	}
	if Count(err) != 4 {
		t.Fatalf("TestAppendNonFatals4() expected 4 errors but got %v", Count(err))
	}
	if IsFatal(err) {
		t.Fatalf("TestAppendNonFatals4() expected non-fatal error but got fatal")
	}
}

func TestAppendFatalFlag1(t *testing.T) {
	var err error
	var errs []error = []error{
		errors.New(ERR1), nil, NonFatal(ERR3), NonFatal(""), nil,
	}
	cnt := 0
	for _, e := range errs {
		err = Append(err, e)
		cnt++
	}
	if cnt != 5 {
		t.Fatalf("TestAppendFatalFlag1() expected to run 5 times but actually %v", cnt)
	}
	if Count(err) != 2 {
		t.Fatalf("TestAppendFatalFlag1() expected 2 errors but got %v", Count(err))
	}
	if !IsFatal(err) {
		t.Fatalf("TestAppendFatalFlag1() expected fatal error but got non-fatal")
	}
	fmt.Printf("TestAppendFatalFlag1()\n%v\n\n", err)
}

func TestAppendFatalFlag2(t *testing.T) {
	var err0 error
	var errs []error = []error{
		errors.New(ERR1), nil, NonFatal(ERR3), NonFatal(""), nil,
	}
	err := Append(err0, errs...)
	if Count(err) != 2 {
		t.Fatalf("TestAppendFatalFlag2() expected 2 errors but got %v", Count(err))
	}
	if IsFatal(err) {
		t.Fatalf("TestAppendFatalFlag2() expected non-fatal error but got fatal")
	}
	fmt.Printf("TestAppendFatalFlag2()\n%v\n\n", err)
}

func TestFlatten(t *testing.T) {
	var err error
	var err2 error = New(false, ERR2, ERR3)
	var errs []error = []error{
		errors.New(ERR0), nil, err2, NonFatal(ERR1), nil,
	}
	cnt := 0
	for _, e := range errs {
		err = Append(err, e)
		cnt++
	}
	if cnt != 5 {
		t.Fatalf("TestFlatten() expected to run 5 times but actually %v", cnt)
	}
	if Count(err) != 4 {
		t.Fatalf("TestFlatten() expected 4 errors but got %v", Count(err))
	}
	if !IsFatal(err) {
		t.Fatalf("TestFlatten() expected fatal error but got non-fatal")
	}
	fmt.Printf("TestFlatten()\n%v\n\n", err)
}

func TestPackage(t *testing.T) {
	err0 := New(false, "yo!")
	var err1 error = errors.New("error 1")
	var err2 error = NonFatal("error 2")
	err := Append(err0, err1)
	err.errors = append(err.errors, err2)
	if Count(err) != 3 {
		t.Fatalf("TestPackage() expected 3 errors but got %v", Count(err))
	}
	if IsFatal(err) {
		t.Fatalf("TestPackage() expected non-fatal error but got fatal")
	}
	fmt.Printf("TestPackage()\n%v\n\n", err)
}
