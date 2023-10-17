package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/gitaepark/pha/dto"
	"github.com/gitaepark/pha/repository"
	mockrepository "github.com/gitaepark/pha/repository/mock"
	"github.com/gitaepark/pha/util"
	"github.com/go-sql-driver/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateProduct(t *testing.T) {
	user, _ := createRandomUser(t)
	product := createRandomProduct(t, user)

	testCases := []struct {
		name          string
		params        CreateProductParams
		buildStubs    func(mockRepository *mockrepository.MockRepository)
		checkResponse func(err CustomErr)
	}{
		{
			name: "성공",
			params: CreateProductParams{
				UserID: user.ID,
				CreateProductRequestBody: dto.CreateProductRequestBody{
					Category:       product.Category,
					Price:          product.Price,
					Cost:           product.Cost,
					Name:           product.Name,
					Description:    product.Description,
					Barcode:        product.Barcode,
					ExpirationDate: product.ExpirationDate.Format(util.DateLayout),
					Size:           string(product.Size),
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(err CustomErr) {
				require.Empty(t, err)
			},
		},
		{
			name: "날짜 타입 변환 실패",
			params: CreateProductParams{
				UserID: user.ID,
				CreateProductRequestBody: dto.CreateProductRequestBody{
					Category:       product.Category,
					Price:          product.Price,
					Cost:           product.Cost,
					Name:           product.Name,
					Description:    product.Description,
					Barcode:        product.Barcode,
					ExpirationDate: util.CreateRandomString(10),
					Size:           string(product.Size),
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(err CustomErr) {
				require.Equal(t, err, errParseDate)
			},
		},
		{
			name: "바코드가 중복된 경우",
			params: CreateProductParams{
				UserID: user.ID,
				CreateProductRequestBody: dto.CreateProductRequestBody{
					Category:       product.Category,
					Price:          product.Price,
					Cost:           product.Cost,
					Name:           product.Name,
					Description:    product.Description,
					Barcode:        product.Barcode,
					ExpirationDate: product.ExpirationDate.Format(util.DateLayout),
					Size:           string(product.Size),
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&mysql.MySQLError{Number: repository.DB_DUPLICATE_ERROR, Message: "barcode"})
			},
			checkResponse: func(err CustomErr) {
				require.Equal(t, err, errDuplicateBarcode)
			},
		},
		{
			name: "회원이 없는 경우",
			params: CreateProductParams{
				UserID: util.CreateRandomInt64(11, 20),
				CreateProductRequestBody: dto.CreateProductRequestBody{
					Category:       product.Category,
					Price:          product.Price,
					Cost:           product.Cost,
					Name:           product.Name,
					Description:    product.Description,
					Barcode:        product.Barcode,
					ExpirationDate: product.ExpirationDate.Format(util.DateLayout),
					Size:           string(product.Size),
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&mysql.MySQLError{Number: repository.DB_FK_ERROR, Message: "user_id"})
			},
			checkResponse: func(err CustomErr) {
				require.Equal(t, err, errNotFoundUser)
			},
		},
		{
			name: "Internal Server Error",
			params: CreateProductParams{
				UserID: user.ID,
				CreateProductRequestBody: dto.CreateProductRequestBody{
					Category:       product.Category,
					Price:          product.Price,
					Cost:           product.Cost,
					Name:           product.Name,
					Description:    product.Description,
					Barcode:        product.Barcode,
					ExpirationDate: product.ExpirationDate.Format(util.DateLayout),
					Size:           string(product.Size),
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(err CustomErr) {
				require.Equal(t, err, NewErrInternalServer(sql.ErrConnDone))
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := mockrepository.NewMockRepository(ctrl)
			service := newTestService(t, repository)

			tc.buildStubs(repository)

			err := service.CreateProduct(context.Background(), tc.params)
			tc.checkResponse(err)
		})
	}
}

func TestGetProductList(t *testing.T) {
	user, _ := createRandomUser(t)
	var productList []repository.Product
	for i := 0; i < 10; i++ {
		productList = append(productList, createRandomProduct(t, user))
	}

	testCases := []struct {
		name          string
		params        GetProductListParams
		buildStubs    func(mockRepository *mockrepository.MockRepository)
		checkResponse func(result dto.GetProductListResponse, err CustomErr)
	}{
		{
			name: "성공",
			params: GetProductListParams{
				UserID: user.ID,
				GetProductListRequestQuery: dto.GetProductListRequestQuery{
					Page: 1,
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetProductList(gomock.Any(), gomock.Any()).
					Times(1).
					Return(productList, nil)
			},
			checkResponse: func(result dto.GetProductListResponse, err CustomErr) {
				require.NotEmpty(t, result)
				require.Empty(t, err)

				for idx, product := range result.List {
					require.Equal(t, product.ID, productList[idx].ID)
					require.Equal(t, product.Category, productList[idx].Category)
					require.Equal(t, product.Price, productList[idx].Price)
					require.Equal(t, product.Cost, productList[idx].Cost)
					require.Equal(t, product.Name, productList[idx].Name)
					require.Equal(t, product.Description, productList[idx].Description)
					require.Equal(t, product.Barcode, productList[idx].Barcode)
					require.Equal(t, product.ExpirationDate, productList[idx].ExpirationDate.Format(util.DateLayout))
					require.Equal(t, product.Size, productList[idx].Size)
					require.WithinDuration(t, product.CreatedAt, productList[idx].CreatedAt, time.Second)
					require.WithinDuration(t, product.UpdatedAt, productList[idx].UpdatedAt, time.Second)
				}
			},
		},
		{
			name: "검색 성공",
			params: GetProductListParams{
				UserID: user.ID,
				GetProductListRequestQuery: dto.GetProductListRequestQuery{
					Page:    1,
					Keyword: productList[0].Name,
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetProductList(gomock.Any(), gomock.Any()).
					Times(1).
					Return(productList[0:1], nil)
			},
			checkResponse: func(result dto.GetProductListResponse, err CustomErr) {
				require.NotEmpty(t, result)
				require.Empty(t, err)

				for _, product := range result.List {
					require.Equal(t, productList[0].ID, product.ID)
					require.Equal(t, productList[0].Category, product.Category)
					require.Equal(t, productList[0].Price, product.Price)
					require.Equal(t, productList[0].Cost, product.Cost)
					require.Equal(t, productList[0].Name, product.Name)
					require.Equal(t, productList[0].Description, product.Description)
					require.Equal(t, productList[0].Barcode, product.Barcode)
					require.Equal(t, productList[0].ExpirationDate.Format(util.DateLayout), product.ExpirationDate)
					require.Equal(t, productList[0].Size, product.Size)
					require.WithinDuration(t, productList[0].CreatedAt, product.CreatedAt, time.Second)
					require.WithinDuration(t, productList[0].UpdatedAt, product.UpdatedAt, time.Second)
				}
			},
		},
		{
			name: "Internal Server Error",
			params: GetProductListParams{
				UserID: user.ID,
				GetProductListRequestQuery: dto.GetProductListRequestQuery{
					Page: 1,
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetProductList(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]repository.Product{}, sql.ErrConnDone)
			},
			checkResponse: func(result dto.GetProductListResponse, err CustomErr) {
				require.Empty(t, result)
				require.Equal(t, err, NewErrInternalServer(sql.ErrConnDone))
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := mockrepository.NewMockRepository(ctrl)
			service := newTestService(t, repository)

			tc.buildStubs(repository)

			result, err := service.GetProductList(context.Background(), tc.params)
			tc.checkResponse(result, err)
		})
	}
}

func TestGetProduct(t *testing.T) {
	user, _ := createRandomUser(t)
	product := createRandomProduct(t, user)

	testCases := []struct {
		name          string
		params        GetProductParams
		buildStubs    func(mockRepository *mockrepository.MockRepository)
		checkResponse func(result dto.GetProductResponse, err CustomErr)
	}{
		{
			name: "성공",
			params: GetProductParams{
				UserID: user.ID,
				GetProductRequestPath: dto.GetProductRequestPath{
					ID: product.ID,
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(product, nil)
			},
			checkResponse: func(result dto.GetProductResponse, err CustomErr) {
				require.NotEmpty(t, result)
				require.Empty(t, err)

				require.Equal(t, result.ID, product.ID)
				require.Equal(t, result.Category, product.Category)
				require.Equal(t, result.Price, product.Price)
				require.Equal(t, result.Cost, product.Cost)
				require.Equal(t, result.Name, product.Name)
				require.Equal(t, result.Description, product.Description)
				require.Equal(t, result.Barcode, product.Barcode)
				require.Equal(t, result.ExpirationDate, product.ExpirationDate.Format(util.DateLayout))
				require.Equal(t, result.Size, product.Size)
				require.WithinDuration(t, result.CreatedAt, product.CreatedAt, time.Second)
				require.WithinDuration(t, result.UpdatedAt, product.UpdatedAt, time.Second)
			},
		},
		{
			name: "상품이 없는 경우",
			params: GetProductParams{
				UserID: user.ID,
				GetProductRequestPath: dto.GetProductRequestPath{
					ID: util.CreateRandomInt64(11, 20),
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(repository.Product{}, sql.ErrNoRows)
			},
			checkResponse: func(result dto.GetProductResponse, err CustomErr) {
				require.Empty(t, result)
				require.Equal(t, err, errNotFoundProduct)
			},
		},
		{
			name: "상품을 등록한 회원이 아닌 경우",
			params: GetProductParams{
				UserID: util.CreateRandomInt64(11, 20),
				GetProductRequestPath: dto.GetProductRequestPath{
					ID: product.ID,
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(product, nil)

				mockRepository.EXPECT().
					DeleteProduct(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(result dto.GetProductResponse, err CustomErr) {
				require.Empty(t, result)
				require.Equal(t, err, errForbiddenProduct)
			},
		},
		{
			name: "Internal Server Error",
			params: GetProductParams{
				UserID: user.ID,
				GetProductRequestPath: dto.GetProductRequestPath{
					ID: product.ID,
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(repository.Product{}, sql.ErrConnDone)
			},
			checkResponse: func(result dto.GetProductResponse, err CustomErr) {
				require.Equal(t, err, NewErrInternalServer(sql.ErrConnDone))
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := mockrepository.NewMockRepository(ctrl)
			service := newTestService(t, repository)

			tc.buildStubs(repository)

			result, err := service.GetProduct(context.Background(), tc.params)
			tc.checkResponse(result, err)
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	user, _ := createRandomUser(t)
	product := createRandomProduct(t, user)

	updatedProduct := repository.Product{
		ID:             product.ID,
		UserID:         user.ID,
		Category:       util.CreateRandomString(15),
		Price:          util.CreateRandomInt32(1000, 10000),
		Cost:           util.CreateRandomInt32(1000, 10000),
		Name:           util.CreateRandomString(10),
		Description:    util.CreateRandomString(50),
		Barcode:        util.CreateRandomString(12),
		ExpirationDate: time.Now().Add(24 * time.Hour),
		Size:           repository.ProductSize(util.CreateRandomProductSize()),
		CreatedAt:      product.CreatedAt,
		UpdatedAt:      time.Now(),
	}
	updatedProductExpirationDate := updatedProduct.ExpirationDate.Format(util.DateLayout)
	invalidExpirationDate := util.CreateRandomString(10)

	testCases := []struct {
		name          string
		params        UpdateProductParams
		buildStubs    func(mockRepository *mockrepository.MockRepository)
		checkResponse func(err CustomErr)
	}{
		{
			name: "카테고리 수정 성공",
			params: UpdateProductParams{
				UserID: user.ID,
				UpdateProductRequestPath: dto.UpdateProductRequestPath{
					ID: product.ID,
				},
				UpdateProductRequestBody: dto.UpdateProductRequestBody{
					Category: &updatedProduct.Category,
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(product, nil)

				mockRepository.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(err CustomErr) {
				require.Empty(t, err)
			},
		},
		{
			name: "가격 수정 성공",
			params: UpdateProductParams{
				UserID: user.ID,
				UpdateProductRequestPath: dto.UpdateProductRequestPath{
					ID: product.ID,
				},
				UpdateProductRequestBody: dto.UpdateProductRequestBody{
					Price: &updatedProduct.Price,
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(product, nil)

				mockRepository.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(err CustomErr) {
				require.Empty(t, err)
			},
		},
		{
			name: "원가 수정 성공",
			params: UpdateProductParams{
				UserID: user.ID,
				UpdateProductRequestPath: dto.UpdateProductRequestPath{
					ID: product.ID,
				},
				UpdateProductRequestBody: dto.UpdateProductRequestBody{
					Cost: &updatedProduct.Cost,
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(product, nil)

				mockRepository.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(err CustomErr) {
				require.Empty(t, err)
			},
		},
		{
			name: "이름 수정 성공",
			params: UpdateProductParams{
				UserID: user.ID,
				UpdateProductRequestPath: dto.UpdateProductRequestPath{
					ID: product.ID,
				},
				UpdateProductRequestBody: dto.UpdateProductRequestBody{
					Name: &updatedProduct.Name,
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(product, nil)

				mockRepository.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(err CustomErr) {
				require.Empty(t, err)
			},
		},
		{
			name: "설명 수정 성공",
			params: UpdateProductParams{
				UserID: user.ID,
				UpdateProductRequestPath: dto.UpdateProductRequestPath{
					ID: product.ID,
				},
				UpdateProductRequestBody: dto.UpdateProductRequestBody{
					Description: &updatedProduct.Description,
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(product, nil)

				mockRepository.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(err CustomErr) {
				require.Empty(t, err)
			},
		},
		{
			name: "바코드 수정 성공",
			params: UpdateProductParams{
				UserID: user.ID,
				UpdateProductRequestPath: dto.UpdateProductRequestPath{
					ID: product.ID,
				},
				UpdateProductRequestBody: dto.UpdateProductRequestBody{
					Barcode: &updatedProduct.Barcode,
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(product, nil)

				mockRepository.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(err CustomErr) {
				require.Empty(t, err)
			},
		},
		{
			name: "유통기한 수정 성공",
			params: UpdateProductParams{
				UserID: user.ID,
				UpdateProductRequestPath: dto.UpdateProductRequestPath{
					ID: product.ID,
				},
				UpdateProductRequestBody: dto.UpdateProductRequestBody{
					ExpirationDate: &updatedProductExpirationDate,
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(product, nil)

				mockRepository.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(err CustomErr) {
				require.Empty(t, err)
			},
		},
		{
			name: "사이즈 수정 성공",
			params: UpdateProductParams{
				UserID: user.ID,
				UpdateProductRequestPath: dto.UpdateProductRequestPath{
					ID: product.ID,
				},
				UpdateProductRequestBody: dto.UpdateProductRequestBody{
					Size: (*string)(&updatedProduct.Size),
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(product, nil)

				mockRepository.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(err CustomErr) {
				require.Empty(t, err)
			},
		},
		{
			name: "상품이 없는 경우",
			params: UpdateProductParams{
				UserID: user.ID,
				UpdateProductRequestPath: dto.UpdateProductRequestPath{
					ID: util.CreateRandomInt64(11, 20),
				},
				UpdateProductRequestBody: dto.UpdateProductRequestBody{
					Category: &updatedProduct.Category,
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(repository.Product{}, sql.ErrNoRows)

				mockRepository.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(err CustomErr) {
				require.Equal(t, err, errNotFoundProduct)
			},
		},
		{
			name: "상품을 등록한 회원이 아닌 경우",
			params: UpdateProductParams{
				UserID: util.CreateRandomInt64(11, 20),
				UpdateProductRequestPath: dto.UpdateProductRequestPath{
					ID: product.ID,
				},
				UpdateProductRequestBody: dto.UpdateProductRequestBody{
					Category: &updatedProduct.Category,
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(product, nil)

				mockRepository.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(err CustomErr) {
				require.Equal(t, err, errForbiddenProduct)
			},
		},
		{
			name: "날짜 변환 실패",
			params: UpdateProductParams{
				UserID: user.ID,
				UpdateProductRequestPath: dto.UpdateProductRequestPath{
					ID: product.ID,
				},
				UpdateProductRequestBody: dto.UpdateProductRequestBody{
					ExpirationDate: &invalidExpirationDate,
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(product, nil)

				mockRepository.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(err CustomErr) {
				require.Equal(t, err, errParseDate)
			},
		},
		{
			name: "Internal Server Error",
			params: UpdateProductParams{
				UserID: user.ID,
				UpdateProductRequestPath: dto.UpdateProductRequestPath{
					ID: product.ID,
				},
				UpdateProductRequestBody: dto.UpdateProductRequestBody{
					Category: &updatedProduct.Category,
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(product, sql.ErrConnDone)

				mockRepository.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(err CustomErr) {
				require.Equal(t, err, NewErrInternalServer(sql.ErrConnDone))
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := mockrepository.NewMockRepository(ctrl)
			service := newTestService(t, repository)

			tc.buildStubs(repository)

			err := service.UpdateProduct(context.Background(), tc.params)
			tc.checkResponse(err)
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	user, _ := createRandomUser(t)
	product := createRandomProduct(t, user)

	testCases := []struct {
		name          string
		params        DeleteProductParams
		buildStubs    func(mockRepository *mockrepository.MockRepository)
		checkResponse func(err CustomErr)
	}{
		{
			name: "성공",
			params: DeleteProductParams{
				UserID: user.ID,
				DeleteProductRequestPath: dto.DeleteProductRequestPath{
					ID: product.ID,
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(product, nil)

				mockRepository.EXPECT().
					DeleteProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(err CustomErr) {
				require.Empty(t, err)
			},
		},
		{
			name: "상품이 없는 경우",
			params: DeleteProductParams{
				UserID: user.ID,
				DeleteProductRequestPath: dto.DeleteProductRequestPath{
					ID: util.CreateRandomInt64(11, 20),
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(repository.Product{}, sql.ErrNoRows)

				mockRepository.EXPECT().
					DeleteProduct(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(err CustomErr) {
				require.Equal(t, err, errNotFoundProduct)
			},
		},
		{
			name: "상품을 등록한 회원이 아닌 경우",
			params: DeleteProductParams{
				UserID: util.CreateRandomInt64(11, 20),
				DeleteProductRequestPath: dto.DeleteProductRequestPath{
					ID: product.ID,
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(product, nil)

				mockRepository.EXPECT().
					DeleteProduct(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(err CustomErr) {
				require.Equal(t, err, errForbiddenProduct)
			},
		},
		{
			name: "Internal Server Error",
			params: DeleteProductParams{
				UserID: user.ID,
				DeleteProductRequestPath: dto.DeleteProductRequestPath{
					ID: product.ID,
				},
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(repository.Product{}, sql.ErrConnDone)
			},
			checkResponse: func(err CustomErr) {
				require.Equal(t, err, NewErrInternalServer(sql.ErrConnDone))
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := mockrepository.NewMockRepository(ctrl)
			service := newTestService(t, repository)

			tc.buildStubs(repository)

			err := service.DeleteProduct(context.Background(), tc.params)
			tc.checkResponse(err)
		})
	}
}

func createRandomProduct(t *testing.T, user repository.User) repository.Product {
	product := repository.Product{
		ID:             util.CreateRandomInt64(1, 10),
		UserID:         user.ID,
		Category:       util.CreateRandomString(15),
		Price:          util.CreateRandomInt32(1000, 10000),
		Cost:           util.CreateRandomInt32(1000, 10000),
		Name:           util.CreateRandomString(10),
		Description:    util.CreateRandomString(50),
		Barcode:        util.CreateRandomString(12),
		ExpirationDate: time.Now().Add(24 * time.Hour),
		Size:           repository.ProductSize(util.CreateRandomProductSize()),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return product
}
