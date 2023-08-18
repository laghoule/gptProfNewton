package AI

import (
	"fmt"
)

const (
	NoError = iota
	MissingEnvKeyErrCode
	InvalidModelErrCode
	FlaggedTermsErrCode
	APIErrCode

	invalidModelErrMsg  = "Modele GPT invalide."
	missingEnvKeyErrMsg = "La variable d'environnement OPENAI_API_KEY est requise."
	flaggedTermsErrMsg  = "Le message que vous avez tenté de soumettre contient des termes inappropriés."
	apiErrMsg           = "Erreur lors de l'appel à l'api d'OpenAI."
)

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
