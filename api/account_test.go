package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mock_db "github.com/julianinsua/the_simp_bank/db/mock"
	"github.com/julianinsua/the_simp_bank/internal/database"
	"github.com/julianinsua/the_simp_bank/token"
	"github.com/julianinsua/the_simp_bank/util"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	user, _ := util.RandomUser()
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		accountID     string
		setupAuthFunc func(t *testing.T, request *http.Request, maker token.PASETOMaker)
		buildStubs    func(store *mock_db.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID.String(),
			setupAuthFunc: func(t *testing.T, request *http.Request, maker token.PASETOMaker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				compareResponse(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID.String(),
			setupAuthFunc: func(t *testing.T, request *http.Request, maker token.PASETOMaker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(database.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "Internal",
			accountID: account.ID.String(),
			setupAuthFunc: func(t *testing.T, request *http.Request, maker token.PASETOMaker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(database.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "BadRequest",
			accountID: "asdf",
			setupAuthFunc: func(t *testing.T, request *http.Request, maker token.PASETOMaker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock_db.NewMockStore(ctrl)

			// Build stubs
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts/%v", tc.accountID)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)

			tc.setupAuthFunc(t, request, server.tokenMaker)
			//check response in the recorder
			tc.checkResponse(t, recorder)
		})
	}

}

func compareResponse(t *testing.T, body *bytes.Buffer, account database.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var getAccount database.Account
	err = json.Unmarshal(data, &getAccount)
	require.NoError(t, err)

	require.Equal(t, getAccount, account)
}

func randomAccount(owner string) database.Account {
	return database.Account{
		ID:       uuid.New(),
		Owner:    owner,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
