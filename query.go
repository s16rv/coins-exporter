package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/s16rv/coins-exporter/types"
)

func SendQueryCoinsDetail(baseApi, id string) (types.ReturnData, error) {
	u := baseApi + "/coins/" + id
	var d types.ReturnData

	client := &http.Client{}

	req, _ := http.NewRequest("GET", u, nil)
	resp, err := client.Do(req)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return d, err
	}
	err = json.Unmarshal(body, &d)
	if err != nil {
		return d, err
	}
	defer resp.Body.Close()

	return d, nil
}
