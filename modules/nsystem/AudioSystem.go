package nsystem

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

// AudioSystem is solely responsible for coordinating audio playback.
// It reads flags from AudioData and executes play/stop commands.
// This completely separates game logic from audio library calls.
type AudioSystem struct {
	// query is created once and reused for high performance.
	query *donburi.Query
}

// NewAudioSystem initializes an AudioSystem with a filter to query only entities containing AudioData.
// Outputs: Returns a pointer to a newly initialized AudioSystem.
func NewAudioSystem() *AudioSystem {
	return &AudioSystem{
		query: donburi.NewQuery(filter.Contains(Audio)),
	}
}

// Update iterates through all entities with AudioData and processes their play/stop flags.
// Inputs: w (donburi.World) - The ECS world containing the entities.
// Purpose: Only entities with the Audio Component are processed (thanks to the Query). It manages playing, pausing, resuming, stopping, and looping audio based on flags in AudioData.
func (this *AudioSystem) Update(w donburi.World) {
	this.query.Each(w, func(entry *donburi.Entry) {
		data := donburi.Get[AudioData](entry, Audio)

		// Skip if no audio name is set
		if data.AudioName == "" {
			return
		}

		// Retrieve IAudioLW from the map by current name
		audioLW, ok := data.Audio[data.AudioName]
		if !ok || audioLW == nil {
			return
		}

		// Update looping state for IAudioLW
		audioLW.SetLooping(data.IsLooping)

		// Continuously update volume if playing (to support real-time volume changes/fading)
		if audioLW.IsPlaying() {
			audioLW.SetVolume(data.Volume)
		}

		// Process STOP command first (highest priority)
		if data.ShouldStop {
			audioLW.Stop()
			data.ShouldStop = false // Reset flag after processing
			return
		}

		// Process PAUSE command
		if data.ShouldPause {
			audioLW.Pause()
			data.ShouldPause = false
			return
		}

		// Process RESUME command
		if data.ShouldResume {
			audioLW.Resume()
			data.ShouldResume = false
			return
		}

		// Process NEW PLAY command
		if data.ShouldPlay {
			// Only play if not already playing (prevent overlapping)
			if !audioLW.IsPlaying() {
				audioLW.Play(data.AudioName, data.Volume, data.Pitch)
			}
			data.ShouldPlay = false // Reset flag after processing
		} else {
			// Process automatic looping
			if data.IsLooping {
				// If not playing, not paused, and NOT forcibly stopped
				if !audioLW.IsPlaying() && !audioLW.IsPaused() && !audioLW.IsStopped() {
					audioLW.Play(data.AudioName, data.Volume, data.Pitch)
				}
			}
		}
	})
}
