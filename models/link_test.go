package models_test

import (
	"github.com/gofrs/uuid"
	"github.com/oxks/mbfl/models"
)

// func (ms *ModelSuite) Test_Link_BeforeValidations() {
// 	link := &models.Link{
// 		UserID: uuid.NewV4(),
// 		Link:   "https://gobuffalo.io",
// 	}

// }

func (ms *ModelSuite) Test_Link_BeforeValidations() {

	usid, err := uuid.NewV4()
	link := &models.Link{
		UserID: usid,
		Link:   "http://gobuffalo.io",
	}
	ms.NoError(err)
	err = link.BeforeValidations(ms.DB)
	ms.NoError(err)
	ms.NotZero(link.Code)
	ms.Len(link.Code, 7)

}
