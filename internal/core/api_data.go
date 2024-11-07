package core

import (
	"bytes"
	"fmt"
	"github.com/immafrady/stock-notifier/utils"
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
	Current         float32   // 当前价
	Opening         float32   // 开盘价
	YesterdayClosed float32   // 昨日收盘价
	Low             float32   // 当日最低
	High            float32   // 当日最高
	Diff            float32   // 当日差价
	Percentage      float32   // 当前涨跌幅
	UpdateAt        time.Time // 更新时间
}

func NewApiData(searchCode string) *ApiData {
	str := utils.Request("http://qt.gtimg.cn/q=" + searchCode)

	// 转换为utf8
	reader := transform.NewReader(bytes.NewReader([]byte(str)), simplifiedchinese.GBK.NewDecoder())
	d, _ := io.ReadAll(reader)
	str = string(d)
	fmt.Println(str)

	if strings.HasPrefix(str, "v_pv_none_match") {
		return nil
	} else {
		data := &ApiData{
			Exchange: str[2:4],
		}

		str = str[strings.Index(str, `"`)+1 : strings.LastIndex(str, `"`)]
		strs := strings.Split(str, "~")

		data.Name = strs[1]
		data.Opening = parseFloat32(strs[5])
		data.YesterdayClosed = parseFloat32(strs[4])
		data.Current = parseFloat32(strs[3])
		data.High = parseFloat32(strs[33])
		data.Low = parseFloat32(strs[34])
		data.Diff = data.Current - data.YesterdayClosed
		data.Percentage = parseFloat32(strs[32])

		var (
			timeTmpl string
			err      error
		)
		switch data.Exchange {
		case "sh", "sz":
			// 20241107151312
			timeTmpl = "20060102150405"
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
	}

}

func parseFloat32(str string) float32 {
	f64, err := strconv.ParseFloat(str, 32)
	if err != nil {
		log.Printf("[error]数字转化错误 - %v: %s", str, err)
	}
	return float32(f64)
}
