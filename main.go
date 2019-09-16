package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"os/exec"

	"github.com/alexflint/go-arg"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/echoview-v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var args struct {
	Port int
}

func main() {
	arg.MustParse(&args)
	if args.Port == 0 {
		args.Port = 1323
	}
	e := echo.New()
	e.Debug = true
	renderer := echoview.New(goview.Config{
		Root:         "views",
		Extension:    ".html",
		Master:       "",
		Partials:     []string{},
		Funcs:        make(template.FuncMap),
		DisableCache: true,
		Delims:       goview.Delims{Left: "{{", Right: "}}"},
	})
	e.Renderer = renderer
	e.Use(middleware.Recover())

	e.GET("/", home)
	e.GET("/list", getProcess)
	e.POST("/send", sendCommand)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", args.Port)))
}

// Handler
func home(c echo.Context) error {
	return c.Render(http.StatusOK, "home", nil)
}

func getProcess(c echo.Context) error {
	cmd := exec.Command("pm2", "jlist")
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, string(out))
}

func sendCommand(c echo.Context) error {
	p := c.FormValue("p")
	plist := strings.Split(p, " ")

	cmd := exec.Command("pm2", plist...)

	out, err := cmd.Output()
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, string(out))
}
