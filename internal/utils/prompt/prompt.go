package prompt

import (
	"fmt"
	"strings"

	"test-ragger/internal/models"
)

func Build(userQ string, hits []models.Hit) string {
	var ctxParts []string
	for i, h := range hits {
		const maxFrag = 800
		txt := h.Text
		if len(txt) > maxFrag {
			txt = txt[:maxFrag] + "…"
		}
		ctxParts = append(ctxParts, fmt.Sprintf("[%d] %s (%s/%s)\n%s", i+1, h.Title, h.DocID, h.ChunkID, txt))
	}
	ctx := strings.Join(ctxParts, "\n\n---\n\n")
	return fmt.Sprintf(`Ты — технический ассистент. Отвечай только по контексту ниже, ссылайся на [номера].
Если ответа в контексте нет — честно скажи об этом.

Вопрос: %s

Контекст:
%s
`, userQ, ctx)
}
