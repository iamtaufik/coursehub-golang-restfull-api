package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/iamtaufik/coursehub-golang-restfull-api/internal/config"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)



func Checkout(c *fiber.Ctx) error {
	var s snap.Client

	s.New(config.Config("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	custAddress := &midtrans.CustomerAddress{
		FName:       "John",
		LName:       "Doe",
		Phone:       "081234567890",
		Address:     "Baker Street 97th",
		City:        "Jakarta",
		Postcode:    "16000",
		CountryCode: "IDN",
	}

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID: "MID-" + uuid.New().String(),
			GrossAmt: 200000,
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: "John",
			LName: "Doe",
			Email: "jhondoe@gmail.com",
			Phone: "08123456789",
			BillAddr: custAddress,
			ShipAddr: custAddress,
		},
		EnabledPayments: snap.AllSnapPaymentType,
		// Expiry: &snap.ExpiryDetails{
		// 	Duration: ,
		// },
		Items: &[]midtrans.ItemDetails{
			{
				ID:    "ITEM1",
				Price: 200000,
				Qty:   1,
				Name:  "Someitem",
			},
		},
	}

	token, err := s.CreateTransactionToken(snapReq)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to create transaction",
		})
	}


	return c.JSON(
		fiber.Map{
			"token": token,
		},
	)

}