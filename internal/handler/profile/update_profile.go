package profile

import (
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
	Children                    string `json:"children"`
	Eligibility                 string `json:"eligibility"`
	Work                        string `json:"work"`
}

func (h *handler) UpdateProfile(c *fiber.Ctx) error {

	var profile model.User
	if c.Params("id") != "" {
		if result := h.DB.First(&profile, c.Params("id")); result.Error != nil {
			// return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
			return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
				fiber.StatusBadRequest,
				"",
				"User Not Found",
			))
		}
	} else {
		profile = c.Locals("profile").(model.User)
	}

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

	//this will update the common fields of the Profile
	if err := UpdateFields(h, c, profile, body); err != nil {
		return err
	}

	///////update education
	if err := UpdateEdu(h, c, profile, body.Education); err != nil {
		return err
	}

	////////update children
	if err := UpdateChildren(h, c, profile, body.Children); err != nil {
		return err
	}

	//////update eligibilitity
	if err := UpdateEli(h, c, profile, body.Eligibility); err != nil {
		return err
	}

	///// update Work
	if err := UpdateWork(h, c, profile, body.Work); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(&profile)

}
