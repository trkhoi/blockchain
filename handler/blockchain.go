package handler

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo"
)

func ReadLogData(c echo.Context) {
	// wss://ropsten.infura.io/ws/v3/4a71ec7b7e324b4b94c4f1dc7811f260
	client, err := ethclient.Dial("wss://ropsten.infura.io/ws")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(client)
}
