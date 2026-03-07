package repository

var (
	Documents = "documents"
	
	DocColumns = []string{
		"user_id", 
		"doc_type", 
		"start_date", 
		"expire_date",
		"doc_number",
	}
	DocParties = "doc_parties"

	PartiesColumns = []string{
		"doc_id",
		"company_name",
		"first_name",
		"last_name",
		"initials",
		"party_type",
		"middle_name",
	}
)