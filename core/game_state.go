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
	PtokensState PTokensState
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
	if wid < 0 || int(wid) >= len(w.States) {
		return fmt.Errorf("wonder id is out of range [0;%d)", len(w.States))
	}
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

// ===== PTokensState =====

type PTokensState struct {
	States [PTokensCount]struct {
		PlayerIndex
		PTokenState
	}
}

func (p PTokensState) IsInRange(pid PTokenID) error {
	if pid < 0 || int(pid) >= len(p.States) {
		return fmt.Errorf("ptoken id is out of range [0;%d)", len(p.States))
	}
	return nil
}

func (p PTokensState) CountBuiltByPlayer(pi PlayerIndex) (count int) {
	for _, ws := range p.States {
		if ws.PTokenState == PTokenTakenByPlayer && ws.PlayerIndex == pi {
			count++
		}
	}
	return count
}

func (p PTokensState) GetByState(s PTokenState) (out []PTokenID) {
	for i, pt := range p.States {
		if pt.PTokenState == s {
			out = append(out, PTokenID(i))
		}
	}
	return out
}

func (p PTokensState) IsTakeble(pid PTokenID) error {
	err := p.IsInRange(pid)
	if err != nil {
		return err
	}

	if p.States[pid].PTokenState != PTokenOnBoard {
		return fmt.Errorf("ptoken id = %d is not on board", pid)
	}
	return nil
}

func (p *PTokensState) setOnBoard(n int, rnd *rand.Rand) {
	for pid, s := range p.States {
		if s.PTokenState == PTokenOnBoard {
			p.States[pid].PTokenState = PTokenDiscarded
		}
	}
	discarded := p.GetByState(PTokenDiscarded)
	for _, id := range TakeNfromM(n, len(discarded), rnd) {
		p.States[discarded[id]].PTokenState = PTokenOnBoard
	}
}

func (p *PTokensState) setForChoose(n int, rnd *rand.Rand) (out []PTokenID) {
	for pid, s := range p.States {
		if s.PTokenState == PTokenChosenFromDiscarded {
			p.States[pid].PTokenState = PTokenDiscarded
		}
	}
	discarded := p.GetByState(PTokenDiscarded)
	for _, id := range TakeNfromM(n, len(discarded), rnd) {
		p.States[discarded[id]].PTokenState = PTokenChosenFromDiscarded
		out = append(out, discarded[id])
	}
	return out
}

func (p *PTokensState) takeFromChosen(pid PTokenID, pi PlayerIndex) error {
	err := p.IsInRange(pid)
	if err != nil {
		return err
	}

	if p.States[pid].PTokenState != PTokenChosenFromDiscarded {
		return fmt.Errorf("ptoken (id = %d): wrong status", pid)
	}

	p.States[pid].PTokenState = PTokenTakenByPlayer
	p.States[pid].PlayerIndex = pi
	return nil
}

func (p *PTokensState) take(pid PTokenID, pi PlayerIndex) error {
	err := p.IsTakeble(pid)
	if err != nil {
		return err
	}
	p.States[pid].PTokenState = PTokenTakenByPlayer
	p.States[pid].PlayerIndex = pi
	return nil
}

type PTokenState uint8

const (
	PTokenDiscarded PTokenState = iota
	PTokenOnBoard
	PTokenTakenByPlayer
	PTokenChosenFromDiscarded
)
