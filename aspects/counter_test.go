package ginmon

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const TestMode string = "test"

const checkMark = "\u2713"
const ballotX = "\u2717"

const testpath = "/foo/bar"

func internalGinCtx() *gin.Context {
	return &gin.Context{
		Request: &http.Request{
			URL: &url.URL{
				Path: testpath,
			},
		},
	}
}

func Test_Inc(t *testing.T) {
	ca := NewCounterAspect()
	expect := 1
	ca.Inc(internalGinCtx())
	ca.reset()
	if assert.Equal(t, ca.RequestsSum, expect, "Incrementation of counter does not work, expect %d but got %d %s",
		expect, ca.RequestsSum, ballotX) {
		t.Logf("Incrementation of counter works, expect %d and git %d %s",
			expect, ca.RequestsSum, checkMark)
	}
}

func Test_GetStats(t *testing.T) {
	ca := NewCounterAspect()
	if assert.NotNil(t, ca.GetStats(), "Return of Getstats() should not be nil") {
		t.Logf("Should be an interface %s", checkMark)
	}

	newCa := ca.GetStats().(CounterAspect)
	expect := 0
	if assert.Equal(t, newCa.RequestsSum, expect, "Return of Getstats() does not work, expect %d but got %d %s",
		expect, newCa.RequestsSum, ballotX) {
		t.Logf("Return of Getstats() works, expect %d and got %d %s",
			expect, newCa.RequestsSum, checkMark)
	}

	ca.Inc(internalGinCtx())
	if assert.Equal(t, newCa.RequestsSum, expect, "Return of Getstats() does not work, expect %d but got %d %s",
		expect, newCa.RequestsSum, ballotX) {
		t.Logf("Return of Getstats() works, expect %d and got %d %s",
			expect, newCa.RequestsSum, checkMark)
	}
	if assert.Equal(t, newCa.Requests[testpath], expect, "Return of Getstats() does not work, expect %d but got %d %s",
		expect, newCa.Requests[testpath], ballotX) {
		t.Logf("Return of Getstats() works, expect %d and got %d %s",
			expect, newCa.Requests[testpath], checkMark)
	}

	ca.reset()
	newCa = ca.GetStats().(CounterAspect)
	expect = 1
	if assert.Equal(t, newCa.RequestsSum, expect, "Return of Getstats() does not work, expect %d but got %d %s",
		expect, newCa.RequestsSum, ballotX) {
		t.Logf("Return of Getstats() works, expect %d and got %d %s",
			expect, newCa.RequestsSum, checkMark)
	}
	if assert.Equal(t, newCa.Requests[testpath], expect, "Return of Getstats() does not work, expect %d but got %d %s",
		expect, newCa.Requests[testpath], ballotX) {
		t.Logf("Return of Getstats() works, expect %d and got %d %s",
			expect, newCa.Requests[testpath], checkMark)
	}
}

func Test_Name(t *testing.T) {
	ca := NewCounterAspect()
	expect := "Counter"
	if assert.Equal(t, ca.Name(), expect, "Return of counter name does not work, expect %s but got %s %s",
		expect, ca.Name(), ballotX) {
		t.Logf("Return of counter name works, expect %s and got %s %s",
			expect, ca.Name(), checkMark)
	}
}

func Test_InRoot(t *testing.T) {
	ca := NewCounterAspect()
	expect := false
	if assert.Equal(t, ca.InRoot(), expect, "Expect %v but got %v %s",
		expect, ca.InRoot(), ballotX) {
		t.Logf("Expect %v and got %v %s",
			expect, ca.InRoot(), checkMark)
	}
}

func Test_CounterHandler(t *testing.T) {
	gin.SetMode(TestMode)
	router := gin.New()
	ca := NewCounterAspect()
	expect := 1
	ca.Inc(internalGinCtx())
	ca.reset()

	router.Use(CounterHandler(ca))
	tryRequest(router, "GET", "/")
	if assert.Equal(t, ca.RequestsSum, expect, "Incrementation of counter does not work, expect %d but got %d %s", expect, ca.RequestsSum, ballotX) {
		t.Logf("CounterHandler works, expect %d and got %d %s", expect, ca.RequestsSum, checkMark)
	}
}

func tryRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
