# barcodeserver
simple webserver that creates ean barcodes 
* this project uses "github.com/boombuler/barcode" to do all the work.

Usage
* http://localhost:9000/ean/<ean13>  results in a barcode (single line of 95px) 
* http://localhost:9000/qr/<text> results in a qr code (100px)

