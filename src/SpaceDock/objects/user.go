/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
 */

package objects

import (
    "SpaceDock"
    "errors"
    "github.com/jameskeane/bcrypt"
    "github.com/jinzhu/gorm"
    "time"
)

type User struct {
    gorm.Model

    Username string `gorm:"size:128;unique_index;not null"`
    Email string `gorm:"size:256;unique_index;not null"`
    Public bool
    Password string `gorm:"size:128"`
    Description string `gorm:"size:10000"`
    Confirmation string `gorm:"size:128"`
    PasswordReset string `gorm:"size:128"`
    PasswordResetExpiry time.Time
    Authed bool

    RoleUsers []RoleUser
}

func NewUser(name string, email string, password string) *User {
    user := &User {
        Username: name,
        Email: email,
        Public: false,
        Description: "",
        Confirmation: "",
        PasswordReset: "",
        PasswordResetExpiry: time.Now(),
        Authed: false,
    }
    user.SetPassword(password)
    return user
}

func (user User) SetPassword(password string) {
    salt, _ := bcrypt.Salt()
    user.Password, _ = bcrypt.Hash(password, salt)
    SpaceDock.Database.Save(user)
}

func (user User) IsAuthenticated() bool {
    return user.Authed
}

func (user User) Login() {
    user.Authed = true
    SpaceDock.Database.Save(user)
}

func (user User) Logout() {
    user.Authed = false
    SpaceDock.Database.Save(user)
}

func (user User) UniqueId() interface{} {
    return user.ID
}

func (user User) GetById(id interface{}) error {
    SpaceDock.Database.First(&user, id)
    if user.Username != "" {
        return errors.New("Invalid user ID")
    }
    return nil
}

func (user User) AddRole(name string) Role {
    role := Role {}
    SpaceDock.Database.Where("name = ?", name).First(&role)
    if role.Name == "" {
        role.Name = name
        role.Params = "{}"
        SpaceDock.Database.Save(&role)
    }
    ru := RoleUser{}
    SpaceDock.Database.Where("roleid = ?", role.ID).Where("userid = ?", user.ID).First(&ru)
    if ru.RoleID != role.ID || ru.UserID != user.ID {
        SpaceDock.Database.Save(NewRoleUser(user, role))
    }
    return role
}

func (user User) RemoveRole(name string) {
    role := Role{}
    SpaceDock.Database.Where("name = ?", name).First(&role)
    if role.Name == "" {
        return
    }
    ru := RoleUser{}
    SpaceDock.Database.Where("roleid = ?", role.ID).Where("userid = ?", user.ID).First(&ru)
    if ru.RoleID == role.ID && ru.UserID == user.ID {
        SpaceDock.Database.Delete(&ru)
    }
}

func (user User) GetRoles() []Role {
    value := make([]Role, len(user.RoleUsers))
    for index,element := range user.RoleUsers {
        role := Role {}
        SpaceDock.Database.First(&role, element.RoleID)
        value[index] = role
    }
    return value
}

func (user User) GetAbilities() []Ability {
    count := 0
    for _,element := range user.GetRoles() {
        count = count + len(element.GetAbilities())
    }
    value := make([]Ability, count)
    c := 0
    for _,element := range user.GetRoles() {
        for _,element2 := range element.GetAbilities() {
            value[c] = element2
            c = c + 1
        }
    }
    return value
}