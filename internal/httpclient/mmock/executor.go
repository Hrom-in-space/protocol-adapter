package mmock

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

import (
	"context"
	mm_httpclient "protocol-adapter/internal/httpclient"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// ExecutorMock implements httpclient.Executor
type ExecutorMock struct {
	t minimock.Tester

	funcDo          func(ctx context.Context, req *mm_httpclient.Request) (rp1 *mm_httpclient.Response, err error)
	inspectFuncDo   func(ctx context.Context, req *mm_httpclient.Request)
	afterDoCounter  uint64
	beforeDoCounter uint64
	DoMock          mExecutorMockDo
}

// NewExecutorMock returns a mock for httpclient.Executor
func NewExecutorMock(t minimock.Tester) *ExecutorMock {
	m := &ExecutorMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.DoMock = mExecutorMockDo{mock: m}
	m.DoMock.callArgs = []*ExecutorMockDoParams{}

	return m
}

type mExecutorMockDo struct {
	mock               *ExecutorMock
	defaultExpectation *ExecutorMockDoExpectation
	expectations       []*ExecutorMockDoExpectation

	callArgs []*ExecutorMockDoParams
	mutex    sync.RWMutex
}

// ExecutorMockDoExpectation specifies expectation struct of the Executor.Do
type ExecutorMockDoExpectation struct {
	mock    *ExecutorMock
	params  *ExecutorMockDoParams
	results *ExecutorMockDoResults
	Counter uint64
}

// ExecutorMockDoParams contains parameters of the Executor.Do
type ExecutorMockDoParams struct {
	ctx context.Context
	req *mm_httpclient.Request
}

// ExecutorMockDoResults contains results of the Executor.Do
type ExecutorMockDoResults struct {
	rp1 *mm_httpclient.Response
	err error
}

// Expect sets up expected params for Executor.Do
func (mmDo *mExecutorMockDo) Expect(ctx context.Context, req *mm_httpclient.Request) *mExecutorMockDo {
	if mmDo.mock.funcDo != nil {
		mmDo.mock.t.Fatalf("ExecutorMock.Do mock is already set by Set")
	}

	if mmDo.defaultExpectation == nil {
		mmDo.defaultExpectation = &ExecutorMockDoExpectation{}
	}

	mmDo.defaultExpectation.params = &ExecutorMockDoParams{ctx, req}
	for _, e := range mmDo.expectations {
		if minimock.Equal(e.params, mmDo.defaultExpectation.params) {
			mmDo.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmDo.defaultExpectation.params)
		}
	}

	return mmDo
}

// Inspect accepts an inspector function that has same arguments as the Executor.Do
func (mmDo *mExecutorMockDo) Inspect(f func(ctx context.Context, req *mm_httpclient.Request)) *mExecutorMockDo {
	if mmDo.mock.inspectFuncDo != nil {
		mmDo.mock.t.Fatalf("Inspect function is already set for ExecutorMock.Do")
	}

	mmDo.mock.inspectFuncDo = f

	return mmDo
}

// Return sets up results that will be returned by Executor.Do
func (mmDo *mExecutorMockDo) Return(rp1 *mm_httpclient.Response, err error) *ExecutorMock {
	if mmDo.mock.funcDo != nil {
		mmDo.mock.t.Fatalf("ExecutorMock.Do mock is already set by Set")
	}

	if mmDo.defaultExpectation == nil {
		mmDo.defaultExpectation = &ExecutorMockDoExpectation{mock: mmDo.mock}
	}
	mmDo.defaultExpectation.results = &ExecutorMockDoResults{rp1, err}
	return mmDo.mock
}

// Set uses given function f to mock the Executor.Do method
func (mmDo *mExecutorMockDo) Set(f func(ctx context.Context, req *mm_httpclient.Request) (rp1 *mm_httpclient.Response, err error)) *ExecutorMock {
	if mmDo.defaultExpectation != nil {
		mmDo.mock.t.Fatalf("Default expectation is already set for the Executor.Do method")
	}

	if len(mmDo.expectations) > 0 {
		mmDo.mock.t.Fatalf("Some expectations are already set for the Executor.Do method")
	}

	mmDo.mock.funcDo = f
	return mmDo.mock
}

// When sets expectation for the Executor.Do which will trigger the result defined by the following
// Then helper
func (mmDo *mExecutorMockDo) When(ctx context.Context, req *mm_httpclient.Request) *ExecutorMockDoExpectation {
	if mmDo.mock.funcDo != nil {
		mmDo.mock.t.Fatalf("ExecutorMock.Do mock is already set by Set")
	}

	expectation := &ExecutorMockDoExpectation{
		mock:   mmDo.mock,
		params: &ExecutorMockDoParams{ctx, req},
	}
	mmDo.expectations = append(mmDo.expectations, expectation)
	return expectation
}

// Then sets up Executor.Do return parameters for the expectation previously defined by the When method
func (e *ExecutorMockDoExpectation) Then(rp1 *mm_httpclient.Response, err error) *ExecutorMock {
	e.results = &ExecutorMockDoResults{rp1, err}
	return e.mock
}

// Do implements httpclient.Executor
func (mmDo *ExecutorMock) Do(ctx context.Context, req *mm_httpclient.Request) (rp1 *mm_httpclient.Response, err error) {
	mm_atomic.AddUint64(&mmDo.beforeDoCounter, 1)
	defer mm_atomic.AddUint64(&mmDo.afterDoCounter, 1)

	if mmDo.inspectFuncDo != nil {
		mmDo.inspectFuncDo(ctx, req)
	}

	mm_params := &ExecutorMockDoParams{ctx, req}

	// Record call args
	mmDo.DoMock.mutex.Lock()
	mmDo.DoMock.callArgs = append(mmDo.DoMock.callArgs, mm_params)
	mmDo.DoMock.mutex.Unlock()

	for _, e := range mmDo.DoMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.rp1, e.results.err
		}
	}

	if mmDo.DoMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmDo.DoMock.defaultExpectation.Counter, 1)
		mm_want := mmDo.DoMock.defaultExpectation.params
		mm_got := ExecutorMockDoParams{ctx, req}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmDo.t.Errorf("ExecutorMock.Do got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmDo.DoMock.defaultExpectation.results
		if mm_results == nil {
			mmDo.t.Fatal("No results are set for the ExecutorMock.Do")
		}
		return (*mm_results).rp1, (*mm_results).err
	}
	if mmDo.funcDo != nil {
		return mmDo.funcDo(ctx, req)
	}
	mmDo.t.Fatalf("Unexpected call to ExecutorMock.Do. %v %v", ctx, req)
	return
}

// DoAfterCounter returns a count of finished ExecutorMock.Do invocations
func (mmDo *ExecutorMock) DoAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmDo.afterDoCounter)
}

// DoBeforeCounter returns a count of ExecutorMock.Do invocations
func (mmDo *ExecutorMock) DoBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmDo.beforeDoCounter)
}

// Calls returns a list of arguments used in each call to ExecutorMock.Do.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmDo *mExecutorMockDo) Calls() []*ExecutorMockDoParams {
	mmDo.mutex.RLock()

	argCopy := make([]*ExecutorMockDoParams, len(mmDo.callArgs))
	copy(argCopy, mmDo.callArgs)

	mmDo.mutex.RUnlock()

	return argCopy
}

// MinimockDoDone returns true if the count of the Do invocations corresponds
// the number of defined expectations
func (m *ExecutorMock) MinimockDoDone() bool {
	for _, e := range m.DoMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.DoMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterDoCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcDo != nil && mm_atomic.LoadUint64(&m.afterDoCounter) < 1 {
		return false
	}
	return true
}

// MinimockDoInspect logs each unmet expectation
func (m *ExecutorMock) MinimockDoInspect() {
	for _, e := range m.DoMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ExecutorMock.Do with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.DoMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterDoCounter) < 1 {
		if m.DoMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to ExecutorMock.Do")
		} else {
			m.t.Errorf("Expected call to ExecutorMock.Do with params: %#v", *m.DoMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcDo != nil && mm_atomic.LoadUint64(&m.afterDoCounter) < 1 {
		m.t.Error("Expected call to ExecutorMock.Do")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *ExecutorMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockDoInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *ExecutorMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *ExecutorMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockDoDone()
}
