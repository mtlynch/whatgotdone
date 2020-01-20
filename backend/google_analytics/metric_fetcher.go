package google_analtyics

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"google.golang.org/api/option"

	ga "google.golang.org/api/analyticsreporting/v4"
)

type (
	// MetricFetcher retrieves metrics from Google Analytics.
	MetricFetcher interface {
		PageViewsByPath(startDate, endDate string) ([]PageViewCount, error)
	}

	// PageViewCount represents the number of pageviews for a given URL path.
	PageViewCount struct {
		Path  string
		Views int
	}

	// defaultMetricFetcher implements MetricFetcher using a real Google Analytics
	// backend.
	defaultMetricFetcher struct {
		svc    *ga.Service
		viewID string
	}
)

// New creates a new MetricFetcher instance.
func New() (mf MetricFetcher, err error) {
	viewID := os.Getenv("GOOGLE_ANALYTICS_VIEW_ID")
	if viewID == "" {
		log.Printf("GOOGLE_ANALYTICS_VIEW_ID is not set, skipping Google Analytics updates")
		return mf, errors.New("Can't create MetricFetcher without Google Analytics View ID")
	}

	const keyFilePath = "google-analytics-service-account.json"
	svc, err := ga.NewService(context.Background(), option.WithCredentialsFile(keyFilePath))
	if err != nil {
		return defaultMetricFetcher{nil, ""}, err
	}
	return defaultMetricFetcher{svc, viewID}, nil
}

// PageViewsByPath retrieves the total pageviews for each URL path over a given
// date range.
func (r defaultMetricFetcher) PageViewsByPath(startDate, endDate string) ([]PageViewCount, error) {
	res, err := getReport(r.svc, r.viewID, startDate, endDate)
	if err != nil {
		return []PageViewCount{}, err
	}

	return extractPageViews(res)
}

func getReport(svc *ga.Service, viewID string, startDate string, endDate string) (*ga.GetReportsResponse, error) {
	req := &ga.GetReportsRequest{
		ReportRequests: []*ga.ReportRequest{
			{
				ViewId: viewID,
				DateRanges: []*ga.DateRange{
					{StartDate: startDate, EndDate: endDate},
				},
				Metrics: []*ga.Metric{
					{Expression: "ga:pageviews"},
				},
				Dimensions: []*ga.Dimension{
					{Name: "ga:pagePath"},
				},
			},
		},
	}

	res, err := svc.Reports.BatchGet(req).Do()
	if err != nil {
		return nil, err
	}

	if res.HTTPStatusCode != 200 {
		return nil, fmt.Errorf("request to Google Analytics failed with %d", res.HTTPStatusCode)
	}

	return res, nil
}

func extractPageViews(res *ga.GetReportsResponse) ([]PageViewCount, error) {
	if len(res.Reports) != 1 {
		return []PageViewCount{}, fmt.Errorf("unexpected report count. wanted %d, got %d", 1, len(res.Reports))
	}
	report := res.Reports[0]
	rows := report.Data.Rows

	if rows == nil {
		log.Println("no data found in report")
		return []PageViewCount{}, nil
	}

	viewCounts := []PageViewCount{}
	for _, row := range rows {
		pagePath := row.Dimensions[0]
		pageViews, err := strconv.Atoi(row.Metrics[0].Values[0])
		if err != nil {
			panic(err)
		}
		viewCounts = append(viewCounts, PageViewCount{pagePath, pageViews})
	}
	return viewCounts, nil
}
