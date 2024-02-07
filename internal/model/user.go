package model

import (
	"strings"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username      string         `gorm:"unique"`
	Email         string         `gorm:"unique"`
	Password      string         `json:"-"`
	Active        string         `gorm:"column:active" json:"active"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Fields        []UserField    `gorm:"foreignKey:user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"fields"`
	Educations    []Edu          `gorm:"foreignKey:user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"educations"`
	Children      []Children     `gorm:"foreignKey:user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"children"`
	Eligibilities []Eli          `gorm:"foreignKey:user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"eligibilities"`
	Works         []Work         `gorm:"foreignKey:user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"works"`
	Voluntaries   []Vol          `gorm:"foreignKey:user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"voluntaries"`
	Learns        []Learn        `gorm:"foreignKey:user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"learns"`
	Skills        []Skill        `gorm:"foreignKey:user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"skills"`
	Associations  []Assoc        `gorm:"foreignKey:user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"associations"`
	References    []Reference    `gorm:"foreignKey:user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"references"`
	Recognitions  []Recognition  `gorm:"foreignKey:user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"recognitions"`
}
type UserResponse struct {
	Id            int                   `json:"id"`
	Username      string                `json:"username"`
	Email         string                `json:"email"`
	Active        string                `json:"active"`
	Fields        map[string]string     `json:"fields"`
	Educations    []interface{}         `json:"educations"`
	Children      []ChildrenResponse    `json:"children"`
	Eligibilities []interface{}         `json:"eligibilities"`
	Works         []interface{}         `json:"works"`
	Voluntaries   []interface{}         `json:"voluntaries"`
	Learns        []interface{}         `json:"learns"`
	Skills        []SkillResponse       `json:"skills"`
	Associations  []AssocResponse       `json:"associations"`
	References    []ReferenceResponse   `json:"references"`
	Recognitions  []RecognitionResponse `json:"recognitions"`
}

func ConvertToUserResponse(db *gorm.DB, user User) UserResponse {
	// db.Preload("Fields").Preload("Educations").Preload("Children").Preload("Eligibilities").Preload("Works").Preload("Voluntaries").Preload("Learns").Preload("Skills").Preload("Associations").Preload("References").Find(&user)

	preloads := []string{"Fields", "Educations", "Children", "Eligibilities", "Works", "Voluntaries", "Learns", "Skills", "Associations", "References", "Recognitions"}

	for _, preload := range preloads {
		db.Preload(preload).Find(&user)
	}

	shortenedFields := make([]UserShortenField, len(user.Fields))

	//user profile fields
	for i, field := range user.Fields {
		shortenedFields[i] = ConvertToUserShortenField(field)
	}

	profileFieldResult := make(map[string]string)

	for _, field := range shortenedFields {
		profileFieldResult[field.Key] = field.Value
	}

	//user education fields
	var eduResponse []interface{}

	for _, eduData := range user.Educations {
		convertedData := ConvertToEduResponse(db, eduData).Fields
		convertedData["id"] = eduData.ID
		eduResponse = append(eduResponse, convertedData)
	}

	//user children fields
	childResponse := ConvertToChildrenResponse(db, user.Children)

	//eligibility fields
	var eliResponse []interface{}

	for _, eliData := range user.Eligibilities {
		eliconvertedData := ConvertToEliResponse(db, eliData).Fields
		eliconvertedData["id"] = eliData.ID
		eliResponse = append(eliResponse, eliconvertedData)
	}

	//Works Fields
	var workResponse []interface{}

	for _, workData := range user.Works {
		workconvertedData := ConvertToWorkResponse(db, workData).Fields
		workconvertedData["id"] = workData.ID
		workResponse = append(workResponse, workconvertedData)
	}

	//Voluntary Works
	var volResponse []interface{}

	for _, volData := range user.Voluntaries {
		volconvertedData := ConvertToVolResponse(db, volData).Fields
		volconvertedData["id"] = volData.ID
		volResponse = append(volResponse, volconvertedData)
	}

	//Learning Works
	var learnResponse []interface{}

	for _, learnData := range user.Learns {
		learnconvertedData := ConvertToLearnResponse(db, learnData).Fields
		learnconvertedData["id"] = learnData.ID
		learnResponse = append(learnResponse, learnconvertedData)
	}

	//user skills fields
	skillResponse := ConvertToSkillResponse(db, user.Skills)

	//user assocs fields
	assocResponse := ConvertToAssocResponse(db, user.Associations)

	//user references fields
	referenceResponse := ConvertToReferenceResponse(db, user.References)

	//user recognitions fields
	recognitionResponse := ConvertToRecognitionResponse(db, user.Recognitions)

	return UserResponse{
		Id:            int(user.ID),
		Username:      user.Username,
		Email:         user.Email,
		Active:        user.Active,
		Fields:        profileFieldResult,
		Educations:    eduResponse,
		Children:      childResponse,
		Eligibilities: eliResponse,
		Works:         workResponse,
		Voluntaries:   volResponse,
		Learns:        learnResponse,
		Skills:        skillResponse,
		Associations:  assocResponse,
		References:    referenceResponse,
		Recognitions:  recognitionResponse,
	}
}

type UserField struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	Key    string `gorm:"column:key" json:"key"`
	Value  string `gorm:"column:value" json:"value"`
}

type UserShortenField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func ConvertToUserShortenField(userfield UserField) UserShortenField {
	return UserShortenField{
		Key:   strings.ToLower(userfield.Key),
		Value: userfield.Value,
	}
}
