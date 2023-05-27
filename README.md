# gptProfNewton

## Description

gptProfNewton est un projet qui utilise le modèle GPT pour simuler un professeur généraliste de niveau élémentaire, le Professeur Newton. Le Professeur Newton est conçu pour utiliser un langage simple et imagé, adapté au niveau de l'élève. Il est toujours enthousiaste et démontre un grand intérêt à transmettre ses connaissances. Le programme est écrit en Go et utilise l'API OpenAI pour générer des réponses.
Caractéristiques

* Langage simple et imagé : Le Professeur Newton utilise un langage simple et imagé pour faciliter la compréhension des concepts par les élèves.

* Ton enthousiaste : Le Professeur Newton est toujours enthousiaste et démontre un grand intérêt à transmettre ses connaissances.

* Références externes : Si le Professeur Newton ne possède pas la réponse à une question, il peut diriger l'élève vers des ressources externes ou suggérer de demander de l'aide à un parent.

* Sécurité : Si le Professeur Newton juge qu'un sujet n'est pas approprié pour l'élève, il le référera à ses parents.

## Installation

### Binaire Go

Pour installer gptProfNewton en tant que binaire Go, vous pouvez télécharger la dernière version depuis la page des releases de notre dépôt GitHub. Nous fournissons des binaires pour Linux, Windows et macOS, à la fois pour les architectures amd64 et arm64.

Une fois que vous avez téléchargé le binaire correspondant à votre système, vous pouvez le rendre exécutable et le déplacer dans un répertoire de votre PATH. Par exemple, pour un binaire Linux, vous pouvez faire :

```bash
chmod +x gptProfNewton-linux-amd64
sudo mv gptProfNewton-linux-amd64 /usr/local/bin/gptProfNewton
```

### Docker

Si vous préférez utiliser Docker, vous pouvez également tirer notre image depuis le GitHub Container Registry ou Docker Hub :

```bash
docker pull ghcr.io/laghoule/gptProfNewton:latest
```

ou

```bash
docker pull laghoule/gptProfNewton:latest
```

### Utilisation

Pour utiliser gptProfNewton, vous devez d'abord définir votre clé API OpenAI comme variable d'environnement :

Pour Linux et macOS :

```bash
export OPENAI_API_KEY=your-api-key
```

Pour Windows :

```powershell
$env:OPENAI_API_KEY="your-api-key"
```

Vous pouvez ensuite exécuter le binaire construit avec la commande suivante :

```bash
./gptProfNewton
```

Ou, si vous utilisez Docker, vous pouvez passer la clé API comme variable d'environnement à Docker :

```bash
docker run -it --rm -e OPENAI_API_KEY=your-api-key ghcr.io/laghoule/gptProfNewton:latest
```

Pour quitter le programme, tapez "quit".

### Contribution

Les contributions sont les bienvenues ! Pour contribuer, veuillez forker le dépôt et créer une pull request.
Licence

gptProfNewton est sous licence GPLv3. Voir le fichier LICENSE pour plus de détails.