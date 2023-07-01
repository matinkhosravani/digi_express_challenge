package openweather

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/matinkhosravani/digi_express_challenge/app"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type Data struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"openweather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

type apiFeature struct {
	response *http.Response
	apiKey   string
}

func (a *apiFeature) resetResponse(*godog.Scenario) {
	a.response = &http.Response{}
}

func (a *apiFeature) iSendrequestTo(endpoint string) (err error) {
	client := &http.Client{}
	endpoint = strings.Replace(endpoint, "YOUR_API_KEY", a.apiKey, 1)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	a.response, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (a *apiFeature) theResponseCodeShouldBe(statusCode int) error {
	if a.response.StatusCode != statusCode {
		return fmt.Errorf("expected status code %d, but got %d", statusCode, a.response.StatusCode)
	}
	return nil
}

func (a *apiFeature) theResponseBodyShouldContainJSONProperties(table *godog.Table) error {
	body, err := io.ReadAll(a.response.Body)
	if err != nil {
		return err
	}

	var actual map[string]interface{}
	if err := json.Unmarshal(body, &actual); err != nil {
		return err
	}

	for _, row := range table.Rows {
		property := row.Cells[0]
		_, ok := getPropertyValue(actual, property.Value)
		if !ok {
			return fmt.Errorf("response body does not contain property: %s", property)
		}
	}

	return nil
}
func (a *apiFeature) theResponseBodyShouldContainXMLProperties(table *godog.Table) error {
	body, err := io.ReadAll(a.response.Body)
	if err != nil {
		return err
	}

	actual, err := unmarshalXMLToMap(body)
	if err != nil {
		return err
	}

	for _, row := range table.Rows {
		property := row.Cells[0]
		_, ok := getPropertyValue(actual, property.Value)
		if !ok {
			return fmt.Errorf("response body does not contain property: %s", property)
		}
	}

	return nil
}

func (a *apiFeature) theResponseShouldMatchJSON(body *godog.DocString) (err error) {
	respBody, err := io.ReadAll(a.response.Body)
	if err != nil {
		return err
	}
	var expected, actual interface{}

	// re-encode expected response
	if err = json.Unmarshal([]byte(body.Content), &expected); err != nil {
		return
	}

	// re-encode actual response too
	if err = json.Unmarshal(respBody, &actual); err != nil {
		return
	}

	// the matching may be adapted per different requirements.
	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected JSON does not match actual, %v vs. %v", expected, actual)
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	app.BootTestApp()
	api := &apiFeature{apiKey: app.GetEnv().OpenWeatherAPIKey}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		api.resetResponse(sc)
		return ctx, nil
	})
	ctx.Step(`^I send GET request to "([^"]*)"$`, api.iSendrequestTo)
	ctx.Step(`^The response status code should be (\d+)$`, api.theResponseCodeShouldBe)
	ctx.Step(`^the response should match json:$`, api.theResponseShouldMatchJSON)
	ctx.Step(`^the response body should contain the following JSON properties:$`, api.theResponseBodyShouldContainJSONProperties)
	ctx.Step(`^the response body should contain the following XML properties:$`, api.theResponseBodyShouldContainXMLProperties)
}

func TestAPIFeature(t *testing.T) {
	status := godog.TestSuite{
		Name:                "Weather API Feature",
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "progress",
			Paths:    []string{"features"},
			NoColors: false,
		},
	}.Run()

	if status != 0 {
		t.Errorf("non-zero status code: %d", status)
	}
}
