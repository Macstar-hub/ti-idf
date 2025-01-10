package httppost

import (
	"net/http"

	mysqlconnector "tf-idf/cmd/mysql"

	"github.com/gin-gonic/gin"
)

func PostLabels(body *gin.Context) {

	// Input valuse sections:
	links := body.PostForm("link")
	name := body.PostForm("name")
	label1 := body.PostForm("label1")
	label2 := body.PostForm("label2")
	label3 := body.PostForm("label3")

	// Insert data to mysql
	mysqlconnector.InsertLabels(links, name, label1, label2, label3)

	body.Redirect(http.StatusFound, "/")
}
