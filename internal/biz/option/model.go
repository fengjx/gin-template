package option

import (
	sysoptionStore "gin-template/internal/store/sysoption"
	"gin-template/pkg/timex"
)

type statusResponse struct {
	Status string `json:"status"`
}

type optionValueResponse struct {
	Value string `json:"value"`
}

type pprofURLResponse struct {
	URL string `json:"url"`
}

type updateOptionRequest struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

type optionPayload struct {
	ID          string `json:"id"`
	OptionKey   string `json:"option_key"`
	OptionValue string `json:"option_value"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
	CTime       int64  `json:"ctime"`
	UTime       int64  `json:"utime"`
}

func itemsToPayload(items []sysoptionStore.Model) []optionPayload {
	resp := make([]optionPayload, 0, len(items))
	for _, item := range items {
		resp = append(resp, toOptionPayload(&item))
	}
	return resp
}

func toOptionPayload(item *sysoptionStore.Model) optionPayload {
	return optionPayload{
		ID:          item.ID,
		OptionKey:   item.OptionKey,
		OptionValue: item.OptionValue,
		Description: item.Description,
		IsPublic:    item.IsPublic,
		CTime:       timex.ToUnixSeconds(item.CTime),
		UTime:       timex.ToUnixSeconds(item.UTime),
	}
}
