package e2e_test

import (
	"github.com/codetent/weasel/e2e/helper"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/yuk7/wsllib-go"
)

var _ = Describe("explore", func() {
	AfterEach(func() {
		wsllib.WslUnregisterDistribution("busybox")
	})

	It("should fail if the environment is undefined", func() {
		cmd, err := helper.NewWeaselCommand("explore", "busybox", "--show-only")
		Expect(err).NotTo(HaveOccurred())

		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Expect(session.Wait(DEFAULT_TIMEOUT)).Should(gexec.Exit(1))
		Expect(session.Err).Should(gbytes.Say("environment busybox not available. Enter it first"))
	})

	It("should show the path of an available environment", func() {
		cmd, err := helper.NewWeaselCommand("enter", "busybox", "--register")
		Expect(err).NotTo(HaveOccurred())

		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Expect(session.Wait(DEFAULT_TIMEOUT)).Should(gexec.Exit(0))

		Expect(wsllib.WslIsDistributionRegistered("busybox")).Should(BeTrue())

		cmd, err = helper.NewWeaselCommand("explore", "busybox", "--show-only")
		Expect(err).NotTo(HaveOccurred())

		session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Expect(session.Wait(DEFAULT_TIMEOUT)).Should(gexec.Exit(0))
		Expect(session.Out).Should(gbytes.Say("\\\\wsl\\$\\\\busybox\\\\root"))
	})
})
