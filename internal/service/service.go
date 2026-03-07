package service

import (
	"ptmk-service/internal/repository"
	"time"
)

type Service struct {
	userRepository *repository.UserRepository
}

type CreateDocument struct {
	UserID         int64      
	DocumentType   string     
	DocumentNumber string     
	DateCreated    time.Time  
	DateEnd        *time.Time 
	Persons        []Person   
}

type Person struct {
	Type          string      
	LegalEntity   *LegalEntity   
	NaturalPerson *NaturalPerson 
}

type LegalEntity struct {
	Name string
}

type NaturalPerson struct {
	FirstName 	*string
	MiddleName 	*string
	LastName 	*string
	Initials 	*string
}


func NewService(userRepository *repository.UserRepository) *Service {
	return &Service{
		userRepository: userRepository,
	}
}
