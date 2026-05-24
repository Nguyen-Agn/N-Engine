# AutoWorld — nAsset Module Reference

**Ý Nghĩa**:
Module `nasset` chịu trách nhiệm tải, nạp các tài nguyên (assets) của game (như hình ảnh, âm thanh) từ tập tin tĩnh và chuyển đổi chúng thành cấu trúc dữ liệu (`ISpriteLW`, `IAudioLW`) mà hệ thống engine có thể hiểu và sử dụng. Ngoài ra, module cung cấp `ManifestLoader` để đọc file JSON khai báo tài nguyên và tự động nạp toàn bộ vào `nGlobal` store.

**Chức năng chính**:
1. **Register.go**:
   Cửa ngõ kết nối với các interface/struct từ core `domain` và `domain/bridge`. Dùng Type Alias để local hóa các file/interface dùng chung, tránh import trực tiếp domain ở các nơi khác.

2. **ImageLoader.go**:
   Cung cấp cấu trúc `SpriteLoader` tuân thủ interface `domain.ISpriteLoader`.
   - **Tải ảnh linh hoạt**: `LoadSingle` tải 1 frame, `LoadSheet` cắt spritesheet thành nhiều frame (hỗ trợ `gapX`, `gapY`).
   - **Strategy Pattern**: Interface `IImageDecoder` — mỗi định dạng (`.png`, `.jpg`, `.gif`) có decoder riêng.
   - **Thay thế nguồn nạp**: `ReadFileFn` có thể override để nạp từ `embed.FS`.

3. **AudioLoader.go**:
   Cung cấp cấu trúc `AudioLoader` tuân thủ interface `domain.IAudioLoader`.
   - **Strategy Pattern**: Interface `IStreamDecoder` — mỗi định dạng (`.ogg`, `.wav`, `.mp3`) có decoder riêng. Có thể đăng ký thêm bằng `RegisterDecoder`.
   - **Flyweight**: Toàn bộ PCM data được đọc vào bộ nhớ một lần, nhiều thực thể có thể dùng chung buffer.
   - **Yêu cầu**: Cần inject `*audio.Context` (của Ebitengine) và `sampleRate` khi khởi tạo.

4. **ManifestLoader.go**:
   Cung cấp cấu trúc `ManifestLoader` triển khai `domain.IManifestLoader`.
   - Đọc file **TOML manifest** khai báo danh sách sprites, audios, vars, constants.
   - Tự động gọi `SpriteLoader` và `AudioLoader` rồi lưu kết quả vào `domain.IGlobal` store.
   - `audioLoader` có thể truyền `nil` nếu chưa cần — module bỏ qua phần audio không báo lỗi.

**Format file Manifest TOML**:
```toml
# Khai báo các sprite
[[sprites]]
key = "hero"
path = "assets/hero.png"
mode = "single"

[[sprites]]
key = "walk"
path = "assets/walk.png"
mode = "sheet"
frameW = 32
frameH = 32
cols = 4
rows = 2
gapX = 0
gapY = 0

# Khai báo audio
[[audios]]
key = "bgm"
path = "assets/bgm.ogg"

# Khai báo các biến và hằng số
[[vars]]
key = "player_score"
value = 0

[[constants]]
key = "gravity"
value = 9.8
```

**Cách sử dụng**:
```go
spriteLoader := nasset.NewSpriteLoader()
loader := nasset.NewManifestLoader(spriteLoader, nil) // nil = bỏ qua audio

store := nglobal.GetInstance()
if err := loader.LoadFromFile("assets/manifest.toml", store); err != nil {
    log.Fatal(err)
}
```

**Mối liên hệ với các thành phần khác**:
- **Domain**: Định nghĩa `ISpriteLoader`, `IAudioLoader`, `IManifestLoader`, `IGlobal`.
- **Domain Bridge**: Dùng `bridge.SpriteLW` để khởi tạo dữ liệu frame, tránh circular dependency với `nsprite`.
- **nGlobal**: Nhận kết quả load và lưu trữ toàn cục. `nasset` chỉ ghi, không đọc lại từ `nGlobal`.
- **nSprite / nAudio**: Phần render/phát âm thanh — `nasset` chỉ chuẩn bị dữ liệu thô.
