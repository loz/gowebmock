package gowebmock_test

import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    "testing"
    )

func TestGowebmock(t *testing.T) {
  RegisterFailHandler(Fail)
    RunSpecs(t, "Gowebmock Suite")
}
