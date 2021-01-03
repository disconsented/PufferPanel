/*
 Copyright 2020 Padduck, LLC
  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at
  	http://www.apache.org/licenses/LICENSE-2.0
  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package services

import (
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type User struct {
	DB *gorm.DB
}

func (us *User) Get(username string) (*models.User, error) {
	model := &models.User{
		Username: username,
	}

	err := us.DB.Where(model).First(model).Error

	if err != nil {
		return nil, err
	}
	return model, nil
}

func (us *User) GetById(id uint) (*models.User, error) {
	model := &models.User{
		ID: id,
	}

	err := us.DB.Where(model).First(model).Error

	if err != nil {
		return nil, err
	}
	return model, nil
}

func (us *User) Login(email string, password string) (user *models.User, sessionToken string, err error) {
	user = &models.User{
		Email: email,
	}

	err = us.DB.Where(user).First(user).Error

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return
	}

	if user.ID == 0 || gorm.IsRecordNotFoundError(err) {
		err = pufferpanel.ErrInvalidCredentials
		return
	}

	if !us.IsValidCredentials(user, password) {
		err = pufferpanel.ErrInvalidCredentials
		return
	}

	sessionToken, err = GenerateSession(user.ID)
	return
}

func (us *User) IsValidCredentials(user *models.User, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)) == nil
}

func (us *User) GetByEmail(email string) (*models.User, error) {
	model := &models.User{
		Email: email,
	}

	err := us.DB.Where(model).First(model).Error

	if err != nil {
		return nil, err
	}
	return model, nil
}

func (us *User) Update(model *models.User) error {
	return us.DB.Save(model).Error
}

func (us *User) Delete(model *models.User) (err error) {
	var trans = us.DB.Begin()
	defer trans.RollbackUnlessCommitted()

	trans.Delete(models.Permissions{}, "user_id = ?", model.ID)
	trans.Delete(models.Client{}, "user_id = ?", model.ID)

	err = trans.Delete(model).Error
	if err != nil {
		return
	}

	return trans.Commit().Error
}

func (us *User) Create(user *models.User) error {
	return us.DB.Create(user).Error
}

func (us *User) ChangePassword(username string, newPass string) error {
	user, err := us.Get(username)

	if err != nil {
		return err
	}

	err = user.SetPassword(newPass)
	if err != nil {
		return err
	}
	return us.Update(user)
}

func (us *User) Search(usernameFilter, emailFilter string, pageSize, page uint) (*models.Users, uint, error) {
	users := &models.Users{}

	query := us.DB

	usernameFilter = strings.Replace(usernameFilter, "*", "%", -1)
	emailFilter = strings.Replace(emailFilter, "*", "%", -1)

	if usernameFilter != "" && usernameFilter != "%" {
		query = query.Where("username LIKE ?", usernameFilter)
	}

	if emailFilter != "" && emailFilter != "%" {
		query = query.Where("email LIKE ?", emailFilter)
	}

	var count uint
	err := query.Model(users).Count(&count).Error

	if err != nil {
		return nil, 0, err
	}

	res := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(users)

	return users, count, res.Error
}
