package domain

// IGlobalConfig đại diện cho trình quản lý cấu hình toàn cục của Engine.
// Triển khai theo Singleton Pattern kết hợp Observer Pattern:
// các module đăng ký IObserver để nhận thông báo khi config thay đổi.
type IGlobalConfig interface {
	// GetConfig trả về instance Singleton của IGlobalConfig.
	// Mọi nơi trong hệ thống đều gọi qua method này thay vì tham chiếu trực tiếp.
	// Purpose: Retrieves the Singleton instance of the global configuration manager.
	// Inputs: None.
	// Outputs: IGlobalConfig - The single active configuration instance.
	GetConfig() IGlobalConfig

	// NotifyChange thông báo đến tất cả IObserver đã đăng ký rằng config vừa thay đổi.
	// Gọi sau khi cập nhật bất kỳ giá trị nào trong config.
	// Purpose: Triggers a notification to all registered observers that the configuration has been updated.
	// Inputs: None.
	// Outputs: None.
	NotifyChange()

	// AddObserver đăng ký một IObserver để nhận thông báo khi config thay đổi.
	// Purpose: Registers an observer to listen for configuration changes.
	// Inputs: observer IObserver - The observer object to add.
	// Outputs: None.
	AddObserver(observer IObserver)

	// RemoveObserver hủy đăng ký một IObserver khỏi danh sách nhận thông báo.
	// Purpose: Unregisters an observer to stop listening for configuration changes.
	// Inputs: observer IObserver - The observer object to remove.
	// Outputs: None.
	RemoveObserver(observer IObserver)

	// ShareGlobal nhúng giao diện đọc/ghi biến key-value có kiểu,
	// cho phép lưu trữ các giá trị cấu hình động (title, width, height, v.v.)
	ShareGlobal
}

// IObserver định nghĩa hợp đồng nhận thông báo khi IGlobalConfig thay đổi.
// Implement interface này và đăng ký qua AddObserver để theo dõi config.
type IObserver interface {
	// NotifyChange được gọi tự động khi IGlobalConfig.NotifyChange() được kích hoạt.
	// config là instance config hiện tại để Observer đọc giá trị mới.
	// Purpose: Callback invoked automatically when the global configuration changes.
	// Inputs: config IGlobalConfig - The updated configuration instance.
	// Outputs: None.
	NotifyChange(config IGlobalConfig)
}
