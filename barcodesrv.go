package main

import (
	"log"
	"net/http"

	"bytes"
	"image/png"
	"strconv"

	"net/url"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/ean"
	"github.com/boombuler/barcode/qr"
)

func eanHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("serving " + r.URL.Path)
	var s string
	if s = r.URL.Path[len("/ean/"):]; len(s) != 13 {
		http.Error(w, s+"is not a correct ean13", http.StatusBadRequest)
		return
	}
	var barcode13 barcode.BarcodeIntCS
	var err error
	if barcode13, err = ean.Encode(s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, barcode13); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	w.Header().Set("ETag", s)
	w.Header().Set("Cache-Control", "max-age=365")

	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("problem writing image to ResponseWriter")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func qrHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("serving " + r.URL.Path)
	var s, q string
	var err error
	s = r.URL.Path[len("/qr/"):]
	if q, err = url.QueryUnescape(s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var qrcode barcode.Barcode
	if qrcode, err = qr.Encode(q, qr.L, qr.Auto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if qrcode, err = barcode.Scale(qrcode, 100, 100); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, qrcode); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	w.Header().Set("ETag", s)
	w.Header().Set("Cache-Control", "max-age=365")

	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("problem writing image to ResponseWriter")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func customNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("notFoundHandler: " + r.URL.Path)
	http.Error(w, r.URL.Path, http.StatusNotFound)
}

func main() {

	log.Println("Starting up barcode webserver at port 9000")
	http.HandleFunc("/ean/", eanHandler)
	http.HandleFunc("/qr/", qrHandler)
	http.HandleFunc("/", customNotFoundHandler)
	http.ListenAndServe(":9000", nil)
}
