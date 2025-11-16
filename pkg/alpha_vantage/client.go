package alphavantage

import (
	"encoding/json"
	"fmt"
	"investbot/pkg/domain"
	"investbot/pkg/errors"
	"net/http"
	"net/url"
)

type AlphaVantageClient struct {
	apiKey string
}

const alphaVantageBaseURL = "https://www.alphavantage.co/query"

func NewAlphaVantageClient(apiKey string) (*AlphaVantageClient, error) {
	return &AlphaVantageClient{apiKey: apiKey}, nil
}

func (c *AlphaVantageClient) GetRealGdpTimeSeries(interval domain.EconomicIndicatorInterval) (domain.EconomicIndicatorTimeSeries, error) {
	// Map domain interval to API interval
	apiInterval := "annual"
	if interval == domain.QuarterlyEconomicIndicatorInterval {
		apiInterval = "quarterly"
	}

	// Build URL with query parameters
	requestUrl, err := url.Parse(alphaVantageBaseURL)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to parse base URL: %v", err),
		}
	}

	q := requestUrl.Query()
	q.Set("function", string(RealGDP))
	q.Set("interval", apiInterval)
	q.Set("apikey", c.apiKey)
	requestUrl.RawQuery = q.Encode()

	// Create HTTP request
	req, err := http.NewRequest("GET", requestUrl.String(), nil)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to create HTTP request: %v", err),
		}
	}

	// Send the request
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to send HTTP request: %v", err),
		}
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: resp.StatusCode,
			Message:    resp.Status,
		}
	}

	// Parse JSON response
	var apiResponse EconomicIndicatorTimeSeriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.JSONMarshalError{
			Message: "failed to decode JSON response",
			Err:     err,
		}
	}

	// Map API response to domain model
	domainInterval := domain.AnnualEconomicIndicatorInterval
	if apiResponse.Interval == "quarterly" {
		domainInterval = domain.QuarterlyEconomicIndicatorInterval
	}

	// Map unit - API returns "billions of dollars" for Real GDP
	domainUnit := domain.BillionsOfDollarsEconomicIndicatorUnit

	// Map data entries
	data := make([]domain.EconomicIndicatortypeTimeSeriesEntry, len(apiResponse.Data))
	for i, entry := range apiResponse.Data {
		data[i] = domain.EconomicIndicatortypeTimeSeriesEntry{
			Date:  entry.Date,
			Value: entry.Value,
		}
	}

	return domain.EconomicIndicatorTimeSeries{
		Name:     domain.RealGDP,
		Interval: domainInterval,
		Unit:     domainUnit,
		Data:     data,
	}, nil
}

// GetTreasuryYieldTimeSeries returns the monthly treasury yield of the given maturity
func (c *AlphaVantageClient) GetTreasuryYieldTimeSeries(
	maturity domain.TreasuryYieldMaturity,
) (domain.EconomicIndicatorTimeSeries, error) {
	apiInterval := "monthly"

	// Build URL with query parameters
	requestUrl, err := url.Parse(alphaVantageBaseURL)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to parse base URL: %v", err),
		}
	}

	var maturityParam string
	switch maturity {
	case domain.ThreeMonthTreasuryYieldMaturity:
		maturityParam = "3month"
	case domain.TwoYearTreasuryYieldMaturity:
		maturityParam = "2year"
	case domain.FiveYearTreasuryYieldMaturity:
		maturityParam = "5year"
	case domain.TenYearTreasuryYieldMaturity:
		maturityParam = "10year"
	case domain.ThirtyYearTreasuryYieldMaturity:
		maturityParam = "30year"
	default:
		maturityParam = "10year"
	}

	q := requestUrl.Query()
	q.Set("function", string(TreasuryYield))
	q.Set("interval", apiInterval)
	q.Set("apikey", c.apiKey)
	q.Set("maturity", maturityParam)
	requestUrl.RawQuery = q.Encode()

	// Create HTTP request
	req, err := http.NewRequest("GET", requestUrl.String(), nil)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to create HTTP request: %v", err),
		}
	}

	// Send the request
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to send HTTP request: %v", err),
		}
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: resp.StatusCode,
			Message:    resp.Status,
		}
	}

	// Parse JSON response
	var apiResponse EconomicIndicatorTimeSeriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.JSONMarshalError{
			Message: "failed to decode JSON response",
			Err:     err,
		}
	}

	domainInterval := domain.MonthlyEconomicIndicatorInterval
	domainUnit := domain.PercentEconomicIndicatorUnit

	// Map data entries
	data := make([]domain.EconomicIndicatortypeTimeSeriesEntry, len(apiResponse.Data))
	for i, entry := range apiResponse.Data {
		data[i] = domain.EconomicIndicatortypeTimeSeriesEntry{
			Date:  entry.Date,
			Value: entry.Value,
		}
	}

	return domain.EconomicIndicatorTimeSeries{
		Name:     domain.TreasuryYield,
		Interval: domainInterval,
		Unit:     domainUnit,
		Data:     data,
	}, nil
}

// GetInterestRatesTimeSeries returns the monthly interest rate time series
func (c *AlphaVantageClient) GetInterestRatesTimeSeries() (domain.EconomicIndicatorTimeSeries, error) {
	apiInterval := "monthly"

	// Build URL with query parameters
	requestUrl, err := url.Parse(alphaVantageBaseURL)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to parse base URL: %v", err),
		}
	}

	q := requestUrl.Query()
	q.Set("function", string(FederalFundsRate))
	q.Set("interval", apiInterval)
	q.Set("apikey", c.apiKey)
	requestUrl.RawQuery = q.Encode()

	// Create HTTP request
	req, err := http.NewRequest("GET", requestUrl.String(), nil)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to create HTTP request: %v", err),
		}
	}

	// Send the request
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to send HTTP request: %v", err),
		}
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: resp.StatusCode,
			Message:    resp.Status,
		}
	}

	// Parse JSON response
	var apiResponse EconomicIndicatorTimeSeriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.JSONMarshalError{
			Message: "failed to decode JSON response",
			Err:     err,
		}
	}

	domainInterval := domain.MonthlyEconomicIndicatorInterval
	domainUnit := domain.PercentEconomicIndicatorUnit

	// Map data entries
	data := make([]domain.EconomicIndicatortypeTimeSeriesEntry, len(apiResponse.Data))
	for i, entry := range apiResponse.Data {
		data[i] = domain.EconomicIndicatortypeTimeSeriesEntry{
			Date:  entry.Date,
			Value: entry.Value,
		}
	}

	return domain.EconomicIndicatorTimeSeries{
		Name:     domain.InterestRate,
		Interval: domainInterval,
		Unit:     domainUnit,
		Data:     data,
	}, nil
}

// GetInflationTimeSeries returns the annual inflation time series
func (c *AlphaVantageClient) GetInflationTimeSeries() (domain.EconomicIndicatorTimeSeries, error) {
	// Build URL with query parameters
	requestUrl, err := url.Parse(alphaVantageBaseURL)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to parse base URL: %v", err),
		}
	}

	q := requestUrl.Query()
	q.Set("function", string(Inflation))
	q.Set("apikey", c.apiKey)
	requestUrl.RawQuery = q.Encode()

	// Create HTTP request
	req, err := http.NewRequest("GET", requestUrl.String(), nil)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to create HTTP request: %v", err),
		}
	}

	// Send the request
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to send HTTP request: %v", err),
		}
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: resp.StatusCode,
			Message:    resp.Status,
		}
	}

	// Parse JSON response
	var apiResponse EconomicIndicatorTimeSeriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.JSONMarshalError{
			Message: "failed to decode JSON response",
			Err:     err,
		}
	}

	domainInterval := domain.AnnualEconomicIndicatorInterval
	domainUnit := domain.PercentEconomicIndicatorUnit

	// Map data entries
	data := make([]domain.EconomicIndicatortypeTimeSeriesEntry, len(apiResponse.Data))
	for i, entry := range apiResponse.Data {
		data[i] = domain.EconomicIndicatortypeTimeSeriesEntry{
			Date:  entry.Date,
			Value: entry.Value,
		}
	}

	return domain.EconomicIndicatorTimeSeries{
		Name:     domain.Inflation,
		Interval: domainInterval,
		Unit:     domainUnit,
		Data:     data,
	}, nil
}

// GetUnemploymentRateTimeSeries returns the monthly unemployment rate time series
func (c *AlphaVantageClient) GetUnemploymentRateTimeSeries() (domain.EconomicIndicatorTimeSeries, error) {
	// Build URL with query parameters
	requestUrl, err := url.Parse(alphaVantageBaseURL)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to parse base URL: %v", err),
		}
	}

	q := requestUrl.Query()
	q.Set("function", string(UnemploymentRate))
	q.Set("apikey", c.apiKey)
	requestUrl.RawQuery = q.Encode()

	// Create HTTP request
	req, err := http.NewRequest("GET", requestUrl.String(), nil)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to create HTTP request: %v", err),
		}
	}

	// Send the request
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to send HTTP request: %v", err),
		}
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return domain.EconomicIndicatorTimeSeries{}, &errors.HTTPError{
			StatusCode: resp.StatusCode,
			Message:    resp.Status,
		}
	}

	// Parse JSON response
	var apiResponse EconomicIndicatorTimeSeriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return domain.EconomicIndicatorTimeSeries{}, &errors.JSONMarshalError{
			Message: "failed to decode JSON response",
			Err:     err,
		}
	}

	domainInterval := domain.MonthlyEconomicIndicatorInterval
	domainUnit := domain.PercentEconomicIndicatorUnit

	// Map data entries
	data := make([]domain.EconomicIndicatortypeTimeSeriesEntry, len(apiResponse.Data))
	for i, entry := range apiResponse.Data {
		data[i] = domain.EconomicIndicatortypeTimeSeriesEntry{
			Date:  entry.Date,
			Value: entry.Value,
		}
	}

	return domain.EconomicIndicatorTimeSeries{
		Name:     domain.UnemploymentRate,
		Interval: domainInterval,
		Unit:     domainUnit,
		Data:     data,
	}, nil
}
