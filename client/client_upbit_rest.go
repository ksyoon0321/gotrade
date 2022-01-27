package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/ksyoon0321/gotrade/config"
)

type UpbitParam struct {
	Apiurl string
	Method string
	Ver    string
	Cmd    string
	Param  interface{}
}

type UpbitRestClient struct {
	*http.Client
	cfg config.ConfigUpbit
}

func NewUpbitRestClient(cfg config.ConfigUpbit) *UpbitRestClient {
	return &UpbitRestClient{Client: http.DefaultClient, cfg: cfg}
}

func (c *UpbitRestClient) GetAllCoins() (interface{}, error) {
	param := UpbitParam{
		Apiurl: c.cfg.Apiurl,
		Method: "GET",
		Ver:    c.cfg.Apiver,
		Cmd:    "/market/all",
		Param: struct {
			IsDetail bool `url:"isDetail"`
		}{false}}

	list, err := c.Call(param)

	return list, err
}

func (c *UpbitRestClient) GetCoinCandles(id string, min int, cnt int) ([]map[string]interface{}, error) {
	param := UpbitParam{
		Apiurl: c.cfg.Apiurl,
		Method: "GET",
		Ver:    c.cfg.Apiver,
		Cmd:    "/candles/minutes/" + strconv.Itoa(min),
		Param: struct {
			Market string `url:"market"`
			Count  int    `url:"count"`
		}{id, cnt}}

	list, err := c.Call(param)

	return list.([]map[string]interface{}), err
}

func (c *UpbitRestClient) GetCoinCurrent(id string) ([]map[string]interface{}, error) {
	param := UpbitParam{
		Apiurl: c.cfg.Apiurl,
		Method: "GET",
		Ver:    c.cfg.Apiver,
		Cmd:    "/ticker",
		Param: struct {
			Markets string `url:"markets"`
		}{id}}

	curr, err := c.Call(param)

	return curr.([]map[string]interface{}), err
}
func (c *UpbitRestClient) Call(param interface{}) (interface{}, error) {
	var upParam = param.(UpbitParam)

	method := upParam.Method
	ver := upParam.Ver
	cmd := upParam.Cmd
	apiUrl := upParam.Apiurl

	values, err := query.Values(upParam.Param)
	if err != nil {
		panic(err)
	}
	encodedQuery := values.Encode()

	//fmt.Println("param = ", param)
	//fmt.Println("values = ", values)
	req, err := http.NewRequest(method, apiUrl+"/"+ver+cmd+"?"+encodedQuery, nil)
	if err != nil {
		return nil, err
	}

	//fmt.Println("=== REQ : ", apiUrl+"/"+ver+cmd+"?"+encodedQuery)
	return getResponse(c.Client, req)
}

func getResponse(client *http.Client, req *http.Request) (interface{}, error) {
	time.Sleep(time.Second / 3)

	//fmt.Println("================== CALL CLIENT")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(" body err 1 = ", body)
		return nil, err
	}

	var r interface{}

	err = json.Unmarshal(body, &r)
	if err != nil {
		fmt.Println(" body err 2 = ", string(body[:]))
		return nil, err
	}

	switch t := r.(type) {
	case []interface{}:
		var a []map[string]interface{}

		for _, item := range t {
			a = append(a, item.(map[string]interface{}))
		}
		r = a
	case map[string]interface{}:
		r = t
	}

	return r, nil
}
