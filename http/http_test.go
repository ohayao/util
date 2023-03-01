package http

import (
	"fmt"
	"net/http"
	"testing"
)

type Res struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

func TestGet(t *testing.T) {
	hc := New("http://localhost:12388/")
	hc.AddQueryParam("name", "zhang san")
	hc.AddQueryParam("age", 10)
	res, _ := hc.Get()
	fmt.Println(res.String())
}

func TestUserJWT(t *testing.T) {
	fmt.Println("第一步登录获取token")
	jwt_token := ""
	hc := New("http://localhost:12388/api/v1/user/login")
	data := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "admin",
		Password: "123456",
	}
	if res, err := hc.PostJSON(data); err != nil {
		t.Error(err)
	} else {
		res_json := &Res{}
		res.JSON(res_json)
		fmt.Println(res_json)
		if tk, ok := res_json.Data["token"]; ok {
			jwt_token = tk.(string)
		}
	}
	fmt.Println("第二步获取用户信息")
	info_req := New("http://localhost:12388/api/v1/admin/userdetail")
	info_req.SetAuthorization("", jwt_token)
	if info_res, err := info_req.Post(); err != nil {
		t.Error(err)
	} else {
		fmt.Println(info_res)
	}
	fmt.Println("第三步获取用户菜单")
	auth_req := New("http://localhost:12388/api/v1/admin/usermenus")
	auth_req.SetAuthorization("", jwt_token)
	if auth_res, err := auth_req.Post(); err != nil {
		t.Error(err)
	} else {
		fmt.Println(auth_res)
	}
}

func TestUserCookie(t *testing.T) {
	fmt.Println("第一步登录获取Cookie")
	cookie := ""
	hc := New("http://localhost:12388/api/v1/user/login")
	data := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "admin",
		Password: "123456",
	}
	if res, err := hc.PostJSON(data); err != nil {
		t.Error(err)
	} else {
		res_json := &Res{}
		res.JSON(res_json)
		fmt.Println(res_json)
		if tk, ok := res_json.Data["token"]; ok {
			cookie = tk.(string)
		}
	}
	fmt.Println("第二步获取用户信息")
	info_req := New("http://localhost:12388/api/v1/admin/userdetail")
	info_req.AddCookie(&http.Cookie{Name: "imadmin", Value: cookie})
	if info_res, err := info_req.Post(); err != nil {
		t.Error(err)
	} else {
		fmt.Println(info_res)
	}
	fmt.Println("第三步获取用户菜单")
	auth_req := New("http://localhost:12388/api/v1/admin/usermenus")
	auth_req.AddCookie(&http.Cookie{Name: "imadmin", Value: cookie})
	if auth_res, err := auth_req.Post(); err != nil {
		t.Error(err)
	} else {
		fmt.Println(auth_res)
	}
}

func BenchmarkGetMenus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		req := New("http://localhost:12388/api/v1/admin/usermenus")
		req.SetAuthorization("", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJUb2tlbiI6IjE2MzA3MjkyODEyOTIyMTAxNzYiLCJleHAiOjE2Nzc2MzgzMzIsImlhdCI6MTY3NzYzMTEzMiwiaXNzIjoiQWRtaW5pc3RyYXRvciJ9.uCbx64EKhtwgxM9NEujhHDEE0wxdRVbrGxWaWZ1JSqo")
		res, _ := req.Post()
		fmt.Println(res)
	}
}

func BenchmarkGetUserDetail(b *testing.B) {
	for i := 0; i < b.N; i++ {
		req := New("http://localhost:12388/api/v1/admin/userdetail")
		req.SetAuthorization("", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJUb2tlbiI6IjE2MzA3MjkyODEyOTIyMTAxNzYiLCJleHAiOjE2Nzc2MzgzMzIsImlhdCI6MTY3NzYzMTEzMiwiaXNzIjoiQWRtaW5pc3RyYXRvciJ9.uCbx64EKhtwgxM9NEujhHDEE0wxdRVbrGxWaWZ1JSqo")
		res, _ := req.Post()
		fmt.Println(res)
	}
}
