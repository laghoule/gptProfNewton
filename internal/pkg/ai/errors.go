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

	genericErrMsg       = "Une erreur inattendue s'est produite."
	invalidModelErrMsg  = "Modele GPT invalide."
	missingEnvKeyErrMsg = "La variable d'environnement OPENAI_API_KEY est requise."
	flaggedTermsErrMsg  = "Le message que vous avez tenté de soumettre contient des termes inappropriés."
	apiErrMsg           = "Erreur lors de l'appel à l'api d'OpenAI."
)

var GenericErr = genericErr()
var InvalidModelErr = invalidModelErr()
var MissingEnvKeyErr = missingEnvKeyErr()
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
	return AIError{Message: genericErrMsg, Code: GenericErrCode}
}

func invalidModelErr() AIError {
	return AIError{Message: invalidModelErrMsg, Code: InvalidModelErrCode}
}

func missingEnvKeyErr() AIError {
	return AIError{Message: missingEnvKeyErrMsg, Code: MissingEnvKeyErrCode}
}

func flaggedTermsErr() AIError {
	return AIError{Message: flaggedTermsErrMsg, Code: FlaggedTermsErrCode}
}

func apiErr() AIError {
	return AIError{Message: apiErrMsg, Code: APIErrCode}
}
