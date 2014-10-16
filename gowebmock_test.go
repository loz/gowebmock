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
      Expect(last[1]).To(Equal("GET " + expected))
    })
  })

  Describe("when an mocked request is made", func () {
    It("does not fail the test", func () {
      expected := webmock.Url("/known")
      webmock.Expect("GET", expected)

      http.Get(expected)

      Expect(mocktest.ErrorCount()).To(Equal(0))
      webmock.Verify(mocktest)
    })
  })
})
