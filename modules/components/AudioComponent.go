package components

import (
	"github.com/Nguyen-Agn/N-Engine/modules/enginetype"

	"github.com/yohamta/donburi"
)

// AudioComponent
type AudioComponent struct {
	IObject
	data *AudioData
}

// interface adready defined at domain

// resgiter new Component
var Audio = enginetype.Audio

// BindComponent binds the base object and its ECS data to this component.
// Inputs:
//   - base: The base IObject to bind to.
func (p *AudioComponent) BindComponent(base IObject) {
	p.IObject = base
	p.data = enginetype.GetComponent(base, Audio)
}

// init initializes the default data for the audio component.
// It registers the "aud" component token with default volume and pitch.
func init() {
	enginetype.RegisterComponentInitializer("aud", func(entry *donburi.Entry) {
		donburi.SetValue(entry, enginetype.Audio, AudioData{
			Audio:  make(map[string]IAudioLW),
			States: make(map[string]*AudioTrackState),
		})
	})
}

// getState is an internal helper to retrieve or lazily initialize the state for an audio track.
func (p AudioComponent) getState(name string) *AudioTrackState {
	if p.data == nil {
		return nil
	}
	if p.data.States == nil {
		p.data.States = make(map[string]*AudioTrackState)
	}
	state, ok := p.data.States[name]
	if !ok {
		state = &AudioTrackState{
			Volume: 1.0,
			Pitch:  1.0,
		}
		p.data.States[name] = state
	}
	return state
}

// Audio retrieves the currently active audio object by name.
// Outputs: Returns the IAudioLW object associated with the current AudioName, or nil if none.
func (p AudioComponent) Audio(name string) IAudioLW {
	if p.data == nil {
		return nil
	}
	return p.data.Audio[name]
}

// SetAudio registers an audio object with a specific name.
// Inputs:
//   - audioName: The name to register the audio under.
//   - audio: The IAudioLW object to associate with the name.
func (p AudioComponent) SetAudio(name string, audio IAudioLW) {
	if p.data != nil {
		p.data.Audio[name] = audio
	}
}

// Volume retrieves the current volume level.
// Outputs: Returns the volume as a float32 (1.0 is default volume).
func (p AudioComponent) Volume(name string) float32 {
	state := p.getState(name)
	if state == nil {
		return 1
	}
	return state.Volume
}

// SetVolume sets the volume level.
// Inputs:
//   - volume: The new volume level (e.g., 1.0 for full volume, 0.0 for muted).
func (p AudioComponent) SetVolume(name string, volume float32) {
	state := p.getState(name)
	if state != nil {
		state.Volume = volume
	}
}

// Pitch retrieves the current pitch level.
// Outputs: Returns the pitch as a float32 (1.0 is default pitch).
func (p AudioComponent) Pitch(name string) float32 {
	state := p.getState(name)
	if state == nil {
		return 1
	}
	return state.Pitch
}

// SetPitch sets the pitch level.
// Inputs:
//   - pitch: The new pitch level (e.g., 1.0 for normal pitch).
func (p AudioComponent) SetPitch(name string, pitch float32) {
	state := p.getState(name)
	if state != nil {
		state.Pitch = pitch
	}
}

// Play starts playing an audio with the specified volume and pitch.
// Inputs:
//   - name: The name of the audio to play.
//   - volume: The volume level to play at.
//   - pitch: The pitch level to play at.
func (p AudioComponent) Play(name string, volume float32, pitch float32) {
	state := p.getState(name)
	if state == nil {
		return
	}
	state.Volume = volume
	state.Pitch = pitch
	state.ShouldPlay = true
	state.ShouldStop = false
}

// StopAudio stops the audio with the specified name.
// Inputs:
//   - name: The name of the audio to stop.
func (p AudioComponent) StopAudio(name string) {
	state := p.getState(name)
	if state != nil {
		state.ShouldStop = true
	}
}

// PauseAudio pauses the audio with the specified name.
// Inputs:
//   - name: The name of the audio to pause.
func (p AudioComponent) PauseAudio(name string) {
	state := p.getState(name)
	if state != nil {
		state.ShouldPause = true
	}
}

// ResumeAudio resumes playback of a paused audio with the specified name.
// Inputs:
//   - name: The name of the audio to resume.
func (p AudioComponent) ResumeAudio(name string) {
	state := p.getState(name)
	if state != nil {
		state.ShouldResume = true
	}
}

// SetLooping sets whether the audio should loop continuously.
// Inputs:
//   - name: The name of the audio.
//   - loop: Boolean indicating whether to loop (true) or not (false).
func (p AudioComponent) SetLooping(name string, loop bool) {
	state := p.getState(name)
	if state != nil {
		state.IsLooping = loop
	}
}

// IsLooping checks if the audio with the given name is currently set to loop.
// Inputs:
//   - name: The name of the audio to check.
//
// Outputs: Returns true if the audio is looping, false otherwise.
func (p AudioComponent) IsLooping(name string) bool {
	state := p.getState(name)
	if state == nil {
		return false
	}
	// Check directly from the internal wrapper if possible, otherwise return the pending flag
	if audioLW, ok := p.data.Audio[name]; ok && audioLW != nil {
		return audioLW.IsLooping()
	}
	return state.IsLooping
}
