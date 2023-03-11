package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/thehaung/simplebank/db/mock"
	db "github.com/thehaung/simplebank/db/sqlc"
	"github.com/thehaung/simplebank/util/randutil"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		Name          string
		AccountID     int64
		BuildStubs    func(store *mockdb.MockStore)
		CheckResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			Name:      "InvalidID",
			AccountID: -1,
			BuildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check resp
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			Name:      "OK",
			AccountID: account.ID,
			BuildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check resp
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			Name:      "NotFound",
			AccountID: account.ID,
			BuildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, sql.ErrNoRows)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check resp
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			Name:      "InternalError",
			AccountID: account.ID,
			BuildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, sql.ErrConnDone)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check resp
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.BuildStubs(store)

			// start server
			server := NewHttpServer(store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts/%d", tc.AccountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.CheckResponse(t, recorder)
		})
	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:       randutil.IntWithRange(1, 1000),
		Owner:    randutil.Owner(),
		Balance:  randutil.Money(),
		Currency: randutil.Currency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)

	require.Equal(t, account, gotAccount)
}

func TestListAccount(t *testing.T) {
	n := 5
	accounts := make([]db.Account, n)
	for i := 0; i < n; i++ {
		accounts[i] = randomAccount()
	}

	type QueryParams struct {
		PageID   int
		PageSize int
	}

	testCases := []struct {
		Name          string
		Query         QueryParams
		BuildStubs    func(store *mockdb.MockStore)
		CheckResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			Name: "InvalidPageID",
			Query: QueryParams{
				PageID:   0,
				PageSize: 9,
			},
			BuildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check resp
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			Name: "InvalidPageSize",
			Query: QueryParams{
				PageID:   1,
				PageSize: 9999,
			},
			BuildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check resp
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			Name: "Ok",
			Query: QueryParams{
				PageID:   1,
				PageSize: 5,
			},
			BuildStubs: func(store *mockdb.MockStore) {
				arg := db.ListAccountsParams{
					Limit:  5,
					Offset: 0,
				}
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(accounts, nil)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check resp
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchListAccount(t, recorder.Body, accounts)
			},
		},
		{
			Name: "InternalServerError",
			Query: QueryParams{
				PageID:   1,
				PageSize: 5,
			},
			BuildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, sql.ErrConnDone)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check resp
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.BuildStubs(store)

			// start server
			server := NewHttpServer(store)
			recorder := httptest.NewRecorder()
			url := "/accounts"
			request, err := http.NewRequest(http.MethodGet, url, nil)

			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.Query.PageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.Query.PageSize))
			request.URL.RawQuery = q.Encode()

			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.CheckResponse(t, recorder)
		})
	}
}

func requireBodyMatchListAccount(t *testing.T, body *bytes.Buffer, listAccount []db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotListAccount []db.Account
	err = json.Unmarshal(data, &gotListAccount)
	require.NoError(t, err)

	require.Equal(t, listAccount, gotListAccount)
}

func TestCreateAccount(t *testing.T) {
	_ = randomAccount()

	testCases := []struct {
		Name          string
		Body          gin.H
		BuildStubs    func(store *mockdb.MockStore)
		CheckResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			Name: "BadRequest",
			Body: gin.H{},
			BuildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			CheckResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			Name: "InvalidCurrency",
			Body: gin.H{
				"currency": "JP",
			},
			BuildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			CheckResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			Name: "Ok",
			Body: gin.H{
				"owner":    "Hehe",
				"currency": "USD",
			},
			BuildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1)
			},
			CheckResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			Name: "InternalServerError",
			Body: gin.H{
				"owner":    "Hehe",
				"currency": "USD",
			},
			BuildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			CheckResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.BuildStubs(store)
			// start server
			server := NewHttpServer(store)
			recorder := httptest.NewRecorder()
			url := "/accounts"
			// Marshal body data to JSON
			data, err := json.Marshal(tc.Body)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.CheckResponse(recorder)
		})
	}
}
