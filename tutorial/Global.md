# Biến và Hằng số Toàn cục / Global Variables & Constants

> **Vision**: Khai báo trong TOML, đọc/ghi bằng một dòng code.
> **Vision**: Declare in TOML, read/write with a single line of code.

---

## 1. Giải thích / Explanation

N-Engine quản lý trạng thái game qua **Store** (napi.Store) — kho dữ liệu toàn cục.
N-Engine manages game state through the **Store** — a global data store.

- **`vars`** — biến toàn cục, có thể đọc và ghi / global variables, readable and writable.
- **`constants`** — hằng số, chỉ đọc / constants, read-only.

---

## 2. Khai báo trong TOML / Declare in TOML

Bạn có thể cấu hình sẵn các giá trị ban đầu trong `manifest.toml`:

```toml
# assets/manifest.toml

# Biến toàn cục / Global variables (read & write)
[[vars]]
key   = "score"
value = 0

[[vars]]
key   = "player_name"
value = "Hero"

[[vars]]
key   = "sound_on"
value = true

# Hằng số / Constants (read-only)
[[constants]]
key   = "gravity"
value = 980.0

[[constants]]
key   = "game_version"
value = "1.0.0"
```

---

## 3. Đọc Biến / Read Variables (napi.Store)

Sử dụng `napi.Store` để lấy giá trị biến một cách an toàn và đúng kiểu:

```go
import "github.com/Nguyen-Agn/N-Engine/modules/napi"

// Số nguyên
score := napi.Store.Int("score")           // int
bigNum := napi.Store.Int64("big_value")    // int64

// Số thực
vol   := napi.Store.Float32("master_volume") // float32
ratio := napi.Store.Float64("ratio")         // float64

// Chuỗi
name := napi.Store.String("player_name")   // string

// Boolean
soundOn := napi.Store.VarBool("sound_on")  // bool
```

---

## 4. Ghi Biến / Write Variables

Để cập nhật giá trị của biến:

```go
// Ghi biến đơn giản / Simple variable write
napi.Store.Value("score", 500)
napi.Store.Value("player_name", "Master")
napi.Store.Value("sound_on", false)
```

> **Lưu ý / Note:** Chỉ có thể ghi vào `vars`. 

---

## 5. Truy cập Hằng số / Read Constants (napi.Assert)

Hằng số được lưu độc lập và truy cập qua `napi.Assert.GetConst`. 

```go
// Trả về interface{}, cần ép kiểu (type assertion)
version := napi.Assert.GetConst("game_version").(string)
gravity := napi.Assert.GetConst("gravity").(float64)
```

**Tạo hằng số mới bằng code:**
```go
napi.Store.NewConst("max_level", 100)
```
