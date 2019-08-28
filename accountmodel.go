package main

type Profile struct {
	Id           int    `json:"-"`
	UserName     string `json:"userName"`
	InviteCode   string `json:"inviteCode"`
	InvitedUsers int    `json:"invitedUsers"`
}

type AccountModel struct {
	ProfilesByID         map[int]*Profile
	ProfilesByInviteCode map[string]*Profile
}

func NewAccountModel() *AccountModel {
	profilesByID := make(map[int]*Profile, 0)
	profilesByInviteCode := make(map[string]*Profile, 0)

	return &AccountModel{
		ProfilesByID:         profilesByID,
		ProfilesByInviteCode: profilesByInviteCode,
	}
}

func (accountModel *AccountModel) AddProfile(userName string, inviteCode string) *Profile {
	id := accountModel.getNextID()
	profile := &Profile{
		Id:           id,
		UserName:     userName,
		InviteCode:   inviteCode,
		InvitedUsers: 0,
	}

	accountModel.ProfilesByID[id] = profile
	accountModel.ProfilesByInviteCode[inviteCode] = profile

	return profile
}

func (accountModel AccountModel) getNextID() int {
	return len(accountModel.ProfilesByID) + 1
}
