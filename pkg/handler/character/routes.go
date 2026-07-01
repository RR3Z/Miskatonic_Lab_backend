package character

import (
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/httpadapter"
	characterService "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	characters characterService.ICharacter
}

func New(characters characterService.ICharacter) *Handler {
	return &Handler{characters: characters}
}

func (h *Handler) RegisterCharacterRoutes(r chi.Router) {
	h.characterRoutes(r)
}

func (h *Handler) characterRoutes(r chi.Router) {
	r.Post("/character", httpadapter.AppHandler(h.createCharacter).ServeHTTP)
	r.Get("/characters", httpadapter.AppHandler(h.getAllCharacters).ServeHTTP)

	r.Route("/character/{characterID}", func(r chi.Router) {
		r.Get("/", httpadapter.AppHandler(h.getCharacter).ServeHTTP)
		r.Put("/", httpadapter.AppHandler(h.updateCharacter).ServeHTTP)
		r.Delete("/", httpadapter.AppHandler(h.deleteCharacter).ServeHTTP)

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

func (h *Handler) characteristicsRoutes(r chi.Router) {
	r.Route("/characteristics", func(r chi.Router) {
		r.Get("/", httpadapter.AppHandler(h.getCharacteristics).ServeHTTP)
		r.Put("/", httpadapter.AppHandler(h.upsertCharacteristics).ServeHTTP)
		r.Delete("/", httpadapter.AppHandler(h.deleteCharacteristics).ServeHTTP)
	})
}

func (h *Handler) derivedStatsRoutes(r chi.Router) {
	r.Route("/derived-stats", func(r chi.Router) {
		r.Get("/", httpadapter.AppHandler(h.getDerivedStats).ServeHTTP)
		r.Put("/", httpadapter.AppHandler(h.upsertDerivedStats).ServeHTTP)
		r.Delete("/", httpadapter.AppHandler(h.deleteDerivedStats).ServeHTTP)
	})
}

func (h *Handler) healthRoutes(r chi.Router) {
	r.Route("/health", func(r chi.Router) {
		r.Get("/", httpadapter.AppHandler(h.getHealth).ServeHTTP)
		r.Put("/", httpadapter.AppHandler(h.upsertHealth).ServeHTTP)
		r.Delete("/", httpadapter.AppHandler(h.deleteHealth).ServeHTTP)
	})
}

func (h *Handler) magicRoutes(r chi.Router) {
	r.Route("/magic", func(r chi.Router) {
		r.Get("/", httpadapter.AppHandler(h.getMagic).ServeHTTP)
		r.Put("/", httpadapter.AppHandler(h.upsertMagic).ServeHTTP)
		r.Delete("/", httpadapter.AppHandler(h.deleteMagic).ServeHTTP)
	})
}

func (h *Handler) sanityRoutes(r chi.Router) {
	r.Route("/sanity", func(r chi.Router) {
		r.Get("/", httpadapter.AppHandler(h.getSanity).ServeHTTP)
		r.Put("/", httpadapter.AppHandler(h.upsertSanity).ServeHTTP)
		r.Delete("/", httpadapter.AppHandler(h.deleteSanity).ServeHTTP)
	})
}

func (h *Handler) luckRoutes(r chi.Router) {
	r.Route("/luck", func(r chi.Router) {
		r.Get("/", httpadapter.AppHandler(h.getLuck).ServeHTTP)
		r.Put("/", httpadapter.AppHandler(h.upsertLuck).ServeHTTP)
		r.Delete("/", httpadapter.AppHandler(h.deleteLuck).ServeHTTP)
	})
}

func (h *Handler) backstoriesRoutes(r chi.Router) {
	r.Route("/backstory", func(r chi.Router) {
		r.Get("/", httpadapter.AppHandler(h.getBackstory).ServeHTTP)
		r.Put("/", httpadapter.AppHandler(h.upsertBackstory).ServeHTTP)
		r.Delete("/", httpadapter.AppHandler(h.deleteBackstory).ServeHTTP)

		r.Route("/items", func(r chi.Router) {
			r.Get("/", httpadapter.AppHandler(h.getBackstoryItems).ServeHTTP)
			r.Post("/", httpadapter.AppHandler(h.createBackstoryItem).ServeHTTP)

			r.Route("/{itemID}", func(r chi.Router) {
				r.Get("/", httpadapter.AppHandler(h.getBackstoryItem).ServeHTTP)
				r.Put("/", httpadapter.AppHandler(h.updateBackstoryItem).ServeHTTP)
				r.Delete("/", httpadapter.AppHandler(h.deleteBackstoryItem).ServeHTTP)
			})
		})
	})
}

func (h *Handler) financesRoutes(r chi.Router) {
	r.Route("/finances", func(r chi.Router) {
		r.Get("/", httpadapter.AppHandler(h.getFinances).ServeHTTP)
		r.Put("/", httpadapter.AppHandler(h.upsertFinances).ServeHTTP)
		r.Delete("/", httpadapter.AppHandler(h.deleteFinances).ServeHTTP)
	})
}

func (h *Handler) skillsRoutes(r chi.Router) {
	r.Route("/skills", func(r chi.Router) {
		r.Get("/", httpadapter.AppHandler(h.getSkills).ServeHTTP)
		r.Post("/", httpadapter.AppHandler(h.createSkill).ServeHTTP)

		r.Route("/{skillID}", func(r chi.Router) {
			r.Get("/", httpadapter.AppHandler(h.getSkill).ServeHTTP)
			r.Put("/", httpadapter.AppHandler(h.updateSkill).ServeHTTP)
			r.Delete("/", httpadapter.AppHandler(h.deleteSkill).ServeHTTP)
		})
	})
}

func (h *Handler) notesRoutes(r chi.Router) {
	r.Route("/notes", func(r chi.Router) {
		r.Get("/", httpadapter.AppHandler(h.getNotes).ServeHTTP)
		r.Post("/", httpadapter.AppHandler(h.createNote).ServeHTTP)

		r.Route("/{noteID}", func(r chi.Router) {
			r.Get("/", httpadapter.AppHandler(h.getNote).ServeHTTP)
			r.Put("/", httpadapter.AppHandler(h.updateNote).ServeHTTP)
			r.Delete("/", httpadapter.AppHandler(h.deleteNote).ServeHTTP)
		})
	})
}
