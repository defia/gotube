package gotube

import (
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	*http.Client
}

func NewClient() *Client {
	tp := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}
	//	jar, err := cookiejar.New(nil)
	//	if err != nil {
	//		log.Fatal(err)
	//		return nil
	//	}
	return &Client{&http.Client{
		Transport: tp,
		//		Jar:       jar,
	}}
}

//get video's raw info
func (c *Client) GetRawInfo(url string) (*VideoRawInfo, error) {
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Println("read finished")
	return GenerateVideoInfo(b)
}

//get info for common format
func (c *Client) GetInfo(url string) (*VideoInfo, error) {
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	v, err := GenerateVideoInfo(b)
	if err != nil {
		return nil, err
	}
	return v.toVideoInfo(), nil
}
