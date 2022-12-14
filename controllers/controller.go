package controllers

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type pagingResult struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	PrevPage  int `json:"prevPage"`
	NextPage  int `json:"nextPage"`
	Count     int `json:"count"`
	TotalPage int `json:"totalPage"`
}
type pagination struct {
	ctx     *gin.Context
	query   *gorm.DB
	records interface{}
}

func (p *pagination) paginate() *pagingResult {
	limit, _ := strconv.Atoi(p.ctx.DefaultQuery("limit", "12"))
	page, _ := strconv.Atoi(p.ctx.DefaultQuery("page", "1"))

	ch := make(chan int)
	go p.countRecores(ch)

	offset := limit * (page - 1) //offset เริ่มมองหาตั้งแต่ตัวที่เท่าไหร่่
	p.query.Limit(limit).Offset(offset).Find(p.records)

	count := <-ch
	totalPage := int(math.Ceil(float64(count) / float64(limit)))

	var nextPage int
	if page == totalPage {
		nextPage = totalPage
	} else {
		nextPage = page - 1
	}

	return &pagingResult{
		Page:      page,
		Limit:     limit,
		Count:     count,
		PrevPage:  (page - 1),
		NextPage:  nextPage,
		TotalPage: totalPage,
	}

}

func (p *pagination) countRecores(ch chan int) {
	var count int
	p.query.Model(p.records).Count(&count)

	ch <- count
}
