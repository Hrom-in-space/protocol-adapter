package mmock

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

import (
	"net/http"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// IHttpClientMock implements smartclient.IHttpClient
type IHttpClientMock struct {
	t minimock.Tester

	funcDo          func(req *http.Request) (rp1 *http.Response, err error)
	inspectFuncDo   func(req *http.Request)
	afterDoCounter  uint64
	beforeDoCounter uint64
	DoMock          mIHttpClientMockDo
}

// NewIHttpClientMock returns a mock for smartclient.IHttpClient
func NewIHttpClientMock(t minimock.Tester) *IHttpClientMock {
	m := &IHttpClientMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.DoMock = mIHttpClientMockDo{mock: m}
	m.DoMock.callArgs = []*IHttpClientMockDoParams{}

	return m
}

type mIHttpClientMockDo struct {
	mock               *IHttpClientMock
	defaultExpectation *IHttpClientMockDoExpectation
	expectations       []*IHttpClientMockDoExpectation

	callArgs []*IHttpClientMockDoParams
	mutex    sync.RWMutex
}

// IHttpClientMockDoExpectation specifies expectation struct of the IHttpClient.Do
type IHttpClientMockDoExpectation struct {
	mock    *IHttpClientMock
	params  *IHttpClientMockDoParams
	results *IHttpClientMockDoResults
	Counter uint64
}

// IHttpClientMockDoParams contains parameters of the IHttpClient.Do
type IHttpClientMockDoParams struct {
	req *http.Request
}

// IHttpClientMockDoResults contains results of the IHttpClient.Do
type IHttpClientMockDoResults struct {
	rp1 *http.Response
	err error
}

// Expect sets up expected params for IHttpClient.Do
func (mmDo *mIHttpClientMockDo) Expect(req *http.Request) *mIHttpClientMockDo {
	if mmDo.mock.funcDo != nil {
		mmDo.mock.t.Fatalf("IHttpClientMock.Do mock is already set by Set")
	}

	if mmDo.defaultExpectation == nil {
		mmDo.defaultExpectation = &IHttpClientMockDoExpectation{}
	}

	mmDo.defaultExpectation.params = &IHttpClientMockDoParams{req}
	for _, e := range mmDo.expectations {
		if minimock.Equal(e.params, mmDo.defaultExpectation.params) {
			mmDo.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmDo.defaultExpectation.params)
		}
	}

	return mmDo
}

// Inspect accepts an inspector function that has same arguments as the IHttpClient.Do
func (mmDo *mIHttpClientMockDo) Inspect(f func(req *http.Request)) *mIHttpClientMockDo {
	if mmDo.mock.inspectFuncDo != nil {
		mmDo.mock.t.Fatalf("Inspect function is already set for IHttpClientMock.Do")
	}

	mmDo.mock.inspectFuncDo = f

	return mmDo
}

// Return sets up results that will be returned by IHttpClient.Do
func (mmDo *mIHttpClientMockDo) Return(rp1 *http.Response, err error) *IHttpClientMock {
	if mmDo.mock.funcDo != nil {
		mmDo.mock.t.Fatalf("IHttpClientMock.Do mock is already set by Set")
	}

	if mmDo.defaultExpectation == nil {
		mmDo.defaultExpectation = &IHttpClientMockDoExpectation{mock: mmDo.mock}
	}
	mmDo.defaultExpectation.results = &IHttpClientMockDoResults{rp1, err}
	return mmDo.mock
}

// Set uses given function f to mock the IHttpClient.Do method
func (mmDo *mIHttpClientMockDo) Set(f func(req *http.Request) (rp1 *http.Response, err error)) *IHttpClientMock {
	if mmDo.defaultExpectation != nil {
		mmDo.mock.t.Fatalf("Default expectation is already set for the IHttpClient.Do method")
	}

	if len(mmDo.expectations) > 0 {
		mmDo.mock.t.Fatalf("Some expectations are already set for the IHttpClient.Do method")
	}

	mmDo.mock.funcDo = f
	return mmDo.mock
}

// When sets expectation for the IHttpClient.Do which will trigger the result defined by the following
// Then helper
func (mmDo *mIHttpClientMockDo) When(req *http.Request) *IHttpClientMockDoExpectation {
	if mmDo.mock.funcDo != nil {
		mmDo.mock.t.Fatalf("IHttpClientMock.Do mock is already set by Set")
	}

	expectation := &IHttpClientMockDoExpectation{
		mock:   mmDo.mock,
		params: &IHttpClientMockDoParams{req},
	}
	mmDo.expectations = append(mmDo.expectations, expectation)
	return expectation
}

// Then sets up IHttpClient.Do return parameters for the expectation previously defined by the When method
func (e *IHttpClientMockDoExpectation) Then(rp1 *http.Response, err error) *IHttpClientMock {
	e.results = &IHttpClientMockDoResults{rp1, err}
	return e.mock
}

// Do implements smartclient.IHttpClient
func (mmDo *IHttpClientMock) Do(req *http.Request) (rp1 *http.Response, err error) {
	mm_atomic.AddUint64(&mmDo.beforeDoCounter, 1)
	defer mm_atomic.AddUint64(&mmDo.afterDoCounter, 1)

	if mmDo.inspectFuncDo != nil {
		mmDo.inspectFuncDo(req)
	}

	mm_params := &IHttpClientMockDoParams{req}

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
		mm_got := IHttpClientMockDoParams{req}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmDo.t.Errorf("IHttpClientMock.Do got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmDo.DoMock.defaultExpectation.results
		if mm_results == nil {
			mmDo.t.Fatal("No results are set for the IHttpClientMock.Do")
		}
		return (*mm_results).rp1, (*mm_results).err
	}
	if mmDo.funcDo != nil {
		return mmDo.funcDo(req)
	}
	mmDo.t.Fatalf("Unexpected call to IHttpClientMock.Do. %v", req)
	return
}

// DoAfterCounter returns a count of finished IHttpClientMock.Do invocations
func (mmDo *IHttpClientMock) DoAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmDo.afterDoCounter)
}

// DoBeforeCounter returns a count of IHttpClientMock.Do invocations
func (mmDo *IHttpClientMock) DoBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmDo.beforeDoCounter)
}

// Calls returns a list of arguments used in each call to IHttpClientMock.Do.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmDo *mIHttpClientMockDo) Calls() []*IHttpClientMockDoParams {
	mmDo.mutex.RLock()

	argCopy := make([]*IHttpClientMockDoParams, len(mmDo.callArgs))
	copy(argCopy, mmDo.callArgs)

	mmDo.mutex.RUnlock()

	return argCopy
}

// MinimockDoDone returns true if the count of the Do invocations corresponds
// the number of defined expectations
func (m *IHttpClientMock) MinimockDoDone() bool {
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
func (m *IHttpClientMock) MinimockDoInspect() {
	for _, e := range m.DoMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to IHttpClientMock.Do with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.DoMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterDoCounter) < 1 {
		if m.DoMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to IHttpClientMock.Do")
		} else {
			m.t.Errorf("Expected call to IHttpClientMock.Do with params: %#v", *m.DoMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcDo != nil && mm_atomic.LoadUint64(&m.afterDoCounter) < 1 {
		m.t.Error("Expected call to IHttpClientMock.Do")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *IHttpClientMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockDoInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *IHttpClientMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *IHttpClientMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockDoDone()
}
