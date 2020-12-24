package common

import "math/rand"

var messages = []string{
	"Text1",
	"Text2",
	"Text3",
	"Text4",
	"Text5",
	"Text6",
}

var yesNo = []string{"yes", "no", "may be", "I won't tell you anything!"}

func GetMessage() string {
	n := rand.Intn(len(messages))
	return messages[n]
}

func GetYesNoMSG() string {
	n := rand.Intn(len(yesNo))
	return yesNo[n]
}
