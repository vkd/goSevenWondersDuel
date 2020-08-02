package core

import (
	"fmt"
	"math/rand"
)

const (
	MaxBuiltWonders = 7
)

type GameState struct {
	WondersState WondersState
}

func InitializeWonders(w WondersState, rnd *rand.Rand) (WondersState, error) {
	for _, ws := range w.States {
		if ws.InGame {
			return w, fmt.Errorf("is already initialized")
		}
	}
	for _, idx := range TakeNfromM(InitialWonders, len(w.States), rnd) {
		w.States[idx].InGame = true
	}
	return w, nil
}

type WonderState struct {
	PlayerIndex
	InGame   bool
	IsChosen bool
	IsBuilt  bool
}

type WondersState struct {
	// States index is a WonderID and value is a WonderState
	States [WondersCount]WonderState
}

func (w WondersState) IsBuildable(wid WonderID, pi PlayerIndex) error {
	if w.CountBuilt() >= MaxBuiltWonders {
		return fmt.Errorf("reach maximum of built wonders: %d", MaxBuiltWonders)
	}
	ws := w.States[wid]
	if !ws.InGame {
		return fmt.Errorf("wonder id = %d is not in game", wid)
	}
	if !ws.IsChosen {
		return fmt.Errorf("wonder id = %d is not chosen", wid)
	}
	if ws.PlayerIndex != pi {
		return fmt.Errorf("wonder id = %d is not related to %d player (actually to %d player)", wid, pi, ws.PlayerIndex)
	}
	if ws.IsBuilt {
		return fmt.Errorf("wonder id = %d is alteady built", wid)
	}
	return nil
}

func (w WondersState) CountBuilt() (out int) {
	for _, ws := range w.States {
		if ws.IsBuilt {
			out++
		}
	}
	return out
}

func (w WondersState) CountBuiltByPlayer(pi PlayerIndex) (out int) {
	for _, ws := range w.States {
		if ws.IsBuilt && ws.PlayerIndex == pi {
			out++
		}
	}
	return out
}

func (w WondersState) AvailableToChoose() (out [InitialWonders]WonderID, _ error) {
	var i = 0
	for wid, s := range w.States {
		if s.InGame && !s.IsChosen {
			out[i] = WonderID(wid)
			i++
		}
	}
	if i != len(out) {
		return out, fmt.Errorf("wrong amount of available wonders (expect: %d): %d", len(out), i)
	}
	return out, nil
}

func (w WondersState) AvailableToBuild(pi PlayerIndex) (out []WonderID) {
	for wid, ws := range w.States {
		if ws.IsChosen && ws.PlayerIndex == pi && !ws.IsBuilt {
			out = append(out, WonderID(wid))
		}
	}
	return out
}

func (w *WondersState) chooseByPlayer(wid WonderID, pi PlayerIndex) error {
	ws := w.States[wid]
	if !ws.InGame {
		return fmt.Errorf("%d wonder is not in game", wid)
	}
	if ws.IsChosen {
		return fmt.Errorf("%d wonder is already choosen by %d player", wid, ws.PlayerIndex)
	}

	w.States[wid].IsChosen = true
	w.States[wid].PlayerIndex = pi
	return nil
}

func (w *WondersState) built(wid WonderID, pi PlayerIndex) error {
	err := w.IsBuildable(wid, pi)
	if err != nil {
		return fmt.Errorf("wonder is not buildable: %w", err)
	}
	w.States[wid].IsBuilt = true
	return nil
}
