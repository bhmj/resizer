package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/bhmj/resizer/storage"
	"github.com/disintegration/imaging"
)

const port = ":8080"

var alive = true // must use mutex or atomic but hey..

func sayBadRequest(w http.ResponseWriter, s string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(s))
}

func sayNotFound(w http.ResponseWriter, s string) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(s))
}

func sayInternal(w http.ResponseWriter, s string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(s))
}

func respond(w http.ResponseWriter, encoded *bytes.Buffer) {
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(encoded.Len()))
	w.Header().Set("Cache-Control", "max-age=3600")
	_, err := io.Copy(w, encoded)
	if err != nil {
		sayInternal(w, err.Error())
		return
	}
}

// MakeHandler ...
func MakeHandler(client *http.Client, cache storage.Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		r.ParseForm()
		url := r.Form["url"]
		width := r.Form["width"]
		height := r.Form["height"]

		// some checks
		if len(url) == 0 || len(width) == 0 || len(height) == 0 {
			sayBadRequest(w, "url, width, height are required")
			return
		}
		if url[0] == "" {
			sayBadRequest(w, "no url specified")
			return
		}

		iwidth, errW := strconv.Atoi(width[0])
		iheight, errH := strconv.Atoi(height[0])
		if errW != nil || errH != nil {
			sayBadRequest(w, "bad width or height")
			return
		}

		fullKey := url[0] + "|" + width[0] + "|" + height[0]

		// try to get cached result
		item, found := cache.Get(fullKey)
		if found {
			fmt.Println("cache level 1 hit")
			encoded := item.(bytes.Buffer)
			respond(w, &encoded)
			return
		}

		// try to get cached source image
		var img image.Image
		im, found := cache.Get(url[0])
		if found {
			fmt.Println("cache level 2 hit")
			img = im.(image.Image)
		} else {
			resp, err := client.Get(url[0])
			if err != nil {
				sayNotFound(w, err.Error())
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				sayNotFound(w, "remote image not received")
				return
			}
			img, _, err = image.Decode(resp.Body)
			if err != nil {
				sayInternal(w, err.Error())
				return
			}
			cache.Put(url[0], img, time.Duration(1*time.Hour))
		}

		dst := imaging.Resize(img, iwidth, iheight, imaging.Lanczos)

		encoded := &bytes.Buffer{}
		err = jpeg.Encode(encoded, dst, nil)
		if err != nil {
			sayInternal(w, err.Error())
			return
		}
		cache.Put(fullKey, *encoded, time.Duration(1*time.Hour))

		respond(w, encoded)
	}
}

// readiness probe (k8s template, do you use it?)
func readiness(w http.ResponseWriter, _ *http.Request) {
	if alive {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func main() {

	fmt.Println("setting up server on", port)

	errs := make(chan error)
	shutdown := make(chan bool)

	cache := storage.NewStore()

	mux := http.NewServeMux()
	cli := &http.Client{
		Timeout: 1 * time.Second,
	}
	mux.HandleFunc("/api/v1/resizer/", MakeHandler(cli, cache))
	mux.HandleFunc("/readiness", readiness)
	srv := http.Server{Addr: port, Handler: mux}
	go func() {
		errs <- srv.ListenAndServe()
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		fmt.Printf("\nshutting down: %s\n", <-errs)
		srv.SetKeepAlivesEnabled(false)
		alive = false
		time.Sleep(2 * time.Second)
		shutdown <- true
	}()

	fmt.Println("started")
	fmt.Println("usage: http://{hostname}:8080/api/v1/resizer/?url={url}&width={width}&height={height}")
	<-shutdown

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() { cancel() }()
	if err := srv.Shutdown(ctxShutDown); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
	fmt.Println("terminated")
}
