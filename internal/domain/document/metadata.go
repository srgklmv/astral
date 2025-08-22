package document

import (
	"github.com/srgklmv/astral/internal/models/apperrors"
	"github.com/srgklmv/astral/internal/models/dto"
)

func ValidateDocumentMetadata(metadata dto.UploadDocumentRequestMetadata) (bool, apperrors.ErrorText) {
	switch {
	case metadata.Name == "":
		return false, apperrors.InvalidFileNameErrorText
	case metadata.IsFile && metadata.Mimetype == "":
		return false, apperrors.InvalidMimeTypeErrorText
	default:
		return true, ""
	}
}
