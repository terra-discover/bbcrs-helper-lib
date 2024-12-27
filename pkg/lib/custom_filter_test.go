package lib

import (
	"fmt"
	"testing"
)

func TestCustomFilters(t *testing.T) {
	// filter single case
	qf := `["id","=",1]`
	queryFilters, whereFilters, _, _ := CustomFilters(qf, "", "")
	fmt.Println("queryFilters :", queryFilters)
	fmt.Println("whereFilters :", whereFilters)
	fmt.Println("filter single case ==================================")

	// filter operator LIKE case
	qf = `["id","LIKE",1]`
	queryFilters, whereFilters, _, _ = CustomFilters(qf, "", "")
	fmt.Println("queryFilters :", queryFilters)
	fmt.Println("whereFilters :", whereFilters)
	fmt.Println("==================================")

	// filter operator NOT LIKE case
	qf = `["name","NOT LIKE","product glass"]`
	queryFilters, whereFilters, _, _ = CustomFilters(qf, "", "")
	fmt.Println("queryFilters :", queryFilters)
	fmt.Println("whereFilters :", whereFilters)
	fmt.Println("filter operator NOT LIKE case ==================================")

	// filter operator IS case
	qf = `["id","IS",null]`
	queryFilters, whereFilters, _, _ = CustomFilters(qf, "", "")
	fmt.Println("queryFilters :", queryFilters)
	fmt.Println("whereFilters :", whereFilters)
	fmt.Println("filter operator IS case ==================================")

	// filter multiple case
	qf = `[["id","=",1],["AND"],["status","=",true],["OR"],["amount","=",20.5]]`
	queryFilters, whereFilters, _, _ = CustomFilters(qf, "", "")
	fmt.Println("queryFilters :", queryFilters)
	fmt.Println("whereFilters :", whereFilters)
	fmt.Println("filter multiple case ==================================")

	// filter operator IN case
	qf = `["payment_status","IN",["paid", "due", "partial"]]`
	queryFilters, whereFilters, _, _ = CustomFilters(qf, "", "")
	fmt.Println("queryFilters :", queryFilters)
	fmt.Println("whereFilters :", whereFilters)
	fmt.Println("filter operator IN case ==================================")

	// filter operator BETWEEN case
	qf = `["DATE(transaction_date)","BETWEEN",["date_start","date_end"]]` // date (YYYY-MM-DD hh:mm:ss)
	queryFilters, whereFilters, _, _ = CustomFilters(qf, "", "")
	fmt.Println("queryFilters :", queryFilters)
	fmt.Println("whereFilters :", whereFilters)
	fmt.Println("filter operator BETWEEN case ==================================")

	// filter nested object field case
	qf = `[["person.id","=","1"],["OR"],["person.status","=","active"]]`
	queryFilters, whereFilters, _, _ = CustomFilters(qf, "", "")
	fmt.Println("queryFilters :", queryFilters)
	fmt.Println("whereFilters :", whereFilters)
	fmt.Println("filter nested object field case ==================================")

	// search column like case
	q := `value`
	col := `["trx_id","id"]`
	_, _, querySearch, columnFilter := CustomFilters("", q, col)
	fmt.Println("querySearch :", querySearch)
	fmt.Println("columnFilter :", columnFilter)
	fmt.Println("search column like case ==================================")

	// filter negative case
	CreateFilter("unexpected json format")
}
