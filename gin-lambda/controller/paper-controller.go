package controller

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	paperservice "github.com/rmocchy/paper-news-backend-v2/gin-lambda/service/paper-service"
)

type PaperController interface {
	RefleshPapers(c *gin.Context)
	ListPapers(c *gin.Context)
}

type paperController struct {
	refleshPaperService paperservice.RefleshPaperService
	listPaperService    paperservice.ListPaperService
}

func NewPaperController(
	refleshPaperService paperservice.RefleshPaperService,
	listPaperService paperservice.ListPaperService,
) PaperController {
	return &paperController{
		refleshPaperService: refleshPaperService,
		listPaperService:    listPaperService,
	}
}

type RefleshPapersRrequest struct {
	MaxSize int `json:"max_size"`
}

func (p *paperController) RefleshPapers(c *gin.Context) {
	// todo: validate by openapi

	// get PapersMaxSize
	var request RefleshPapersRrequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("error to get maxsize: %w", err.Error())
		request.MaxSize = 20 // default
	}
	if request.MaxSize > 100 {
		log.Println("cannnot get more than 100 papers")
		request.MaxSize = 100 // max
	}

	res, err := p.refleshPaperService.RefleshPaper(c, &paperservice.RefleshPaperInput{
		PapersMaxSize: int64(request.MaxSize),
	})
	if err != nil {
		log.Println("error at RefleshPapers: %w", err.Error())
		c.JSON(500, gin.H{"message": fmt.Sprintln("RefleshFailed: %w", err.Error())})
		return
	}

	c.JSON(200, gin.H{"message": fmt.Sprint("RefleshSuccess: refleshed newest %d papers", res.RefleshedPaperSize)})
}

func (p *paperController) ListPapers(c *gin.Context) {
	// todo: validate by openapi

	result, err := p.listPaperService.ListPapers(c, &paperservice.ListPapersInput{})
	if err != nil {
		c.JSON(500, gin.H{"message": fmt.Sprintln("ListFailed: %w", err.Error())})
		return
	}
	c.JSON(200, gin.H{"message": "ListSuccess", "result": result.Papers})
}
