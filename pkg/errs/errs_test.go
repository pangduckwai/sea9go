package errs

import (
	"errors"
	"fmt"
	"testing"
)

const ERR0 = "unit-tests (errs) 0"
const ERR1 = "unit-tests (errs) 1"
const ERR2 = "unit-tests (errs) 2"
const ERR3 = "unit-tests (errs) 3"

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
		errx = Append(errx, NonFatalf("unit-tests (errs) %v", i))
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
	for _, e := range errs {
		err = Append(err, e)
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
	for _, e := range errs {
		err = Append(err, e)
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

func TestAppendFatalFlag(t *testing.T) {
	var errs []error = []error{
		errors.New(ERR1), nil, NonFatal(ERR3), NonFatal(""), nil,
	}

	var err1, err2 error
	for _, e := range errs {
		err1 = Append(err1, e)
	}
	err2 = Append(err2, errs...)

	if Count(err1) != 2 {
		t.Fatalf("TestAppendFatalFlag() expected 2 errors but got %v", Count(err1))
	}
	if !IsFatal(err1) {
		t.Fatalf("TestAppendFatalFlag() expected fatal error but got non-fatal")
	}

	if Count(err2) != 2 {
		t.Fatalf("TestAppendFatalFlag() expected 2 errors but got %v", Count(err2))
	}
	if IsFatal(err2) {
		t.Fatalf("TestAppendFatalFlag() expected non-fatal error but got fatal")
	}
}

func TestFlatten(t *testing.T) {
	var err error
	var errs []error = []error{
		errors.New(ERR0),
		nil,
		New(false, ERR2, ERR3),
		NonFatal(ERR1),
		nil,
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

func TestAppendBaseBehavior(t *testing.T) {
	var err error = New(true)
	var errs []error = []error{
		Fatal(""), NonFatal(ERR1), nil, errors.New(ERR3), Fatal(""), nil,
	}
	cnt := 0
	for _, e := range errs {
		err = Append(err, e)
		cnt++
	}
	if cnt != 6 {
		t.Fatalf("TestAppendBaseBehavior() expected to run 6 times but actually %v", cnt)
	}
	if Count(err) != 2 {
		t.Fatalf("TestAppendBaseBehavior() expected 2 errors but got %v", Count(err))
	}
	if !IsFatal(err) {
		t.Fatalf("TestAppendBaseBehavior() expected fatal error but got non-fatal")
	}
	x := err.(*Err)
	fmt.Printf("TestAppendBaseBehavior() count: %v (%v)\n%v\n\n", Count(err), len(x.errors), err)
}

func TestAppendBaseBehavior2(t *testing.T) {
	var errs []error = []error{
		Fatal(""), nil, nil, NonFatal(ERR3), Fatal(""), nil,
	}

	var err1 error = New(false)
	var err2 error
	for _, e := range errs {
		err1 = Append(err1, e)
		err2 = Append(err2, e)
	}

	x1 := err1.(*Err)
	x2 := err2.(*Err)

	if Count(err1) != 1 && len(x1.errors) != 2 {
		t.Fatalf("TestAppendBaseBehavior2() expected 1(2) errors but got %v(%v)", Count(err1), len(x1.errors))
	}
	if IsFatal(err1) {
		t.Fatalf("TestAppendBaseBehavior2() expected non-fatal error but got fatal")
	}
	fmt.Printf("TestAppendBaseBehavior2() count: %v (%v)\n%v\n\n", Count(err1), len(x1.errors), err1)

	if Count(err2) != 1 && len(x2.errors) != 2 {
		t.Fatalf("TestAppendBaseBehavior3() expected 1(2) errors but got %v(%v)", Count(err2), len(x2.errors))
	}
	if IsFatal(err2) {
		t.Fatalf("TestAppendBaseBehavior3() expected non-fatal error but got fatal")
	}
	fmt.Printf("TestAppendBaseBehavior3() count: %v (%v)\n%v\n\n", Count(err2), len(x2.errors), err2)
}

func TestWithinPackage(t *testing.T) {
	err0 := New(false, "yo!")
	var err1 error = NonFatal(ERR1)
	var err2 error = NonFatal(ERR2)
	err := Append(err0, err1)
	err.errors = append(err.errors, err2)
	if Count(err) != 3 {
		t.Fatalf("TestWithinPackage() expected 3 errors but got %v", Count(err))
	}
	if IsFatal(err) {
		t.Fatalf("TestWithinPackage() expected non-fatal error but got fatal")
	}
	if err.errors[1].Error() != ERR1 {
		t.Fatalf("TestWithinPackage() expected \"%v\" but got \"%v\"", ERR1, err.errors[1])
	}
	if err.errors[2].Error() != "[non-fatal] "+ERR2 {
		t.Fatalf("TestWithinPackage() expected \"%v\" but got \"%v\"", "[non-fatal] "+ERR2+".", err.errors[2])
	}
	fmt.Printf("TestWithinPackage()\n%v\n\n", err)
}
