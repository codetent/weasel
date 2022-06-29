package e2e_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Run", func() {
	var weaselPath string

	BeforeEach(func() {
		var err error

		weaselPath, err = gexec.Build("github.com/codetent/weasel")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		// err := store.UnregisterDistribution("test_hub_image")
		// Expect(err).NotTo(HaveOccurred())

		gexec.CleanupBuildArtifacts()
	})

	Describe("Run instance", func() {
		It("from existing distribution", func() {
			cmd := exec.Command(weaselPath, "build", "--tag", "test_hub_image", "hub:busybox")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			session.Wait(30)
			Expect(session).Should(gexec.Exit(0))

			cmd = exec.Command(weaselPath, "run", "test_hub_image", "echo", "works!")
			session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			session.Wait(30)
			Expect(session).Should(gexec.Exit(0))
			Expect(session.Out).Should(gbytes.Say("works!"))
		})

		It("from non-existing distribution", func() {
			cmd := exec.Command(weaselPath, "run", "_", "echo", "works!")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			session.Wait(30)
			Expect(session).Should(gexec.Exit(1))
			Expect(session.Out).Should(gbytes.Say(".* Distribution with id _ not found"))
		})
	})
})
