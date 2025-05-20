# gptProfNewton

[![Go Report Card](https://goreportcard.com/badge/github.com/laghoule/gptProfNewton)](https://goreportcard.com/report/github.com/laghoule/gptProfNewton)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=laghoule_gptProfNewton&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=laghoule_gptProfNewton)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=laghoule_gptProfNewton&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=laghoule_gptProfNewton)

## Description

gptProfNewton est un projet qui utilise le modèle GPT d'OpenAI pour simuler un tuteur scolaire, le Professeur Newton. Le tuteur virtuel est conçu pour utiliser un langage simple et imagé, adapté au niveau de l'élève. Il est toujours enthousiaste et démontre un grand intérêt à transmettre ses connaissances. Le programme est écrit en Go et utilise l'API OpenAI pour générer des réponses.

### Caractéristiques

* Mode texte seulement, permet d'éliminer les distractions.

* Langage simple et imagé : Le Professeur Newton utilise un langage simple et imagé pour faciliter la compréhension des concepts par les élèves.

* Ton enthousiaste : Le Professeur Newton est toujours enthousiaste et démontre un grand intérêt à transmettre ses connaissances.

* Références externes : Si le Professeur Newton ne possède pas la réponse à une question, il peut diriger l'élève vers des ressources externes ou suggérer de demander de l'aide à un parent.

* Sécurité : Si le Professeur Newton juge qu'un sujet n'est pas approprié pour l'élève, il le référera à ses parents.

* Modération : Utilise l'api de modération d'OpenAI, afin d'assurer une conversation appropriée.

## Installation

### Binaire Go

Pour installer gptProfNewton, vous pouvez télécharger la dernière version depuis la page des releases de notre dépôt GitHub. Nous fournissons des binaires pour Linux, Windows et macOS, à la fois pour les architectures amd64 et arm64.

Une fois que vous avez téléchargé le binaire correspondant à votre système, vous pouvez le rendre exécutable et le déplacer dans un répertoire de votre PATH. Par exemple, pour un binaire Linux, vous pouvez faire :

```bash
chmod +x gptProfNewton-linux-amd64
sudo mv gptProfNewton-linux-amd64 /usr/local/bin/gptProfNewton
```

### Utilisation

Pour utiliser gptProfNewton, vous devez d'abord posseder une clef d'utilisation d'OpenAI. Vous devez par la suite creer un fichier de configuration:

config.yaml

```yaml
eleve:
  nom: Bob
  niveau: 1
  details: "Bob est un eleve qui a des difficultes en francais, porte un attention particuliere a l'orthographe et a la grammaire, il est tres curieux et adore les maths."

openai:
  creatif: true
  modele: gpt-4o
  clef_api: "sk-..."
```

Vous pouvez ensuite exécuter le binaire avec la commande suivante :

Linux & MacOS:

```bash
./gptProfNewton -config config.yaml
```

Pour Windows :

```powershell
gptProfNewton -config config.yaml
```

Ou, si vous preferer utilisez Docker:

```bash
docker run -it --rm -v $(pwd)/config.yaml:/config.yaml ghcr.io/laghoule/gptprofnewton -config /config.yaml
```

Pour quitter le programme, tapez `/quit`, pour reinitialiser un conversation, tapez `/reset`.

### Contribution

Les contributions sont les bienvenues ! Pour contribuer, veuillez forker le dépôt et créer une pull request.

### Licence

gptProfNewton est sous licence GPLv3. Voir le fichier LICENSE pour plus de détails.
