package data

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/xuri/excelize/v2"
	"strconv"
)

type member struct {
	Name string
	CF   string
	Room string
	Seat string
}
type Accept struct {
	CF      string
	Problem string
}

type volunteer struct {
	It   int
	Info []Info
}
type Info struct {
	QQ   int64
	Name string
}

var (
	Member    map[string]member
	Problem   map[Accept]bool
	Volunteer map[string]volunteer
)

func LoadData() {
	f, err := excelize.OpenFile("data.xlsx")
	if err != nil {
		log.Error(err.Error())
	}
	rows, err := f.GetRows("candidate")
	if err != nil {
		log.Error(err.Error())
	}

	Member = make(map[string]member)
	for _, row := range rows {
		Member[row[1]] = member{
			Name: row[0],
			CF:   row[1],
			Room: row[2],
			Seat: row[3],
		}
	}

	rows, err = f.GetRows("volunteer")
	if err != nil {
		log.Error(err.Error())
	}

	Volunteer = make(map[string]volunteer)
	for _, row := range rows {
		_, vis := Volunteer[row[2]]
		if !vis {
			Volunteer[row[2]] = volunteer{
				It: 0,
			}
		}
		v := Volunteer[row[2]]
		qq, err := strconv.ParseInt(row[1], 10, 64)
		if err != nil {
			log.Error(err.Error())
		}
		v.Info = append(v.Info, Info{
			Name: row[0],
			QQ:   qq,
		})
		Volunteer[row[2]] = v
	}
	fmt.Println(Volunteer)

	Problem = make(map[Accept]bool)
}
