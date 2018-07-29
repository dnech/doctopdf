package main

import (
	"net/http"
	"html/template"
	"strings"
	"log"
	"github.com/labstack/echo"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/satori/go.uuid"
	"os"
	"fmt"
	"encoding/json"
	"bytes"
)

type Pdf struct {
	Dpi       uint `json:"dpi"`
	Orientation string `json:"orientation"`
	Grayscale bool `json:"grayscale"`
	PageSize string `json:"pagesize"`
}

type Config struct {
	Host string `json:"host"`
	Static string `json:"static"`
	Pdf Pdf `json:"pdf"`
}

func LoadConfiguration(file string) Config {

	// Default
	config := Config{}
	config.Host = ":5555"
	config.Static = "http://localhost:5555/static"
	config.Pdf.Dpi = 300
	config.Pdf.Grayscale = false
	config.Pdf.Orientation = wkhtmltopdf.OrientationPortrait
	config.Pdf.PageSize = wkhtmltopdf.PageSizeA4


	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
		panic(err);
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

var config Config

func main() {
	config = LoadConfiguration("./config.json");

	e := echo.New()

	e.Static("static", "./static")

	e.GET("/temp/:id", func(c echo.Context) error {
		id := c.Param("id")
		return c.File("./temp/" + id + ".pdf")
	})

	e.GET("/doc/:id", getDoc)
	e.POST("/doc/:id", getDoc)

	e.GET("/pdf/:id", getPdf)
	e.POST("/pdf/:id", getPdf)

	e.Logger.Fatal(e.Start(config.Host))
}

type TemplateData struct {
	Id string
	Config Config
	Method string
	Data string
	Query map[string]string
	Form map[string]string
	Cookies map[string]string
}

func convertParams(data map[string][]string) map[string]string {
	ret := make(map[string]string)
	for k, v := range data {
		str := ""
		for _, a := range v {
			str = str + string(a)
		}
		ret[k] = str
	}
	return ret
}

func convertCookie(data []*http.Cookie) map[string]string {
	ret := make(map[string]string)
	for _, v := range data {
		ret[v.Name] = v.Value
	}
	return ret
}

func getHtml(c echo.Context) (string, error) {
	id := c.Param("id")

	t, err := template.ParseFiles("./template/" + id + ".html")
	if err != nil {
		return "", err
	}

	data := TemplateData{
		Id: id,
		Config: config,
		Method: c.Request().Method,
		Query: convertParams(c.Request().URL.Query()),
		Form: convertParams(c.Request().URL.Query()),
		Cookies: convertCookie(c.Cookies()),
	}

	params_json := "{}"
	params, ok := c.Request().URL.Query()["data"]
	if ok {
		if len(params[0]) > 0 {
			params_json= params[0]
		}
	}

	err = c.Request().ParseForm()
	if err == nil {
		data.Form = convertParams(c.Request().Form)
		form_data := c.Request().Form.Get("data")
		if form_data != "" {
			params_json = form_data
		}
	}


	fmt.Println(params_json)
	fmt.Println(c.Request().Method)
	fmt.Println(c.ParamValues())
	fmt.Println(c.Cookies())

	data.Data = params_json

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

func getDoc(c echo.Context) error {
	html, err := getHtml(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.HTML(http.StatusOK, html)
}

func getPdf(c echo.Context) error {
	html, err := getHtml(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	id := c.Param("id")
	guid := uuid.Must(uuid.NewV4()).String()

	if err := convertToPdf(html, id + "_" + guid); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.File("./temp/" + id + "_" + guid + ".pdf")
}

func convertToPdf(html string, filename string) error {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Println(err)
		return err
	}

	// Set global options
	pdfg.Dpi.Set(config.Pdf.Dpi)
	pdfg.Orientation.Set(config.Pdf.Orientation)
	pdfg.Grayscale.Set(config.Pdf.Grayscale)
	pdfg.PageSize.Set(config.Pdf.PageSize)

	// Create a new input page from HTML
	pdfg.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(html)))

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		log.Println(err)
		return err
	}

	// Write buffer contents to file on disk
	err = pdfg.WriteFile("./temp/" + filename + ".pdf")
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Created file './temp/" + filename + ".pdf'")
	return nil
}