package components

import (
	"autoworld/modules/enginetype"

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
			Volume: 1.0,
			Pitch:  1.0,
		})
	})
}

// Audio retrieves the currently active audio object.
// Outputs: Returns the IAudioLW object associated with the current AudioName, or nil if none.
func (p AudioComponent) Audio() IAudioLW {
	if p.data == nil || p.data.AudioName == "" {
		return nil
	}
	return p.data.Audio[p.data.AudioName]
}

// SetAudio registers an audio object with a specific name.
// Inputs:
//   - audioName: The name to register the audio under.
//   - audio: The IAudioLW object to associate with the name.
func (p AudioComponent) SetAudio(audioName string, audio IAudioLW) {

	if p.data != nil {
		p.data.Audio[audioName] = audio
	}
}

// AudioName retrieves the name of the currently active audio.
// Outputs: Returns the name as a string, or an empty string if none is active.
func (p AudioComponent) AudioName() string {

	if p.data == nil {
		return ""
	}
	return p.data.AudioName
}

// SetAudioName sets the name of the currently active audio.
// Inputs:
//   - audioName: The name of the audio to set as active.
func (p AudioComponent) SetAudioName(audioName string) {
	if p.data != nil {
		p.data.AudioName = audioName
	}
}

// AudioSpeed retrieves the playback speed of the audio.
// Outputs: Returns the speed multiplier as a float32 (1 is default speed).
func (p AudioComponent) AudioSpeed() float32 {

	if p.data == nil {
		return 1
	}
	return p.data.AudioSpeed
}

// SetAudioSpeed sets the playback speed of the audio.
// Inputs:
//   - audioSpeed: The new speed multiplier (e.g., 1.0 for normal, 2.0 for double speed).
func (p AudioComponent) SetAudioSpeed(audioSpeed float32) {

	if p.data != nil {
		p.data.AudioSpeed = audioSpeed
	}
}

// Volume retrieves the current volume level.
// Outputs: Returns the volume as a float32 (1.0 is default volume).
func (p AudioComponent) Volume() float32 {

	if p.data == nil {
		return 1
	}
	return p.data.Volume
}

// SetVolume sets the volume level.
// Inputs:
//   - volume: The new volume level (e.g., 1.0 for full volume, 0.0 for muted).
func (p AudioComponent) SetVolume(volume float32) {

	if p.data != nil {
		p.data.Volume = volume
	}
}

// Pitch retrieves the current pitch level.
// Outputs: Returns the pitch as a float32 (1.0 is default pitch).
func (p AudioComponent) Pitch() float32 {

	if p.data == nil {
		return 1
	}
	return p.data.Pitch
}

// SetPitch sets the pitch level.
// Inputs:
//   - pitch: The new pitch level (e.g., 1.0 for normal pitch).
func (p AudioComponent) SetPitch(pitch float32) {

	if p.data != nil {
		p.data.Pitch = pitch
	}
}

// Play starts playing an audio with the specified volume and pitch.
// Inputs:
//   - name: The name of the audio to play.
//   - volume: The volume level to play at.
//   - pitch: The pitch level to play at.
func (p AudioComponent) Play(name string, volume float32, pitch float32) {

	if p.data == nil {
		return
	}
	p.data.AudioName = name
	p.data.Volume = volume
	p.data.Pitch = pitch
	p.data.ShouldPlay = true
	p.data.ShouldStop = false
}

// PlayDefault starts playing an audio with the currently set default volume and pitch.
// Inputs:
//   - name: The name of the audio to play.
func (p AudioComponent) PlayDefault(name string) {
	p.Play(name, p.data.Volume, p.data.Pitch)
}

// StopAudio stops the audio with the specified name.
// Inputs:
//   - name: The name of the audio to stop.
func (p AudioComponent) StopAudio(name string) {

	if p.data != nil {
		p.data.AudioName = name
		p.data.ShouldStop = true
	}
}

// PauseAudio pauses the audio with the specified name.
// Inputs:
//   - name: The name of the audio to pause.
func (p AudioComponent) PauseAudio(name string) {

	if p.data != nil {
		p.data.AudioName = name
		p.data.ShouldPause = true
	}
}

// ResumeAudio resumes playback of a paused audio with the specified name.
// Inputs:
//   - name: The name of the audio to resume.
func (p AudioComponent) ResumeAudio(name string) {

	if p.data != nil {
		p.data.AudioName = name
		p.data.ShouldResume = true
	}
}

// SetLooping sets whether the audio should loop continuously.
// Inputs:
//   - name: The name of the audio.
//   - loop: Boolean indicating whether to loop (true) or not (false).
func (p AudioComponent) SetLooping(name string, loop bool) {

	if p.data != nil {
		p.data.AudioName = name
		p.data.IsLooping = loop
	}
}

// IsLooping checks if the audio with the given name is currently set to loop.
// Inputs:
//   - name: The name of the audio to check.
// Outputs: Returns true if the audio is looping, false otherwise.
func (p AudioComponent) IsLooping(name string) bool {

	if p.data == nil {
		return false
	}
	// Check directly from the internal wrapper if possible, otherwise return the pending flag
	if audioLW, ok := p.data.Audio[name]; ok && audioLW != nil {
		return audioLW.IsLooping()
	}
	return p.data.IsLooping && p.data.AudioName == name
}
