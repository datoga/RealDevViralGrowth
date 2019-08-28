package main

import (
	"math/rand"
	"strings"
	"time"
)

type AccountService struct {
	randomNumbers   *rand.Rand
	inviteCodeChars string
	accountModel    *AccountModel
}

func NewAccountService(inviteCodeChars string, accountModel *AccountModel) *AccountService {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	return &AccountService{
		randomNumbers:   r,
		inviteCodeChars: inviteCodeChars,
		accountModel:    accountModel,
	}
}

func (service AccountService) GetProfile(id int) (*Profile, bool) {
	profile, found := service.accountModel.ProfilesByID[id]

	return profile, found
}

func (service *AccountService) GetReferer(inviteCode string) (*Profile, bool) {
	referer, found := service.accountModel.ProfilesByInviteCode[inviteCode]

	return referer, found
}

func (service *AccountService) AddProfile(userName string, referer *Profile) *Profile {
	newCodeFree := false
	var newInviteCode string

	for newCodeFree == false {
		newInviteCode = service.generateInviteCode()

		_, ok := service.accountModel.ProfilesByInviteCode[newInviteCode]

		if !ok {
			newCodeFree = true
		}
	}

	newProfile := service.accountModel.AddProfile(userName, newInviteCode)

	if referer != nil {
		referer.InvitedUsers = referer.InvitedUsers + 1
	}

	return newProfile
}

func (service AccountService) generateInviteCode() string {

	var str strings.Builder

	for range service.inviteCodeChars {
		char := service.generateRandomChar()

		str.WriteString(char)
	}

	return str.String()
}

func (service AccountService) generateRandomChar() string {
	tam := len(service.inviteCodeChars)

	pos := service.randomNumbers.Intn(tam)

	return (string)(service.inviteCodeChars[pos])
}
