package httppost

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// test
type userInfoTemplate struct {
	FirstName    string
	LastName     string
	Email        string
	TicketNumber int
}

var remainingTickets int = 50

func UserInfoPost(body *gin.Context) {

	firstName := body.PostForm("firstnames")
	lastName := body.PostForm("lastnames")
	email := body.PostForm("emails")
	ticketNumberString := body.PostForm("ticketnumbers")

	ticketNumber, _ := strconv.Atoi(ticketNumberString)
	UserInfo := userInfoTemplate{
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		TicketNumber: ticketNumber,
	}
	userData := userhandaler.UserHandelerStruct(UserInfo.FirstName, UserInfo.LastName, UserInfo.Email, UserInfo.TicketNumber)
	firstName, lastName, email, ticketNumber = userData.FirstName, userData.LastName, userData.Email, userData.TicketNumber
	allUsers, allFirsnames := userhandaler.MakeSlicesd(firstName, lastName, email, ticketNumber)
	fmt.Printf("All users are: %v And all just firstnames: %v\n", allUsers, allFirsnames)

	isNameValid, isMailValid, isTicketNumberValid := uservalidation.UserInputValidations(firstName, lastName, email, ticketNumber)
	if isNameValid && isMailValid && isTicketNumberValid {
		remainingTickets = remainingtickets.AvailableTickets(remainingTickets, ticketNumber)

		// Produce test info to rabbitmq
		rabbitmq.RabbitProducer(firstName, lastName, email, ticketNumber)
		// go rabbitmq.RabbitConsumer()
		//

		notEnoughTickets := remainingTickets <= 0
		if notEnoughTickets {
			fmt.Printf("Booking Failed With Ticket Remaining: %v \n", remainingTickets)
			// break
		} else {
			fmt.Printf("User %v %v With Email %v Booked Successfuly %v Tickets \n", firstName, lastName, email, ticketNumber)
		}

		// Make user feedback guid
	} else {
		if !isNameValid {
			fmt.Println("Please Enter Correct FirstName And LastName ... ")
			fmt.Printf("FisrtName that enterd: %v And lastName that enterd: %v", firstName, lastName)
		}
		if !isMailValid {
			fmt.Println("Please Enter Correct Email Address ... ")
		}
		if !isTicketNumberValid {
			fmt.Printf("Please Select Ticket Number In Range Remaining Tickets: %v \n", remainingTickets)
		}
	}
	body.Redirect(http.StatusFound, "/api/v1/getusers")
}

func BookedUsers(body *gin.Context) {
	// Authentication Section
	state := "someState"
	if body.Request.URL.Query().Get("state") != state {
		body.Redirect(308, "http://192.168.1.100:8091/?uri=/api/v1/getusers")
	}
	// End Of authentication section.

	allUsersInfos := mysqlconnector.SelectQury()
	userListNumber := len(allUsersInfos.Firstnames)
	var users []gin.H
	for i := 0; i < userListNumber; i++ {
		users = append(users, gin.H{
			"Firstname":    allUsersInfos.Firstnames[i],
			"LastName":     allUsersInfos.LastName[i],
			"Email":        allUsersInfos.Email[i],
			"TicketNumber": allUsersInfos.TicketNumber[i],
		})
	}
	body.HTML(http.StatusOK, "adminArea.html", gin.H{
		"Booking": users,
	})
	body.Redirect(http.StatusFound, "/")
}
