package AI

import (
	"fmt"
)

const (
	NoError = iota
	GenericErrCode
	MissingEnvKeyErrCode
	InvalidModelErrCode
	FlaggedTermsErrCode
	APIErrCode

	GenericErrMsg      = "Une erreur inattendue s'est produite."
	InvalidModelErrMsg = "Modele GPT invalide."
	MissingKeyErrMsg   = "La variable d'environnement OPENAI_API_KEY ou la clef_api est requise."
	FlaggedTermsErrMsg = "Le message que vous avez tenté de soumettre contient des termes inappropriés."
	ApiErrMsg          = "Erreur lors de l'appel à l'api d'OpenAI."
)

var GenericErr = genericErr()
var InvalidModelErr = invalidModelErr()
var MissingEnvKeyErr = missingKeyErr()
var FlaggedTermsErr = flaggedTermsErr()
var ApiErr = apiErr()

type AIError struct {
	Message string
	Code    int
}

func (e AIError) Error() string {
	return fmt.Sprintf("GptProfNewton code d'erreur %d: %s", e.Code, e.Message)
}

func genericErr() AIError {
	return AIError{Message: GenericErrMsg, Code: GenericErrCode}
}

func invalidModelErr() AIError {
	return AIError{Message: InvalidModelErrMsg, Code: InvalidModelErrCode}
}

func missingKeyErr() AIError {
	return AIError{Message: MissingKeyErrMsg, Code: MissingEnvKeyErrCode}
}

func flaggedTermsErr() AIError {
	return AIError{Message: FlaggedTermsErrMsg, Code: FlaggedTermsErrCode}
}

func apiErr() AIError {
	return AIError{Message: ApiErrMsg, Code: APIErrCode}
}
