package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/Xinnan-Alex/simplebank/db/mock"
	db "github.com/Xinnan-Alex/simplebank/db/sqlc"
	"github.com/Xinnan-Alex/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAccount(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name         string
		accountID    int64
		buildStubs   func(store *mockdb.MockStore)
		checkReponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).            //expect function to return 1 time
					Return(account, nil) //expect funtion to return specific value
			},
			checkReponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check response
				//checking the returned status code
				require.Equal(t, http.StatusOK, recorder.Code)
				//checking the content in the body
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		// TODO: add more cases
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).                           //expect function to return 1 time
					Return(db.Account{}, sql.ErrNoRows) //expect funtion to return specific value
			},
			checkReponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check response
				//checking the returned status code
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).                             //expect function to return 1 time
					Return(db.Account{}, sql.ErrConnDone) //expect funtion to return specific value
			},
			checkReponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check response
				//checking the returned status code
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0) //expect function to return 1 time
			},
			checkReponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check response
				//checking the returned status code
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			//start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkReponse(t, recorder)
		})

	}

}

func TestCreateAccount(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name           string
		body           gin.H
		buildStubs     func(store *mockdb.MockStore)
		expectedStatus int
	}{
		{
			name: "OK",
			body: gin.H{
				"owner":    account.Owner,
				"currency": account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    account.Owner,
					Balance:  0,
					Currency: account.Currency,
				}

				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "InvalidOwner",
			body: gin.H{
				"owner":    "",
				"currency": account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    account.Owner,
					Balance:  0,
					Currency: account.Currency,
				}

				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "InvalidCurrency",
			body: gin.H{
				"owner":    account.Owner,
				"currency": "YEN",
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    account.Owner,
					Balance:  0,
					Currency: account.Currency,
				}

				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "InternalError",
			body: gin.H{
				"owner":    account.Owner,
				"currency": account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    account.Owner,
					Balance:  0,
					Currency: account.Currency,
				}

				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/accounts"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			require.Equal(t, recorder.Code, tc.expectedStatus)

		})
	}

}

func TestListAccounts(t *testing.T) {
	accounts_size := 5
	accounts := make([]db.Account, accounts_size)

	//create a list of accounts with the size of accounts_size
	for i := 0; i < accounts_size; i++ {
		accounts[i] = randomAccount()
	}

	type Query struct {
		pageSize int
		pageID   int
	}

	testCases := []struct {
		name         string
		query        Query
		buildStubs   func(store *mockdb.MockStore)
		checkReponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				pageSize: accounts_size,
				pageID:   1,
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.ListAccountsParams{
					Limit:  int32(accounts_size),
					Offset: 0,
				}
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(accounts, nil)
			},
			checkReponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccounts(t, recorder.Body, accounts)
			},
		},
		{
			name: "Invalid Page ID",
			query: Query{
				pageSize: accounts_size,
				pageID:   -1,
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.ListAccountsParams{
					Limit:  int32(accounts_size),
					Offset: 0,
				}
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkReponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Page Size",
			query: Query{
				pageSize: -1,
				pageID:   1,
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.ListAccountsParams{
					Limit:  int32(accounts_size),
					Offset: 0,
				}
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkReponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Internal Error",
			query: Query{
				pageSize: accounts_size,
				pageID:   1,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Account{}, sql.ErrConnDone)
			},
			checkReponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/accounts"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkReponse(recorder)
		})

	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, gotAccount, account)
}

func requireBodyMatchAccounts(t *testing.T, body *bytes.Buffer, accounts []db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []db.Account
	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)
	require.Equal(t, gotAccounts, accounts)

}
