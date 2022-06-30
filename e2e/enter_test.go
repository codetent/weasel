package e2e_test

import (
	"os"
	"os/exec"

	"github.com/codetent/weasel/e2e/helper"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("enter", func() {
	var weaselPath string

	BeforeEach(func() {
		var err error
		weaselPath, err = gexec.Build("github.com/codetent/weasel")
		Expect(err).NotTo(HaveOccurred())
	})

	It("should fail if there is no configuration file", func() {
		dir, err := helper.CreateEmptyWorkspace()
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(dir)

		cmd := exec.Command(weaselPath, "enter", "foo")
		cmd.Dir = dir
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)

		Expect(err).NotTo(HaveOccurred())
		Expect(session.Wait()).Should(gexec.Exit(1))
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
		Expect(session.Wait()).Should(gexec.Exit(1))
		Expect(session.Err).Should(gbytes.Say("undefined environment foo"))
	})

	It("should successfully register a defined environment", func() {
		dir, err := helper.CreateConfigWorkspace("testdata/config/v1alpha1/foo.yml")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(dir)

		cmd := exec.Command(weaselPath, "enter", "foo", "--register")
		cmd.Dir = dir
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)

		Expect(err).NotTo(HaveOccurred())
		Expect(session.Wait()).Should(gexec.Exit(0))
	})
})
