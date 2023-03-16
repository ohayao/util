package http

import (
	// origin http library
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	oh "net/http"
	"net/url"
	"strings"
)

type authorization struct {
	scheme string
	value  string
}

type Response struct {
	StatusCode    int
	Body          []byte
	Proto         string
	Header        map[string][]string
	ContentLength int64
}
type param struct {
	key   string
	value interface{}
}

type Client struct {
	cli           *oh.Client
	url           string
	user_agent    string
	headers       map[string]string
	cookies       []*oh.Cookie
	query         []param
	form          []param
	authorization *authorization
}

func New(url string) *Client {
	cli := &Client{
		cli:     &oh.Client{},
		url:     url,
		query:   make([]param, 0),
		form:    make([]param, 0),
		headers: make(map[string]string),
		cookies: make([]*oh.Cookie, 0),
	}
	return cli
}

func (that *Client) AddHeader(key, value string) *Client {
	that.headers[key] = value
	return that
}

func (that *Client) AddCookie(cookie *oh.Cookie) *Client {
	that.cookies = append(that.cookies, cookie)
	return that
}

// scheme Eg: "Bearer " or "Basic " or ""
func (that *Client) SetAuthorization(scheme, value string) *Client {
	that.authorization = &authorization{
		scheme: scheme, value: value,
	}
	return that
}

// 设置代理头
func (that *Client) UserAgent(value string) *Client {
	that.user_agent = value
	return that
}

// 添加URL参数
func (that *Client) AddQueryParam(name string, value interface{}) *Client {
	that.query = append(that.query, param{key: name, value: value})
	return that
}

// 添加FormData参数
func (that *Client) AddFormData(name string, value interface{}) *Client {
	that.form = append(that.form, param{name, value})
	return that
}

// Get 请求
func (that *Client) Get() (*Response, error) {
	url := that.warpperURL()
	req, err := oh.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	that.wrapperHeader(req).wrapperCookie(req).wrapperOther(req)
	response, err := that.cli.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	result := &Response{
		StatusCode:    response.StatusCode,
		Body:          body,
		Proto:         response.Proto,
		ContentLength: response.ContentLength,
		Header:        response.Header,
	}
	return result, nil
}

// Post 请求
func (that *Client) Post() (*Response, error) {
	url := that.warpperURL()
	datas := that.convertData()
	reader := strings.NewReader(datas.Encode())
	req, err := oh.NewRequest("POST", url, reader)
	if err != nil {
		return nil, err
	}
	that.wrapperHeader(req).wrapperCookie(req).wrapperOther(req)
	response, err := that.cli.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	result := &Response{
		StatusCode:    response.StatusCode,
		Body:          body,
		Proto:         response.Proto,
		ContentLength: response.ContentLength,
		Header:        response.Header,
	}
	return result, nil
}

// PostJSON 请求
func (that *Client) PostJSON(data interface{}) (*Response, error) {
	url := that.warpperURL()
	that.AddHeader("Content-Type", "application/json")
	bs, _ := json.Marshal(data)
	bnr := bytes.NewReader(bs)
	req, err := oh.NewRequest("POST", url, bnr)
	if err != nil {
		return nil, err
	}
	that.wrapperHeader(req).wrapperCookie(req).wrapperOther(req)
	response, err := that.cli.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	result := &Response{
		StatusCode:    response.StatusCode,
		Body:          body,
		Proto:         response.Proto,
		ContentLength: response.ContentLength,
		Header:        response.Header,
	}
	return result, nil
}

func (that *Client) wrapperHeader(req *oh.Request) *Client {
	for k, v := range that.headers {
		req.Header.Add(k, v)
	}
	return that
}

func (that *Client) wrapperCookie(req *oh.Request) *Client {
	for _, c := range that.cookies {
		req.AddCookie(c)
	}
	return that
}

func (that *Client) wrapperOther(req *oh.Request) *Client {
	req.Header.Set("User-Agent", that.user_agent)
	if that.authorization != nil {
		req.Header.Set("Authorization", fmt.Sprintf("%s%s", that.authorization.scheme, that.authorization.value))
	}
	return that
}

func (that *Client) warpperURL() string {
	if len(that.query) == 0 {
		return that.url
	}
	params := strings.Builder{}
	for _, v := range that.query {
		params.WriteString(fmt.Sprintf("&%s=%v", v.key, v.value))
	}
	if strings.Index(that.url, "?") > 0 {
		return that.url + "&" + url.PathEscape(params.String())
	} else {
		return that.url + "?" + url.PathEscape(strings.TrimPrefix(params.String(), "&"))
	}
}
func (that *Client) convertData() url.Values {
	res := url.Values{}
	for _, v := range that.form {
		if vi, vv := v.value.(int); vv {
			res.Add(v.key, fmt.Sprintf("%d", vi))
		} else if vs, vv := v.value.(string); vv {
			res.Add(v.key, vs)
		}
	}
	return res
}

func (that *Response) JSON(result interface{}) error {
	return json.Unmarshal(that.Body, result)
}

func (that *Response) String() string {
	return string(that.Body)
}
