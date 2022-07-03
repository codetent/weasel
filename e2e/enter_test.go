package e2e_test

import (
	"os"
	"os/exec"

	"github.com/codetent/weasel/e2e/helper"
	"github.com/codetent/weasel/pkg/weasel/wsl"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/yuk7/wsllib-go"
)

var _ = Describe("enter", Ordered, func() {
	var weaselPath string

	BeforeEach(func() {
		var err error
		weaselPath, err = gexec.Build("github.com/codetent/weasel")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		wsllib.WslUnregisterDistribution("weasel-foo")
	})

	It("should fail if there is no configuration file", func() {
		dir, err := helper.CreateEmptyWorkspace()
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(dir)

		cmd := exec.Command(weaselPath, "enter", "foo")
		cmd.Dir = dir
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)

		Expect(err).NotTo(HaveOccurred())
		Expect(session.Wait(DEFAULT_TIMEOUT)).Should(gexec.Exit(1))
		Expect(session.Err).Should(gbytes.Say("configuration not found"))
	})

	It("should fail if the environment is undefined", func() {
		dir, err := helper.CreateConfigWorkspace("testdata/config/v1alpha1/empty.yml")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(dir)

		cmd := exec.Command(weaselPath, "enter", "foo")
		cmd.Dir = dir
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)

		Expect(err).NotTo(HaveOccurred())
		Expect(session.Wait(DEFAULT_TIMEOUT)).Should(gexec.Exit(1))
		Expect(session.Err).Should(gbytes.Say("undefined environment foo"))
	})

	It("should fail registering an environment with an invalid image", func() {
		dir, err := helper.CreateConfigWorkspace("testdata/config/v1alpha1/invalid.yml")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(dir)

		cmd := exec.Command(weaselPath, "enter", "foo", "--register")
		cmd.Dir = dir
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)

		Expect(err).NotTo(HaveOccurred())
		Expect(session.Wait(DEFAULT_TIMEOUT)).Should(gexec.Exit(1))
		Expect(session.Err).Should(gbytes.Say("requested image _invalid_:foo not found"))
	})

	Context("When a defined environment is successfully registered", func() {
		var workspace string

		BeforeEach(func() {
			var err error

			workspace, err = helper.CreateConfigWorkspace("testdata/config/v1alpha1/foo.yml")
			Expect(err).NotTo(HaveOccurred())

			cmd := exec.Command(weaselPath, "enter", "foo", "--register")
			cmd.Dir = workspace
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)

			Expect(err).NotTo(HaveOccurred())
			Expect(session.Wait(DEFAULT_TIMEOUT)).Should(gexec.Exit(0))
			Expect(wsllib.WslIsDistributionRegistered("weasel-foo")).Should(BeTrue())
		})

		AfterEach(func() {
			os.RemoveAll(workspace)
		})

		It("should be accessible", func() {
			_, err := wsl.ExecuteSilently("weasel-foo", "echo foo")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should not be recreated", func() {
			_, err := wsl.ExecuteSilently("weasel-foo", "touch /marker")
			Expect(err).NotTo(HaveOccurred())

			cmd := exec.Command(weaselPath, "enter", "foo", "--register")
			cmd.Dir = workspace
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)

			Expect(err).NotTo(HaveOccurred())
			Expect(session.Wait(DEFAULT_TIMEOUT)).Should(gexec.Exit(0))

			_, err = wsl.ExecuteSilently("weasel-foo", "test -f /marker")
			Expect(err).NotTo(HaveOccurred())
		})

		It("can be recreated", func() {
			_, err := wsl.ExecuteSilently("weasel-foo", "touch /marker")
			Expect(err).NotTo(HaveOccurred())

			cmd := exec.Command(weaselPath, "enter", "foo", "--register", "--recreate")
			cmd.Dir = workspace
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)

			Expect(err).NotTo(HaveOccurred())
			Expect(session.Wait(DEFAULT_TIMEOUT)).Should(gexec.Exit(0))

			_, err = wsl.ExecuteSilently("weasel-foo", "test ! -f /marker")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
