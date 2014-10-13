package gowebmock_test

import (
    . "github.com/loz/gowebmock"
    "net/http"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    )

type MockTest struct {
  errors []([]interface{})
}

func NewMockTest() *MockTest {
  t := new(MockTest)
  t.errors = make([]([]interface{}), 0)
  return t
}

func (t *MockTest) Error(args ...interface{}) {
  t.errors = append(t.errors, args)
}

func (t *MockTest) ErrorCount() int {
  return len(t.errors)
}

func (t *MockTest) LastError() []interface{} {
  return t.errors[0]
}

var _ = Describe("Gowebmock", func() {
  var webmock *WebMock
  var mocktest *MockTest

  BeforeEach(func () {
    mocktest = NewMockTest()
    webmock = NewWebMock(mocktest)
  })

  Describe("when an unknown request is made", func () {
    It("fails test with error message", func () {
      webmock.Expect("GET", "something")
      expected := webmock.Url("/unknown")
      http.Get(expected)

      last := mocktest.LastError()
      Expect(last[0]).To(Equal("Unexpected Request"))
      Expect(last[1]).To(Equal("GET /unknown"))
    })
  })

  Describe("when an mocked request is made", func () {
    It("does not fail the test", func () {
      expected := webmock.Url("/known")
      webmock.Expect("GET", "/known")

      http.Get(expected)

      Expect(mocktest.ErrorCount()).To(Equal(0))
      webmock.Verify(mocktest)
    })

    It("returns supplied body", func () {
      body := make([]byte, 1024)
      expected := "Some HTTP Body"
      webmock.Expect("GET", "/known").WithBody(expected)

      res, _ := http.Get(webmock.Url("/known"))

      bytes, _ := res.Body.Read(body)
      actual := string(body[:bytes])
      Expect(actual).To(BeEquivalentTo(expected))
    })

    It("gives supplied headers", func () {
      headers := map[string]string{"X-Header":"Value"}

      webmock.Expect("GET", "/known").WithHeaders(headers)
      res, _ := http.Get(webmock.Url("/known"))

      header := res.Header["X-Header"][0]
      Expect(header).To(Equal("Value"))
    })
  })
})
