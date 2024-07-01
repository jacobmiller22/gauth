package service

import (
	"errors"
	"fmt"
	"math/rand"
)

type AuthStore struct {
	Requests map[string]AuthReq
	Codes    map[string]string
}

type AuthReq struct {
	Id           string
	ClientId     string
	RedirectUri  string
	ResponseType string
	Scope        string
	State        string
}

var store *AuthStore

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randSeq(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GetLoginUrl() string {
	return "http://localhost:8080/login"
}

func MakeStore() (bool, error) {
	if store != nil {
		return false, errors.New("Store already exists")
	}
	store = &AuthStore{
		Requests: map[string]AuthReq{},
		Codes:    map[string]string{},
	}
	return true, nil
}

func (store *AuthStore) SaveRequest(id string, req AuthReq) {
	fmt.Printf("Before storing: %+v\n", store.Requests)
	store.Requests[id] = req
	fmt.Printf("after storing: %+v\n", store.Requests)
}

func (store *AuthStore) SaveRequestCode(code string, id string) {
	store.Codes[code] = id
}

func (store *AuthStore) GetRequestById(id string) (AuthReq, error) {
	if val, ok := store.Requests[id]; ok {
		return val, nil
	}
	return AuthReq{}, fmt.Errorf("Error while getting request for id: %s", id)
}

func (store *AuthStore) GetRequestByCode(code string) (AuthReq, error) {
	if val, ok := store.Codes[code]; ok {
		return store.GetRequestById(val)
	}
	return AuthReq{}, fmt.Errorf("Error while getting request for code: %s", code)
}

func MakeAuthReq(clientId string, redirectUri string, responseType string, scope string, state string) *AuthReq {
	id := randSeq(12)
	return &AuthReq{
		Id:           id,
		ClientId:     clientId,
		RedirectUri:  redirectUri,
		ResponseType: responseType,
		State:        state,
	}
}

func (aq *AuthReq) ForcefullyValidate() error {
	return nil
}

func (aq *AuthReq) Save() string {
	store.SaveRequest(aq.Id, *aq)
	return aq.Id
}

func (aq *AuthReq) GenerateCode() string {
	code := randSeq(12)
	store.SaveRequestCode(code, aq.Id)
	return code
}

func GetAuthReqById(id string) (AuthReq, error) {
	authReq, err := store.GetRequestById(id)
	if err != nil {
		return AuthReq{}, errors.New("bad auth req id")
	}
	return authReq, nil
}

func GetAuthReqByCode(code string) (AuthReq, error) {
	authReq, err := store.GetRequestByCode(code)
	if err != nil {
		return AuthReq{}, errors.New("bad auth req id")
	}
	return authReq, nil
}
