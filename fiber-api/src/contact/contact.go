package contact

import (
	"database/sql"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
)

type Contact struct {
	Id        uint64 `json:"id"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Email     string `json:"email"`
	Telephone string `json:"telephone"`
}

type ContactDao interface {
	Insert(contact Contact) Contact
	Delete(contact Contact)
	GetAll() ([]Contact, error)
}

type ContactDaoS struct {
	Connection *sql.DB
}

func RegisterRouterHook(app *fiber.App) func(s *ContactDaoS) []fiber.Router {
	return func(contactDao *ContactDaoS) []fiber.Router {
		return []fiber.Router{
			app.Get("/", func(ctx *fiber.Ctx) error {
				all, err := contactDao.GetAll()
				if err != nil {
					ctx.Status(500)
					return err
				}
				res, jsonErr := json.Marshal(all)
				if jsonErr != nil {
					ctx.Status(500)
					return err
				}
				return ctx.Send(res)
			}),
			app.Post("/contact", func(ctx *fiber.Ctx) error {
				var entity Contact
				err := json.Unmarshal(ctx.Body(), &entity)
				if err == nil {
					_, err = contactDao.Insert(entity)
				}
				return err
			}),
			app.Delete("/contact/:id", func(ctx *fiber.Ctx) error {
				id, _ := strconv.ParseInt(ctx.Params("id"), 10, 0)
				deletedUser, err := contactDao.Delete(uint64(id))
				if err != nil {
					return err
				}
				if deletedUser == 0 {
					return ctx.SendStatus(404)
				} else {
					return ctx.SendStatus(200)
				}
			}),
		}
	}
}

func (dao *ContactDaoS) Insert(contact Contact) (*Contact, error) {
	err := dao.Connection.QueryRow(`INSERT INTO contact (name, surname, email, telephone) VALUES ($1, $2, $3, $4) RETURNING id`,
		contact.Name, contact.Surname, contact.Email, contact.Telephone).Scan(&contact.Id)
	if err != nil {
		return nil, err
	}
	return &contact, err
}

func (dao *ContactDaoS) Delete(id uint64) (int64, error) {
	res, err := dao.Connection.Exec(`DELETE FROM contact where id = $1`,
		id)
	if err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func (dao *ContactDaoS) GetAll() ([]Contact, error) {
	rows, err := dao.Connection.Query(`SELECT id, name, surname, email, telephone FROM contact`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cnts []Contact
	var tempContact Contact
	for rows.Next() {
		err = rows.Scan(&tempContact.Id, &tempContact.Name, &tempContact.Surname, &tempContact.Email, &tempContact.Telephone)
		if err != nil {
			log.Fatal("Error during row mapping")
			return nil, err
		}
		cnts = append(cnts, tempContact)
	}
	return cnts, nil
}
