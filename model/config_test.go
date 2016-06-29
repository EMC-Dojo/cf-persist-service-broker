package model_test

import (
	. "github.com/EMC-Dojo/cf-persist-service-broker/model"

	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	It("Should know how to get a storage host", func() {
		c := GetConfig()
		Expect(c["libstorage.uri"]).To(Equal(os.Getenv("LIBSTORAGE_URI")))
	})
})
