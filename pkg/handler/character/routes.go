package character

import (
	httpAdapter "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/httpadapter"
	characterService "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character"
	"github.com/go-chi/chi/v5"
)

type CharacterHandler struct {
	service characterService.ICharacter
}

func New(service characterService.ICharacter) *CharacterHandler {
	return &CharacterHandler{service: service}
}

func (h *CharacterHandler) RegisterRoutes(r chi.Router) {
	h.characterRoutes(r)
}

func (h *CharacterHandler) characterRoutes(r chi.Router) {
	r.Post("/", httpAdapter.AppHandler(h.createCharacter).ServeHTTP)
	r.Get("/", httpAdapter.AppHandler(h.getAllCharacters).ServeHTTP)

	r.Route("/{characterID}", func(r chi.Router) {
		r.Get("/", httpAdapter.AppHandler(h.getCharacter).ServeHTTP)
		r.Put("/", httpAdapter.AppHandler(h.updateCharacter).ServeHTTP)
		r.Patch("/", httpAdapter.AppHandler(h.replacePortrait).ServeHTTP)
		r.Delete("/", httpAdapter.AppHandler(h.deleteCharacter).ServeHTTP)

		h.characteristicsRoutes(r)
		h.derivedStatsRoutes(r)
		h.skillsRoutes(r)

		h.healthRoutes(r)
		h.magicRoutes(r)
		h.sanityRoutes(r)
		h.luckRoutes(r)

		h.backstoriesRoutes(r)
		h.financesRoutes(r)

		h.notesRoutes(r)
	})
}

func (h *CharacterHandler) characteristicsRoutes(r chi.Router) {
	r.Route("/characteristics", func(r chi.Router) {
		r.Get("/", httpAdapter.AppHandler(h.getCharacteristics).ServeHTTP)
		r.Put("/", httpAdapter.AppHandler(h.upsertCharacteristics).ServeHTTP)
		r.Delete("/", httpAdapter.AppHandler(h.deleteCharacteristics).ServeHTTP)
	})
}

func (h *CharacterHandler) derivedStatsRoutes(r chi.Router) {
	r.Route("/derived-stats", func(r chi.Router) {
		r.Get("/", httpAdapter.AppHandler(h.getDerivedStats).ServeHTTP)
		r.Put("/", httpAdapter.AppHandler(h.upsertDerivedStats).ServeHTTP)
		r.Delete("/", httpAdapter.AppHandler(h.deleteDerivedStats).ServeHTTP)
	})
}

func (h *CharacterHandler) healthRoutes(r chi.Router) {
	r.Route("/health", func(r chi.Router) {
		r.Get("/", httpAdapter.AppHandler(h.getHealth).ServeHTTP)
		r.Put("/", httpAdapter.AppHandler(h.upsertHealth).ServeHTTP)
		r.Delete("/", httpAdapter.AppHandler(h.deleteHealth).ServeHTTP)
	})
}

func (h *CharacterHandler) magicRoutes(r chi.Router) {
	r.Route("/magic", func(r chi.Router) {
		r.Get("/", httpAdapter.AppHandler(h.getMagic).ServeHTTP)
		r.Put("/", httpAdapter.AppHandler(h.upsertMagic).ServeHTTP)
		r.Delete("/", httpAdapter.AppHandler(h.deleteMagic).ServeHTTP)
	})
}

func (h *CharacterHandler) sanityRoutes(r chi.Router) {
	r.Route("/sanity", func(r chi.Router) {
		r.Get("/", httpAdapter.AppHandler(h.getSanity).ServeHTTP)
		r.Put("/", httpAdapter.AppHandler(h.upsertSanity).ServeHTTP)
		r.Delete("/", httpAdapter.AppHandler(h.deleteSanity).ServeHTTP)
	})
}

func (h *CharacterHandler) luckRoutes(r chi.Router) {
	r.Route("/luck", func(r chi.Router) {
		r.Get("/", httpAdapter.AppHandler(h.getLuck).ServeHTTP)
		r.Put("/", httpAdapter.AppHandler(h.upsertLuck).ServeHTTP)
		r.Delete("/", httpAdapter.AppHandler(h.deleteLuck).ServeHTTP)
	})
}

func (h *CharacterHandler) backstoriesRoutes(r chi.Router) {
	r.Route("/backstory", func(r chi.Router) {
		r.Get("/", httpAdapter.AppHandler(h.getBackstory).ServeHTTP)
		r.Put("/", httpAdapter.AppHandler(h.upsertBackstory).ServeHTTP)
		r.Delete("/", httpAdapter.AppHandler(h.deleteBackstory).ServeHTTP)

		r.Route("/items", func(r chi.Router) {
			r.Get("/", httpAdapter.AppHandler(h.getBackstoryItems).ServeHTTP)
			r.Post("/", httpAdapter.AppHandler(h.createBackstoryItem).ServeHTTP)

			r.Route("/{itemID}", func(r chi.Router) {
				r.Get("/", httpAdapter.AppHandler(h.getBackstoryItem).ServeHTTP)
				r.Put("/", httpAdapter.AppHandler(h.updateBackstoryItem).ServeHTTP)
				r.Delete("/", httpAdapter.AppHandler(h.deleteBackstoryItem).ServeHTTP)
			})
		})
	})
}

func (h *CharacterHandler) financesRoutes(r chi.Router) {
	r.Route("/finances", func(r chi.Router) {
		r.Get("/", httpAdapter.AppHandler(h.getFinances).ServeHTTP)
		r.Put("/", httpAdapter.AppHandler(h.upsertFinances).ServeHTTP)
		r.Delete("/", httpAdapter.AppHandler(h.deleteFinances).ServeHTTP)
	})
}

func (h *CharacterHandler) skillsRoutes(r chi.Router) {
	r.Route("/skills", func(r chi.Router) {
		r.Get("/", httpAdapter.AppHandler(h.getSkills).ServeHTTP)
		r.Post("/", httpAdapter.AppHandler(h.createSkill).ServeHTTP)

		r.Route("/{skillID}", func(r chi.Router) {
			r.Get("/", httpAdapter.AppHandler(h.getSkill).ServeHTTP)
			r.Put("/", httpAdapter.AppHandler(h.updateSkill).ServeHTTP)
			r.Delete("/", httpAdapter.AppHandler(h.deleteSkill).ServeHTTP)
		})
	})
}

func (h *CharacterHandler) notesRoutes(r chi.Router) {
	r.Route("/notes", func(r chi.Router) {
		r.Get("/", httpAdapter.AppHandler(h.getNotes).ServeHTTP)
		r.Post("/", httpAdapter.AppHandler(h.createNote).ServeHTTP)

		r.Route("/{noteID}", func(r chi.Router) {
			r.Get("/", httpAdapter.AppHandler(h.getNote).ServeHTTP)
			r.Put("/", httpAdapter.AppHandler(h.updateNote).ServeHTTP)
			r.Delete("/", httpAdapter.AppHandler(h.deleteNote).ServeHTTP)
		})
	})
}
