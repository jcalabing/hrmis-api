package profile

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jcalabing/hrmis-api/internal/common/errors"
	"github.com/jcalabing/hrmis-api/internal/common/functions"
	"github.com/jcalabing/hrmis-api/internal/model"
)

type CreateRequestBody struct {
	Firstname                   string `json:"firstname" validate:"required,nonempty"`
	Surname                     string `json:"surname" validate:"required,nonempty"`
	Civilstatus                 string `json:"civilstatus" validate:"required,nonempty"`
	Sex                         string `json:"sex" validate:"required,nonempty"`
	Dateofbirth                 string `json:"dateofbirth" validate:"required,nonempty"`
	Email                       string `json:"email" validate:"required,nonempty"`
	Mobilenumber                string `json:"mobilenumber" validate:"required,nonempty"`
	Issuedid                    string `json:"issuedid" validate:"required,nonempty"`
	Issuedidnumber              string `json:"issuedidnumber" validate:"required,nonempty"`
	Issuedidissuance            string `json:"issuedidissuance" validate:"required,nonempty"`
	Middlename                  string `json:"middlename"`
	Extname                     string `json:"extname"`
	Placeofbirth                string `json:"placeofbirth"`
	Height                      string `json:"height"`
	Weight                      string `json:"weight"`
	Bloodtype                   string `json:"bloodtype"`
	Gsis                        string `json:"gsis"`
	Pagibig                     string `json:"pagibig"`
	Philhealth                  string `json:"philhealth"`
	Sss                         string `json:"sss"`
	Tin                         string `json:"tin"`
	Agencyemployeenumber        string `json:"agencyemployeenumber"`
	Country                     string `json:"country"`
	Residentialhouse            string `json:"residentialhouse"`
	Residentialstreet           string `json:"residentialstreet"`
	Residentialsubdivision      string `json:"residentialsubdivision"`
	Residentialbarangay         string `json:"residentialbarangay"`
	Residentialmunicipality     string `json:"residentialmunicipality"`
	Residentialprovince         string `json:"residentialprovince"`
	Residentialzip              string `json:"residentialzip"`
	Permanenthouse              string `json:"permanenthouse"`
	Permanentstreet             string `json:"permanentstreet"`
	Permanentsubdivision        string `json:"permanentsubdivision"`
	Permanentbarangay           string `json:"permanentbarangay"`
	Permanentmunicipality       string `json:"permanentmunicipality"`
	Permanentprovince           string `json:"permanentprovince"`
	Permanentzip                string `json:"permanentzip"`
	Telephonenumber             string `json:"telephonenumber"`
	Spousesurname               string `json:"spousesurname"`
	Spousefirstname             string `json:"spousefirstname"`
	Spousemiddlename            string `json:"spousemiddlename"`
	Spousextname                string `json:"spousextname"`
	Spouseoccupation            string `json:"spouseoccupation"`
	Spouseemployer              string `json:"spouseemployer"`
	Spousebusinessaddress       string `json:"spousebusinessaddress"`
	Fathersurname               string `json:"fathersurname"`
	Fatherfirstname             string `json:"fatherfirstname"`
	Fathermiddlename            string `json:"fathermiddlename"`
	Fatherextname               string `json:"fatherextname"`
	Mothersurname               string `json:"mothersurname"`
	Motherfirstname             string `json:"motherfirstname"`
	Mothermiddlename            string `json:"mothermiddlename"`
	Relatedthirddegree          string `json:"relatedthirddegree"`
	Relatedfourthdegree         string `json:"relatedfourthdegree"`
	Relatedfourthdegreedetails  string `json:"relatedfourthdegreedetails"`
	Guiltyofoffense             string `json:"guiltyofoffense"`
	Guiltyofoffensedetail       string `json:"guiltyofoffensedetail"`
	Criminallycharged           string `json:"criminallycharged"`
	Criminallychargeddetails    string `json:"criminallychargeddetails"`
	Criminalconvicted           string `json:"criminalconvicted"`
	Criminalconvicteddetails    string `json:"criminalconvicteddetails"`
	Seperatedfromservice        string `json:"seperatedfromservice"`
	Seperatedfromservicedetails string `json:"seperatedfromservicedetails"`
	Candidate                   string `json:"candidate"`
	Candidatedetails            string `json:"candidatedetails"`
	Electionresigned            string `json:"electionresigned"`
	Electionresigneddetails     string `json:"electionresigneddetails"`
	Immigrantstatus             string `json:"immigrantstatus"`
	Immigrantstatusdetails      string `json:"immigrantstatusdetails"`
	Indigenousmember            string `json:"indigenousmember"`
	Indigenousmemberdetails     string `json:"indigenousmemberdetails"`
	Pwd                         string `json:"pwd"`
	Pwddetails                  string `json:"pwddetails"`
	Soloparent                  string `json:"soloparent"`
	Soloparentdetails           string `json:"soloparentdetails"`
	Filipino                    string `json:"filipino"`
	Dualcitizenship             string `json:"dualcitizenship"`
	Bybirth                     string `json:"bybirth"`
	Naturalized                 string `json:"naturalized"`
	Education                   string `json:"education"`
}

func (h *handler) UpdateProfile(c *fiber.Ctx) error {
	body := new(CreateRequestBody)

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
			fiber.StatusBadRequest,
			"",
			"",
		))
	}

	validate := validator.New()
	validate.RegisterValidation("nonempty", functions.ValidateNonEmpty)

	if err := validate.Struct(body); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Field())
		}
		errorMsg := "Kindly check the following fields: " + strings.Join(validationErrors, ", ")

		return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
			fiber.StatusBadRequest,
			errorMsg,
			validationErrors,
		))
	}

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
		"Filipino",
		"Dualcitizenship",
		"Bybirth",
		"Naturalized",
	}

	profile := c.Locals("profile").(model.User)
	for _, fieldName := range fieldlist {
		fieldValue := reflect.ValueOf(body).Elem().FieldByName(fieldName)
		if fieldValue.IsValid() {
			userfield := model.UserField{
				Key:    fieldName,
				Value:  fieldValue.Interface().(string),
				UserID: profile.ID,
			}
			// Look for User
			var oldUserField model.UserField
			h.DB.Where("user_id = ? AND key = ?", userfield.UserID, userfield.Key).First(&oldUserField)

			// If key does not exist, create a new UserField
			if oldUserField.ID == 0 {
				if result := h.DB.Create(&userfield); result.Error != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
						fiber.StatusInternalServerError,
						"",
						"Error occurred while creating the new user field.",
					))
				}
			} else {
				// If key exists, update the existing UserField
				if result := h.DB.Model(&oldUserField).Updates(&userfield); result.Error != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
						fiber.StatusInternalServerError,
						"",
						"Error occurred while updating the user field.",
					))
				}
			}
		}

	}

	///////education array/////////////////////////////////////////////////////////////////////////////////////
	var eduArray []map[string]interface{}

	// Unmarshal the JSON array into the map
	if err := json.Unmarshal([]byte(body.Education), &eduArray); err != nil {
		fmt.Println("Error unmarshaling Education:", err)
	} else {
		fmt.Printf("%+v\n", eduArray)
	}

	return c.Status(fiber.StatusOK).JSON(&profile)

}
