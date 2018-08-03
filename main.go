package main

import (
	"net/http"
	"strings"
	"log"
	"github.com/labstack/echo"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/satori/go.uuid"
	"os"
	"fmt"
	"encoding/json"
)

type Pdf struct {
	Dpi       uint `json:"dpi"`
	Orientation string `json:"orientation"`
	Grayscale bool `json:"grayscale"`
	PageSize string `json:"pagesize"`
	MarginTop uint `json:"margin_top"`
	MarginBottom uint `json:"margin_bottom"`
	MarginRight uint `json:"margin_right"`
	MarginLeft uint `json:"margin_left"`
}

type Config struct {
	Host string `json:"host"`
	Static string `json:"static"`
	Public string `json:"public"`
	Pdf Pdf `json:"pdf"`
}

func LoadConfiguration(file string) Config {

	// Default
	config := Config{}
	config.Host = ":8888"
	config.Static = "http://localhost:8888/static"
	config.Public = "http://localhost:8888"
	config.Pdf.Dpi = 300
	config.Pdf.Grayscale = false
	config.Pdf.Orientation = wkhtmltopdf.OrientationPortrait
	config.Pdf.PageSize = wkhtmltopdf.PageSizeA4
	config.Pdf.MarginTop = 5
	config.Pdf.MarginBottom = 5
	config.Pdf.MarginRight  = 5
	config.Pdf.MarginLeft = 20

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
	e.Static("/", "./public")

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
	Id string `json:"id"`
	Config Config `json:"config"`
	Data json.RawMessage `json:"data"`
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

func getDataParams(c echo.Context) string {
	params_json := "{}"
	params, ok := c.Request().URL.Query()["data"]
	if ok {
		if len(params[0]) > 0 {
			params_json= params[0]
		}
	}

	err := c.Request().ParseForm()
	if err == nil {
		form_data := c.Request().Form.Get("data")
		if form_data != "" {
			params_json = form_data
		}
	}
	return params_json
}


func loadPdfConfig(c echo.Context) Pdf {
	id := c.Param("id")

	// Default
	pdf := Pdf{}
	pdf.Dpi = config.Pdf.Dpi
	pdf.Grayscale = config.Pdf.Grayscale
	pdf.Orientation = config.Pdf.Orientation
	pdf.PageSize = config.Pdf.PageSize
	pdf.MarginTop = config.Pdf.MarginTop
	pdf.MarginBottom = config.Pdf.MarginBottom
	pdf.MarginRight = config.Pdf.MarginRight
	pdf.MarginLeft = config.Pdf.MarginLeft

	file, err := os.Open("./template/" + id + ".json")
	if err != nil {
		return pdf
	}

	jsonParser := json.NewDecoder(file)
	jsonParser.Decode(&pdf)

	return pdf
}

func getReplaceHtml(c echo.Context) (string, error) {
	id := c.Param("id")

	file, err := os.Open("./template/" + id + ".html")
	if err != nil {
		return "", err
	}
	defer file.Close()

	// получить размер файла
	stat, err := file.Stat()
	if err != nil {
		return "", err
	}
	// чтение файла
	bs := make([]byte, stat.Size())
	_, err = file.Read(bs)
	if err != nil {
		return "", err
	}

	str := string(bs)


	data := TemplateData{
		Id: id,
		Config: config,
		Data: json.RawMessage(getDataParams(c)),
	}

	params_json := ""
	b, err := json.Marshal(data)
	if err != nil {
		params_json = `{"message":"`+ string(err.Error()) +`"}`
	} else {
		params_json = string(b)
	}

	str = strings.Replace(str, "${template.id}", id, -1)
	str = strings.Replace(str, "${template.static}", config.Static, -1)
	str = strings.Replace(str, "${template.public}", config.Public, -1)

	params_json = strings.Replace(params_json, "\\", "\\\\", -1)
	params_json = strings.Replace(params_json, "\"", "\\\"", -1)

	str = strings.Replace(str, "${template.data}", params_json, -1)

	return str, nil
}

/*
*   CONTENT DOC
*/
func getDoc(c echo.Context) error {
	html, err := getReplaceHtml(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.HTML(http.StatusOK, html)
}

/*
*   CONTENT PDF
*/
func getPdf(c echo.Context) error {
	html, err := getReplaceHtml(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	pdf_config := loadPdfConfig(c)

	id := c.Param("id")
	guid := uuid.Must(uuid.NewV4()).String()

	if err := convertToPdf(html, id + "_" + guid, pdf_config); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.File("./temp/" + id + "_" + guid + ".pdf")
}

/*
* CONVERT TO PDF
*/
func convertToPdf(html string, filename string, config Pdf) error {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Println(err)
		return err
	}

	// Set global options
	pdfg.Dpi.Set(config.Dpi)
	pdfg.Orientation.Set(config.Orientation)
	pdfg.Grayscale.Set(config.Grayscale)
	pdfg.PageSize.Set(config.PageSize)
	pdfg.MarginTop.Set(config.MarginTop)
	pdfg.MarginBottom.Set(config.MarginBottom)
	pdfg.MarginRight.Set(config.MarginRight)
	pdfg.MarginLeft.Set(config.MarginLeft)

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