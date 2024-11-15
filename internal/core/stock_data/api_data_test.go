package stock_data

import "testing"

func TestApiData(t *testing.T) {
	NewApiData("sh600036")
	NewApiData("hk03968")
	NewApiData("usCIHKY")
}
