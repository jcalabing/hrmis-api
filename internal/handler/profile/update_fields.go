package profile

import (
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/jcalabing/hrmis-api/internal/common/errors"
	"github.com/jcalabing/hrmis-api/internal/model"
	"gorm.io/gorm"
)

func UpdateFields(h *handler, c *fiber.Ctx, profile model.User, body interface{}) error {
	fieldlist := []string{
		"Firstname",
		"Surname",
		"Civilstatus",
		"Sex",
		"Dateofbirth",
		"Email",
		"Mobilenumber",
		"Issuedid",
		"Issuedidnumber",
		"Issuedidissuance",
		"Middlename",
		"Extname",
		"Placeofbirth",
		"Height",
		"Weight",
		"Bloodtype",
		"Gsis",
		"Pagibig",
		"Philhealth",
		"Sss",
		"Tin",
		"Agencyemployeenumber",
		"Country",
		"Residentialhouse",
		"Residentialstreet",
		"Residentialsubdivision",
		"Residentialbarangay",
		"Residentialmunicipality",
		"Residentialprovince",
		"Residentialzip",
		"Permanenthouse",
		"Permanentstreet",
		"Permanentsubdivision",
		"Permanentbarangay",
		"Permanentmunicipality",
		"Permanentprovince",
		"Permanentzip",
		"Telephonenumber",
		"Spousesurname",
		"Spousefirstname",
		"Spousemiddlename",
		"Spousextname",
		"Spousenumber",
		"Spouseoccupation",
		"Spouseemployer",
		"Spousebusinessaddress",
		"Fathersurname",
		"Fatherfirstname",
		"Fathermiddlename",
		"Fatherextname",
		"Mothersurname",
		"Motherfirstname",
		"Mothermiddlename",
		"Relatedthirddegree",
		"Relatedfourthdegree",
		"Relatedfourthdegreedetails",
		"Guiltyofoffense",
		"Guiltyofoffensedetail",
		"Criminallycharged",
		"Criminallychargeddate",
		"Criminallychargeddetails",
		"Criminalconvicted",
		"Criminalconvicteddetails",
		"Seperatedfromservice",
		"Seperatedfromservicedetails",
		"Candidate",
		"Candidatedetails",
		"Electionresigned",
		"Electionresigneddetails",
		"Immigrantstatus",
		"Immigrantstatusdetails",
		"Indigenousmember",
		"Indigenousmemberdetails",
		"Pwd",
		"Pwddetails",
		"Soloparent",
		"Soloparentdetails",
		"Citizenship",
		"Dualcitizenship",
	}

	for _, fieldName := range fieldlist {
		fieldValue := reflect.ValueOf(body).Elem().FieldByName(fieldName)
		if fieldValue.IsValid() {
			userfield := model.UserField{
				Key:    fieldName,
				Value:  fmt.Sprintf("%v", fieldValue.Interface()),
				UserID: profile.ID,
			}
			// Look for UserField
			var oldUserField model.UserField
			if err := h.DB.Where("user_id = ? AND key = ?", userfield.UserID, userfield.Key).First(&oldUserField).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					// If key does not exist, create a new UserField
					if err := h.DB.Create(&userfield).Error; err != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
							fiber.StatusInternalServerError,
							"",
							"Error occurred while creating the new user field.",
						))
					}
				} else {
					return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
						fiber.StatusInternalServerError,
						"",
						"Error occurred while querying the database.",
					))
				}
			} else {
				// If key exists, update the existing UserField
				if err := h.DB.Model(&oldUserField).Updates(&userfield).Error; err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
						fiber.StatusInternalServerError,
						"",
						"Error occurred while updating the user field.",
					))
				}
			}
		}
	}

	return nil
}
