/*
 Copyright 2019 Padduck, LLC
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
)

type Template struct {
	DB *gorm.DB
}

func (t *Template) GetAll() (*models.Templates, error) {
	templates := &models.Templates{}
	err := t.DB.Find(&templates).Error
	if err != nil {
		return nil, err
	}
	//because we don't want to return a ton of data, we'll only return a few select fields
	replacement := make(models.Templates, len(*templates))

	for k, v := range *templates {
		replacement[k] = &models.Template{
			Name: v.Name,
			Template: pufferpanel.Template{
				Server: pufferpanel.Server{
					Display: v.Template.Server.Display,
					Type:    v.Template.Server.Type,
				},
			},
		}
	}
	return &replacement, err
}

func (t *Template) Get(name string) (*models.Template, error) {
	template := &models.Template{
		Name:  name,
	}
	err := t.DB.Find(&template).Error
	if err != nil {
		return nil, err
	}
	return template, err
}


func (t *Template) Save(template *models.Template) error {
	return t.DB.Save(template).Error
}