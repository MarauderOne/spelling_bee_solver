package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sort"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func main() {

	//Parse the argument flags (like the ones in Heroku's Procfile)
	flag.Parse()

	webServer := gin.Default()

	//Serve static files
	webServer.Static("/spellingbeesolver", "./frontend")

	//Define routes
	webServer.GET("/", func(c *gin.Context) {
		c.File("./frontend/index.html")
	})

	//Endpoint to handle Spelling Bee solving
	webServer.POST("/letters", parseLetters)

	//Define the port
	port := os.Getenv("PORT")
	//Define default port (for local testing)
	if port == "" {
		port = "8080"
	}

	//Run the webserver
	err := webServer.Run(":" + port)
	if err != nil {
		glog.Fatalf("Web server initialisation failed: %v", err)
	}
}

//Define function for parse user's guesses
func parseLetters(c *gin.Context) {
	var gridData []CellData
	if err := c.ShouldBindJSON(&gridData); err != nil {
		glog.Errorf("Unable to bind JSON from page: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Convert Character values to uppercase
	for i := range gridData {
			gridData[i].Character = strings.ToUpper(gridData[i].Character)
	}

	//Call the Spelling Bee solving function
	glog.Infof("Calling solveSpellingBee function with gridData: %v", gridData)
	result, countOfResults, solvingError, httpStatus := solveSpellingBee(gridData)

	if solvingError != "" {
		//Respond with the http status code and error message returned by the solveSpellingBee function
		glog.Warningf("Responding to page with error response: %v", solvingError)
		glog.Flush()
		c.JSON(httpStatus, gin.H{"error": solvingError})
	} else {
		// Respond with the result
		glog.Info("Responding to the page with expected results")
		glog.Flush()
		c.JSON(httpStatus, gin.H{"result": result, "resultCount": countOfResults})
	}
}

//Spelling Bee solving logic
func solveSpellingBee(gridData []CellData) (result string, countOfResults int, solvingError string, httpStatus int) {

	//Set default HTTP response code (will be updated if there is an error)
	httpStatus = http.StatusOK

	//Initialise answer list
	glog.Info("Calling createNewAnswerList function")
	answerList := createNewAnswerList()

	//Define variables for regex patterns
	var yellowRegex string
	var greyRegex string
	

	//Start looping through the user boxes
	for _, box := range gridData {

		//Check for non-alphabetic characters
		if nonAlpha(box.Character) {
			glog.Errorf("Invalid character recieved: %v", box.Character)
			solvingError = fmt.Sprintf("Invalid character: %v", box.Character)
			httpStatus = http.StatusBadRequest
			glog.Error("Breaking revision loop")
			glog.Info("Writing results count")
			var resultCount int = answerList.Count()
			glog.Info("Writing results")
			results := strings.Join(answerList.Words, " ")
			glog.Info("Returning solveSpellingBee function")
			return results, resultCount, solvingError, httpStatus
		}

		//Skip over boxes which have either no character or no color
		if (box.Character == "") || (box.Color == "") {
			glog.Info("Character or Color value is missing, skipping this iteration of revision loop")
			continue
		}

		switch box.Color {
		case "yellow":
			//Set regex pattern for yellow box
			glog.Info("Define regex pattern for character in yellow box")
			yellowRegex = fmt.Sprintf(".*%v.*", box.Character)
			glog.Info("Add yellow box character to greybox regex pattern")
			greyRegex = fmt.Sprintf("%v%v", greyRegex, box.Character)
		case "grey":
			//Set regex pattern for grey box
			glog.Info("Define regex pattern for character in grey box")
			greyRegex = fmt.Sprintf("%v%v", greyRegex, box.Character)
		default:
			//Invalid color, this should never be reached (except in the tests)
			glog.Errorf("Invalid color recieved: %v", box.Color)
			solvingError = fmt.Sprintf("Invalid color: %v", box.Color)
			httpStatus = http.StatusBadRequest
			glog.Error("Breaking loop")
			glog.Info("Writing results count")
			var resultCount int = answerList.Count()
			glog.Info("Writing results")
			results := strings.Join(answerList.Words, " ")
			glog.Info("Returning solveSpellingBee function")
			return results, resultCount, solvingError, httpStatus
		}
	}

	glog.Info("Revising answer list for yellow box character")
	answerList = reviseAnswerList(answerList, yellowRegex)
	glog.Info("Revising answer list for grey box character(s)")
	greyRegex = fmt.Sprintf("[%v]*", greyRegex)
	answerList = reviseAnswerList(answerList, greyRegex)
	
	glog.Info("Writing results count")
	var resultCount int = answerList.Count()
	glog.Info("Sorting results")
	sort.Slice(answerList.Words, func(i, j int) bool {
        return len(answerList.Words[i]) > len(answerList.Words[j])
    })
	glog.Info("Writing results")
	results := strings.Join(answerList.Words, " ")
	glog.Info("Returning solveSpellingBee function")
	return results, resultCount, solvingError, httpStatus
}
