package apierr

// ===== INTERNAL SERVER ERROR =====

const InternalServerErrorCode = "INTERNAL_SERVER"
const InternalServerErrorMessage = "Sorry, something went wrong on our end. Please try again later."

// ===== DATA NOT FOUND ERROR =====

const DataNotFoundErrorCode = "DATA_NOT_FOUND"
const DataNotFoundErrorMessage = "{entity} not found."

// ===== ID NOT FOUND ERROR =====

const IdNotFoundErrorCode = "ID_NOT_FOUND"
const IdNotFoundErrorMessage = "{field} '{id}' not found."

// ===== EMPTY FIELD ERROR =====

const EmptyFieldErrorCode = "EMPTY_FIELD"
const EmptyFieldErrorMessage = "{field} is required."

// ===== INVALID REQUEST ERROR =====

const InvalidRequestErrorCode = "INVALID_REQUEST"
const InvalidRequestErrorMessage = "Request body format is invalid. Please check the data structure and try again."

// ===== VALIDATION ERROR =====

const ValidationErrorCode = "VALIDATION_ERROR"
const ValidationErrorMessage = "Validation error"
