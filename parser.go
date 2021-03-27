package parser

import (
	"encoding/xml"
	"fmt"
	"io"
)

type Share struct {
	ID        string
	Name      string
	ShortName string
	LatName   string
	Status    string
	LotSize   int
}

type Trade struct {
	ID string

	PriceOpen float64
	PriceLow  float64
	PriceHigh float64
	PriceLast float64

	WaPrice   float64
	NumTrades int64
	NumOffers int64
	NumBids   int64
	VolToday  int64
}

type document struct {
	XMLName xml.Name `xml:"document"`
	Datas   []data   `xml:"data"`
}

type data struct {
	XMLName xml.Name `xml:"data"`
	ID      string   `xml:"id,attr"`
	Rows    rows     `xml:"rows"`
}

type rows struct {
	XMLName xml.Name `xml:"rows"`
	Rows    []row    `xml:"row"`
}

type row struct {
	XMLName   xml.Name `xml:"row"`
	ID        string   `xml:"SECID,attr"`
	Name      string   `xml:"SECNAME,attr"`
	ShortName string   `xml:"SHORTNAME,attr"`
	LotSize   int      `xml:"LOTSIZE,attr"`
	Status    string   `xml:"STATUS,attr"`
	LatName   string   `xml:"LATNAME,attr"`

	PriceOpen float64 `xml:"OPEN,attr"` // цена открытия
	PriceLow  float64 `xml:"LOW,attr"`  // минимальная цена сделки
	PriceHigh float64 `xml:"HIGH,attr"` // максимальная цена сделки
	PriceLast float64 `xml:"LAST,attr"` // цена последней сделки

	WaPrice   float64 `xml:"WAPRICE,attr"`   // средневзвешенная цена
	NumTrades int64   `xml:"NUMTRADES,attr"` // число сделок
	NumOffers int64   `xml:"NUMOFFERS,attr"` // число заявок на продажу
	NumBids   int64   `xml:"NUMBIDS,attr"`   // число заявок на покупку
	VolToday  int64   `xml:"VOLTODAY,attr"`  // число инструментов в сделках
	QTY       int64   `xml:"QTY,attr"`       // число инструментов в последней сделке
}

func parse(in io.Reader) {

	var doc document

	xml.NewDecoder(in).Decode(&doc)

	shares := map[string]*Share{}
	trades := map[string]*Trade{}

	for _, d := range doc.Datas {
		if d.ID == "securities" {
			for _, row := range d.Rows.Rows {
				shares[row.ID] = &Share{
					ID:        row.ID,
					Name:      row.Name,
					LatName:   row.LatName,
					ShortName: row.ShortName,
					LotSize:   row.LotSize,
					Status:    row.Status,
				}
			}
		} else if d.ID == "marketdata" {
			for _, row := range d.Rows.Rows {
				trades[row.ID] = &Trade{
					ID:        row.ID,
					PriceOpen: row.PriceOpen,
					PriceLow:  row.PriceLow,
					PriceHigh: row.PriceHigh,
					PriceLast: row.PriceLast,
					WaPrice:   row.WaPrice,
					NumTrades: row.NumTrades,
					NumOffers: row.NumOffers,
					NumBids:   row.NumBids,
					VolToday:  row.VolToday,
				}
			}
		}
	}

	fmt.Println(shares["SBERP"])

	fmt.Println(trades["SBERP"])
}
