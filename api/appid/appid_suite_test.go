package appid

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAppID(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AppID Suite")
}
