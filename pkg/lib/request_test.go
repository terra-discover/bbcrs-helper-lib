package lib

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func TestGetXUserID(t *testing.T) {
	id := uuid.New().String()
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"id": GetXUserID(c),
		})
	})

	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("X-User-ID", id)
	response, err := app.Test(request)

	utils.AssertEqual(t, nil, err, "Getting response data")
	defer response.Body.Close()
	bte, err := ioutil.ReadAll(response.Body)
	utils.AssertEqual(t, nil, err, "Reading response data")

	var result map[string]interface{}
	err = json.Unmarshal(bte, &result)
	utils.AssertEqual(t, nil, err, "Parsing response data")
	utils.AssertEqual(t, id, result["id"], "same id")

	request2 := httptest.NewRequest("GET", "/", nil)
	request2.Header.Add("Accept", "application/json")
	response2, err := app.Test(request2)

	utils.AssertEqual(t, nil, err, "Getting response data")
	defer response2.Body.Close()
	bte2, err := ioutil.ReadAll(response2.Body)
	utils.AssertEqual(t, nil, err, "Reading response data")

	var result2 map[string]interface{}
	err = json.Unmarshal(bte2, &result2)
	utils.AssertEqual(t, nil, err, "Parsing response data")
	utils.AssertEqual(t, nil, result2["id"], "null id")
}

func TestGetXAgentID(t *testing.T) {
	id := uuid.New().String()
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"id": GetXAgentID(c),
		})
	})

	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("X-Agent-ID", id)
	response, err := app.Test(request)

	utils.AssertEqual(t, nil, err, "Getting response data")
	defer response.Body.Close()
	bte, err := ioutil.ReadAll(response.Body)
	utils.AssertEqual(t, nil, err, "Reading response data")

	var result map[string]interface{}
	err = json.Unmarshal(bte, &result)
	utils.AssertEqual(t, nil, err, "Parsing response data")
	utils.AssertEqual(t, id, result["id"], "same id")

	request2 := httptest.NewRequest("GET", "/", nil)
	request2.Header.Add("Accept", "application/json")
	response2, err := app.Test(request2)

	utils.AssertEqual(t, nil, err, "Getting response data")
	defer response2.Body.Close()
	bte2, err := ioutil.ReadAll(response2.Body)
	utils.AssertEqual(t, nil, err, "Reading response data")

	var result2 map[string]interface{}
	err = json.Unmarshal(bte2, &result2)
	utils.AssertEqual(t, nil, err, "Parsing response data")
	utils.AssertEqual(t, nil, result2["id"], "null id")

	viper.Set("AGENT_ID", uuid.New().String())
	request3 := httptest.NewRequest("GET", "/", nil)
	request.Header.Add("Accept", "application/json")
	response3, err := app.Test(request3)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 200, response3.StatusCode)
	viper.Set("AGENT_ID", "")

}

func TestGetXCorporateID(t *testing.T) {
	id := uuid.New().String()
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"id": GetXCorporateID(c),
		})
	})

	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("X-Corporate-ID", id)
	response, err := app.Test(request)

	utils.AssertEqual(t, nil, err, "Getting response data")
	defer response.Body.Close()
	bte, err := ioutil.ReadAll(response.Body)
	utils.AssertEqual(t, nil, err, "Reading response data")

	var result map[string]interface{}
	err = json.Unmarshal(bte, &result)
	utils.AssertEqual(t, nil, err, "Parsing response data")
	utils.AssertEqual(t, id, result["id"], "same id")

	request2 := httptest.NewRequest("GET", "/", nil)
	request2.Header.Add("Accept", "application/json")
	response2, err := app.Test(request2)

	utils.AssertEqual(t, nil, err, "Getting response data")
	defer response2.Body.Close()
	bte2, err := ioutil.ReadAll(response2.Body)
	utils.AssertEqual(t, nil, err, "Reading response data")

	var result2 map[string]interface{}
	err = json.Unmarshal(bte2, &result2)
	utils.AssertEqual(t, nil, err, "Parsing response data")
	utils.AssertEqual(t, nil, result2["id"], "null id")
}

func TestGetLanguage(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"lang": GetLanguage(c),
		})
	})

	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Add("Accept-Language", "fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5")
	response, err := app.Test(request, 500)
	utils.AssertEqual(t, nil, err, "Sending request")
	utils.AssertEqual(t, 200, response.StatusCode, "Getting status code")

	defer response.Body.Close()
	bte, err := ioutil.ReadAll(response.Body)
	utils.AssertEqual(t, nil, err, "Reading response data")

	var result map[string]interface{}
	err = json.Unmarshal(bte, &result)
	utils.AssertEqual(t, nil, err, "Parsing response data")
	utils.AssertEqual(t, "fr", result["lang"], "same language")

	request2 := httptest.NewRequest("GET", "/", nil)
	request2.Header.Add("Accept-Language", "a")
	response2, err := app.Test(request2, 500)
	utils.AssertEqual(t, nil, err, "Sending request")
	utils.AssertEqual(t, 200, response2.StatusCode, "Getting status code")

	defer response2.Body.Close()
	bte2, err := ioutil.ReadAll(response2.Body)
	utils.AssertEqual(t, nil, err, "Reading response data")

	var result2 map[string]interface{}
	err = json.Unmarshal(bte2, &result2)
	utils.AssertEqual(t, nil, err, "Parsing response data")
	utils.AssertEqual(t, "en", result2["lang"], "same language")
}

func TestBodyParser(t *testing.T) {
	type sample struct {
		Name *string `json:"name" validate:"required,specialcharacter,gte=8"`
	}

	app := fiber.New()
	app.Post("/validate", func(c *fiber.Ctx) error {
		data := new(sample)
		if err := BodyParser(c, data); nil != err {
			return ErrorBadRequest(c, err)
		}

		return OK(c)
	})

	res, body, err := PostTest(app, "/validate", nil, "")
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, false, nil == body)
	utils.AssertEqual(t, 400, res.StatusCode)

	res, body, err = PostTest(app, "/validate", nil, `{"name":"john"}`)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 400, res.StatusCode)

	res, body, err = PostTest(app, "/validate", nil, `{"name":"john doe"}`)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 200, res.StatusCode)
}

func TestQueryParser(t *testing.T) {
	type sample struct {
		Age *int `query:"age" validate:"required"`
	}

	app := fiber.New()
	app.Get("/validate", func(c *fiber.Ctx) error {
		data := new(sample)
		if err := QueryParser(c, data); nil != err {
			return ErrorBadRequest(c, err)
		}

		return OK(c)
	})

	res, body, err := GetTest(app, "/validate", nil)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, false, nil == body)
	utils.AssertEqual(t, 400, res.StatusCode)

	res, body, err = GetTest(app, "/validate", nil)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 400, res.StatusCode)

	res, body, err = GetTest(app, "/validate?age=twenty", nil)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 400, res.StatusCode)

	res, body, err = GetTest(app, "/validate?age=12", nil)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 200, res.StatusCode)
}
