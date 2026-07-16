package characterDTO

import (
	"bytes"
	"encoding/json"
)

type CharacterRequest struct {
	Name       string  `json:"name"`
	Occupation *string `json:"occupation"`
	Age        *int16  `json:"age"`
	Sex        *string `json:"sex"`
	Residence  *string `json:"residence"`
	Birthplace *string `json:"birthplace"`
}

type PatchValue[T any] struct {
	Set   bool
	Value *T
}

func (v *PatchValue[T]) UnmarshalJSON(data []byte) error {
	v.Set = true
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		v.Value = nil
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	v.Value = &value
	return nil
}

type PatchCharacterRequest struct {
	Name       PatchValue[string] `json:"name"`
	Occupation PatchValue[string] `json:"occupation"`
	Age        PatchValue[int16]  `json:"age"`
	Sex        PatchValue[string] `json:"sex"`
	Residence  PatchValue[string] `json:"residence"`
	Birthplace PatchValue[string] `json:"birthplace"`
}

func (r PatchCharacterRequest) HasChanges() bool {
	return r.Name.Set ||
		r.Occupation.Set ||
		r.Age.Set ||
		r.Sex.Set ||
		r.Residence.Set ||
		r.Birthplace.Set
}
