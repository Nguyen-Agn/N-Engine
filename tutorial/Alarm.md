# Bộ đếm thời gian / Alarms (ncom.Alrm)

> **Vision**: Hẹn giờ gọi hàm chỉ với một dòng code thay vì đếm frame thủ công.
> **Vision**: Schedule function calls with a single line of code instead of manual frame counting.

---

## 1. Giải thích / Explanation

Alarm Component (`ncom.Alrm`) giúp bạn thực thi một đoạn code (callback) sau một khoảng thời gian nhất định (tính bằng frame).
The Alarm Component (`ncom.Alrm`) helps you execute a piece of code (callback) after a certain amount of time (measured in frames).

Thay vì phải tạo biến `timer` và trừ dần trong `StepUpdate()`, bạn chỉ cần gắn ID cho Alarm và gọi `SetAlarm()`. Engine sẽ tự động kích hoạt callback và xóa Alarm sau khi hoàn thành.
Instead of creating a `timer` variable and subtracting it in `StepUpdate()`, you just assign an ID to the Alarm and call `SetAlarm()`. The engine will automatically trigger the callback and clear the Alarm upon completion.

---

## 2. Ví dụ / Code Example

```go
package objects

import (
	"fmt"

	"autoworld/modules/napi"
	"autoworld/modules/napi/ncom"
)

type Spawner struct {
	ncom.Object
	ncom.Alrm // Thêm Alarm component
}

func NewSpawner() *Spawner {
	s := &Spawner{}
	// "alr" là token của Alarm
	napi.Obj.NewObject(s, "MonsterSpawner", "alr sce-main")
	return s
}

func (s *Spawner) Create() {
	// Hẹn giờ Alarm số 0 chạy sau 120 frames (~2 giây nếu 60FPS)
	// Schedule Alarm ID 0 to run after 120 frames (~2 seconds at 60FPS)
	s.SetAlarm(0, 120, s.spawnMonster)
}

// Callback của Alarm
func (s *Spawner) spawnMonster() {
	fmt.Println("Đã sinh ra một quái vật! / A monster has spawned!")

	// Lặp lại chu kỳ (Tự động đặt lại Alarm)
	// Repeat the cycle (Re-schedule the Alarm)
	s.SetAlarm(0, 120, s.spawnMonster)
}

func (s *Spawner) StopSpawning() {
	// Hủy Alarm nếu không muốn nó kích hoạt nữa
	// Cancel the Alarm if you don't want it to trigger anymore
	s.StopAlarm(0)
}
```

## 3. Các API Khác / Other APIs
- `s.GetAlarm(id)`: Lấy số frame còn lại của Alarm tương ứng. Trả về `-1` nếu Alarm không tồn tại hoặc đã dừng. / Get remaining frames of the Alarm. Returns `-1` if not active.
- ID của Alarm là kiểu `int`, bạn có thể chạy nhiều Alarm song song bằng cách truyền vào các ID khác nhau (0, 1, 2...).
