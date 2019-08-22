package main

import (
	"github.com/ValleyZw/comgo"
	"html/template"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Channels and points
type Entry struct {
	AnalogIds []IDs `json:"analog_ids"`
	DigitIds  []IDs `json:"digital_ids"`
}

// Channel id
type IDs struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Points Points `json:"points"`
}

// Data point format
type Points struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Point Point  `json:"point"`
}

// Data point
type Point struct {
	X []string  `json:"x"`
	Y []float64 `json:"y"`
}

var entry Entry
var temp *template.Template

const AxisFormat = "2006-01-02 15:04:05.000"

func init() {
	temp = template.Must(template.ParseGlob("templates/*")) // need to cd to current folder
}

func main() {
	m := http.NewServeMux()
	m.HandleFunc("/", handleIndex)
	m.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8000", m)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Set a lower memory limit for multipart forms (default is 32 MiB)
		err := r.ParseMultipartForm(100 << 20) //100MiB
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//get a ref to the parsed multipart form
		files := r.MultipartForm.File
		fileNames := make(map[string]*multipart.FileHeader)
		for _, v := range files {
			fileNames[strings.ToLower(filepath.Ext(v[0].Filename))] = v[0]
		}

		cfg := comgo.New()
		entry = Entry{}

		if cfgFile, ok := fileNames[".cfg"]; ok {
			file, err := cfgFile.Open()
			defer file.Close()
			if err != nil {
				log.Println(err)
				return
			}
			err = cfg.ReadCFG(file)
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			http.Error(w, ".cfg file not exist", http.StatusInternalServerError)
			return
		}

		if dat, ok := fileNames[".dat"]; ok {
			file, err := dat.Open()
			defer file.Close()
			if err != nil {
				log.Println(err)
				return
			}
			cfg.ReadDAT(file)
		} else {
			http.Error(w, ".dat file not exist", http.StatusInternalServerError)
			return
		}

		ra := cfg.GetSamplingRate()
		ti := cfg.GetStartTime()

		var t []string
		for i := 0; i < cfg.GetSamplingNumber(); i++ {
			s, _ := time.ParseDuration(strconv.FormatFloat(float64(i)/float64(ra), 'g', 1, 64) + "s")
			t = append(t, ti.Add(s).Format(AxisFormat))
		}

		wg := sync.WaitGroup{}
		names := cfg.GetAnalogChannelNames()
		wg.Add(len(names))
		for k, v := range names {
			go func(k int, v string) {
				defer wg.Done()
				points, err := cfg.GetAnalogChannelData(uint16(k + 1))
				if err != nil {
					log.Println(err)
					return
				}
				anaPoints := Points{v, "line", Point{t, points}}
				entry.AnalogIds = append(entry.AnalogIds, IDs{v, v, anaPoints})
			}(k, v)
		}
		wg.Wait()
	}
	err := temp.ExecuteTemplate(w, "index.html", &entry)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error refreshing page", http.StatusInternalServerError)
		return
	}
}
