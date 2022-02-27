package downloader

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
)

type Bw struct {
	Name     string
	Bwr      int64
	Fb       int64
	Paused   bool
	Finished bool
	Cancel   bool
	Error    error
}

func DownloadFile(filepath string, url string, out chan Bw) {
	defer close(out)
	head, err := http.Head(url)
	if err != nil {
		out <- Bw{Name: filepath, Bwr: -1, Fb: -1, Paused: false, Finished: true, Cancel: false, Error: err}
		return
	}
	size := head.ContentLength
	resp, err := http.Get(url)
	fmt.Println("downloading from " + url)
	if err != nil {
		out <- Bw{Name: filepath, Bwr: -1, Fb: -1, Paused: false, Finished: true, Cancel: false, Error: err}
		return
	}
	defer resp.Body.Close()

	if filepath == "" {
		fl := strings.Split(resp.Header.Get("content-disposition"), "filename=")
		if len(fl) >= 2 {
			filepath = fl[1]
		} else {
			filepath = uuid.NewString()
		}
	}

	of, err := os.Create(filepath)
	if err != nil {
		out <- Bw{Name: filepath, Bwr: -1, Fb: -1, Paused: false, Finished: true, Cancel: false, Error: err}
		return
	}
	defer of.Close()

	buffer := make([]byte, 1024)
	var totalbytes int64 = 0
	for {
		select {
		case command := <-out:
			if command.Paused {
				for {
					command = <-out
					if !command.Paused {
						break
					}
				}
			} else if command.Cancel {
				out <- Bw{Name: filepath, Bwr: totalbytes, Fb: size, Paused: false, Finished: false, Cancel: true, Error: nil}
				os.Remove(filepath)
				return
			}
		default:
		}
		read, err := resp.Body.Read(buffer)
		totalbytes += int64(read)
		if read == 0 || err != nil {
			out <- Bw{Name: filepath, Bwr: totalbytes, Fb: size, Paused: false, Finished: true, Cancel: false, Error: nil}
			break
		}
		_, err = of.Write(buffer)
		if err != nil {
			out <- Bw{Name: filepath, Bwr: totalbytes, Fb: size, Paused: false, Finished: false, Cancel: false, Error: err}
			break
		}
		out <- Bw{Name: filepath, Bwr: totalbytes, Fb: size, Paused: false, Finished: false, Cancel: false, Error: err}
	}
}
