package main

import (
	"flag"
	"fmt"
	"gotube"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cheggaaa/pb"
)

//var wg = &sync.WaitGroup{}
//var ch = make(chan bool, 5)

func main() {
	log.SetFlags(log.Lshortfile)
	u := flag.String("u", "https://www.youtube.com/watch?v=QS7lN7giXXc", "youtube url")
	fn := flag.String("f", "", "downloadfilename")
	flag.Parse()

	url := *u
	url = strings.Replace(url, `https://`, `http://`, -1)

	client := gotube.NewClient()
	info, err := client.GetInfo(url)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(info.Title)

	for i, v := range info.AvailableDownloadInfo {
		log.Println(i, v.Type)
	}
	var selection int
	fmt.Print("Please select download type:")
	fmt.Scanf("%d", &selection)
	fmt.Print("You typed:", selection, "\nAre you sure(y/n):")
	var sure string
	fmt.Scan(&sure)
	if sure != "y" {
		log.Fatal("Bye.")
	}
	filename := ""
	if *fn == "" {
		filename = info.Title + info.AvailableDownloadInfo[selection].FileExtension
	} else {
		filename = *fn + info.AvailableDownloadInfo[selection].FileExtension
	}
	download(client, info.AvailableDownloadInfo[selection].Url, *fn+"."+filename)

}

func download(c *gotube.Client, url, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	resp, err := c.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	log.Println("size:", resp.ContentLength)
	mycopy(f, resp.Body, resp.ContentLength)

}

func mycopy(dst io.Writer, src io.Reader, total int64) {
	s := 0
	b := make([]byte, 4096)
	//start := time.Now()
	stop := false
	bar := pb.StartNew(int(total))

	bar.ShowSpeed = true

	go func() {
		for !stop {
			time.Sleep(time.Second)
			bar.Set(s)

			//log.Printf("%d/%d  %dk/s\n", s, total, int(float64(s)/1024.0/time.Now().Sub(start).Seconds()))
		}

	}()
	defer func() {
		stop = true
		if int64(s) == total {
			bar.FinishPrint("Finished!")
		} else {
			bar.FinishPrint("something wrong!")
		}
	}()
	for {
		n, err := src.Read(b)
		if n > 0 {
			s += n
			_, err := dst.Write(b[:n])
			if err != nil {
				//log.Println(err)
				return
			}
		}
		if err != nil {
			//log.Println(err)
			return
		}
	}

}
