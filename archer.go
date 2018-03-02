package archer

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Archer - inteface to http.client
type Archer interface {
	Get(url string) ([]byte, error)
	Post(url string, body io.Reader) ([]byte, error)
	Put(url string) ([]byte, error)
	Delete(url string) ([]byte, error)
	Upload(url string, file string) ([]byte, error)
}

type archer struct {
	c *http.Client
}

// NewArcher - create a new requester which provides http.client
func NewArcher(timeout time.Duration) Archer {
	c := &http.Client{
		Timeout: timeout,
	}
	return &archer{
		c: c,
	}
}

func (a *archer) Get(url string) ([]byte, error) {
	return a.call(http.MethodGet, url, "", nil)
}

func (a *archer) Post(url string, body io.Reader) ([]byte, error) {
	return a.call(http.MethodPost, url, "", body)
}

func (a *archer) Put(url string) ([]byte, error) {
	// #TODO add body for put
	return a.call(http.MethodPut, url, "", nil)
}

func (a *archer) Delete(url string) ([]byte, error) {
	return a.call(http.MethodDelete, url, "", nil)
}

func (a *archer) Upload(url string, file string) ([]byte, error) {
	return a.call(http.MethodPut, url, file, nil)
}

func (a *archer) call(method string, url string, filePath string, body io.Reader) ([]byte, error) {

	var request *http.Request
	var err error

	if filePath != "" {
		// Trying to upload a file
		file, err := os.Open(filePath)
		if err != nil {
			return []byte{}, err
		}
		defer file.Close()

		pr, pw := io.Pipe()
		bufin := bufio.NewReader(file)

		go func() {
			_, err := bufin.WriteTo(pw)
			if err != nil {
				fmt.Println(err)
			}
			pw.Close()
		}()

		request, err = http.NewRequest(method, url, pr)
	} else if body != nil {
		request, err = http.NewRequest(method, url, body)
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	} else {
		request, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return []byte{}, err
	}

	response, err := a.c.Do(request)
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}
