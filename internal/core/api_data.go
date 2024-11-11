package core

import (
	"bytes"
	"github.com/immafrady/stock-notifier/utils"
	"github.com/shopspring/decimal"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

type ApiData struct {
	Name            string    // 名字
	Exchange        string    // 市场
	Classify        string    // 分类
	Current         float64   // 当前价
	Opening         float64   // 开盘价
	YesterdayClosed float64   // 昨日收盘价
	Low             float64   // 当日最低
	High            float64   // 当日最高
	Diff            float64   // 当日差价
	Percentage      float64   // 当前涨跌幅
	UpdateAt        time.Time // 更新时间
}

// ParsePrice 格式化价格
func (a ApiData) ParsePrice(v float64) string {
	prec := 2
	if a.Classify == "LOF" || a.Classify == "ETF" {
		prec = 3
	}
	return strconv.FormatFloat(v, 'f', prec, 64)
}

func NewApiData(searchCode string) *ApiData {
	str := utils.Request("http://qt.gtimg.cn/q=" + searchCode)

	// 转换为utf8
	reader := transform.NewReader(bytes.NewReader([]byte(str)), simplifiedchinese.GBK.NewDecoder())
	d, _ := io.ReadAll(reader)
	str = string(d)

	if strings.HasPrefix(str, "v_pv_none_match") {
		return nil
	} else if len(str) > 5 {
		data := &ApiData{
			Exchange: str[2:4],
		}

		str = str[strings.Index(str, `"`)+1 : strings.LastIndex(str, `"`)]
		strs := strings.Split(str, "~")

		data.Name = strs[1]
		data.Opening = parseFloat64(strs[5])
		data.YesterdayClosed = parseFloat64(strs[4])
		data.Current = parseFloat64(strs[3])
		data.High = parseFloat64(strs[33])
		data.Low = parseFloat64(strs[34])
		data.Diff = data.Current - data.YesterdayClosed
		data.Percentage = parseFloat64(strs[32])

		var (
			timeTmpl string
			err      error
		)
		switch data.Exchange {
		case "sh", "sz":
			// 20241107151312
			timeTmpl = "20060102150405"
			data.Classify = strs[61] // 补充一下分类，如果是ETF、LOF做特殊处理
		case "hk":
			// 2024/11/07 14:58:20
			timeTmpl = "2006/01/02 15:04:05"
		case "us":
			// 2024-11-06 15:55:34
			timeTmpl = "2006-01-02 15:04:05"
		default:
			return nil
		}
		data.UpdateAt, err = time.Parse(timeTmpl, strs[30])
		if err != nil {
			log.Printf("[error]日期格式错误 - %v: %s", strs[30], err)
			return nil
		}
		return data
	} else {
		log.Println("[error]接收字符串不合要求：", str)
		return nil
	}
}

func parseFloat64(str string) float64 {
	d, err := decimal.NewFromString(str)
	if err != nil {
		log.Printf("[error]数字转化错误 - %v: %s", str, err)
	}
	res, _ := d.Float64()
	return res
}
