package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// div id="box_26783" -> a href data-img-img="image url" class="box-img-load"
// https://shot.cafe/images/o/for-a-few-dollars-more-1965-310-26783.jpg

// GET request at a url like this: https://shot.cafe/server.php?c=26783
// gives all image and tag information
func main() {
	siteUrl := os.Args[1]
	const (
		host   = "localhost"
		port   = "5432"
		user   = "postgres"
		dbname = "personal"
	)
	conn, err := NewPostgres(user, host, port, dbname)
	if err != nil {
		log.Fatalf("Error : %v", err)
	}
	defer conn.Close(context.Background())

	fmt.Println("Successfully connected!")

	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("shot.cafe"),
	)

	c.OnXML(`//*[@id="main"]/div[1]/div[2]/p[1]/strong/em`, func(e *colly.XMLElement) {
		totalFramesTxt := e.Text
		fmt.Println("Total Frames: ", totalFramesTxt)
		totalFrames, err := strconv.Atoi(totalFramesTxt)
		if err != nil {
			panic(err)
		}
		totalPages := (totalFrames + 19) / 20

		linkParts := strings.Split(siteUrl, "/")
		movieTitle := linkParts[len(linkParts)-1]

		for page := 1; page <= totalPages; page++ {
			pageUrl := fmt.Sprintf("https://shot.cafe/results.php?j=1&tz=movie&movie=%s&page=%d", movieTitle, page)
			fmt.Println("Page URL: ", pageUrl)
			ImagesScraper(conn, pageUrl)
		}

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	fmt.Println("Scraping", siteUrl)

	//TODO:
	//1. Use the new url below for grabbing the image details
	//2. Target the small totes p with number in <em> to get total number of images and calculate number of pages. I believe I got 21 images on one page.
	//https://shot.cafe/results.php?j=1&tz=movie&movie=for-a-few-dollars-more-1965-310&page=1
	//Start scraping on https://shot.cafe
	c.Visit(siteUrl)
}
