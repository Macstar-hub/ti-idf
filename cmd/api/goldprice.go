package httppost

import (
	"fmt"
	"io"
	"net/http"
)

func GetURL() string {
	return fmt.Sprintf("https://brsapi.ir/FreeTsetmcBourseApi/Api_Free_Gold_Currency_v2.json")
}

func main() {
	getAllMessage, err := http.Get(GetURL())
	if err != nil {
		fmt.Println("Get all udpate failed with: ", err)
	}

	body, err := io.ReadAll(getAllMessage.Body)
	if err != nil {
		fmt.Println("Cannot read body with error: ", err)
	}

	/*
		Make Persian Format.
		// bodyString := persian.ToPersianDigits(string(body))
	*/

	bodyStringFormat, err := fmt.Printf("%q\n", body)
	if err != nil {
		fmt.Println("Cannot convert to utf64 with error: ", err)
	}
	fmt.Println(bodyStringFormat)
	fmt.Println("================================================================================")
	fmt.Println(string(body))
}
