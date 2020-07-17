package client

import (
	"errors"
	"fmt"
)

type Credentials struct {
	Permissions []string `json:"permissions,omitempty"`
}

type Login struct {
	Credentials Credentials `json:"credentials,omitempty"`
	PersonalGit bool        `json:"personalGit,omitempty"`
	Permissions []string    `json:"permissions,omitempty"`
	IDP         IDP         `json:"idp,omitempty"`
}

type ShortProfile struct {
	UserName string `json:"userName,omitempty"`
}

type Personal struct {
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	CompanyName string `json:"companyName,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	Country     string `json:"country,omitempty"`
}

type User struct {
	ID             string              `json:"_id,omitempty"`
	UserName       string              `json:"userName"`
	Email          string              `json:"email"`
	Personal       *Personal           `json:"personal,omitempty"`
	Roles          []string            `json:"roles,omitempty"`
	DefaultAccount int                 `json:"defaultAccount,omitempty"`
	Account        []Account           `json:"account,omitempty"`
	Status         string              `json:"status,omitempty"`
	RegisterDate   string              `json:"register_date,omitempty"`
	HasPassword    bool                `json:"hasPassword,omitempty"`
	Notifications  []NotificationEvent `json:"notifications,omitempty"`
	ShortProfile   ShortProfile        `json:"shortProfile,omitempty"`
	Logins         []Login             `json:"logins,omitempty"`
	InviteURL      string              `json:"inviteUrl,omitempty"`
}

type NewUser struct {
	UserName string   `json:"userName"`
	Email    string   `json:"email"`
	Logins   []Login  `json:"logins,omitempty"`
	Roles    []string `json:"roles,omitempty"`
	Account  []string `json:"account,omitempty"`
}

type UserAccounts struct {
	UserName string    `json:"userName`
	Account  []Account `json:"account`
}

func (client *Client) AddNewUserToAccount(accountId, userName, userEmail string) (*User, error) {

	userDetails := fmt.Sprintf(`{"userName": "%s", "email": "%s"}`, userName, userEmail)

	fullPath := fmt.Sprintf("/accounts/%s/adduser", accountId)

	opts := RequestOptions{
		Path:   fullPath,
		Method: "POST",
		Body:   []byte(userDetails),
	}

	resp, err := client.RequestAPI(&opts)
	if err != nil {
		return nil, err
	}

	var user User

	err = DecodeResponseInto(resp, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (client *Client) AddPendingUser(user *NewUser) (*User, error) {

	body, err := EncodeToJSON(user)
	if err != nil {
		return nil, err
	}
	opts := RequestOptions{
		Path:   "/admin/accounts/addpendinguser",
		Method: "POST",
		Body:   body,
	}

	resp, err := client.RequestAPI(&opts)
	if err != nil {
		return nil, err
	}

	var respUser User

	err = DecodeResponseInto(resp, &respUser)
	if err != nil {
		return nil, err
	}

	return &respUser, nil
}

func (client *Client) ActivateUser(userId string) (*User, error) {

	opts := RequestOptions{
		Path:   fmt.Sprintf("/admin/user/%s/activate", userId),
		Method: "POST",
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var user User

	err = DecodeResponseInto(resp, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (client *Client) SetUserAsAccountAdmin(accountId, userId string) error {

	opts := RequestOptions{
		Path:   fmt.Sprintf("/accounts/%s/%s/admin", accountId, userId),
		Method: "POST",
	}

	_, err := client.RequestAPI(&opts)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) DeleteUserAsAccountAdmin(accountId, userId string) error {

	opts := RequestOptions{
		Path:   fmt.Sprintf("/accounts/%s/%s/admin", accountId, userId),
		Method: "DELETE",
	}

	_, err := client.RequestAPI(&opts)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) GetAllUsers() (*[]User, error) {

	opts := RequestOptions{
		Path:   "/admin/user",
		Method: "GET",
	}

	resp, err := client.RequestAPI(&opts)
	if err != nil {
		return nil, err
	}

	var users []User

	err = DecodeResponseInto(resp, &users)
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func (client *Client) GetUserByID(userId string) (*User, error) {

	users, err := client.GetAllUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range *users {
		if user.ID == userId {
			return &user, nil
		}
	}

	return nil, errors.New(fmt.Sprint("[ERROR] User with ID %s wasn't found.", userId))
}

func (client *Client) DeleteUser(userName string) error {

	opts := RequestOptions{
		Path:   fmt.Sprintf("/admin/user/%s", userName),
		Method: "DELETE",
	}

	_, err := client.RequestAPI(&opts)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) DeleteUserFromAccount(accountId, userId string) error {

	opts := RequestOptions{
		Path:   fmt.Sprintf("/accounts/%s/%s", accountId, userId),
		Method: "DELETE",
	}

	_, err := client.RequestAPI(&opts)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) UpdateUserAccounts(userId string, accounts []Account) error {

	// API call '/accounts/{accountId}/{userId}/adduser' doesn't work

	user, err := client.GetUserByID(userId)
	if err != nil {
		return err
	}

	postUser := UserAccounts{
		UserName: user.UserName,
		Account:  accounts,
	}

	body, err := EncodeToJSON(postUser)
	if err != nil {
		return err
	}

	opts := RequestOptions{
		Path:   "/admin/user/account",
		Method: "POST",
		Body:   body,
	}

	_, err = client.RequestAPI(&opts)
	if err != nil {
		return err
	}

	return nil
}
