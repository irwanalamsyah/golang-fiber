package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var app = fiber.New()

func TestRoutingHelloWorld(t *testing.T) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello World!")
	})

	request := httptest.NewRequest("GET", "/", nil)
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello World!", string(bytes))
}

func TestCtx(t *testing.T) {
	app.Get("/hello", func(ctx *fiber.Ctx) error {
		name := ctx.Query("name", "Guest")
		return ctx.SendString("Hello " + name)
	})

	request := httptest.NewRequest("GET", "/hello?name=Irwan", nil)
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello Irwan", string(bytes))

	request = httptest.NewRequest("GET", "/hello", nil)
	response, err = app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err = io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello Guest", string(bytes))
}

func TestHttpRequest(t *testing.T) {
	app.Get("/request", func(ctx *fiber.Ctx) error {
		first := ctx.Get("firstname")   // Header
		last := ctx.Cookies("lastname") // Cookie
		return ctx.SendString("Hello " + first + " " + last)

	})

	request := httptest.NewRequest("GET", "/request", nil)
	request.Header.Set("firstname", "Irwan")
	request.AddCookie(&http.Cookie{Name: "lastname", Value: "Alamsyah"})
	response, err := app.Test(request)
	assert.Nil(t, err)
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello Irwan Alamsyah", string(bytes))
}

func TestRouteParameter(t *testing.T) {
	app.Get("/users/:userId/orders/:orderId", func(ctx *fiber.Ctx) error {
		userId := ctx.Params("userId")
		orderId := ctx.Params("orderId")
		return ctx.SendString("Get Order " + orderId + " From User " + userId)
	})

	request := httptest.NewRequest("GET", "/users/irwan/orders/10", nil)
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	bytes, err := io.ReadAll(response.Body)
	assert.Equal(t, "Get Order 10 From User irwan", string(bytes))
}

func TestFormRequest(t *testing.T) {
	app.Post("/hello", func(ctx *fiber.Ctx) error {
		name := ctx.FormValue("name")
		return ctx.SendString("Hello " + name)
	})

	body := strings.NewReader("name=Irwan")
	request := httptest.NewRequest("POST", "/hello", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello Irwan", string(bytes))
}

//go:embed source/contoh.txt
var contohFile []byte

func TestMultipartForm(t *testing.T) {
	app.Post("/upload", func(ctx *fiber.Ctx) error {
		file, err := ctx.FormFile("file")
		if err != nil {
			return err
		}

		err = ctx.SaveFile(file, "./target/"+file.Filename)
		if err != nil {
			return err
		}

		return ctx.SendString("Upload Success")
	})

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	file, _ := writer.CreateFormFile("file", "contoh.txt")
	file.Write(contohFile)
	writer.Close()

	request := httptest.NewRequest("POST", "/upload", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response, err := app.Test(request)
	assert.Nil(t, err)
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Upload Success", string(bytes))
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestRequestBody(t *testing.T) {
	app.Post("/login", func(ctx *fiber.Ctx) error {
		body := ctx.Body()
		request := new(LoginRequest)
		err := json.Unmarshal(body, request)
		if err != nil {
			return err
		}
		return ctx.SendString("Login Success " + request.Username)
	})

	body := strings.NewReader(`{"username": "Irwan", "password":"rahasia"}`)
	request := httptest.NewRequest("POST", "/login", body)
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Login Success Irwan", string(bytes))
}

type RegisterRequest struct {
	Username string `json:"username" xml:"username" form:"username"`
	Password string `json:"password" xml:"password" form:"password"`
	Name     string `json:"name" xml:"name" form:"name"`
}

func TestBodyParser(t *testing.T) {
	app.Post("/register", func(ctx *fiber.Ctx) error {
		request := new(RegisterRequest)
		err := ctx.BodyParser(request)
		if err != nil {
			return err
		}

		return ctx.SendString("Register Success " + request.Username)
	})
}

func TestBodyParserJSON(t *testing.T) {
	TestBodyParser(t)

	body := strings.NewReader(`{"username": "Irwan", "password":"rahasia", "name":"Irwan Alamsyah"}`)
	request := httptest.NewRequest("POST", "/register", body)
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)
	assert.Nil(t, err)
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Register Success Irwan", string(bytes))

}

func TestBodyParserForm(t *testing.T) {
	TestBodyParser(t)

	body := strings.NewReader(`username=Irwan&password=rahasia&name=Irwan+Alamsyah`)
	request := httptest.NewRequest("POST", "/register", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := app.Test(request)
	assert.Nil(t, err)
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Register Success Irwan", string(bytes))
}

func TestBodyParserXML(t *testing.T) {
	TestBodyParser(t)

	body := strings.NewReader(
		`<RegisterRequest>
				<username>Irwan</username>
				<password>rahasia</password>
				<name>Irwan Alamsyah</name>
			</RegisterRequest>`)

	request := httptest.NewRequest("POST", "/register", body)
	request.Header.Set("Content-Type", "application/xml")
	response, err := app.Test(request)
	assert.Nil(t, err)
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Register Success Irwan", string(bytes))
}

func TestResponseJSON(t *testing.T) {
	app.Get("/user", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"username": "irwanalamsyah",
			"name": "Irwan Alamsyah",
		})
	})

	request := httptest.NewRequest("GET", "/user", nil)
	response, err := app.Test(request)
	assert.Nil(t, err)
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"name":"Irwan Alamsyah","username":"irwanalamsyah"}`, string(bytes))
}

func TestDownloadFile(t *testing.T) {
	app.Get("/download", func(ctx *fiber.Ctx) error {
		return ctx.Download("./source/contoh_download.txt", "contoh_download.txt")
	})

	request := httptest.NewRequest("GET", "/download", nil)
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, `attachment; filename="contoh_download.txt"`, response.Header.Get("Content-Disposition"))
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "this is sample file for upload", string(bytes))
}

func TestRoutingGroup(t *testing.T) {
	helloWorld := func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello World")
	}

	api := app.Group("/api")
	api.Get("/hello", helloWorld)		// -> /api/hello
	api.Get("/world", helloWorld)		// -> /api/world

	web := app.Group("/web")
	web.Get("/hello", helloWorld)		// -> /web/hello
	web.Get("/world", helloWorld)		// -> /web/world

	request := httptest.NewRequest("GET", "/api/hello", nil)
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello World", string(bytes))
}

func TestStatic(t *testing.T) {
	app.Static("/public", "./source")

	request := httptest.NewRequest("GET", "/public/contoh_download.txt", nil)
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "this is sample file for upload", string(bytes))
}
