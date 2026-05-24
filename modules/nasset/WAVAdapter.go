package nasset

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// ─── WAV Header structs ───────────────────────────────────────────────────────

// wavHeader là phần đầu RIFF của file WAV.
type wavHeader struct {
	ChunkID   [4]byte // "RIFF"
	ChunkSize uint32
	Format    [4]byte // "WAVE"
}

// fmtChunk là chunk định dạng âm thanh (fmt ).
type fmtChunk struct {
	Subchunk1ID   [4]byte // "fmt "
	Subchunk1Size uint32
	AudioFormat   uint16 // 1 = PCM
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
}

// ─── WAVAdapter ──────────────────────────────────────────────────────────────

// WAVAdapter là adapter xử lý file WAV với các bit depth khác nhau.
// Nếu file là 24-bit PCM (không được Ebitengine hỗ trợ), adapter tự động
// chuyển đổi xuống 16-bit để tương thích. Các bit depth khác (8/16-bit) được
// chuyển thẳng mà không thay đổi.
type WAVAdapter struct{}

// NewWAVAdapter tạo một WAVAdapter mới.
func NewWAVAdapter() *WAVAdapter {
	return &WAVAdapter{}
}

// Adapt nhận dữ liệu WAV thô từ r, kiểm tra bit depth và convert nếu cần.
// Trả về io.ReadSeeker chứa WAV hợp lệ (8 hoặc 16-bit) sẵn sàng cho Ebitengine.
func (a *WAVAdapter) Adapt(r io.ReadSeeker) (io.ReadSeeker, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("wav adapter: không đọc được dữ liệu: %w", err)
	}

	bits, err := a.readBitsPerSample(data)
	if err != nil {
		return nil, err
	}

	// 8-bit và 16-bit tương thích trực tiếp với Ebitengine — chuyển thẳng
	if bits == 8 || bits == 16 {
		return bytes.NewReader(data), nil
	}

	// 24-bit: cần convert xuống 16-bit
	if bits == 24 {
		converted, err := a.convert24to16(data)
		if err != nil {
			return nil, fmt.Errorf("wav adapter: lỗi convert 24→16 bit: %w", err)
		}
		return bytes.NewReader(converted), nil
	}

	return nil, fmt.Errorf("wav adapter: bit depth %d không được hỗ trợ (chỉ hỗ trợ 8, 16, 24)", bits)
}

// readBitsPerSample parse header WAV và đọc giá trị BitsPerSample.
func (a *WAVAdapter) readBitsPerSample(data []byte) (uint16, error) {
	r := bytes.NewReader(data)

	var riff wavHeader
	if err := binary.Read(r, binary.LittleEndian, &riff); err != nil {
		return 0, fmt.Errorf("wav adapter: parse RIFF header lỗi: %w", err)
	}
	if string(riff.ChunkID[:]) != "RIFF" || string(riff.Format[:]) != "WAVE" {
		return 0, fmt.Errorf("wav adapter: không phải file WAV hợp lệ")
	}

	// Duyệt qua các chunk để tìm chunk "fmt "
	for {
		var chunkID [4]byte
		var chunkSize uint32
		if err := binary.Read(r, binary.LittleEndian, &chunkID); err != nil {
			return 0, fmt.Errorf("wav adapter: không tìm thấy chunk fmt")
		}
		if err := binary.Read(r, binary.LittleEndian, &chunkSize); err != nil {
			return 0, fmt.Errorf("wav adapter: đọc chunk size lỗi")
		}

		if string(chunkID[:]) == "fmt " {
			var audioFormat, numChannels, blockAlign, bitsPerSample uint16
			var sampleRate, byteRate uint32
			binary.Read(r, binary.LittleEndian, &audioFormat)
			binary.Read(r, binary.LittleEndian, &numChannels)
			binary.Read(r, binary.LittleEndian, &sampleRate)
			binary.Read(r, binary.LittleEndian, &byteRate)
			binary.Read(r, binary.LittleEndian, &blockAlign)
			binary.Read(r, binary.LittleEndian, &bitsPerSample)
			return bitsPerSample, nil
		}

		// Bỏ qua chunk này (align to even byte boundary)
		skip := int64(chunkSize)
		if chunkSize%2 != 0 {
			skip++
		}
		r.Seek(skip, io.SeekCurrent)
	}
}

// convert24to16 chuyển đổi toàn bộ WAV 24-bit PCM sang 16-bit PCM.
// Cách làm: lấy 2 byte cao của mỗi sample 3-byte (bỏ byte thấp nhất).
// Rebuild file WAV hoàn chỉnh với header cập nhật.
func (a *WAVAdapter) convert24to16(data []byte) ([]byte, error) {
	r := bytes.NewReader(data)

	// Đọc RIFF header
	var riff wavHeader
	binary.Read(r, binary.LittleEndian, &riff)

	var out bytes.Buffer

	// Duyệt từng chunk, xử lý riêng "fmt " và "data"
	for {
		var chunkID [4]byte
		var chunkSize uint32

		if err := binary.Read(r, binary.LittleEndian, &chunkID); err != nil {
			break // hết file
		}
		if err := binary.Read(r, binary.LittleEndian, &chunkSize); err != nil {
			break
		}

		chunkData := make([]byte, chunkSize)
		io.ReadFull(r, chunkData)

		// Align: WAV chunks phải có kích thước chẵn
		if chunkSize%2 != 0 {
			r.Seek(1, io.SeekCurrent)
		}

		id := string(chunkID[:])

		switch id {
		case "fmt ":
			// Ghi chunk fmt đã sửa: 24-bit → 16-bit
			cr := bytes.NewReader(chunkData)
			var audioFmt, numCh, blockAlign, bitsPS uint16
			var sampleRate, byteRate uint32
			binary.Read(cr, binary.LittleEndian, &audioFmt)
			binary.Read(cr, binary.LittleEndian, &numCh)
			binary.Read(cr, binary.LittleEndian, &sampleRate)
			binary.Read(cr, binary.LittleEndian, &byteRate)
			binary.Read(cr, binary.LittleEndian, &blockAlign)
			binary.Read(cr, binary.LittleEndian, &bitsPS)

			// Tính lại BlockAlign và ByteRate cho 16-bit
			newBitsPS := uint16(16)
			newBlockAlign := numCh * (newBitsPS / 8)
			newByteRate := sampleRate * uint32(newBlockAlign)

			var fmtBuf bytes.Buffer
			binary.Write(&fmtBuf, binary.LittleEndian, audioFmt)
			binary.Write(&fmtBuf, binary.LittleEndian, numCh)
			binary.Write(&fmtBuf, binary.LittleEndian, sampleRate)
			binary.Write(&fmtBuf, binary.LittleEndian, newByteRate)
			binary.Write(&fmtBuf, binary.LittleEndian, newBlockAlign)
			binary.Write(&fmtBuf, binary.LittleEndian, newBitsPS)

			out.WriteString("fmt ")
			binary.Write(&out, binary.LittleEndian, uint32(fmtBuf.Len()))
			out.Write(fmtBuf.Bytes())

		case "data":
			// Convert mỗi sample 3-byte → 2-byte (lấy byte 1 và 2, bỏ byte 0)
			converted := make([]byte, 0, len(chunkData)*2/3)
			for i := 0; i+2 < len(chunkData); i += 3 {
				converted = append(converted, chunkData[i+1], chunkData[i+2])
			}
			out.WriteString("data")
			binary.Write(&out, binary.LittleEndian, uint32(len(converted)))
			out.Write(converted)

		default:
			// Các chunk khác (LIST, INFO...) giữ nguyên
			out.Write(chunkID[:])
			binary.Write(&out, binary.LittleEndian, chunkSize)
			out.Write(chunkData)
		}
	}

	// Ghép RIFF header + body
	body := out.Bytes()
	var final bytes.Buffer
	final.WriteString("RIFF")
	binary.Write(&final, binary.LittleEndian, uint32(len(body)+4))
	final.WriteString("WAVE")
	final.Write(body)

	return final.Bytes(), nil
}
