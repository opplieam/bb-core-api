package utils

import (
	"database/sql"
	"fmt"
	"math/rand/v2"

	"github.com/jaswdr/faker/v2"

	"github.com/opplieam/bb-core-api/.gen/buy-better-core/public/model"
	. "github.com/opplieam/bb-core-api/.gen/buy-better-core/public/table"
)

var (
	numberOfUsers        = 20
	numberOfSellers      = 35
	numberOfGroupProduct = 100
)

func adr(s string) *string {
	return &s
}

func randRange(min, max int) int {
	return rand.IntN(max+1-min) + min
}

func SeedUsers(db *sql.DB) error {
	fake := faker.New()

	var users []model.Users
	for _ = range numberOfUsers {
		user := model.Users{
			Email:     fmt.Sprintf("%s@gmail.com", fake.Internet().User()),
			FirstName: adr(fake.Person().FirstName()),
			LastName:  adr(fake.Person().LastName()),
			Role:      "basic",
		}
		users = append(users, user)
	}

	stmt := Users.INSERT(Users.Email, Users.FirstName, Users.LastName, Users.Role).MODELS(users)
	if _, err := stmt.Exec(db); err != nil {
		return err
	}

	fmt.Println("Seeded Users")
	return nil
}

func SeedSellers(db *sql.DB) error {
	fake := faker.New()
	var sellers []model.Sellers
	for _ = range numberOfSellers {
		seller := model.Sellers{
			Name: fake.Company().Name(),
			URL:  adr(fake.Internet().URL()),
		}
		sellers = append(sellers, seller)
	}
	stmt := Sellers.INSERT(Sellers.Name, Sellers.URL).MODELS(sellers)
	if _, err := stmt.Exec(db); err != nil {
		return err
	}
	fmt.Println("Seeded Sellers")
	return nil
}

func SeedProducts(db *sql.DB) error {
	fake := faker.New()

	// Add group_product
	for _ = range numberOfGroupProduct {
		groupProductName := fmt.Sprintf("%s %s", fake.App().Name(), fake.Color().ColorName())
		groupProduct := model.GroupProduct{
			Name: groupProductName,
		}
		groupProductStmt := GroupProduct.
			INSERT(GroupProduct.Name).
			MODEL(groupProduct).
			RETURNING(GroupProduct.ID)
		var groupProductDest model.GroupProduct
		if err := groupProductStmt.Query(db, &groupProductDest); err != nil {
			return err
		}
		// Add product
		productVariant := randRange(1, 5)
		for _ = range productVariant {
			product := model.Products{
				Name:     groupProductName,
				URL:      fake.Internet().URL(),
				SellerID: int32(randRange(1, numberOfSellers)),
			}
			var productDest model.Products
			productStmt := Products.
				INSERT(Products.Name, Products.URL, Products.SellerID).
				MODEL(product).
				RETURNING(Products.ID)
			if err := productStmt.Query(db, &productDest); err != nil {
				return err
			}
			// Add match_group_product
			matchProductGroupStmt := MatchProductGroup.
				INSERT(MatchProductGroup.GroupID, MatchProductGroup.ProductID).
				MODEL(model.MatchProductGroup{
					GroupID:   groupProductDest.ID,
					ProductID: productDest.ID,
				})
			if _, err := matchProductGroupStmt.Exec(db); err != nil {
				return err
			}
			// Add product image
			var productImages []model.ImageProduct
			for _ = range randRange(1, 3) {
				productImage := model.ImageProduct{
					ImageURL:  fake.Internet().URL(),
					ProductID: productDest.ID,
				}
				productImages = append(productImages, productImage)
			}
			imageStmt := ImageProduct.
				INSERT(ImageProduct.ImageURL, ImageProduct.ProductID).
				MODELS(productImages)
			if _, err := imageStmt.Exec(db); err != nil {
				return err
			}
			// Add price_now
			var prices []model.PriceNow
			priceVariant := randRange(1, 5)
			currency := fake.Currency().Code()
			for _ = range priceVariant {
				price := model.PriceNow{
					Price:     fake.Float64(2, 15, 1500),
					Currency:  currency,
					ProductID: productDest.ID,
				}
				prices = append(prices, price)
			}
			priceStmt := PriceNow.
				INSERT(PriceNow.Price, PriceNow.Currency, PriceNow.ProductID).
				MODELS(prices)
			if _, err := priceStmt.Exec(db); err != nil {
				return err
			}
		}
	}
	fmt.Println("Seeded Product")
	return nil
}

func SeedSubscribeProduct(db *sql.DB) error {
	var userSubs []model.UserSubProduct
	for i := range numberOfUsers {
		subVariant := randRange(1, 3)
		for _ = range subVariant {
			groupID := randRange(1, numberOfGroupProduct)
			userSub := model.UserSubProduct{
				UserID:         int32(i + 1),
				GroupProductID: int32(groupID),
			}
			userSubs = append(userSubs, userSub)
		}
	}
	stmt := UserSubProduct.
		INSERT(UserSubProduct.UserID, UserSubProduct.GroupProductID).
		MODELS(userSubs)
	if _, err := stmt.Exec(db); err != nil {
		return err
	}
	fmt.Println("Seeded User Subscribed Product")
	return nil
}
