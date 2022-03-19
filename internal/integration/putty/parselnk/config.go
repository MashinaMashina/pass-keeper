package parselnk

import "pass-keeper/internal/config"

func fillHandler(h *config.Part) {
	d := h.DefaultValues()
	d["lnk.replace"] = "PUTTY -=>"
	h.SetDefaultValues(d)
}
