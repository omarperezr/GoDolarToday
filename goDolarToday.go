package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// DolarTodayAPI struct has all the info provided by the Dolar Today API
type DolarTodayAPI struct {
	Antibloqueo struct {
		Mobile                   string `json:"mobile"`
		Video                    string `json:"video"`
		CortoAlternativo         string `json:"corto_alternativo"`
		EnableIads               string `json:"enable_iads"`
		EnableAdmobbanners       string `json:"enable_admobbanners"`
		EnableAdmobinterstitials string `json:"enable_admobinterstitials"`
		Alternativo              string `json:"alternativo"`
		Alternativo2             string `json:"alternativo2"`
		Notifications            string `json:"notifications"`
		ResourceID               string `json:"resource_id"`
	} `json:"_antibloqueo"`
	Labels struct {
		A  string `json:"a"`
		A1 string `json:"a1"`
		B  string `json:"b"`
		C  string `json:"c"`
		D  string `json:"d"`
		E  string `json:"e"`
	} `json:"_labels"`
	Timestamp struct {
		Epoch       string `json:"epoch"`
		Fecha       string `json:"fecha"`
		FechaCorta  string `json:"fecha_corta"`
		FechaCorta2 string `json:"fecha_corta2"`
		FechaNice   string `json:"fecha_nice"`
		Dia         string `json:"dia"`
		DiaCorta    string `json:"dia_corta"`
	} `json:"_timestamp"`
	USD struct {
		Transferencia  float64 `json:"transferencia"`
		TransferCucuta float64 `json:"transfer_cucuta"`
		Efectivo       float64 `json:"efectivo"`
		EfectivoReal   float64 `json:"efectivo_real"`
		EfectivoCucuta float64 `json:"efectivo_cucuta"`
		Promedio       float64 `json:"promedio"`
		PromedioReal   float64 `json:"promedio_real"`
		Cencoex        float64 `json:"cencoex"`
		Sicad1         float64 `json:"sicad1"`
		Sicad2         float64 `json:"sicad2"`
		BitcoinRef     float64 `json:"bitcoin_ref"`
		Dolartoday     float64 `json:"dolartoday"`
	} `json:"USD"`
	EUR struct {
		Transferencia  float64 `json:"transferencia"`
		TransferCucuta float64 `json:"transfer_cucuta"`
		Efectivo       float64 `json:"efectivo"`
		EfectivoReal   float64 `json:"efectivo_real"`
		EfectivoCucuta float64 `json:"efectivo_cucuta"`
		Promedio       float64 `json:"promedio"`
		PromedioReal   float64 `json:"promedio_real"`
		Cencoex        float64 `json:"cencoex"`
		Sicad1         float64 `json:"sicad1"`
		Sicad2         float64 `json:"sicad2"`
		Dolartoday     float64 `json:"dolartoday"`
	} `json:"EUR"`
	COL struct {
		Efectivo float64 `json:"efectivo"`
		Transfer float64 `json:"transfer"`
		Compra   float64 `json:"compra"`
		Venta    float64 `json:"venta"`
	} `json:"COL"`
	GOLD struct {
		Rate float64 `json:"rate"`
	} `json:"GOLD"`
	USDVEF struct {
		Rate float64 `json:"rate"`
	} `json:"USDVEF"`
	USDCOL struct {
		Setfxsell     float64 `json:"setfxsell"`
		Setfxbuy      float64 `json:"setfxbuy"`
		Rate          float64 `json:"rate"`
		Ratecash      float64 `json:"ratecash"`
		Ratetrm       float64 `json:"ratetrm"`
		Trmfactor     float64 `json:"trmfactor"`
		Trmfactorcash float64 `json:"trmfactorcash"`
	} `json:"USDCOL"`
	EURUSD struct {
		Rate float64 `json:"rate"`
	} `json:"EURUSD"`
	BCV struct {
		Fecha     string `json:"fecha"`
		FechaNice string `json:"fecha_nice"`
		Liquidez  string `json:"liquidez"`
		Reservas  string `json:"reservas"`
	} `json:"BCV"`
	MISC struct {
		Petroleo string `json:"petroleo"`
		Reservas string `json:"reservas"`
	} `json:"MISC"`
}

var opciones = map[string]bool{
	os.Args[0]: true,
	"all":      true,
	"dtoday":   true,
	"dcucuta":  true,
	"dbitcoin": true,
	"ddicom":   true,
	"dimpli":   true,
	"etoday":   true,
	"ecucuta":  true,
	"ebitcoin": true,
	"edicom":   true,
	"eimpli":   true,
}

// Dotf turns a regular float and turns it into a human readable number (string)
func Dotf(v float64) string {
	buf := &bytes.Buffer{}
	if v < 0 {
		buf.Write([]byte{'-'})
		v = 0 - v
	}

	comma := []byte{'.'}

	parts := strings.Split(strconv.FormatFloat(v, 'f', -1, 64), ".")
	pos := 0
	if len(parts[0])%3 != 0 {
		pos += len(parts[0]) % 3
		buf.WriteString(parts[0][:pos])
		buf.Write(comma)
	}
	for ; pos < len(parts[0]); pos += 3 {
		buf.WriteString(parts[0][pos : pos+3])
		buf.Write(comma)
	}
	buf.Truncate(buf.Len() - 1)

	if len(parts) > 1 {
		buf.Write([]byte{','})
		buf.WriteString(parts[1])
	}
	return buf.String()
}

func help(x string) {
	if !opciones[x] && x != "help" {
		fmt.Printf("\n%s: no es un argumento valido\n", x)
		fmt.Println("uso: dolarToday [help] [all] [dtoday] [dcucuta] [dbitcoin] [dimpli] [etoday] [ecucuta] [dbitcoin] [eimpli]")
	} else {
		fmt.Println(`
uso: dolarToday [help] [all] [dtoday] [dcucuta] [dbitcoin] [ddicom] [dimpli] [etoday] [ecucuta] [dbitcoin] [ddicom] [eimpli]

argumentos opcionales:
  help	Muestra este mensaje de ayuda
  all	Muestra el precio del Dolar y del Euro de Dolartoday y Cucuta
  dtoday	Muestra el precio del Dolar de DolarToday
  dcucuta	Muestra el precio del Dolar de Cucuta
  dbitcoin	Muestra el precio del Dolar Bitcoin
  ddicom	Muestra el precio del Dolar de Dicom
  dimpli	Muestra el precio del Dolar Implicito (Liquidez Monetaria/Reservas Internacionales) Datos BCV
  etoday	Muestra el precio del Euro de DolarToday
  ecucuta	Muestra el precio del Euro de Cucuta
  ebitcoin	Muestra el precio del Euro Bitcoin
  edicom	Muestra el precio del Euro de Dicom
  eimpli	Muestra el precio del Euro Implicito (Liquidez Monetaria/Reservas Internacionales) Datos BCV`)
	}
}

func main() {
	var (
		args   = os.Args
		leArg  = args[0]
		argLen = len(args)
	)

	if argLen > 1 {
		leArg = args[1]
	}

	if opciones[leArg] && argLen < 3 {
		response, err := http.Get("https://s3.amazonaws.com/dolartoday/data.json")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var responseObject DolarTodayAPI
		json.Unmarshal(responseData, &responseObject)

		if argLen == 1 {
			var temp string
			fmt.Printf("\nDOLAR:\nToday: %sbsf\nCucuta: %sbsf\nBitcoin: %sbsf\nDicom: %s\n\nEURO:\nToday: %sbsf\nCucuta: %sbsf\nBitcoin: %s\nDicom: %s\n\nPresione Enter para salir...",
				Dotf(responseObject.USD.Transferencia), Dotf(responseObject.USD.EfectivoCucuta), Dotf(responseObject.USD.BitcoinRef), Dotf(responseObject.USD.PromedioReal),
				Dotf(responseObject.EUR.Transferencia), Dotf(responseObject.EUR.EfectivoCucuta), Dotf(responseObject.USD.BitcoinRef*responseObject.EURUSD.Rate), Dotf(responseObject.EUR.PromedioReal))

			fmt.Scanln(&temp)
		} else {
			switch leArg {
			case "all":
				fmt.Printf("\nDOLAR:\nToday: %sbsf\nCucuta: %sbsf\nBitcoin: %sbsf\nDicom: %sbsf\nImplicito: %s\n\nEURO:\nToday: %sbsf\nCucuta: %sbsf\nBitcoin: %sbsf\nDicom: %s\nImplicito: %sbsf\n",
					Dotf(responseObject.USD.Transferencia), Dotf(responseObject.USD.EfectivoCucuta), Dotf(responseObject.USD.BitcoinRef), Dotf(responseObject.USD.PromedioReal), Dotf(responseObject.USD.Efectivo),
					Dotf(responseObject.EUR.Transferencia), Dotf(responseObject.EUR.EfectivoCucuta), Dotf(responseObject.USD.BitcoinRef*responseObject.EURUSD.Rate), Dotf(responseObject.EUR.PromedioReal), Dotf(responseObject.EUR.Efectivo))
			case "dtoday":
				fmt.Printf("\nDolar Today: %sbsf\n\n", Dotf(responseObject.USD.Transferencia))
			case "dcucuta":
				fmt.Printf("\nDolar Cucuta: %sbsf\n\n", Dotf(responseObject.USD.EfectivoCucuta))
			case "dbitcoin":
				fmt.Printf("\nDolar Bitcoin: %sbsf\n\n", Dotf(responseObject.USD.BitcoinRef))
			case "ddicom":
				fmt.Printf("\nDolar Dicom: %sbsf\n\n", Dotf(responseObject.USD.PromedioReal))
			case "dimpli":
				fmt.Printf("\nDolar Implicito: %sbsf\n\n", Dotf(responseObject.USD.Efectivo))
			case "etoday":
				fmt.Printf("\nEuro Today: %sbsf\n\n", Dotf(responseObject.EUR.Transferencia))
			case "ecucuta":
				fmt.Printf("\nEuro Cucuta: %sbsf\n\n", Dotf(responseObject.EUR.EfectivoCucuta))
			case "ebitcoin":
				fmt.Printf("\nEuro Bitcoin: %sbsf\n\n", Dotf(responseObject.USD.BitcoinRef*responseObject.EURUSD.Rate))
			case "edicom":
				fmt.Printf("\nEuro Dicom: %sbsf\n\n", Dotf(responseObject.EUR.PromedioReal))
			case "eimpli":
				fmt.Printf("\nEuro Implicito: %sbsf\n\n", Dotf(responseObject.EUR.Efectivo))
			}
		}
	} else {
		help(leArg)
	}
}
