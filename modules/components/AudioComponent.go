package components

import (
	"autoworld/modules/enginetype"

	"github.com/yohamta/donburi"
)

// AudioComponent là Mixin để nhúng vào Custom Object.
type AudioComponent struct {
	IObject
	data *AudioData
}

// interface adready defined at domain

// resgiter new Component
var Audio = enginetype.Audio

func (p *AudioComponent) BindComponent(base IObject) {
	p.IObject = base
	p.data = enginetype.GetComponent(base, Audio)
}

func init() {
	enginetype.RegisterComponentInitializer("aud", func(entry *donburi.Entry) {
		donburi.SetValue(entry, enginetype.Audio, AudioData{
			Audio:  make(map[string]IAudioLW),
			Volume: 1.0,
			Pitch:  1.0,
		})
	})
}

func (p AudioComponent) Audio() IAudioLW {

	if p.data == nil || p.data.AudioName == "" {
		return nil
	}
	return p.data.Audio[p.data.AudioName]
}

func (p AudioComponent) SetAudio(audioName string, audio IAudioLW) {

	if p.data != nil {
		p.data.Audio[audioName] = audio
	}
}

func (p AudioComponent) AudioName() string {

	if p.data == nil {
		return ""
	}
	return p.data.AudioName
}

func (p AudioComponent) SetAudioName(audioName string) {

	if p.data != nil {
		p.data.AudioName = audioName
	}
}

func (p AudioComponent) AudioSpeed() float32 {

	if p.data == nil {
		return 1
	}
	return p.data.AudioSpeed
}

func (p AudioComponent) SetAudioSpeed(audioSpeed float32) {

	if p.data != nil {
		p.data.AudioSpeed = audioSpeed
	}
}

func (p AudioComponent) Volume() float32 {

	if p.data == nil {
		return 1
	}
	return p.data.Volume
}

func (p AudioComponent) SetVolume(volume float32) {

	if p.data != nil {
		p.data.Volume = volume
	}
}

func (p AudioComponent) Pitch() float32 {

	if p.data == nil {
		return 1
	}
	return p.data.Pitch
}

func (p AudioComponent) SetPitch(pitch float32) {

	if p.data != nil {
		p.data.Pitch = pitch
	}
}

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

func (p AudioComponent) PlayDefault(name string) {
	p.Play(name, 1.0, 1.0)
}

func (p AudioComponent) StopAudio(name string) {

	if p.data != nil {
		p.data.AudioName = name
		p.data.ShouldStop = true
	}
}

func (p AudioComponent) PauseAudio(name string) {

	if p.data != nil {
		p.data.AudioName = name
		p.data.ShouldPause = true
	}
}

func (p AudioComponent) ResumeAudio(name string) {

	if p.data != nil {
		p.data.AudioName = name
		p.data.ShouldResume = true
	}
}

func (p AudioComponent) SetLooping(name string, loop bool) {

	if p.data != nil {
		p.data.AudioName = name
		p.data.IsLooping = loop
	}
}

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
