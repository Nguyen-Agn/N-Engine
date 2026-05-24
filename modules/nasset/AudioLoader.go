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
	// Decode nhận dữ liệu thô và trả về stream PCM cùng số lượng kênh (channel) và sample rate.
	Decode(r io.ReadSeeker) (stream io.ReadSeeker, err error)
}

// ─── Decoder implementations ────────────────────────────────────────────────

// OGGDecoder giải mã file .ogg (Vorbis)
type OGGDecoder struct{ SampleRate int }

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

func NewWAVDecoder(sampleRate int) *WAVDecoder {
	return &WAVDecoder{
		SampleRate: sampleRate,
		adapter:    NewWAVAdapter(),
	}
}

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

// NewAudioLoader tạo một AudioLoader mới với các decoder mặc định (ogg, wav, mp3).
// ctx là Ebitengine audio.Context đã được khởi tạo sẵn (thường lấy từ AudioSystem).
// sampleRate: sample rate chung của game (ví dụ: 44100).
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

// RegisterDecoder cho phép đăng ký thêm decoder tùy chỉnh cho định dạng mới.
func (l *AudioLoader) RegisterDecoder(ext string, decoder IStreamDecoder) {
	l.decoders[ext] = decoder
}

// Load triển khai domain.IAudioLoader.
// Đọc file tại path, giải mã thành PCM bytes và tạo naudio.AudioLW sẵn sàng phát.
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
