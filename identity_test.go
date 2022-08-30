package kratox_test

import (
	"encoding/json"

	"github.com/jaswdr/faker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	kratox "github.com/w6d-io/kratox"
)

var _ = Describe("Identity", func() {
	Context("CRUD Identity", func() {
		BeforeEach(func() {
		})
		AfterEach(func() {
			kratox.Kratox = nil
		})
		It("succeeds to Create, Get, Update And Delete Identity", func() {

			// RAMDOM FAKER
			faker := faker.New()
			name := faker.Person().FirstName()

			// INPUT SCHEMA TYPE
			schemaId := "default"

			// INPUT CREATION
			jsonStrCreate := `{
    "email": "` + name + `@wildcard.io",
    "name": {
      "first": "` + name + `",
      "last": "` + name + `"
    },
    "roles": {
      "organizations": [],
      "private_projects": [
        {
          "key": 666,
          "value": "admin"
        }
      ],
      "scopes": [],
      "affiliate_projects": []
    },
    "projects":[]
  }`
			// INPUT UPDATED
			jsonStrUpdate := `{
    "email": "` + name + `@wildcard.io",
    "name": {
      "first": "` + name + `",
      "last": "` + name + `"
    },
    "roles": {
      "organizations": [],
      "private_projects": [
        {
          "key": 666,
          "value": "admin"
        },
        {
          "key": 777,
          "value": "user"
        },
        {
          "key": 888,
          "value": "222"
        }
      ],
      "scopes": [],
      "affiliate_projects": []
    },
    "projects":[]
  }`

			// CONNECT TO SERVER
			kratox.SetAddressDetails("127.0.0.1", Verbose, 4434)

			// TRAIT
			var traitCreation map[string]interface{}
			var traitUpdate map[string]interface{}

			// CREATE
			json.Unmarshal([]byte(jsonStrCreate), &traitCreation)
			CreateIdentity, _ := kratox.Kratox.CreateIdentity(ctx, schemaId, traitCreation)
			idCreated := CreateIdentity.Id
			Expect(CreateIdentity.State.IsValid()).To(Equal(true))

			// GET
			gettedIdentity, _ := kratox.Kratox.GetIdentity(ctx, idCreated)
			Expect(gettedIdentity.GetId()).To(Equal(CreateIdentity.GetId()))

			// UPDATE
			json.Unmarshal([]byte(jsonStrUpdate), &traitUpdate)
			newUpdatedIdentity, _ := kratox.Kratox.UpdateIdentity(ctx, idCreated, schemaId, traitUpdate)
			lastUpdateIdentity := gettedIdentity.Traits
			Expect(newUpdatedIdentity.Traits).To(Not(Equal(lastUpdateIdentity)))

			// DELETE
			e := kratox.Kratox.DeleteIdentity(ctx, idCreated)
			Expect(e).To(BeNil())
		})
	})
})
