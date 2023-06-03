// main.go
package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"html/template"
	"io/ioutil"
	"net/http"
)

//go:embed templates/*.html
var templateFiles embed.FS

type Data struct {
	originJson       map[string]interface{}
	originResultSeq  []map[string]interface{}
	currentResultSeq []map[string]interface{}
}

func NewData(o map[string]interface{}) (*Data, error) {
	if origin, ok := o["results"].([]interface{}); ok {
		resultSeq := make([]map[string]interface{}, len(origin))
		for i, r := range origin {
			if res, ok := r.(map[string]interface{}); ok {
				resultSeq[i] = res
			} else {
				return nil, fmt.Errorf("init json error, result is not a map[string]interface{}")
			}
		}
		return &Data{originJson: o, originResultSeq: resultSeq, currentResultSeq: resultSeq}, nil
	}
	return nil, fmt.Errorf("init json error, json needs \"results\" key")
}

func (d *Data) FilterKeyValue(key, value string) bool {
	current := make([]map[string]interface{}, 0)
	for _, data := range d.currentResultSeq {
		for k, v := range data {
			if key == k {
				if containsString(v, value) {
					current = append(current, data)
				}
			}
		}
	}
	if len(current) == 0 {
		return false
	}
	d.currentResultSeq = current
	return true
}

func (d *Data) FilterValue(value string) bool {
	current := make([]map[string]interface{}, 0)
	for _, data := range d.currentResultSeq {
		for _, v := range data {
			if containsString(v, value) {
				current = append(current, data)
			}
		}
	}
	if len(current) == 0 {
		return false
	}
	d.currentResultSeq = current
	return true
}

func containsString(data interface{}, target string) bool {
	switch v := data.(type) {
	case string:
		return v == target
	case []interface{}:
		for _, item := range v {
			if containsString(item, target) {
				return true
			}
		}
	case map[string]interface{}:
		for _, item := range v {
			if containsString(item, target) {
				return true
			}
		}
	}

	return false
}

func main() {
	var GlobalData *Data
	router := gin.Default()
	templates, err := template.ParseFS(templateFiles, "templates/*.html")
	if err != nil {
		log.Fatal("Failed to parse templates: ", err)
	}

	router.SetHTMLTemplate(templates)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			log.Error("Failed to get uploaded file: ", err)
			c.String(http.StatusBadRequest, "Failed to get uploaded file")
			return
		}

		src, err := file.Open()
		if err != nil {
			log.Error("Failed to open uploaded file: ", err)
			c.String(http.StatusInternalServerError, "Failed to open uploaded file")
			return
		}
		defer src.Close()

		data, err := ioutil.ReadAll(src)
		if err != nil {
			log.Error("Failed to read uploaded file: ", err)
			c.String(http.StatusInternalServerError, "Failed to read uploaded file")
			return
		}

		var jsonData map[string]interface{}
		err = json.Unmarshal(data, &jsonData)
		if err != nil {
			log.Error("Failed to unmarshal JSON data: ", err)
			c.String(http.StatusBadRequest, "Invalid JSON data")
			return
		}

		GlobalData, err = NewData(jsonData)
		if err != nil {
			log.Error("Failed to create Data instance: ", err)
			c.String(http.StatusBadRequest, "Failed to create Data instance")
			return
		}

		//c.String(http.StatusOK, "JSON data uploaded successfully")
		c.HTML(http.StatusOK, "results.html", gin.H{
			"Data": GlobalData,
		})
	})

	router.POST("/filter_key_value", func(c *gin.Context) {
		key := c.PostForm("key")
		value := c.PostForm("value")

		if GlobalData == nil {
			log.Error("Data is not initialized")
			c.String(http.StatusBadRequest, "Data is not initialized")
			return
		}

		filtered := GlobalData.FilterKeyValue(key, value)
		if filtered {
			c.String(http.StatusOK, "Filter applied successfully")
			c.HTML(http.StatusOK, "results.html", gin.H{
				"Data": GlobalData,
			})
		} else {
			c.String(http.StatusOK, "No matching results.html found")
			c.HTML(http.StatusOK, "results.html", gin.H{
				"Data": GlobalData,
			})
		}
	})

	router.POST("/filter_value", func(c *gin.Context) {
		value := c.PostForm("value")

		if GlobalData == nil {
			log.Error("Data is not initialized")
			c.String(http.StatusBadRequest, "Data is not initialized")
			return
		}

		filtered := GlobalData.FilterValue(value)
		if filtered {
			c.HTML(http.StatusOK, "results.html", gin.H{
				"Data": GlobalData,
			})
		} else {
			c.HTML(http.StatusOK, "results.html", gin.H{
				"Data": GlobalData,
			})
		}
	})

	router.GET("/download", func(c *gin.Context) {
		if GlobalData == nil {
			log.Error("Data is not initialized")
			c.String(http.StatusBadRequest, "Data is not initialized")
			return
		}

		filename := "filtered_data.json"
		fileData, err := json.Marshal(GlobalData.currentResultSeq)
		if err != nil {
			log.Error("Failed to marshal JSON data: ", err)
			c.String(http.StatusInternalServerError, "Failed to generate JSON file")
			return
		}

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		c.Data(http.StatusOK, "application/json", fileData)
	})

	router.Run(":8080")
}
