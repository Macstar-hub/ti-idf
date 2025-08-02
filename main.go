package main

import (
	httppost "tf-idf/cmd/api"

	"github.com/gin-gonic/gin"
)

const dockPath string = "./test/sampleFile-2.txt"
const stopWordFilePath string = "./configs/stopWords.txt"

func main() {
	LandigPage()
}

func LandigPage() {
	server := gin.Default()
	server.LoadHTMLGlob("./web/**/*")
	server.StaticFile("/", "./web/landingPages/linkSubmit.html")
	server.StaticFile("/assetcalc", "./web/landingPages/assetCalc.html")
	server.StaticFile("/uploadfile", "./web/landingPages/upload.html")
	server.POST("/api/v1/postLinks", httppost.PostLabels)
	server.POST("/api/v1/assetcalc", httppost.CalcAsset)
	server.GET("/api/v1/linkslist", httppost.ShowLinks)
	server.POST("/api/v1/search", httppost.Search)
	server.POST("/api/v1/uploadfile", httppost.UploadFile)
	server.Run(":9080")
}
