package nasset

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"autoworld/domain"
	"autoworld/modules/naudio"

	ebitenAudio "github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

// IStreamDecoder định nghĩa chiến lược (Strategy) giải mã từng định dạng âm thanh.
// Mỗi decoder nhận một io.ReadSeeker và trả về một io.ReadSeeker đã được decode
// sẵn sàng để Ebitengine đọc dữ liệu PCM.
type IStreamDecoder interface {
	// Purpose: Decodes raw audio data into a PCM stream.
	// Inputs: r - io.ReadSeeker containing the raw audio data.
	// Outputs: stream - io.ReadSeeker containing the decoded PCM data, err - error if decoding fails.
	Decode(r io.ReadSeeker) (stream io.ReadSeeker, err error)
}

// ─── Decoder implementations ────────────────────────────────────────────────

// OGGDecoder giải mã file .ogg (Vorbis)
type OGGDecoder struct{ SampleRate int }

// Purpose: Decodes OGG format audio data.
// Inputs: r (io.ReadSeeker) - The raw audio data stream.
// Outputs: (io.ReadSeeker) - The decoded PCM audio stream, (error) - Error if decoding fails.
func (d *OGGDecoder) Decode(r io.ReadSeeker) (io.ReadSeeker, error) {
	s, err := vorbis.DecodeWithSampleRate(d.SampleRate, r)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// WAVDecoder giải mã file .wav — tự động xử lý 24-bit qua WAVAdapter.
type WAVDecoder struct {
	SampleRate int
	adapter    *WAVAdapter
}

// Purpose: Initializes a WAVDecoder with the given sample rate and a new WAVAdapter.
// Inputs: sampleRate (int) - The sample rate to use for decoding.
// Outputs: (*WAVDecoder) - A pointer to the newly created WAVDecoder.
func NewWAVDecoder(sampleRate int) *WAVDecoder {
	return &WAVDecoder{
		SampleRate: sampleRate,
		adapter:    NewWAVAdapter(),
	}
}

// Purpose: Adapts the WAV file (e.g., converting 24-bit to 16-bit) and decodes it.
// Inputs: r (io.ReadSeeker) - The raw audio data stream.
// Outputs: (io.ReadSeeker) - The decoded PCM audio stream, (error) - Error if decoding or adapting fails.
func (d *WAVDecoder) Decode(r io.ReadSeeker) (io.ReadSeeker, error) {
	// Chạy qua adapter trước: nếu là 24-bit sẽ được convert sang 16-bit
	adapted, err := d.adapter.Adapt(r)
	if err != nil {
		return nil, err
	}
	s, err := wav.DecodeWithSampleRate(d.SampleRate, adapted)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// MP3Decoder giải mã file .mp3
type MP3Decoder struct{ SampleRate int }

// Purpose: Decodes MP3 format audio data.
// Inputs: r (io.ReadSeeker) - The raw audio data stream.
// Outputs: (io.ReadSeeker) - The decoded PCM audio stream, (error) - Error if decoding fails.
func (d *MP3Decoder) Decode(r io.ReadSeeker) (io.ReadSeeker, error) {
	s, err := mp3.DecodeWithSampleRate(d.SampleRate, r)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// ─── AudioLoader ─────────────────────────────────────────────────────────────

// AudioLoader triển khai domain.IAudioLoader.
// Đọc file âm thanh từ đĩa, giải mã theo định dạng và tạo ra IAudioLW
// chứa toàn bộ dữ liệu PCM sẵn sàng phát qua Ebitengine audio.Context.
type AudioLoader struct {
	ctx      *ebitenAudio.Context
	decoders map[string]IStreamDecoder
}

// Purpose: Initializes an AudioLoader with default decoders (OGG, WAV, MP3).
// Inputs: ctx (*ebitenAudio.Context) - The Ebitengine audio context, sampleRate (int) - The default sample rate for the game.
// Outputs: (*AudioLoader) - A pointer to the newly created AudioLoader.
func NewAudioLoader(ctx *ebitenAudio.Context, sampleRate int) *AudioLoader {
	return &AudioLoader{
		ctx: ctx,
		decoders: map[string]IStreamDecoder{
			".ogg": &OGGDecoder{SampleRate: sampleRate},
			".wav": NewWAVDecoder(sampleRate),
			".mp3": &MP3Decoder{SampleRate: sampleRate},
		},
	}
}

// Purpose: Allows extending the supported audio formats by adding custom decoders.
// Inputs: ext (string) - The file extension (e.g., ".flac"), decoder (IStreamDecoder) - The decoder strategy for the format.
func (l *AudioLoader) RegisterDecoder(ext string, decoder IStreamDecoder) {
	l.decoders[ext] = decoder
}

// Purpose: Reads an audio file, decodes it using the registered decoder for its extension, and creates an IAudioLW object.
// Inputs: path (string) - The file path to the audio file.
// Outputs: (domain.IAudioLW) - The loaded audio object ready for playback, (error) - Error if the format is unsupported, file cannot be read, or decoding fails.
func (l *AudioLoader) Load(path string) (domain.IAudioLW, error) {
	ext := filepath.Ext(path)
	decoder, ok := l.decoders[ext]
	if !ok {
		return nil, fmt.Errorf("audio loader: định dạng không hỗ trợ '%s'", ext)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("audio loader: không thể mở file '%s': %w", path, err)
	}
	defer f.Close()

	stream, err := decoder.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("audio loader: lỗi giải mã '%s': %w", path, err)
	}

	// Đọc toàn bộ PCM data vào bộ nhớ (buffer dùng chung — Flyweight pattern).
	data, err := io.ReadAll(stream)
	if err != nil {
		return nil, fmt.Errorf("audio loader: lỗi đọc stream '%s': %w", path, err)
	}

	return naudio.NewAudioLW(l.ctx, data), nil
}
