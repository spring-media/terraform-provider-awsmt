package models

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

func TestSchedulerModels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "awsmt/models")
}
