package main

import (
	"context"
	"github.com/gocarina/gocsv"
	"github.com/kennykarnama/school-adminstration-api/config"
	user3 "github.com/kennykarnama/school-adminstration-api/domain/entity/user"
	"github.com/kennykarnama/school-adminstration-api/domain/repository/user"
	user2 "github.com/kennykarnama/school-adminstration-api/domain/service/user"
	"github.com/kennykarnama/school-adminstration-api/util/dbconn"
	"log"
	"os"
)

type Credential struct {
	AlternativeId string `csv:"alternative_id"`
	Name          string `csv:"name"`
	Password      string `csv:"password"`
	SchoolID      string `csv:"school_id"`
}

type Credentials []*Credential

func (credentials Credentials) ToTeachers() []*user3.Teacher {
	var teachers []*user3.Teacher
	for _, c := range credentials {
		schoolID := c.SchoolID
		if schoolID == "" {
			schoolID = "00000000-0000-4000-8000-000000000001"
		}
		teachers = append(teachers, &user3.Teacher{
			AlternativeId:      c.AlternativeId,
			Name:               c.Name,
			Password:           c.Password,
			SchoolID:           &schoolID,
			Role:               user3.RoleTeacher,
			Active:             true,
			MustChangePassword: true,
		})
	}
	return teachers
}

func main() {
	clientsFile, err := os.OpenFile("input.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer clientsFile.Close()

	var clients Credentials

	if err := gocsv.UnmarshalFile(clientsFile, &clients); err != nil { // Load clients from file
		panic(err)
	}

	cfg := config.Get()

	dbconn := dbconn.InitGorm()

	userRepo := user.NewMySqlRepository(dbconn)
	userSvc := user2.NewService(userRepo, cfg)

	err = userSvc.RegisterTeachers(context.Background(), clients.ToTeachers())
	if err != nil {
		log.Fatalf("%v", err)
	}
}
