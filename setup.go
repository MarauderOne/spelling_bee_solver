package main

import (
	"github.com/MarauderOne/spelling_bee_solver/dictionary_tools"
	"github.com/golang/glog"
	"strings"
)

//Define a data structure for user inputs
type CellData struct {
	Character string `json:"character"`
	Color     string `json:"color"`
}

//Define a function to initialise the answerList
func createNewAnswerList() (dictionary *dictionary_tools.MySimpleDict) {
	glog.Info("Creating d variable as NewSimpleDict struct")
	newDictionary := dictionary_tools.NewSimpleDict()
	glog.Info("Loading dictionary_tools/initialList.dict into d variable")
	newDictionary.Load("dictionary_tools/initialList.dict")
	glog.Info("Returning createNewAnswerList function")
	return newDictionary
}

//Define a function to revise the answerList based on given regex patterns
func reviseAnswerList(answersList *dictionary_tools.MySimpleDict, regexPattern string) (revisedDictionary *dictionary_tools.MySimpleDict) {
	glog.Info("Filtering list of potential answers using regex to create new list of potential answers in newAnswerList variable")
	newAnswerList := answersList.Lookup(regexPattern, 0, 80000)
	glog.Info("Creating d variable as NewSimpleDict struct")
	d := dictionary_tools.NewSimpleDict()
	glog.Info("Loading newAnswerList variable into d variable")
	d.AddWordsList(newAnswerList)
	return d
}

//Define a function to determine if a character is not an uppercase alphabetic character
func nonAlpha(char string) bool {

	//Define a list of uppercase alphabetic characters
	glog.Info("Defining list of uppercase alphabetic characters")
	const alpha = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	//If the character in the box is in the alpha string or is ""
	glog.Info("Checking if box character is alphabetic")
	if (strings.Contains(alpha, strings.ToUpper(string(char)))) || (string(char) == "") {
		//Then return no error
		glog.Info("Box character is alphabetic")
		return false
	} else {
		//Else return an error
		glog.Errorf("Box character is not alphabetic: %v", char)
		return true
	}
}
