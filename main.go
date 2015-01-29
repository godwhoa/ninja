package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"

	"net/http"
	"os"
)

const domain = "http://104.131.80.165:8080/u/"

//const domain = "http://localhost:8080/u/"

func GetVal(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("urlInput")
	var alias = r.FormValue("alias")
	if Exist(alias) {
		alias = randSeq()
		err := ioutil.WriteFile("urls/"+alias, []byte(url), 0644)
		perror(err)
		fmt.Fprintf(w, "Alias has been taken so you get a random one: "+domain+alias)
	}
	if alias == "val" {
		fmt.Fprintf(w, "/val is the post url so can't use it.")
		return
	}
	err := ioutil.WriteFile("urls/"+alias, []byte(url), 0644)
	perror(err)
	fmt.Fprintf(w, domain+alias)

}
func short(w http.ResponseWriter, r *http.Request) {
	parm := r.URL.String()[3:]

	url, err := ioutil.ReadFile("urls/" + parm)
	perror(err)
	http.Redirect(w, r, string(url), 301)
	return
}

func main() {
	http.HandleFunc("/val", GetVal)
	http.HandleFunc("/u/", short)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.ListenAndServe(":8080", nil)
}

//Error handling
func perror(err error) {
	if err != nil {
		panic(err)
	}
}

//File exists
func Exist(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}

	return true
}

//Random string
//by Imgur Team
func randSeq() string {
	hashGetter := make(chan string)
	length := 7

	go func() {
		for {
			str := ""

			for len(str) < length {
				c := 10
				bArr := make([]byte, c)
				_, err := rand.Read(bArr)
				if err != nil {
					log.Println("error:", err)
					break
				}

				for _, b := range bArr {
					if len(str) == length {
						break
					}

					/**
					 * Each byte will be in [0, 256), but we only care about:
					 *
					 * [48, 57]     0-9
					 * [65, 90]     A-Z
					 * [97, 122]    a-z
					 *
					 * Which means that the highest bit will always be zero, since the last byte with high bit
					 * zero is 01111111 = 127 which is higher than 122. Lower our odds of having to re-roll a byte by
					 * dividing by two (right bit shift of 1).
					 */

					b = b >> 1

					// The byte is any of        0-9                  A-Z                      a-z
					byteIsAllowable := (b >= 48 && b <= 57) || (b >= 65 && b <= 90) || (b >= 97 && b <= 122)

					if byteIsAllowable {
						str += string(b)
					}
				}

			}

			hashGetter <- str
		}
	}()

	return <-hashGetter
}
