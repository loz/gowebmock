package gowebmock

import (
  "net/http"
  "net/http/httptest"
)

type ErrorReporter interface {
  Error(args ...interface{})
}

type Expectation struct {
  Verb string
  Path string
  body string
  headers map[string]string
}

func (e *Expectation) Matches(verb string, path string) bool {
  return (e.Verb == verb && e.Path == path)
}

func (e *Expectation) Respond(w http.ResponseWriter) {
  for h,value := range e.headers {
    w.Header()[h] = []string{value}
  }
  w.Write([]byte(e.body))
}

func (e *Expectation) WithBody(body string) *Expectation {
  e.body = body
  return e
}

func (e *Expectation) WithHeaders(headers map[string]string) *Expectation {
  e.headers = headers
  return e
}


type WebMock struct {
  testServer *httptest.Server
  t ErrorReporter
  expectations []*Expectation
}

func NewWebMock(t ErrorReporter) *WebMock {
  mock := new(WebMock)
  mock.t = t
  mock.expectations = make([]*Expectation, 0)
  mock.testServer = httptest.NewServer(http.HandlerFunc(mock.HttpHandler))
  return mock
}

func (m *WebMock) Url(path string) string {
  return m.testServer.URL + path
}

func (m *WebMock) Expect(verb string, path string) *Expectation {
  expected := Expectation{Verb: verb, Path: path}
  m.expectations = append(m.expectations, &expected)
  return &expected
}

func (m *WebMock) FindExpected(verb string, path string) *Expectation {
  for _, e := range m.expectations {
    if e.Matches(verb, path) {
      return e
    }
  }
  return nil
}

func (m *WebMock) Verify(t ErrorReporter) {
  m.testServer.Close()
}

func (m *WebMock) HttpHandler(w http.ResponseWriter, r *http.Request) {
  path := r.URL.String()
  expected := m.FindExpected(r.Method, path)
  if expected == nil {
    m.t.Error("Unexpected Request", r.Method + " " + path)
  } else {
    expected.Respond(w)
  }
}
