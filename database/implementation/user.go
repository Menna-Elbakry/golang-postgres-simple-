package database

import (
	"log"
	table "task/database/table"
	"task/model"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-pg/pg/orm"
	uuid "github.com/satori/go.uuid"
)

//GetAllUsers to select all users exist
func GetAllUsers() ([]model.User, error) {
	var userEntity []table.User
	selErr := db.Model(&userEntity).Select()
	if selErr != nil {
		log.Printf("Error While Getting Users Reason:  %v\n", selErr)
		return nil, selErr
	}
	var users []model.User
	for _, usr := range userEntity {
		users = append(users, usr.MapToModule())
		log.Printf("Get Users Successful . \n Name: %v ,Email:%v", usr.Name, usr.Email)
	}
	return users, nil
}

//GetUserByID to search for user
func GetUserById(user *model.User) error {
	var userEntity table.User
	getErr := db.Model(&userEntity).
		Where("id=?", user.ID).
		Returning("name,email").
		Select()
	if getErr != nil {
		log.Printf("Error While Getting User By Id Reason: \n %v", getErr)
		return getErr
	}
	log.Printf("Get User By Id Successful For \n Name:%v , Email: %v\n ", userEntity.Name, userEntity.Email)
	return nil
}

func GetUserRoleById(user *model.User) (model.RoleEnum, error) {
	var userEntity table.User
	getErr := db.Model(&userEntity).
		Where("id=?", user.ID).
		Returning("Role").
		Select()
	if getErr != nil {
		log.Printf("Error While Getting User By Id Reason: \n %v", getErr)
		return " ", getErr
	}
	log.Printf("Get User By Id Successful For \n Role:%v \n ", userEntity.Role)
	return userEntity.Role, nil
}

//GetUserByRole to search for user
func GetUsersByRole(user *model.User) ([]model.User, error) {
	var userEntity []table.User
	getErr := db.Model(&userEntity).
		Where("role=?", user.Role).
		Returning("name,id").
		Select()
	if getErr != nil {
		log.Printf("Error While Getting Users Reason: \n %v", getErr)
		return nil, getErr
	}
	var users []model.User
	for _, usr := range userEntity {
		users = append(users, usr.MapToModule())
		log.Printf("Get Users Successful .\n ID: %v, Name: %v", usr.ID, usr.Name)
	}
	return users, nil
}

//AddNewUser to database
func AddNewUser(user *model.User) (uuid.UUID, error) {
	_, insertErr := db.Model(table.User{}.Fill(user)).Insert()
	if insertErr != nil {
		log.Printf("Error While Adding New User Reason:  %v\n", insertErr)
		return uuid.Nil, insertErr
	}
	log.Printf("New User Added Successfully .\n User Id: %s", user.ID)
	return user.ID, nil
}

//DeleteUser to delete user
func DeleteUser(user *model.User) error {
	var userEntity table.User
	_, deleteErr := db.Model(&userEntity).Where("id=?", user.ID).
		Returning("id").
		Delete()
	if deleteErr != nil {
		log.Printf("Error While Deleting User. Reason: %v \n", deleteErr)
		return deleteErr
	}
	log.Printf("Successfully Deleted User With ID: %v \n", userEntity.ID)
	return nil
}

//UpdateUserPassword to update the password
func UpdateUserPassword(user *model.User) error {
	var userEntity table.User
	_, updateErr := db.Model(&userEntity).
		Set("password=?,updated_at=?", user.Password, time.Now()).
		Where("id=?", user.ID).
		Update()
	if updateErr != nil {
		log.Printf("Error While Updating Password. Reason:  %v\n", updateErr)
		return updateErr
	}
	log.Printf("Password Updated Successfully For User: %v, Email: %v\n ", userEntity.Name, userEntity.Email)
	return nil
}

//UpdateUsername to modify the user name
func UpdateUsername(user *model.User) error {
	var userEntity table.User
	_, updateErr := db.Model(&userEntity).
		Set("name=?", user.Name).
		Where("id=?,updated_at=?", user.ID, time.Now()).
		Returning("name,email").
		Update()
	if updateErr != nil {
		log.Printf("Error While Updating Name. Reason:  %v\n", updateErr)
		return updateErr
	}
	log.Printf("Name Updated Successfully For Email %v, Username : %v\n ", userEntity.Email, userEntity.Name)
	return nil
}

//UpdateUserMail to modify user email
func UpdateUserMail(user *model.User) error {
	var userEntity table.User
	_, updateErr := db.Model(&userEntity).
		Set("email=?,,updated_at=?", user.Email, time.Now()).
		Where("id=?", user.ID).
		Returning("name,email").
		Update()
	if updateErr != nil {
		log.Printf("Error While Updating Reason:  %v\n", updateErr)
		return updateErr
	}
	log.Printf("Email Updated Successfully For User: %v, Email: %v\n ", userEntity.Name, userEntity.Email)
	return nil
}

//Login func
func Login(user *model.User) (uuid.UUID, error) {
	var userEntity table.User
	getErr := db.Model(&userEntity).Where("Email=?", user.Email).
		Returning("id,password").Select()
	if getErr != nil {
		log.Printf("Error While Login :  %v\n", getErr)
		return uuid.Nil, getErr
	}
	err := bcrypt.CompareHashAndPassword([]byte(userEntity.Password), []byte(user.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		log.Println(err)
		return uuid.Nil, nil
	}
	log.Printf("Successfully LogedIn ID .%v\n ", userEntity.ID)

	return userEntity.ID, nil
}

//CreateUserTable to add users table to database
func CreateUserTable() error {
	opts := &orm.CreateTableOptions{
		IfNotExists: true,
	}
	createErr := db.CreateTable(&table.User{}, opts)
	if createErr != nil {
		log.Printf("Error While Creating Table Users Reason: \n %v", createErr)
		return createErr
	}
	log.Printf("Table Users Created Successfully.\n")
	return nil
}
