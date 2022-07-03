package e2e_test

import (
	"github.com/codetent/weasel/e2e/helper"
	"github.com/codetent/weasel/pkg/weasel/wsl"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/yuk7/wsllib-go"
)

var _ = Describe("enter", Ordered, func() {
	AfterEach(func() {
		wsllib.WslUnregisterDistribution("busybox")
	})

	It("should fail registering an environment with an invalid image", func() {
		cmd, err := helper.NewWeaselCommand("enter", "_invalid_:foo", "--register")
		Expect(err).NotTo(HaveOccurred())

		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)

		Expect(err).NotTo(HaveOccurred())
		Expect(session.Wait(DEFAULT_TIMEOUT)).Should(gexec.Exit(1))
		Expect(session.Err).Should(gbytes.Say("requested image _invalid_:foo not found"))
	})

	Context("When a defined environment is successfully registered", func() {
		BeforeEach(func() {
			cmd, err := helper.NewWeaselCommand("enter", "busybox", "--register")
			Expect(err).NotTo(HaveOccurred())

			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Expect(session.Wait(DEFAULT_TIMEOUT)).Should(gexec.Exit(0))

			Expect(wsllib.WslIsDistributionRegistered("busybox")).Should(BeTrue())
		})

		It("should be accessible", func() {
			_, err := wsl.ExecuteSilently("busybox", "echo foo")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should not be recreated", func() {
			_, err := wsl.ExecuteSilently("busybox", "touch /marker")
			Expect(err).NotTo(HaveOccurred())

			cmd, err := helper.NewWeaselCommand("enter", "busybox", "--register")
			Expect(err).NotTo(HaveOccurred())

			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Expect(session.Wait(DEFAULT_TIMEOUT)).Should(gexec.Exit(0))

			_, err = wsl.ExecuteSilently("busybox", "test -f /marker")
			Expect(err).NotTo(HaveOccurred())
		})

		It("can be recreated", func() {
			_, err := wsl.ExecuteSilently("busybox", "touch /marker")
			Expect(err).NotTo(HaveOccurred())

			cmd, err := helper.NewWeaselCommand("enter", "busybox", "--register", "--recreate")
			Expect(err).NotTo(HaveOccurred())

			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Expect(session.Wait(DEFAULT_TIMEOUT)).Should(gexec.Exit(0))

			_, err = wsl.ExecuteSilently("busybox", "test ! -f /marker")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
