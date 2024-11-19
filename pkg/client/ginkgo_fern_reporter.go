package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/guidewire-oss/fern-ginkgo-client/pkg/models"

	gt "github.com/onsi/ginkgo/v2/types"
)

var (
	aggregatedTestRun *models.TestRun
	mu                sync.Mutex // to handle concurrent access if needed
)

// InitializeTestRun sets up the test run object before starting the suites.
func (f *FernApiClient) InitializeTestRun(projectName string) {
	aggregatedTestRun = &models.TestRun{
		TestProjectName: projectName,
		StartTime:       time.Now(),
		SuiteRuns:       []models.SuiteRun{},
	}
}

func (f *FernApiClient) Report(testName string, report gt.Report) error {
	//ayo????
	mu.Lock()
	defer mu.Unlock()

	// Create a SuiteRun for the current suite
	suiteRun := models.SuiteRun{
		SuiteName: report.SuiteDescription,
		StartTime: report.StartTime,
		EndTime:   report.EndTime,
	}

	// Gather SpecRuns for this suite
	var specRuns []models.SpecRun
	for _, spec := range report.SpecReports {
		specRun := models.SpecRun{
			SpecDescription: spec.FullText(),
			Status:          spec.State.String(),
			Message:         spec.Failure.Message,
			StartTime:       spec.StartTime,
			EndTime:         spec.EndTime,
			Tags:            convertTags(report.SuiteLabels),
		}
		specRuns = append(specRuns, specRun)
	}
	suiteRun.SpecRuns = specRuns

	// Append the SuiteRun to the aggregated TestRun
	aggregatedTestRun.SuiteRuns = append(aggregatedTestRun.SuiteRuns, suiteRun)

	return nil

}

// SubmitFinalReport sends the consolidated report after all suites are complete.
func (f *FernApiClient) SubmitFinalReport() error {
	mu.Lock()
	defer mu.Unlock()

	// Finalize end time for the entire test run
	aggregatedTestRun.EndTime = time.Now()

	testJson, err := json.Marshal(aggregatedTestRun)
	if err != nil {
		return err
	}

	fmt.Printf("%s", string(testJson))

	bodyReader := bytes.NewReader(testJson)
	reportURL, err := url.JoinPath(f.baseURL, "api/testrun")
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, reportURL, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := f.httpClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making HTTP request: %s\n", err)
		return err
	}
	defer resp.Body.Close()

	return nil
}

func convertTags(specLabels []string) []models.Tag {
	// Convert Ginkgo tags to Tag struct
	var tags []models.Tag
	for _, label := range specLabels {
		tags = append(tags, models.Tag{
			Name: label, // Or however you want to define the tag
		})
	}
	return tags
}
