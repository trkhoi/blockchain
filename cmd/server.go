package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	MarketABI = `[
		{
		  "inputs": [],
		  "stateMutability": "nonpayable",
		  "type": "constructor"
		},
		{
		  "anonymous": false,
		  "inputs": [
			{
			  "indexed": true,
			  "internalType": "uint256",
			  "name": "itemId",
			  "type": "uint256"
			},
			{
			  "indexed": true,
			  "internalType": "address",
			  "name": "nftContract",
			  "type": "address"
			},
			{
			  "indexed": true,
			  "internalType": "uint256",
			  "name": "tokenId",
			  "type": "uint256"
			},
			{
			  "indexed": false,
			  "internalType": "address",
			  "name": "seller",
			  "type": "address"
			},
			{
			  "indexed": false,
			  "internalType": "address",
			  "name": "owner",
			  "type": "address"
			},
			{
			  "indexed": false,
			  "internalType": "uint256",
			  "name": "price",
			  "type": "uint256"
			},
			{
			  "indexed": false,
			  "internalType": "bool",
			  "name": "sold",
			  "type": "bool"
			}
		  ],
		  "name": "MarketItemCreated",
		  "type": "event"
		},
		{
		  "inputs": [
			{
			  "internalType": "address",
			  "name": "nftContract",
			  "type": "address"
			},
			{
			  "internalType": "uint256",
			  "name": "tokenId",
			  "type": "uint256"
			},
			{
			  "internalType": "uint256",
			  "name": "price",
			  "type": "uint256"
			}
		  ],
		  "name": "createMarketItem",
		  "outputs": [],
		  "stateMutability": "payable",
		  "type": "function"
		},
		{
		  "inputs": [
			{
			  "internalType": "address",
			  "name": "nftContract",
			  "type": "address"
			},
			{
			  "internalType": "uint256",
			  "name": "itemId",
			  "type": "uint256"
			}
		  ],
		  "name": "createMarketSale",
		  "outputs": [],
		  "stateMutability": "payable",
		  "type": "function"
		},
		{
		  "inputs": [],
		  "name": "fetchItemsCreated",
		  "outputs": [
			{
			  "components": [
				{
				  "internalType": "uint256",
				  "name": "itemId",
				  "type": "uint256"
				},
				{
				  "internalType": "address",
				  "name": "nftContract",
				  "type": "address"
				},
				{
				  "internalType": "uint256",
				  "name": "tokenId",
				  "type": "uint256"
				},
				{
				  "internalType": "address payable",
				  "name": "seller",
				  "type": "address"
				},
				{
				  "internalType": "address payable",
				  "name": "owner",
				  "type": "address"
				},
				{
				  "internalType": "uint256",
				  "name": "price",
				  "type": "uint256"
				},
				{
				  "internalType": "bool",
				  "name": "sold",
				  "type": "bool"
				}
			  ],
			  "internalType": "struct NFTMarket.MarketItem[]",
			  "name": "",
			  "type": "tuple[]"
			}
		  ],
		  "stateMutability": "view",
		  "type": "function"
		},
		{
		  "inputs": [],
		  "name": "fetchMarketItems",
		  "outputs": [
			{
			  "components": [
				{
				  "internalType": "uint256",
				  "name": "itemId",
				  "type": "uint256"
				},
				{
				  "internalType": "address",
				  "name": "nftContract",
				  "type": "address"
				},
				{
				  "internalType": "uint256",
				  "name": "tokenId",
				  "type": "uint256"
				},
				{
				  "internalType": "address payable",
				  "name": "seller",
				  "type": "address"
				},
				{
				  "internalType": "address payable",
				  "name": "owner",
				  "type": "address"
				},
				{
				  "internalType": "uint256",
				  "name": "price",
				  "type": "uint256"
				},
				{
				  "internalType": "bool",
				  "name": "sold",
				  "type": "bool"
				}
			  ],
			  "internalType": "struct NFTMarket.MarketItem[]",
			  "name": "",
			  "type": "tuple[]"
			}
		  ],
		  "stateMutability": "view",
		  "type": "function"
		},
		{
		  "inputs": [],
		  "name": "fetchMyNFTs",
		  "outputs": [
			{
			  "components": [
				{
				  "internalType": "uint256",
				  "name": "itemId",
				  "type": "uint256"
				},
				{
				  "internalType": "address",
				  "name": "nftContract",
				  "type": "address"
				},
				{
				  "internalType": "uint256",
				  "name": "tokenId",
				  "type": "uint256"
				},
				{
				  "internalType": "address payable",
				  "name": "seller",
				  "type": "address"
				},
				{
				  "internalType": "address payable",
				  "name": "owner",
				  "type": "address"
				},
				{
				  "internalType": "uint256",
				  "name": "price",
				  "type": "uint256"
				},
				{
				  "internalType": "bool",
				  "name": "sold",
				  "type": "bool"
				}
			  ],
			  "internalType": "struct NFTMarket.MarketItem[]",
			  "name": "",
			  "type": "tuple[]"
			}
		  ],
		  "stateMutability": "view",
		  "type": "function"
		},
		{
		  "inputs": [],
		  "name": "getListingPrice",
		  "outputs": [
			{
			  "internalType": "uint256",
			  "name": "",
			  "type": "uint256"
			}
		  ],
		  "stateMutability": "view",
		  "type": "function"
		}
	  ]`
	NftABI = `
	[
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "marketplaceAddress",
          "type": "address"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "constructor"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "owner",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "approved",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "Approval",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "owner",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "operator",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "bool",
          "name": "approved",
          "type": "bool"
        }
      ],
      "name": "ApprovalForAll",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "from",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "Transfer",
      "type": "event"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "approve",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "owner",
          "type": "address"
        }
      ],
      "name": "balanceOf",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "tokenURI",
          "type": "string"
        }
      ],
      "name": "createToken",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "getApproved",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "owner",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "operator",
          "type": "address"
        }
      ],
      "name": "isApprovedForAll",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "name",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "ownerOf",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "from",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "safeTransferFrom",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "from",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "tokenId",
          "type": "uint256"
        },
        {
          "internalType": "bytes",
          "name": "_data",
          "type": "bytes"
        }
      ],
      "name": "safeTransferFrom",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "operator",
          "type": "address"
        },
        {
          "internalType": "bool",
          "name": "approved",
          "type": "bool"
        }
      ],
      "name": "setApprovalForAll",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes4",
          "name": "interfaceId",
          "type": "bytes4"
        }
      ],
      "name": "supportsInterface",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "symbol",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "tokenURI",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "from",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "transferFrom",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ]
	`
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf("cannot read .env, err: %v", err.Error()))
	}

	// init logger
	// undo := Log()
	// defer zap.L().Sync()
	// defer undo()

	// db, err := gorm.Open(postgres.New(postgres.Config{
	// 	DSN: fmt.Sprintf("user=%v password=%v dbname=%v host=%v port=%v sslmode=disable",
	// 		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
	// 	PreferSimpleProtocol: true, // disables implicit prepared statement usage
	// }), &gorm.Config{
	// 	Logger: glog.Default.LogMode(glog.Info),
	// })
	// if err != nil {
	// 	zap.L().Panic("cannot connect to db", zap.Error(err))
	// }

	// if err := db.AutoMigrate(
	// 	&model.Transaction{},
	// ); err != nil {
	// 	zap.L().Panic("cannot migrate db", zap.Error(err))
	// }

	// transactionSvc := transaction.NewPGService(db)

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	e.GET("/test", ReadLogData)

	// handler.NewTransactionHandler(e, transactionSvc)

	e.Logger.Fatal(e.Start(":8080"))
}

type Response struct {
	TotalVolume big.Int `json:"totalVolume"`
}

func ReadLogData(c echo.Context) error {
	// wss://ropsten.infura.io/ws/v3/4a71ec7b7e324b4b94c4f1dc7811f260
	client, err := ethclient.Dial("wss://ropsten.infura.io/ws/v3/4a71ec7b7e324b4b94c4f1dc7811f260")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(client)

	contractAddress := common.HexToAddress("0x6a5ad6704a511d8B1e953076A63A6b1077814C32")
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(11307119),
		ToBlock:   big.NewInt(11347669),
		Addresses: []common.Address{
			contractAddress,
		},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(MarketABI))
	if err != nil {
		log.Fatal(err)
	}

	totalVolume := big.NewInt(0)
	for _, vLog := range logs {
		a, err := contractAbi.Unpack("MarketItemCreated", vLog.Data)
		if err != nil {
			log.Fatal(err)
		}

		totalVolume = totalVolume.Add(totalVolume, a[2].(*big.Int))
	}
	fmt.Println(totalVolume)

	return c.JSON(http.StatusOK, &Response{TotalVolume: *totalVolume})
}
