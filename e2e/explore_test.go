package e2e_test

import (
	"os"
	"os/exec"

	"github.com/codetent/weasel/e2e/helper"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/yuk7/wsllib-go"
)

var _ = Describe("explore", func() {
	var weaselPath string

	BeforeEach(func() {
		var err error
		weaselPath, err = gexec.Build("github.com/codetent/weasel")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		wsllib.WslUnregisterDistribution("weasel-foo")
	})

	It("should fail if the environment is undefined", func() {
		dir, err := helper.CreateConfigWorkspace("testdata/config/v1alpha1/empty.yml")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(dir)

		cmd := exec.Command(weaselPath, "explore", "foo", "--show-only")
		cmd.Dir = dir
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)

		Expect(err).NotTo(HaveOccurred())
		Expect(session.Wait(DEFAULT_TIMEOUT)).Should(gexec.Exit(1))
		Expect(session.Err).Should(gbytes.Say("undefined environment foo"))
	})

	It("should show the path of an available environment", func() {
		workspace, err := helper.CreateConfigWorkspace("testdata/config/v1alpha1/foo.yml")
		Expect(err).NotTo(HaveOccurred())

		cmd := exec.Command(weaselPath, "enter", "foo", "--register")
		cmd.Dir = workspace
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)

		Expect(err).NotTo(HaveOccurred())
		Expect(session.Wait(DEFAULT_TIMEOUT)).Should(gexec.Exit(0))
		Expect(wsllib.WslIsDistributionRegistered("weasel-foo")).Should(BeTrue())

		cmd = exec.Command(weaselPath, "explore", "foo", "--show-only")
		cmd.Dir = workspace
		session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)

		Expect(err).NotTo(HaveOccurred())
		Expect(session.Wait(DEFAULT_TIMEOUT)).Should(gexec.Exit(0))
		Expect(session.Out).Should(gbytes.Say("\\\\wsl\\$\\\\weasel-foo\\\\root"))
	})
})
