package e2e_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// var _ = Describe("Build", func() {
// 	var weaselPath string

// 	BeforeEach(func() {
// 		var err error

// 		weaselPath, err = gexec.Build("github.com/codetent/weasel")
// 		Expect(err).NotTo(HaveOccurred())

// 		err = store.UnregisterDistribution("test_hub_image")
// 		Expect(err).NotTo(HaveOccurred())

// 		path, err := store.GetRegisteredDistribution("test_hub_image")
// 		Expect(err).NotTo(HaveOccurred())
// 		Expect(path).Should(BeEmpty())
// 	})

// 	AfterEach(func() {
// 		err := store.UnregisterDistribution("test_hub_image")
// 		Expect(err).NotTo(HaveOccurred())

// 		gexec.CleanupBuildArtifacts()
// 	})

// 	Describe("Build distribution", func() {
// 		It("from docker hub image", func() {
// 			cmd := exec.Command(weaselPath, "build", "--tag", "test_hub_image", "hub:busybox")
// 			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
// 			Expect(err).NotTo(HaveOccurred())

// 			session.Wait(30)
// 			Expect(session).Should(gexec.Exit(0))

// 			path, err := store.GetRegisteredDistribution("test_hub_image")
// 			Expect(err).NotTo(HaveOccurred())
// 			Expect(path).Should(BeARegularFile())
// 		})

// 		It("from unknown docker hub image", func() {
// 			cmd := exec.Command(weaselPath, "build", "--tag", "test_hub_image", "hub:_")
// 			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
// 			Expect(err).NotTo(HaveOccurred())

// 			session.Wait(30)
// 			Expect(session).Should(gexec.Exit(1))
// 			Expect(session.Out).Should(gbytes.Say(".* Error pulling image"))
// 		})

// 		It("from context with dockerfile", func() {
// 			cmd := exec.Command(weaselPath, "build", "--tag", "test_hub_image", "--file", "testdata/dummy.dockerfile", "context:testdata")
// 			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
// 			Expect(err).NotTo(HaveOccurred())

// 			session.Wait(30)
// 			Expect(session).Should(gexec.Exit(0))

// 			path, err := store.GetRegisteredDistribution("test_hub_image")
// 			Expect(err).NotTo(HaveOccurred())
// 			Expect(path).Should(BeARegularFile())
// 		})

// 		It("from non-existing context", func() {
// 			cmd := exec.Command(weaselPath, "build", "--tag", "test_hub_image", "context:_")
// 			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
// 			Expect(err).NotTo(HaveOccurred())

// 			session.Wait(30)
// 			Expect(session).Should(gexec.Exit(1))
// 			Expect(session.Out).Should(gbytes.Say(".* Error building image"))
// 		})

// 		It("from unknown specifier", func() {
// 			cmd := exec.Command(weaselPath, "build", "--tag", "test_hub_image", "_")
// 			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
// 			Expect(err).NotTo(HaveOccurred())

// 			session.Wait(30)
// 			Expect(session).Should(gexec.Exit(1))
// 			Expect(session.Out).Should(gbytes.Say(".* Unknown specifier _"))
// 		})
// 	})
// })
