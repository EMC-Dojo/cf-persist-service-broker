package libstoragewrapper_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLibstoragewrapper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Libstoragewrapper Suite")
}
