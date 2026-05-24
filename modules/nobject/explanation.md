# AutoWorld — nObject Module Reference

**Ý Nghĩa**:
    Module `nobject` chịu trách nhiệm định nghĩa kiến trúc Object cơ bản và điều phối vòng đời logic (`StepUpdate`) cũng như xử lý Input callback của các đối tượng trong game.

**Chức năng chính**:
    1. **Register.go**:
        Nơi duy nhất kết nối với `domain`, thực hiện Type Alias cho `IObject`. Giúp toàn bộ module sử dụng kiểu dữ liệu nội bộ mà vẫn đảm bảo tính tương thích với Engine.
    2. **Object (Object.go)**:
        Lớp thủy tổ của mọi đối tượng trong game, triển khai giao diện `IObject`. Chứa các thuộc tính vật lý cơ bản (X, Y, Rotation, Scale) và các thành phần liên quan đến Sprite, Audio.
    3. **LogicSystem (LogicSystem.go)**:
        Hệ thống điều phối logic. Nó không chứa dữ liệu mà chỉ thực hiện việc duyệt qua danh sách các Object và kích hoạt sự kiện `StepUpdate`.
    4. **InputSystem (InputSystem.go)**:
        Hệ thống điều phối Input Callback. Mỗi frame, nó duyệt qua danh sách Object có `InputData`, kiểm tra trạng thái phím qua `IInputManager` và gọi hàm Handler tương ứng. Nhận `domain.IInputManager` qua interface để tránh phụ thuộc trực tiếp vào module `core`.

**Quy trình hoạt động**:
    1. **Update Flow** (trong Scene.Update):
       InputSystem.Update(objectList) → kích hoạt callbacks
       LogicSystem.Update(objectList) → gọi StepUpdate trên từng Object
    2. **Duyệt Object**:
       Hệ thống nhận danh sách `IObject` và lặp qua từng phần tử.
    3. **Kích hoạt Logic**:
       Gọi hàm `StepUpdate()` trên từng Object để thực thi các kịch bản logic cụ thể của đối tượng đó.