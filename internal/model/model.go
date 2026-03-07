package model

import "time"

type Document struct {
	UserID         uint64
	DocumentType   string
	DocumentNumber string
	DateCreated    time.Time
	DateEnd        *time.Time
	Persons        []Party
}

type Party struct {
	PartyType   string 
	CompanyName *string
	FirstName   *string
	MiddleName  *string
	LastName    *string
	Initials    *string
}

type ListDocuments struct {
	UserID uint64
	Sort string
	Limit *int
}

type DocumentResponse struct {
	Total int64
	Documents []Document
}

var SortType = []string {
	"company_name",
	"first_name",
	"middle_name",
	"last_name",
	"initials",
	"start_date",
	"expire_date",
}