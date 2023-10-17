package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/gitaepark/pha/util"
	"github.com/gitaepark/pha/util/validator"
	"github.com/stretchr/testify/require"
)

func TestCreateProduct(t *testing.T) {
	user := getRandomUser(t)
	createRandomProduct(t, user)
}

func TestGetProductList(t *testing.T) {
	user := getRandomUser(t)
	for i := 0; i < 10; i++ {
		createRandomProduct(t, user)
	}

	arg := GetProductListParams{
		UserID:        user.ID,
		Searchchosung: "",
		Offset:        5,
	}

	productList, err := testQueries.GetProductList(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, productList)
	require.Equal(t, len(productList), int(10-arg.Offset))

	for _, product := range productList {
		require.NotZero(t, product.ID)
		require.NotEmpty(t, product.Category)
		require.NotZero(t, product.Price)
		require.NotZero(t, product.Cost)
		require.NotEmpty(t, product.Name)
		require.NotEmpty(t, product.Description)
		require.NotEmpty(t, product.Barcode)
		require.NotZero(t, product.ExpirationDate)
		require.NotEmpty(t, product.Size)
		require.True(t, validator.IsSupportedProductSize(string(product.Size)))
		require.NotZero(t, product.CreatedAt)
		require.NotZero(t, product.UpdatedAt)
	}
}

func TestGetProductListWithKeyword(t *testing.T) {
	user := getRandomUser(t)
	for i := 0; i < 9; i++ {
		createRandomProduct(t, user)
	}
	testQueries.CreateProduct(context.Background(), CreateProductParams{
		UserID:         user.ID,
		Category:       util.CreateRandomString(15),
		Price:          util.CreateRandomInt32(1000, 10000),
		Cost:           util.CreateRandomInt32(1000, 10000),
		Name:           "슈크림 라떼",
		Description:    util.CreateRandomString(50),
		Barcode:        util.CreateRandomString(12),
		ExpirationDate: time.Now().Add(24 * time.Hour),
		Size:           ProductSize(util.CreateRandomProductSize()),
	})

	arg := GetProductListParams{
		UserID:        user.ID,
		Searchchosung: "슈크림",
		Offset:        0,
	}

	productList, err := testQueries.GetProductList(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, productList)

	require.NotZero(t, productList[0].ID)
	require.NotEmpty(t, productList[0].Category)
	require.NotZero(t, productList[0].Price)
	require.NotZero(t, productList[0].Cost)
	require.Equal(t, productList[0].Name, "슈크림 라떼")
	require.NotEmpty(t, productList[0].Description)
	require.NotEmpty(t, productList[0].Barcode)
	require.NotZero(t, productList[0].ExpirationDate)
	require.NotEmpty(t, productList[0].Size)
	require.True(t, validator.IsSupportedProductSize(string(productList[0].Size)))
	require.NotZero(t, productList[0].CreatedAt)
	require.NotZero(t, productList[0].UpdatedAt)
}

func TestGetProductListWithChosungKeyword(t *testing.T) {
	user := getRandomUser(t)
	for i := 0; i < 9; i++ {
		createRandomProduct(t, user)
	}
	testQueries.CreateProduct(context.Background(), CreateProductParams{
		UserID:         user.ID,
		Category:       util.CreateRandomString(15),
		Price:          util.CreateRandomInt32(1000, 10000),
		Cost:           util.CreateRandomInt32(1000, 10000),
		Name:           "슈크림 라떼",
		Description:    util.CreateRandomString(50),
		Barcode:        util.CreateRandomString(12),
		ExpirationDate: time.Now().Add(24 * time.Hour),
		Size:           ProductSize(util.CreateRandomProductSize()),
	})

	arg := GetProductListParams{
		UserID:        user.ID,
		Searchchosung: "ㅅㅋㄹ",
		Offset:        0,
	}

	productList, err := testQueries.GetProductList(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, productList)

	require.NotZero(t, productList[0].ID)
	require.NotEmpty(t, productList[0].Category)
	require.NotZero(t, productList[0].Price)
	require.NotZero(t, productList[0].Cost)
	require.Equal(t, productList[0].Name, "슈크림 라떼")
	require.NotEmpty(t, productList[0].Description)
	require.NotEmpty(t, productList[0].Barcode)
	require.NotZero(t, productList[0].ExpirationDate)
	require.NotEmpty(t, productList[0].Size)
	require.True(t, validator.IsSupportedProductSize(string(productList[0].Size)))
	require.NotZero(t, productList[0].CreatedAt)
	require.NotZero(t, productList[0].UpdatedAt)
}

func TestGetProduct(t *testing.T) {
	user := getRandomUser(t)
	createRandomProduct(t, user)
	productList, _ := testQueries.GetProductList(context.Background(), GetProductListParams{
		UserID:        user.ID,
		Searchchosung: "",
		Offset:        0,
	})

	product, err := testQueries.GetProduct(context.Background(), productList[0].ID)
	require.NoError(t, err)
	require.NotEmpty(t, product)

	require.Equal(t, product.ID, productList[0].ID)
	require.Equal(t, product.UserID, productList[0].UserID)
	require.Equal(t, product.Category, productList[0].Category)
	require.Equal(t, product.Price, productList[0].Price)
	require.Equal(t, product.Cost, productList[0].Cost)
	require.Equal(t, product.Name, productList[0].Name)
	require.Equal(t, product.Description, productList[0].Description)
	require.Equal(t, product.Barcode, productList[0].Barcode)
	require.WithinDuration(t, product.ExpirationDate, productList[0].ExpirationDate, 24*time.Hour)
	require.Equal(t, product.Size, productList[0].Size)
	require.WithinDuration(t, product.CreatedAt, productList[0].CreatedAt, time.Second)
	require.WithinDuration(t, product.UpdatedAt, productList[0].UpdatedAt, time.Second)
}

func TestUpdateProduct(t *testing.T) {
	user := getRandomUser(t)
	createRandomProduct(t, user)
	productList, _ := testQueries.GetProductList(context.Background(), GetProductListParams{
		UserID:        user.ID,
		Searchchosung: "",
		Offset:        0,
	})

	arg := UpdateProductParams{
		Category:       util.CreateRandomString(15),
		Price:          util.CreateRandomInt32(1000, 10000),
		Cost:           util.CreateRandomInt32(1000, 10000),
		Name:           util.CreateRandomString(10),
		Description:    util.CreateRandomString(50),
		Barcode:        util.CreateRandomString(12),
		ExpirationDate: time.Now().Add(24 * time.Hour),
		Size:           ProductSize(util.CreateRandomProductSize()),
		ID:             productList[0].ID,
	}

	time.Sleep(time.Second)

	err := testQueries.UpdateProduct(context.Background(), arg)
	require.NoError(t, err)

	product, _ := testQueries.GetProduct(context.Background(), productList[0].ID)
	require.Equal(t, product.ID, productList[0].ID)
	require.Equal(t, product.Category, arg.Category)
	require.Equal(t, product.Price, arg.Price)
	require.Equal(t, product.Cost, arg.Cost)
	require.Equal(t, product.Name, arg.Name)
	require.Equal(t, product.Description, arg.Description)
	require.Equal(t, product.Barcode, arg.Barcode)
	require.WithinDuration(t, product.ExpirationDate, arg.ExpirationDate, 24*time.Hour)
	require.Equal(t, product.Size, arg.Size)
	require.WithinDuration(t, product.CreatedAt, productList[0].CreatedAt, time.Second)
	fmt.Println(product.UpdatedAt)
	fmt.Println(productList[0].UpdatedAt)
	require.NotEqual(t, product.UpdatedAt, productList[0].UpdatedAt)
}

func TestDeleteProduct(t *testing.T) {
	user := getRandomUser(t)
	createRandomProduct(t, user)
	productList, _ := testQueries.GetProductList(context.Background(), GetProductListParams{
		UserID:        user.ID,
		Searchchosung: "",
		Offset:        0,
	})

	err := testQueries.DeleteProduct(context.Background(), productList[0].ID)
	require.NoError(t, err)

	_, err = testQueries.GetProduct(context.Background(), productList[0].ID)
	require.ErrorIs(t, err, sql.ErrNoRows)
}

func createRandomProduct(t *testing.T, user User) {
	arg := CreateProductParams{
		UserID:         user.ID,
		Category:       util.CreateRandomString(15),
		Price:          util.CreateRandomInt32(1000, 10000),
		Cost:           util.CreateRandomInt32(1000, 10000),
		Name:           util.CreateRandomString(10),
		Description:    util.CreateRandomString(50),
		Barcode:        util.CreateRandomString(12),
		ExpirationDate: time.Now().Add(24 * time.Hour),
		Size:           ProductSize(util.CreateRandomProductSize()),
	}

	err := testQueries.CreateProduct(context.Background(), arg)
	require.NoError(t, err)
}
