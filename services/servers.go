package services

import (
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/database"
	"github.com/pufferpanel/pufferpanel/models"
	uuid2 "github.com/satori/go.uuid"
	"strings"
)

type ServerService struct {
	db *gorm.DB
}

func GetServerService() (*ServerService, error) {
	db, err := database.GetConnection()
	if err != nil {
		return nil, err
	}

	service := &ServerService{
		db: db,
	}

	return service, nil
}

func (ss *ServerService) Search(nodeId uint, nameFilter string, pageSize, page uint) (*models.Servers, error) {
	servers := &models.Servers{}

	query := ss.db.Offset((page - 1) * pageSize).Limit(pageSize)

	if nodeId != 0 {
		query = query.Where(&models.Server{NodeID: nodeId})
	}

	nameFilter = strings.Replace(nameFilter, "*", "%", -1)

	if nameFilter != "" && nameFilter != "%" {
		query = query.Where("name LIKE ?", nameFilter)
	}

	res := query.Find(servers)

	return servers, res.Error
}

//TODO: Waiting on user objects with rights to implement correctly
func (ss *ServerService) GetForUser(userId uint) (*models.Servers, error) {
	servers := &models.Servers{}

	res := ss.db.Find(servers)

	return servers, res.Error
}

func (ss *ServerService) Get(id string) (*models.Server, bool, error) {
	model := &models.Server{
		Identifier: id,
	}

	res := ss.db.First(model)

	return model, model.ID != 0, res.Error
}

func (ss *ServerService) Update(model *models.Server) error {
	res := ss.db.Save(model)
	return res.Error
}

func (ss *ServerService) Delete(id uint) error {
	model := &models.Server{
		ID: id,
	}

	res := ss.db.Delete(model)
	return res.Error
}

func (ss *ServerService) Create(model *models.Server, serverData interface{}) (err error) {
	uuid := uuid2.NewV4()
	generatedId := strings.ToUpper(uuid.String())[0:8]
	model.Identifier = generatedId

	conn := ss.db.Begin()
	successful := false

	defer func() {
		if successful && err == nil {
			conn.Commit()
		} else {
			conn.Rollback()
		}
	}()

	res := conn.Create(model)
	if res.Error != nil {
		err = res.Error
		return
	}

	return
}