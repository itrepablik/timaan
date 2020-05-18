
![timaan](https://user-images.githubusercontent.com/58651329/82109578-beff8b00-9769-11ea-9f38-826f38901ff1.png)
The **timaan** package is a token generator for your user's authentication process in your app whether it's a WEB, CLI, or Mobile applications.

# Installation
```
go get -u github.com/itrepablik/timaan
```

# Usage
This is the sample usage for the timaan package.
```
package main

import (
	"fmt"
	"time"

	"github.com/itrepablik/itrlog"
	"github.com/itrepablik/timaan"
)

func main() {
	//*************************************************************************************
	// Add new token request for token key e.g username: 'politz' after successful login
	// This will serve as your session token stored in memory
	//*************************************************************************************
	tokenPayload := timaan.TP{
		"LOG_LEVEL":  "INFO",
		"API_KEY":    timaan.RandomToken(),
		"USERNAME":   "politz",
		"EMAIL":      "info@email.com",
		"FIRST_NAME": "Juan",
		"LAST_NAME":  "Dela Cruz",
	}
	tok := timaan.TK{
		TokenKey: "politz",
		Payload:  tokenPayload,
		ExpireOn: time.Now().Add(time.Hour * 30).Unix(),
	}
	newToken, err := timaan.GenerateToken("politz", tok)
	if err != nil {
		itrlog.Fatal(err)
	}
	fmt.Println("newToken: ", newToken)

	//*****************************************************************************
	// Extract the token payload for specific user, e.g username: 'politz'
	//*****************************************************************************
	tok, err = timaan.DecodePayload("politz")
	if err != nil {
		itrlog.Fatal(err)
	}
	fmt.Println("TokenKey: ", tok.TokenKey)

	payLoad := tok.Payload
	for field, val := range payLoad {
		fmt.Println(field, " : ", val)
	}

	//************************************************
	// Example for email confirmation token
	//************************************************
	rt := timaan.RandomToken()
	emailConfirmPayload := timaan.TP{
		"USERNAME": "juan",
		"EMAIL":    "juan@email.com",
	}
	tok = timaan.TK{
		TokenKey: rt,
		Payload:  emailConfirmPayload,
		ExpireOn: time.Now().Add(time.Hour * 30).Unix(),
	}
	newToken, err = timaan.GenerateToken(rt, tok)
	if err != nil {
		itrlog.Fatal(err)
	}
	confirmURL := "https://itrepablik.com/confirm/" + newToken
	fmt.Println(confirmURL)

	//*******************************************************************
	// Remove token for any specific user, e.g username: 'politz'
	//*******************************************************************
	isTokenRemove, err := timaan.UT.Remove("politz")
	if err != nil {
		itrlog.Fatal(err)
	}
	fmt.Println("isTokenRemove: ", isTokenRemove)
}
```
