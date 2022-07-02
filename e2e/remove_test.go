package e2e_test

import (
	"os/exec"

	"github.com/codetent/weasel/e2e/helper"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/yuk7/wsllib-go"
)

var _ = Describe("remove", func() {
	var weaselPath string

	BeforeEach(func() {
		var err error
		weaselPath, err = gexec.Build("github.com/codetent/weasel")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		wsllib.WslUnregisterDistribution("weasel-foo")
	})

	It("should fail for unavailable environments", func() {
		workspace, err := helper.CreateConfigWorkspace("testdata/config/v1alpha1/empty.yml")
		Expect(err).NotTo(HaveOccurred())

		cmd := exec.Command(weaselPath, "remove", "foo")
		cmd.Dir = workspace
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)

		Expect(err).NotTo(HaveOccurred())
		Expect(session.Wait(DEFAULT_TIMEOUT)).Should(gexec.Exit(1))
		Expect(session.Err).Should(gbytes.Say("undefined environment foo"))
	})

	It("should fail for unavailable environments", func() {
		workspace, err := helper.CreateConfigWorkspace("testdata/config/v1alpha1/foo.yml")
		Expect(err).NotTo(HaveOccurred())

		cmd := exec.Command(weaselPath, "remove", "foo")
		cmd.Dir = workspace
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)

		Expect(err).NotTo(HaveOccurred())
		Expect(session.Wait(DEFAULT_TIMEOUT)).Should(gexec.Exit(1))
		Expect(session.Err).Should(gbytes.Say("environment foo not available"))
	})

	It("should remove a registered environment", func() {
		workspace, err := helper.CreateConfigWorkspace("testdata/config/v1alpha1/foo.yml")
		Expect(err).NotTo(HaveOccurred())

		cmd := exec.Command(weaselPath, "enter", "foo", "--register")
		cmd.Dir = workspace
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)

		Expect(err).NotTo(HaveOccurred())
		Expect(session.Wait(DEFAULT_TIMEOUT)).Should(gexec.Exit(0))
		Expect(wsllib.WslIsDistributionRegistered("weasel-foo")).Should(BeTrue())

		cmd = exec.Command(weaselPath, "remove", "foo")
		cmd.Dir = workspace
		session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)

		Expect(err).NotTo(HaveOccurred())
		Expect(session.Wait(DEFAULT_TIMEOUT)).Should(gexec.Exit(0))
		Expect(wsllib.WslIsDistributionRegistered("weasel-foo")).Should(BeFalse())
	})
})
