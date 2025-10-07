package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

/*
'Name': name,
'Tarih': "",
'Maktülün yaşı': "",
'İl/ilçe': "",
'Neden öldürüldü': "",
'Kim tarafından öldürüldü': "",
'Korunma talebi': "",
'Öldürülme şekli': "",
'Failin durumu': "",
'Kaynak': "",
'image': img_src,
'url': href
*/

type Incident struct {
	Id         int      `json:"id"`
	Name       string   `json:"name"`
	FullName   string   `json:"fullname"`
	Age        string   `json:"age"`
	Location   string   `json:"location"`
	Date       string   `json:"date"`
	Reason     string   `json:"reason"`
	By         string   `json:"by"`
	Protection string   `json:"protection"`
	Method     string   `json:"method"`
	Status     string   `json:"status"`
	Source     []string `json:"source"`
	Image      string   `json:"image"`
	Url        string   `json:"url"`
}

type Detail struct {
	Name       string   `json:"name"`
	Age        string   `json:"age"`
	Location   string   `json:"location"`
	Date       string   `json:"date"`
	Reason     string   `json:"reason"`
	By         string   `json:"by"`
	Protection string   `json:"protection"`
	Method     string   `json:"method"`
	Status     string   `json:"status"`
	Source     []string `json:"source"`
	Image      string   `json:"image"`
}

// Static and dynamic Variables
var (
	baseUrl      = "https://anitsayac.com"
	jsonFileName = "data.json"
	csvFileName  = "data.csv"
)

func ReplaceAll(s, old, new string, n int) string {
	// Replace All for Regexp
	re := regexp.MustCompile(old)
	return re.ReplaceAllString(s, new)
}

// validateFiles checks if the generated files are valid and have reasonable content
func validateFiles(jsonFile, csvFile string, expectedCount int) bool {
	// Check if files exist
	if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
		log.Printf("JSON file %s does not exist", jsonFile)
		return false
	}

	if _, err := os.Stat(csvFile); os.IsNotExist(err) {
		log.Printf("CSV file %s does not exist", csvFile)
		return false
	}

	// Check file sizes
	jsonInfo, err := os.Stat(jsonFile)
	if err != nil {
		log.Printf("Failed to get JSON file info: %s", err)
		return false
	}

	csvInfo, err := os.Stat(csvFile)
	if err != nil {
		log.Printf("Failed to get CSV file info: %s", err)
		return false
	}

	// Files should not be empty
	if jsonInfo.Size() == 0 {
		log.Printf("JSON file is empty")
		return false
	}

	if csvInfo.Size() == 0 {
		log.Printf("CSV file is empty")
		return false
	}

	// Basic JSON validation - try to parse
	jsonFileContent, err := os.ReadFile(jsonFile)
	if err != nil {
		log.Printf("Failed to read JSON file: %s", err)
		return false
	}

	var incidents []Incident
	if err := json.Unmarshal(jsonFileContent, &incidents); err != nil {
		log.Printf("Invalid JSON format: %s", err)
		return false
	}

	// Check if we have reasonable number of incidents
	if len(incidents) < expectedCount/2 {
		log.Printf("Too few incidents in JSON: got %d, expected at least %d", len(incidents), expectedCount/2)
		return false
	}

	// Basic CSV validation - count lines
	csvFileContent, err := os.ReadFile(csvFile)
	if err != nil {
		log.Printf("Failed to read CSV file: %s", err)
		return false
	}

	lines := strings.Split(string(csvFileContent), "\n")
	// Should have header + data lines (allowing for empty last line)
	if len(lines) < expectedCount {
		log.Printf("Too few lines in CSV: got %d, expected at least %d", len(lines), expectedCount)
		return false
	}

	return true
}

// Get article content
func getArticleContent(url string) Detail {
	/*
		<b>Ad Soyad:</b> Fidan Çakır<br><b>Maktülün yaşı: </b>Reşit<br><b>İl/ilçe: </b>İzmir<br><b>Tarih: </b>16/10/2024<br><b>Neden öldürüldü:</b>  Tespit Edilemeyen<br><b>Kim tarafından öldürüldü:</b>  Tespit Edilemeyen<br><b>Korunma talebi:</b>  Yok<br><b>Öldürülme şekli:</b>  Kesici Alet<br><b>Failin durumu: </b>Soruşturma Sürüyor<br><b>Kaynak:</b>  <a target=_blank href='https://www.t24.com.tr/haber/supheli-bir-kadin-olumu-daha-izmir-de-bir-kadin-evinde-vucudunda-kesi-izleriyle-olu-bulundu,1190726'><u>https://www.t24.com.tr/haber/supheli-bir-kadin-olumu-daha-izmir-de-bir-kadin-evinde-vucudunda-kesi-izleriyle-olu-bulundu,1190726</u></a><br><img width=750 style='margin-top:10px' src=ii/3202024.jpg>

		<b>Ad Soyad:</b> Beyhan Yavuz<br><b>Tarih: </b>05/02/2008<br><b>Neden öldürüldü:</b>  Reddetme<br><b>Kim tarafından öldürüldü:</b>  Dini nikahlı kocası<br><b>Korunma talebi:</b>  Tespit Edilemeyen<br><b>Öldürülme şekli:</b>  Ateşli Silah<br><b>Kaynak:</b>  <a target=_blank href='http://hurarsiv.hurriyet.com.tr/goster/haber.aspx?id=8170561&tarih=2008-02-05'><u>http://hurarsiv.hurriyet.com.tr/goster/haber.aspx?id=8170561&tarih=2008-02-05</u></a>

		<b>Ad Soyad:</b> Ayşenur Halil<br><b>Maktülün yaşı: </b>Reşit<br><b>İl/ilçe: </b>İstanbul<br><b>Tarih: </b>04/10/2024<br><b>Neden öldürüldü:</b>  Tespit Edilemeyen<br><b>Kim tarafından öldürüldü:</b>  Eski Sevgilisi<br><b>Korunma talebi:</b>  Yok<br><b>Öldürülme şekli:</b>  Kesici Alet<br><b>Failin durumu: </b>İntihar<br><b>Kaynak:</b>  <br><a target=_blank href='https://www.t24.com.tr/haber/fatih-te-vahset-kadinin-kafasini-kesip-surlardan-atladi,1187812'><u>https://www.t24.com.tr/haber/fatih-te-vahset-kadinin-kafasini-kesip-surlardan-atladi,1187812</u></a><br><a target=_blank href='https://www.t24.com.tr/haber/yarim-saat-arayla-iki-kadini-katleden-semih-celik-5-kez-psikolojik-tedavi-gormus,1187904'><u>https://www.t24.com.tr/haber/yarim-saat-arayla-iki-kadini-katleden-semih-celik-5-kez-psikolojik-tedavi-gormus,1187904</u></a><br><img width=750 style='margin-top:10px' src=ii/2892024.jpg>

		<b>Ad Soyad:</b> İkbal Uzuner<br><b>Maktülün yaşı: </b>Reşit<br><b>İl/ilçe: </b>İstanbul<br><b>Tarih: </b>04/10/2024<br><b>Neden öldürüldü:</b>  Tespit Edilemeyen<br><b>Kim tarafından öldürüldü:</b>  Tanımadığı Birisi <br><b>Korunma talebi:</b>  Yok<br><b>Öldürülme şekli:</b>  Kesici Alet<br><b>Failin durumu: </b>İntihar<br><b>Kaynak:</b>  <br><a target=_blank href='https://www.t24.com.tr/haber/fatih-te-vahset-kadinin-kafasini-kesip-surlardan-atladi,1187812'><u>https://www.t24.com.tr/haber/fatih-te-vahset-kadinin-kafasini-kesip-surlardan-atladi,1187812</u></a><br><a target=_blank href='https://www.t24.com.tr/haber/yarim-saat-arayla-iki-kadini-katleden-semih-celik-5-kez-psikolojik-tedavi-gormus,1187904'><u>https://www.t24.com.tr/haber/yarim-saat-arayla-iki-kadini-katleden-semih-celik-5-kez-psikolojik-tedavi-gormus,1187904</u></a><br><img width=750 style='margin-top:10px' src=ii/2902024.jpg>
	*/

	detail := Detail{}
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("anitsayac.com"),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./anitsayac_cache"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting Detail: ", r.URL.String())
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		nameMatches := regexp.MustCompile(`(?i)<b>Ad Soyad:\s*</b>\s*(.+?)<br>`).FindStringSubmatch(string(e.Response.Body))
		if len(nameMatches) > 1 {
			detail.Name = nameMatches[1]
		}

		ageMatches := regexp.MustCompile(`(?i)<b>Maktülün yaşı:\s*</b>\s*(.+?)<br>`).FindStringSubmatch(string(e.Response.Body))
		if len(ageMatches) > 1 {
			detail.Age = ageMatches[1]
		}

		locationMatches := regexp.MustCompile(`(?i)<b>İl/ilçe:\s*</b>\s*(.+?)<br>`).FindStringSubmatch(string(e.Response.Body))
		if len(locationMatches) > 1 {
			detail.Location = locationMatches[1]
		}

		dateMatches := regexp.MustCompile(`(?i)<b>Tarih:\s*</b>\s*(.+?)<br>`).FindStringSubmatch(string(e.Response.Body))
		if len(dateMatches) > 1 {
			detail.Date = dateMatches[1]
		}

		reasonMatches := regexp.MustCompile(`(?i)<b>Neden öldürüldü:\s*</b>\s*(.+?)<br>`).FindStringSubmatch(string(e.Response.Body))
		if len(reasonMatches) > 1 {
			detail.Reason = reasonMatches[1]
		}

		byMatches := regexp.MustCompile(`(?i)<b>Kim tarafından öldürüldü:\s*</b>\s*(.+?)<br>`).FindStringSubmatch(string(e.Response.Body))
		if len(byMatches) > 1 {
			detail.By = byMatches[1]
		}

		protectionMatches := regexp.MustCompile(`(?i)<b>Korunma talebi:\s*</b>\s*(.+?)<br>`).FindStringSubmatch(string(e.Response.Body))
		if len(protectionMatches) > 1 {
			detail.Protection = protectionMatches[1]
		}

		methodMatches := regexp.MustCompile(`(?i)<b>Öldürülme şekli:\s*</b>\s*(.+?)<br>`).FindStringSubmatch(string(e.Response.Body))
		if len(methodMatches) > 1 {
			detail.Method = methodMatches[1]
		}

		statusMatches := regexp.MustCompile(`(?i)<b>Failin durumu:\s*</b>\s*(.+?)<br>`).FindStringSubmatch(string(e.Response.Body))
		if len(statusMatches) > 1 {
			detail.Status = statusMatches[1]
		}

		sources := e.ChildAttrs("a", "href")
		detail.Source = sources

		detail.Image = baseUrl + "/" + e.ChildAttr("img", "src")
	})

	c.Visit(url)

	// Return detail object
	return detail
}

func main() {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("anitsayac.com"),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./anitsayac_cache"),
	)

	// Create another collector to scrape additional details
	// detailCollector := c.Clone()

	incidents := make([]Incident, 0, 100000)

	/*
		<span class='xxy'>
			<a href='details.aspx?id=38931'  data-width='800' data-height='380'  class='html5lightbox'  adata-group='mygroup' >Burçin Sevgi T.</a>
			</span>
		<span class='xxy'>
			<a href='details.aspx?id=38933'  data-width='800' data-height='380'  class='html5lightbox'  adata-group='mygroup' >Rojin Kabaiş</a>
		</span>
	*/
	// Blog post page scraper
	c.OnHTML("div#divcounter", func(e *colly.HTMLElement) {
		e.ForEach("span.xxy", func(i int, e *colly.HTMLElement) {
			// Detail url: https://anitsayac.com/details.aspx?id=38931
			detail := getArticleContent(baseUrl + "/" + e.ChildAttr("span.xxy > a", "href"))
			incident := Incident{
				Id: func() int {
					id, _ := strconv.Atoi(strings.Split(e.ChildAttr("span.xxy > a", "href"), "=")[1])
					return id
				}(),
				Name:       e.ChildText("span.xxy > a"),
				FullName:   detail.Name,
				Age:        detail.Age,
				Location:   detail.Location,
				Date:       detail.Date,
				Reason:     detail.Reason,
				By:         detail.By,
				Protection: detail.Protection,
				Method:     detail.Method,
				Status:     detail.Status,
				Source:     detail.Source,
				Image:      detail.Image,
				Url:        baseUrl + "/" + e.ChildAttr("span.xxy > a", "href"),
			}

			incidents = append(incidents, incident)
		})
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL.String())
	})

	// Url: https://anitsayac.com/?year=2000
	c.Visit(baseUrl + "/?year=2000")

	// Check if we have valid data before writing
	if len(incidents) == 0 {
		log.Println("No incidents found, not updating files")
		return
	}

	// Write to temporary files first to avoid data loss
	tempJsonFile := jsonFileName + ".tmp"
	tempCsvFile := csvFileName + ".tmp"

	// Write JSON to temporary file
	file, err := os.Create(tempJsonFile)
	if err != nil {
		log.Fatalf("Cannot create temp file %q: %s\n", tempJsonFile, err)
		return
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	// Encode json to temporary file
	if err := enc.Encode(incidents); err != nil {
		log.Fatalf("Failed to encode JSON: %s\n", err)
		return
	}
	file.Close()

	// Write CSV to temporary file
	csvTempFile, err := os.Create(tempCsvFile)
	if err != nil {
		log.Fatalf("Cannot create temp CSV file %q: %s\n", tempCsvFile, err)
		return
	}
	defer csvTempFile.Close()

	csvTempFile.WriteString("Id,Name,FullName,Age,Location,Date,Reason,By,Protection,Method,Status,Source,Image,Url\n")
	for _, incident := range incidents {
		csvTempFile.WriteString(fmt.Sprintf("%d,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n",
			incident.Id,
			incident.Name,
			incident.FullName,
			incident.Age,
			incident.Location,
			incident.Date,
			strings.Join(strings.Split(incident.Reason, ","), " "),
			strings.Join(strings.Split(incident.By, ","), " "),
			strings.Join(strings.Split(incident.Protection, ","), " "),
			strings.Join(strings.Split(incident.Method, ","), " "),
			strings.Join(strings.Split(incident.Status, ","), " "),
			ReplaceAll(ReplaceAll(strings.Join(incident.Source, " "), ",", "%2C", 1), "\n", "", 1),
			incident.Image,
			incident.Url,
		))
	}
	csvTempFile.Close()

	// Validate temporary files before replacing originals
	if !validateFiles(tempJsonFile, tempCsvFile, len(incidents)) {
		log.Println("Validation failed, not updating original files")
		// Clean up temporary files
		os.Remove(tempJsonFile)
		os.Remove(tempCsvFile)
		return
	}

	// If validation passes, replace original files
	if err := os.Rename(tempJsonFile, jsonFileName); err != nil {
		log.Fatalf("Failed to replace JSON file: %s\n", err)
		return
	}

	if err := os.Rename(tempCsvFile, csvFileName); err != nil {
		log.Fatalf("Failed to replace CSV file: %s\n", err)
		return
	}

	fmt.Printf("Successfully updated files with %d incidents\n", len(incidents))
}
