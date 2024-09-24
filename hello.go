package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/gocolly/colly"
	"github.com/jackc/pgx/v5"
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
		Zero     []string `json:"0"`
	} `json:"tags"`
}

func saveImages(conn *pgx.Conn, imgUrl string, imageDetails ImageDetails) {
	// Store image to hard drive
	// fullImgUrl := "https://shot.cafe/" + imgUrl
	// fileName := imgUrl[strings.LastIndex(imgUrl, "/")+1:]
	// response, err := resty.New().R().Get(fullImgUrl)
	// if err != nil {
	// 	log.Fatalf("Failed to download image: %v", err)
	// }

	// err = os.WriteFile(fileName, response.Body(), 0644)
	// if err != nil {
	// 	log.Fatalf("Failed to save image: %v", err)
	// }

	// fmt.Printf("Image saved to %s\n", fileName)
	// absPath, err := os.Getwd()
	// if err != nil {
	// 	log.Fatalf("Failed to get current directory: %v", err)
	// }
	// fullFilePath := absPath + "/" + fileName

	//Save movie details to db
	//Use a map to ensure movie hasn't already been saved
	movieId := InsertMovie(conn, imageDetails)

	//Save image details to db
	InsertMovieImage(conn, imageDetails, movieId, "test")
}

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
	// count := 0

	// Add CSS selector for a href elements with the class box-img-load
	c.OnHTML("a.box-img-load[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		imgUrl := e.Attr("data-img-img")
		fmt.Printf("Link found: %q -> %s -> %s\n", e.Text, link, imgUrl)

		// Do a GET request to get all image and tag information
		linkParts := strings.Split(link, "-")
		imgNum := linkParts[len(linkParts)-1]
		// fmt.Println("Img Pars: ", linkParts)
		fmt.Printf("Img Num: %s\n", imgNum)

		var imgDetails ImageDetails

		imgGetUrl := "https://shot.cafe/server.php?c=" + imgNum
		fmt.Println("ImgUrl: ", imgGetUrl)
		response, err := resty.New().R().EnableTrace().SetResult(imgDetails).Get(imgGetUrl)
		if err != nil {
			fmt.Println("Error: ", err)
			log.Fatal(err)
		}
		fmt.Println("Response: ", response)
		saveImages(conn, imgUrl, imgDetails)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	fmt.Println("Scraping", siteUrl)

	//Start scraping on https://shot.cafe
	c.Visit(siteUrl)
}
