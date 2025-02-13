package lib

import (
	"testing"
	"time"

	"github.com/gofiber/fiber/v2/utils"
)

func TestObjectToSingleLevel(t *testing.T) {
	fields := []string{
		"username",
		"password",
		"contact.phone",
		"contact.address",
		"contact.city.id",
		"contact.city.name",
		"contact.city.country.id",
		"contact.city.country.name",
		"member.name",
		"member.join_date",
	}

	source := map[string]interface{}{
		"username": "johndoe",
		"password": "s3cr3t",
		"contact": map[string]interface{}{
			"phone":   123456,
			"address": "address line",
			"city": map[string]interface{}{
				"id":   1,
				"name": "city name",
				"country": map[string]interface{}{
					"id":   2,
					"name": "country name",
				},
			},
		},
		"member": map[string]interface{}{
			"name":      "john doe",
			"join_date": time.Now().UTC(),
		},
	}

	target := map[string]interface{}{}

	ObjectToSingleLevel(source, fields, &target)

	utils.AssertEqual(t, source["username"], target["username"])
	utils.AssertEqual(t, source["password"], target["password"])

	contact := source["contact"].(map[string]interface{})
	utils.AssertEqual(t, float64(contact["phone"].(int)), target["contact.phone"])
	utils.AssertEqual(t, contact["address"], target["contact.address"])

	city := contact["city"].(map[string]interface{})
	utils.AssertEqual(t, float64(city["id"].(int)), target["contact.city.id"])
	utils.AssertEqual(t, city["name"], target["contact.city.name"])

	country := city["country"].(map[string]interface{})
	utils.AssertEqual(t, float64(country["id"].(int)), target["contact.city.country.id"])
	utils.AssertEqual(t, country["name"], target["contact.city.country.name"])

	member := source["member"].(map[string]interface{})
	joinDate, _ := time.Parse(time.RFC3339Nano, target["member.join_date"].(string))
	utils.AssertEqual(t, member["name"], target["member.name"])
	utils.AssertEqual(t, member["join_date"].(time.Time), joinDate)
}
