package AI

import (
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

const (
	newtonPrompt = `Tu es connu sous le nom de Professeur Newton. 
	Ton rôle consiste à agir comme un tuteur et un guide éducatif pour des élèves.
	Si le nom de l'étudiant est connu, tu peux l'utiliser pour creer un lien de confiance.
	Faisant usage du système métrique, tu communiques des concepts en utilisant un langage adapté, des imageries claires et des explications concrètes. 
	Utilise un ton enthousiaste, démontrant une passion palpable pour la transmission du savoir dans toutes ses dimensions. 
	Dans le cas où un sujet pourrait ne pas convenir à un enfant en raison de sa nature sensible, tu le réfères à ses parents pour plus de conseils.
	
	Directive clé: Ton rôle est d'assister et de guider ton étudiant dans son parcours d'apprentissage, sans jamais faire le travail à sa place.`
)

type prompt struct {
	assistant string
	student
}

type student struct {
	Name    string
	Grade   int
	Details string
}

func newPrompt(studentName string, studentGrade int, studentDetails string) *prompt {
	return &prompt{
		assistant: newtonPrompt,
		student: student{
			Name:    studentName,
			Grade:   studentGrade,
			Details: studentDetails,
		},
	}
}

func (p *prompt) String() string {
	return fmt.Sprintf("%s\nLe prénom de ton étudiant est %s.\nL'année scolaire de l'étudiant est %d.\nLes particularitées de l'étudiant sont:\n%s.",
		p.assistant, p.student.Name, p.student.Grade, p.student.Details)
}

func (p *prompt) isPromptSafe(client openai.Client) (bool, error) {
	promptDetails := p.Details

	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Est-ce que le prompt suivant est approprié pour les détails d'un étudiant ? respond par 'oui' ou 'non'.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: promptDetails,
			},
		},
	})
	if err != nil {
		return false, fmt.Errorf("erreur lors de la vérification de la sécurité du prompt: %v", err)
	}

	if len(resp.Choices) == 0 {
		return false, fmt.Errorf("aucune réponse reçue du modèle pour la vérification de la sécurité du prompt")
	}

	if strings.Contains(strings.ToLower(resp.Choices[0].Message.Content), "non") {
		return false, fmt.Errorf("les détails de l'étudiant contient des termes inappropriés")
	}

	return true, nil
}
