package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mock_db "github.com/julianinsua/the_simp_bank/db/mock"
	"github.com/julianinsua/the_simp_bank/internal/database"
	"github.com/julianinsua/the_simp_bank/token"
	"github.com/julianinsua/the_simp_bank/util"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// Create the necessary structs to match the parameters passed to the mock database for user creation
type eqCreateUserParamsMatcher struct {
	arg      database.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(database.CreateUserParams)
	if !ok {
		return false
	}
	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches args %v with password %v", e.arg, e.password)
}

func usrChk(args database.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{args, password}
}

func TestCreateUser(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuthFunc func(t *testing.T, request *http.Request, maker token.PASETOMaker)
		buildStubs    func(store *mock_db.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username": user.Username,
				"password": password,
				"fullName": user.FullName,
				"email":    user.Email,
			},

			setupAuthFunc: func(t *testing.T, request *http.Request, maker token.PASETOMaker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				usrParams := database.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().CreateUser(gomock.Any(), usrChk(usrParams, password)).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username": user.Username,
				"password": password,
				"fullName": user.FullName,
				"email":    user.Email,
			},
			setupAuthFunc: func(t *testing.T, request *http.Request, maker token.PASETOMaker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(database.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DuplicateUsername",
			body: gin.H{
				"username": user.Username,
				"password": password,
				"fullName": user.FullName,
				"email":    user.Email,
			},
			setupAuthFunc: func(t *testing.T, request *http.Request, maker token.PASETOMaker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.
					EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(database.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username": "invalid-username#1",
				"password": password,
				"fullName": user.FullName,
				"email":    user.Email,
			},
			setupAuthFunc: func(t *testing.T, request *http.Request, maker token.PASETOMaker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.
					EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"username": user.Username,
				"password": password,
				"fullName": user.FullName,
				"email":    "invalid-email",
			},
			setupAuthFunc: func(t *testing.T, request *http.Request, maker token.PASETOMaker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.
					EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "PasswordTooShort",
			body: gin.H{
				"username": user.Username,
				"password": "123",
				"fullName": user.FullName,
				"email":    user.Email,
			},
			setupAuthFunc: func(t *testing.T, request *http.Request, maker token.PASETOMaker) {
				addAuthorization(t, request, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.
					EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock_db.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data)) // Turn JSON into bytes
			require.NoError(t, err)

			tc.setupAuthFunc(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomUser(t *testing.T) (user database.User, password string) {
	password = util.RandomString(8)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = database.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user database.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser database.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}

func TestLoginUser(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mock_db.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(user, nil)
				store.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "badUsername",
			body: gin.H{"username": "user-name#1", "password": password},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "badPassword",
			body: gin.H{"username": user.Username, "password": "badPassword"},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(1).Return(user, nil)
				store.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "userNotFound",
			body: gin.H{"username": user.Username, "password": password},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(1).Return(database.User{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "internalError",
			body: gin.H{"username": user.Username, "password": password},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(1).Return(database.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock_db.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users/login"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)

		})
	}
}
