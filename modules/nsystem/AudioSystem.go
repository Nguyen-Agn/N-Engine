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

		if data.States == nil {
			return
		}

		for audioName, state := range data.States {
			// Retrieve IAudioLW from the map by current name
			audioLW, ok := data.Audio[audioName]
			if !ok || audioLW == nil {
				continue
			}

			// Update looping state for IAudioLW
			audioLW.SetLooping(state.IsLooping)

			// Continuously update volume if playing (to support real-time volume changes/fading)
			if audioLW.IsPlaying() {
				audioLW.SetVolume(state.Volume)
			}

			// Process STOP command first (highest priority)
			if state.ShouldStop {
				audioLW.Stop()
				state.ShouldStop = false // Reset flag after processing
				continue
			}

			// Process PAUSE command
			if state.ShouldPause {
				audioLW.Pause()
				state.ShouldPause = false
				continue
			}

			// Process RESUME command
			if state.ShouldResume {
				audioLW.Resume()
				state.ShouldResume = false
				continue
			}

			// Process NEW PLAY command
			if state.ShouldPlay {
				// Only play if not already playing (prevent overlapping)
				if !audioLW.IsPlaying() {
					audioLW.Play(audioName, state.Volume, state.Pitch)
				}
				state.ShouldPlay = false // Reset flag after processing
			} else {
				// Process automatic looping
				if state.IsLooping {
					// If not playing, not paused, and NOT forcibly stopped
					if !audioLW.IsPlaying() && !audioLW.IsPaused() && !audioLW.IsStopped() {
						audioLW.Play(audioName, state.Volume, state.Pitch)
					}
				}
			}
		}
	})
}
