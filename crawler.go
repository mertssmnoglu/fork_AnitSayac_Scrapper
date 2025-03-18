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
	baseUrl = "https://anitsayac.com"
)

func ReplaceAll(s, old, new string, n int) string {
	// Replace All for Regexp
	re := regexp.MustCompile(old)
	return re.ReplaceAllString(s, new)
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
	fName := "data.json"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()

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

	incidents := make([]Incident, 0, 20000)

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

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	// Dump json to the standard output
	enc.Encode(incidents)

	// CSV conversion
	csvFile, err := os.Create("data.csv")
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", "data.csv", err)
		return
	}
	defer csvFile.Close()

	csvFile.WriteString("Id,Name,FullName,Age,Location,Date,Reason,By,Protection,Method,Status,Source,Image,Url\n")
	for _, incident := range incidents {
		csvFile.WriteString(fmt.Sprintf("%d,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n",
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
}
