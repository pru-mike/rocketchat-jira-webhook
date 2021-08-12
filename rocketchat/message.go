package rocketchat

import (
	"github.com/pru-mike/rocketchat-jira-webhook/logger"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

type translate struct {
	tag string
	key string
	msg catalog.Message
}

var translates = []translate{
	{
		"en", "Found %d issue",
		plural.Selectf(1, "%d",
			"=1", "Found 1 task.",
			plural.Other, "Found %[1]d tasks.",
		),
	},
	{
		"ru", "Found %d issue",
		plural.Selectf(1, "%d",
			plural.One, "Найдена %[1]d задача.",
			plural.Few, "Найдено %[1]d задачи.",
			plural.Other, "Найдено %[1]d задач.",
		),
	},
	{
		"fr", "Found %d issue",
		plural.Selectf(1, "%d",
			"=1", "1 tâche trouvée.",
			plural.Other, "Trouvé %[1]d tâches.",
		),
	},
	{
		"de", "Found %d issue",
		plural.Selectf(1, "%d",
			"=1", "1 Aufgabe gefunden.",
			plural.Other, "%[1]d Aufgaben gefunden.",
		),
	},

	{
		"en", "Found %d doc",
		plural.Selectf(1, "%d",
			"=1", "Found 1 document.",
			plural.Other, "Found %[1]d documents.",
		),
	},
	{
		"ru", "Found %d doc",
		plural.Selectf(1, "%d",
			plural.One, "Найден %[1]d документ.",
			plural.Few, "Найдено %[1]d документа.",
			plural.Other, "Найдено %[1]d документов.",
		),
	},
	{
		"fr", "Found %d doc",
		plural.Selectf(1, "%d",
			"=1", "1 document trouvée.",
			plural.Other, "Trouvé %[1]d documents.",
		),
	},
	{
		"de", "Found %d pages",
		plural.Selectf(1, "%d",
			"=1", "1 dokument gefunden.",
			plural.Few, "%[1]d dokumente gefunden.",
			plural.Other, "%[1]d dokument gefunden.",
		),
	},
}

func init() {
	for _, e := range translates {
		err := message.Set(language.MustParse(e.tag), e.key, e.msg)
		if err != nil {
			logger.Error(err)
		}
	}
}
