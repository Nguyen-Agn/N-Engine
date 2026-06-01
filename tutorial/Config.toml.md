# Cấu hình Tài nguyên / Resource Config (TOML)

> **Vision**: Khai báo tài nguyên một lần trong TOML, dùng ở khắp nơi.
> **Vision**: Declare resources once in TOML, use everywhere.

---

## 1. Giải thích / Explanation

N-Engine dùng file **manifest.toml** để khai báo tài nguyên:
N-Engine uses **manifest.toml** to declare resources:

- `[[sprites]]` — hình ảnh / images
- `[[audios]]` — âm thanh / audio files
- `[[vars]]` — biến toàn cục / global variables
- `[[constants]]` — hằng số / constants

Nạp một lần với `napi.Game.LoadFromFile(path)`, dùng mọi nơi qua `napi.Assert.*` và `napi.Store.*`.
Load once with `napi.Game.LoadFromFile(path)`, use anywhere via `napi.Assert.*` and `napi.Store.*`.

---

## 2. Cấu trúc File / File Structure

```toml
# assets/manifest.toml

# ── Sprites ──────────────────────────────────────────

# Sprite đơn / Single sprite
[[sprites]]
key  = "player_idle"
path = "assets/images/player_idle.png"

# Sprite sheet (cắt nhiều frame) / Sprite sheet (multiple frames)
[[sprites]]
key         = "player_run"
path        = "assets/images/player_run.png"
cols        = 8
rows        = 1

# ── Audios ───────────────────────────────────────────

[[audios]]
key  = "bgm_main"
path = "assets/audio/main_theme.ogg"

[[audios]]
key  = "sfx_jump"
path = "assets/audio/jump.wav"

# ── Variables (có thể thay đổi) / Mutable variables ──

[[vars]]
key   = "score"
value = 0       # int

[[vars]]
key   = "player_name"
value = "Hero"  # string

[[vars]]
key   = "sound_on"
value = true    # bool

# ── Constants (chỉ đọc) / Read-only constants ─────────

[[constants]]
key   = "gravity"
value = 980.0

[[constants]]
key   = "max_speed"
value = 400
```

---

## 3. Nạp và Sử dụng / Load and Use

**Nạp tài nguyên:**
```go
import "autoworld/modules/napi"

func main() {
	// 1. Phải gọi Init trước
	napi.Game.Init(napi.GameConfig{...})

	// 2. Load manifest
	napi.Game.LoadFromFile("assets/manifest.toml")
}
```

**Truy cập tài nguyên:**
```go
// Truy cập Sprite
spr := napi.Assert.GetSprite("player_run")

// Truy cập Audio
aud := napi.Assert.GetAudio("sfx_jump")

// Lấy/Ghi Biến (Variables)
score := napi.Store.Int("score")
napi.Store.Value("score", score + 10) // Set value

// Lấy Hằng số (Constants)
grav := napi.Assert.GetConst("gravity").(float64)
```
