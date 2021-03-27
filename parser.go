package parser

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"time"

	"github.com/wmentor/ua"
)

var (
	ErrFetchFailed  error = errors.New("fetch data error")
	ErrDecodeFailed error = errors.New("decode data failed")
)

type (
	Shares map[string]*Share
	Trades map[string]*Trade
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

func parse(in io.Reader) (Shares, Trades, error) {

	var doc document

	if err := xml.NewDecoder(in).Decode(&doc); err != nil {
		return nil, nil, ErrDecodeFailed
	}

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

	return shares, trades, nil
}

func Get() (Shares, Trades, error) {
	client := ua.New()

	client.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"
	client.Timeout = time.Second * 10

	resp, err := client.Request("GET", "https://iss.moex.com/iss/engines/stock/markets/shares/boards/TQBR/securities.xml", nil, nil)

	if err != nil {
		return nil, nil, ErrFetchFailed
	}

	if resp == nil || resp.StatusCode != 200 {
		return nil, nil, ErrFetchFailed
	}

	return parse(bytes.NewReader(resp.Content))
}
