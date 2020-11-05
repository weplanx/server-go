package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-curd/operates"
	"github.com/kainonly/gin-helper/hash"
	"github.com/kainonly/gin-helper/res"
	"gorm.io/gorm"
	"taste-api/application/cache"
	"taste-api/application/common"
	"taste-api/application/model"
)

type Controller struct {
}

type originListsBody struct {
	operates.OriginListsBody
}

func (c *Controller) OriginLists(ctx *gin.Context, i interface{}) interface{} {
	app := common.Inject(i)
	var body originListsBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return res.Error(err)
	}
	return app.Curd.
		Originlists(model.Admin{}, body.OriginListsBody).
		OrderBy([]string{"create_time desc"}).
		Field([]string{"id", "username", "role", "call", "email", "phone", "avatar", "status"}).
		Exec()
}

type listsBody struct {
	operates.ListsBody
}

func (c *Controller) Lists(ctx *gin.Context, i interface{}) interface{} {
	app := common.Inject(i)
	var body listsBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return res.Error(err)
	}
	return app.Curd.
		Lists(model.Admin{}, body.ListsBody).
		OrderBy([]string{"create_time desc"}).
		Field([]string{"id", "username", "role", "call", "email", "phone", "avatar", "status"}).
		Exec()
}

type getBody struct {
	operates.GetBody
}

func (c *Controller) Get(ctx *gin.Context, i interface{}) interface{} {
	app := common.Inject(i)
	var body getBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return res.Error(err)
	}
	return app.Curd.
		Get(model.Admin{}, body.GetBody).
		Field([]string{"id", "username", "role", "call", "email", "phone", "avatar", "status"}).
		Exec()
}

type addBody struct {
	Username string `binding:"required,min=4,max=20"`
	Password string `binding:"required,min=12,max=20"`
	Role     string `binding:"required"`
	Email    string
	Phone    string
	Call     string
	Avatar   string
	Status   bool
}

func (c *Controller) Add(ctx *gin.Context, i interface{}) interface{} {
	app := common.Inject(i)
	var body addBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return res.Error(err)
	}
	var password string
	password, err = hash.Make(body.Password, hash.Option{})
	if err != nil {
		return res.Error(err)
	}
	data := model.AdminBasic{
		Username: body.Username,
		Password: password,
		Email:    body.Email,
		Phone:    body.Phone,
		Call:     body.Call,
		Avatar:   body.Avatar,
		Status:   body.Status,
	}
	return app.Curd.
		Add().
		After(func(tx *gorm.DB) error {
			roleData := model.AdminRoleAssoc{
				Username: body.Username,
				RoleKey:  body.Role,
			}
			err = tx.Create(&roleData).Error
			if err != nil {
				return err
			}
			clearcache(app.Cache)
			return nil
		}).
		Exec(&data)
}

type editBody struct {
	operates.EditBody
	Username string
	Password string `binding:"min=12,max=20"`
	Role     string `binding:"required_if=switch false"`
	Email    string
	Phone    string
	Call     string
	Avatar   string
	Status   bool
}

func (c *Controller) Edit(ctx *gin.Context, i interface{}) interface{} {
	app := common.Inject(i)
	var body editBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return res.Error(err)
	}
	var password string
	if !body.Switch {
		if body.Password != "" {
			password, err = hash.Make(body.Password, hash.Option{})
			if err != nil {
				return res.Error(err)
			}
		}
	}
	data := model.AdminBasic{
		Username: body.Username,
		Password: password,
		Email:    body.Email,
		Phone:    body.Phone,
		Call:     body.Call,
		Avatar:   body.Avatar,
		Status:   body.Status,
	}
	return app.Curd.
		Edit(model.AdminBasic{}, body.EditBody).
		After(func(tx *gorm.DB) error {
			if !body.Switch {
				roleData := model.AdminRoleAssoc{
					Username: body.Username,
					RoleKey:  body.Role,
				}
				err = tx.Create(&roleData).Error
				if err != nil {
					return err
				}
			}
			clearcache(app.Cache)
			return nil
		}).
		Exec(data)
}

type deleteBody struct {
	operates.DeleteBody
}

func (c *Controller) Delete(ctx *gin.Context, i interface{}) interface{} {
	app := common.Inject(i)
	var body deleteBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return res.Error(err)
	}
	return app.Curd.
		Delete(model.AdminBasic{}, body.DeleteBody).
		After(func(tx *gorm.DB) error {
			clearcache(app.Cache)
			return nil
		}).
		Exec()
}

func clearcache(cache *cache.Model) {
	cache.AdminClear()
}
