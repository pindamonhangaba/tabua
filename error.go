package tabua

import (
	pq "github.com/lib/pq"
)

// ColumnValidationError custom error for column validation failure
type ColumnValidationError struct {
	ColumnName   string
	ErrorMessage string
}

func (cve ColumnValidationError) Error() string {
	return cve.ColumnName + ": " + cve.ErrorMessage
}

// QueryGenerationError custom error for use when erring in query generation
type QueryGenerationError struct {
	Message string
}

func (cve QueryGenerationError) Error() string {
	return "Error generation query: " + cve.Message
}

func handleError(err error) error {
	if driverErr, ok := err.(*pq.Error); ok { // Now the error number is accessible directly
		if driverErr.Code == "1045" {
			// Handle the permission-denied error
		}
	} /*else if driverErr, ok := err.(*sql.ErrNoRows); ok { // Now the error number is accessible directly
		// not found
	}*/
	return err
}
