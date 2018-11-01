package generic_test

import (
	"fmt"
	"os"
	"strconv"
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

const defaultHPATimeout = "210s"

var (
	hpaDeployment = test_helpers.PathFromRoot("specs/hpa-php-apache.yml")
)

var _ = Describe("Horizontal Pod Autoscaling", func() {
	BeforeEach(func() {
		createHPADeployment()
		autoscaleHPA()
	})

	AfterEach(func() {
		deleteHPADeployment()
	})

	It("scales the pods accordingly", func() {
		HPATimeout := os.Getenv("HPA_TIMEOUT")
		if HPATimeout == "" {
			HPATimeout = defaultHPATimeout
		}

		By("creating more pods when the CPU load increases")

		increaseCPULoad()
		Eventually(func() int {
			session := runner.RunKubectlCommand("get", "hpa/php-apache", "-o", "jsonpath={.status.currentReplicas}")
			Eventually(session, "10s").Should(gexec.Exit(0))
			replicas, _ := strconv.Atoi(string(session.Out.Contents()))
			return replicas
		}, HPATimeout).Should(BeNumerically(">", 1))

		By("decreasing the number of pods when the CPU load decreases")

		session := runner.RunKubectlCommand("delete", "deployment/load-generator")
		Eventually(session, "10s").Should(gexec.Exit(0))

		Eventually(func() int {
			session := runner.RunKubectlCommand("get", "hpa/php-apache", "-o", "jsonpath={.status.currentReplicas}")
			Eventually(session, "10s").Should(gexec.Exit(0))
			replicas, _ := strconv.Atoi(string(session.Out.Contents()))
			return replicas
		}, HPATimeout).Should(BeNumerically("==", 1))
	})
})

func deleteHPADeployment() {
	runner.RunKubectlCommand("delete", "-f", hpaDeployment)
}

func createHPADeployment() {
	session := runner.RunKubectlCommand("create", "-f", hpaDeployment)
	Eventually(session, "10s").Should(gexec.Exit(0))

	Eventually(func() string {
		return runner.GetPodStatusBySelector(runner.Namespace(), "run=php-apache")
	}, "120s").Should(Equal("Running"))
}

func autoscaleHPA() {
	session := runner.RunKubectlCommand("autoscale", "deployment/php-apache", "--cpu-percent=25", "--min=1", "--max=2")
	Eventually(session, "10s").Should(gexec.Exit(0))
}

func increaseCPULoad() {
	remoteCommand := fmt.Sprintf("while true; do wget -q -O- http://php-apache.%s.svc.cluster.local; done", runner.Namespace())

	session := runner.RunKubectlCommand("run", "-i", "--tty", "load-generator", "--image=busybox", "--", "/bin/sh", "-c", remoteCommand)
	Eventually(session, "10s").Should(gexec.Exit(0))

	Eventually(func() string {
		return runner.GetPodStatusBySelector(runner.Namespace(), "run=load-generator")
	}, "120s").Should(Equal("Running"))
}
