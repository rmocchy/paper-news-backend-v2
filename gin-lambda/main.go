package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	controller "github.com/rmocchy/paper-news-backend-v2/gin-lambda/controller"
	paperRepo "github.com/rmocchy/paper-news-backend-v2/gin-lambda/repository/paper-repository"
	paperServ "github.com/rmocchy/paper-news-backend-v2/gin-lambda/service/paper-service"
)

var ginLambda *ginadapter.GinLambdaV2

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	if ginLambda == nil {
		log.Printf("Gin cold start")

		hs, err := buildHandlers()
		if err != nil {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       "Failed to build handlers",
			}, err
		}

		r := gin.Default()
		r.GET("/greet", sayHello)
		r.GET("/papers", hs.paperController.ListPapers)
		r.POST("/papers/reflesh", hs.paperController.RefleshPapers)
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			c.JSON(http.StatusNotFound, gin.H{"message": "Not Method Found at " + path + "req was " + req.RawPath})
		})

		ginLambda = ginadapter.NewV2(r)
	}

	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}

func sayHello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
}

type MyHandlers struct {
	paperController controller.PaperController
}

func buildHandlers() (*MyHandlers, error) {
	db, err := sqlx.Open("mysql", "user:password@tcp(local_db)/paper-news_local?parseTime=true")
	if err != nil {
		return nil, err
	}

	// paper service
	paperRepository := paperRepo.NewPaperRepository(db)
	arxivClient := paperRepo.NewArxivFetchClient()
	refleshPaperService := paperServ.NewRefleshPaperService(paperRepository, arxivClient)
	listPaperService := paperServ.NewListPaperService(paperRepository)
	paperController := controller.NewPaperController(refleshPaperService, listPaperService)

	return &MyHandlers{
		paperController: paperController,
	}, nil
}
