package tests_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
  fern "github.com/guidewire/fern-ginkgo-client/pkg/client"

	. "github.com/guidewire/fern-ginkgo-client/tests"
)

var f *fern.FernApiClient

// Initialize the Fern client and the aggregated test run at the start
var _ = BeforeSuite(func() {
	f = fern.New("Example Test",
		fern.WithBaseURL("http://localhost:8080/"),
	)
	f.InitializeTestRun("Example Project") // Initializes the aggregated report
})


var _ = Describe("Adder", func() {

		Describe("Add", func() {

		It("adds two numbers", func() {
			sum := Add(2, 3)
			Expect(sum).To(Equal(5))
		})
	})

})

// Report each suite's results as they complete
var _ = ReportAfterSuite("", func(report Report) {
	err := f.Report("example test", report)
	Expect(err).To(BeNil(), "Unable to report suite results")
})

// After all tests are complete, submit the final aggregated report
var _ = AfterSuite(func() {
	err := f.SubmitFinalReport() // Sends the final aggregated report
	Expect(err).To(BeNil(), "Unable to submit final report")
})