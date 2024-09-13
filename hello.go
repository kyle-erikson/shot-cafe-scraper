package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/gocolly/colly"
)

// div id="box_26783" -> a href data-img-img="image url" class="box-img-load"
// https://shot.cafe/images/o/for-a-few-dollars-more-1965-310-26783.jpg

// GET request at a url like this: https://shot.cafe/server.php?c=26783
// gives all image and tag information

type ImageDetails struct {
	Image struct {
		Image   string `json:"image"`
		Project string `json:"project"`
		Date    string `json:"date"`
		Slug    string `json:"slug"`
		Hits    string `json:"hits"`
		// T_H              string   `json:"t_h"`
		// T_W              string   `json:"t_w"`
		Alt            string `json:"alt"`
		Wfm            string `json:"wfm"`
		ImTags         string `json:"im_tags"`
		QTag           string `json:"q_tag"`
		TagSubmissions string `json:"tag_submissions"`
		Type           string `json:"type"`
		Anb            string `json:"anb"`
		Title          string `json:"title"`
		Year           string `json:"year"`
		Aspect         string `json:"aspect"`
		Camera         string `json:"camera"`
		Lens           string `json:"lens"`
		D_ID           string `json:"d_id"`
		Director       string `json:"director"`
		DP_ID          string `json:"dp_id"`
		DP             string `json:"dp"`
		PSlug          string `json:"pslug"`
		Fave           string `json:"fave"`
		Coll           string `json:"coll"`
		TypeText       string `json:"typetext"`
		TypeLink       string `json:"typelink"`
	} `json:"image"`
	Tags struct {
		Two      []string `json:"2"`
		Three    []string `json:"3"`
		MinusOne []string `json:"-1"`
		One      []string `json:"1"`
	} `json:"tags"`
}

func main() {
	siteUrl := os.Args[1]
	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("shot.cafe"),
	)
	// count := 0

	// Add CSS selector for a href elements with the class box-img-load
	c.OnHTML("a.box-img-load[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		imgUrl := e.Attr("data-img-img")
		fmt.Printf("Link found: %q -> %s -> %s\n", e.Text, link, imgUrl)

		// if count > 1 {
		// 	return
		// }

		// Do a GET request to get all image and tag information
		linkParts := strings.Split(link, "-")
		imgNum := linkParts[len(linkParts)-1]
		// fmt.Println("Img Pars: ", linkParts)
		fmt.Printf("Img Num: %s\n", imgNum)

		var imgDetails ImageDetails

		imgGetUrl := "https://shot.cafe/server.php?c=" + imgNum
		fmt.Println("ImgUrl: ", imgGetUrl)
		response, err := resty.New().R().EnableTrace().SetResult(&imgDetails).Get(imgGetUrl)
		if err != nil {
			fmt.Println("Error: ", err)
			log.Fatal(err)
		}
		fmt.Println("Response: ", response)
		// count++
	})

	// On every a element which has href attribute call callback
	// c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	// 	link := e.Attr("href")
	// 	// Print link
	// 	fmt.Printf("Link found: %q -> %s\n", e.Text, link)
	// 	// Visit link found on page
	// 	// Only those links are visited which are in AllowedDomains
	// 	c.Visit(e.Request.AbsoluteURL(link))
	// })

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	fmt.Println("Scraping", siteUrl)

	//Start scraping on https://shot.cafe
	c.Visit(siteUrl)
}
