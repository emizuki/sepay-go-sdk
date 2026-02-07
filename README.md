# SePay Go SDK

Thư viện Go SDK không chính thức của cổng thanh toán SePay. Hỗ trợ các hình thức tích hợp thanh toán quét mã chuyển khoản ngân hàng VietQr, quét mã chuyển khoản Napas, thanh toán qua thẻ quốc tế/nội địa Visa/Master Card/JCB.

## Yêu cầu

- Go 1.21 hoặc cao hơn

## Cài đặt

```
go get github.com/emizuki/sepay-go-sdk
```

## Bắt đầu nhanh

Thư viện cần cấu hình các thông tin từ đơn vị bán hàng của bạn (merchant) được lấy từ my.sepay.vn.

```go
import "github.com/emizuki/sepay-go-sdk"

client, err := sepay.NewClient(sepay.Config{
	Env:        sepay.Sandbox,
	MerchantID: "YOUR_MERCHANT_ID",
	SecretKey:  "YOUR_MERCHANT_SECRET_KEY",
})
if err != nil {
	log.Fatal(err)
}
```

Tạo biểu mẫu thanh toán cho đơn hàng có mã `DH0001`, số tiền thanh toán là 10.000đ. Sau khi thanh toán thành công sẽ tự động chuyển hướng về đường dẫn: https://example.com/order/DH0001

```go
signed := client.Checkout.InitOneTimePaymentFields(sepay.OnetimePaymentFields{
	Operation:          sepay.OperationPurchase,
	OrderInvoiceNumber: "DH0001",
	OrderAmount:        10000,
	Currency:           "VND",
	SuccessURL:         sepay.String("https://example.com/order/DH0001"),
	OrderDescription:   "Thanh toan don hang DH0001",
})

fmt.Println(signed)
```

Kết quả trả về kiểu `*SignedCheckoutFields`:

```go
&sepay.SignedCheckoutFields{
	Merchant:           "YOUR_MERCHANT_ID",
	Operation:          "PURCHASE",
	OrderInvoiceNumber: "DH0001",
	OrderAmount:        10000,
	Currency:           "VND",
	OrderDescription:   "Thanh toan don hang DH0001",
	SuccessURL:         sepay.String("https://example.com/order/DH0001"),
	Signature:          "AUTO_GENERATED_SIGNATURE",
}
```

Kết hợp với `client.Checkout.InitCheckoutURL()` để lấy URL và tạo `form` xử lý thanh toán sau khi đã tạo đơn hàng. Sử dụng `signed.FormValues()` để lấy các trường dạng `map[string]string`:

```go
checkoutURL := client.Checkout.InitCheckoutURL()
formValues := signed.FormValues()
```

```html
<form action="{{.CheckoutURL}}" method="POST">
  {{range $key, $value := .FormValues}}
  <input type="hidden" name="{{$key}}" value="{{$value}}" />
  {{end}}
  <button type="submit">Thanh toán</button>
</form>
```

## Cấu hình

```go
client, err := sepay.NewClient(sepay.Config{
	Env:        sepay.Sandbox,
	MerchantID: "YOUR_MERCHANT_ID",
	SecretKey:  "YOUR_MERCHANT_SECRET_KEY",
})
```

| Tham số             | Mô tả                                                                                    |
| ------------------- | ---------------------------------------------------------------------------------------- |
| **Env**             | Môi trường hiện tại, giá trị hỗ trợ: `sepay.Sandbox`, `sepay.Production`                 |
| **MerchantID**      | Mã đơn vị merchant                                                                       |
| **SecretKey**       | Khóa bảo mật merchant                                                                    |
| **APIVersion**      | Phiên bản API sử dụng, giá trị hỗ trợ: `sepay.APIVersionV1` (mặc định)                   |
| **CheckoutVersion** | Phiên bản trang thanh toán sử dụng, giá trị hỗ trợ: `sepay.CheckoutVersionV1` (mặc định) |

## Khởi tạo đối tượng cho biểu mẫu thanh toán

Sử dụng `client.Checkout.InitCheckoutURL()` để tạo URL thanh toán theo thông tin đã cấu hình.

| Môi trường     | URL thanh toán                                |
| -------------- | --------------------------------------------- |
| **sandbox**    | https://pay-sandbox.sepay.vn/v1/checkout/init |
| **production** | https://pay.sepay.vn/v1/checkout/init         |

### Đơn hàng thanh toán 1 lần

```go
signed := client.Checkout.InitOneTimePaymentFields(sepay.OnetimePaymentFields{
	Operation:          sepay.OperationPurchase,
	PaymentMethod:      sepay.BankTransfer,
	OrderInvoiceNumber: "DH0001",
	OrderAmount:        10000,
	Currency:           "VND",
	OrderDescription:   "Thanh toan don hang DH0001",
	CustomerID:         sepay.String("KH001"),
	SuccessURL:         sepay.String("https://example.com/success"),
	ErrorURL:           sepay.String("https://example.com/error"),
	CancelURL:          sepay.String("https://example.com/cancel"),
	CustomData:         sepay.String("custom_value"),
})
```

Tham số đầu vào (`OnetimePaymentFields`):

| Tham số                | Bắt buộc | Mô tả                                                                   |
| ---------------------- | -------- | ----------------------------------------------------------------------- |
| **Operation**          | ✔︎        | Loại giao dịch, hiện chỉ hỗ trợ: `sepay.OperationPurchase`              |
| **PaymentMethod**      | ✔︎        | Phương thức thanh toán: `sepay.BankTransfer`, `sepay.NapasBankTransfer` |
| **OrderInvoiceNumber** | ✔︎        | Mã đơn hàng/hoá đơn (duy nhất)                                          |
| **OrderAmount**        | ✔︎        | Số tiền giao dịch                                                       |
| **Currency**           | ✔︎        | Đơn vị tiền tệ (VD: `"VND"`, `"USD"`)                                   |
| **OrderDescription**   |          | Mô tả đơn hàng                                                          |
| **OrderTaxAmount**     |          | Thuế đơn hàng, sử dụng `sepay.Float64()`                                |
| **CustomerID**         |          | Mã khách hàng (nếu có), sử dụng `sepay.String()`                        |
| **SuccessURL**         |          | URL callback khi thanh toán thành công, sử dụng `sepay.String()`        |
| **ErrorURL**           |          | URL callback khi xảy ra lỗi, sử dụng `sepay.String()`                   |
| **CancelURL**          |          | URL callback khi người dùng hủy thanh toán, sử dụng `sepay.String()`    |
| **CustomData**         |          | Dữ liệu tuỳ chỉnh (merchant tự định nghĩa), sử dụng `sepay.String()`    |

Dữ liệu trả về (`*SignedCheckoutFields`):

| Tham số                | Bắt buộc | Mô tả                                                              |
| ---------------------- | -------- | ------------------------------------------------------------------ |
| **Merchant**           | ✔︎        | Mã merchant                                                        |
| **Operation**          | ✔︎        | Loại giao dịch, hiện chỉ hỗ trợ: `"PURCHASE"`                      |
| **PaymentMethod**      | ✔︎        | Phương thức thanh toán: `"BANK_TRANSFER"`, `"NAPAS_BANK_TRANSFER"` |
| **OrderInvoiceNumber** | ✔︎        | Mã đơn hàng/hoá đơn (duy nhất)                                     |
| **OrderAmount**        | ✔︎        | Số tiền giao dịch                                                  |
| **Currency**           | ✔︎        | Đơn vị tiền tệ (VD: `"VND"`, `"USD"`)                              |
| **OrderDescription**   | ✔︎        | Mô tả đơn hàng                                                     |
| **OrderTaxAmount**     |          | Thuế đơn hàng                                                      |
| **CustomerID**         |          | Mã khách hàng (nếu có)                                             |
| **SuccessURL**         |          | URL callback khi thanh toán thành công                             |
| **ErrorURL**           |          | URL callback khi xảy ra lỗi                                        |
| **CancelURL**          |          | URL callback khi người dùng hủy thanh toán                         |
| **CustomData**         |          | Dữ liệu tuỳ chỉnh (merchant tự định nghĩa)                         |
| **Signature**          | ✔︎        | Chữ ký bảo mật (HMAC SHA256) để xác thực dữ liệu trả về            |

## API

SDK cung cấp các phương thức để gọi Open API cho cổng thanh toán SePay. Tất cả phương thức API đều nhận `context.Context` làm tham số đầu tiên.

### Tra cứu danh sách đơn hàng

```go
resp, err := client.Order.All(ctx, &sepay.OrderQueryParams{
	PerPage:       sepay.Int(10),
	Q:             sepay.String("keyword"),
	OrderStatus:   sepay.String("COMPLETED"),
	CreatedAt:     sepay.String("2024-01-01"),
	FromCreatedAt: sepay.String("2024-01-01"),
	ToCreatedAt:   sepay.String("2024-12-31"),
	CustomerID:    sepay.String("KH001"),
	SortCreatedAt: sepay.String("desc"),
})
```

### Xem chi tiết đơn hàng

```go
resp, err := client.Order.Retrieve(ctx, "DH0001")
```

### Hủy giao dịch đơn hàng (dành cho thanh toán bằng thẻ tín dụng)

```go
resp, err := client.Order.VoidTransaction(ctx, "DH0001")
```

### Hủy đơn hàng (dành cho thanh toán bằng quét mã QR)

```go
resp, err := client.Order.Cancel(ctx, "DH0001")
```

### Xử lý phản hồi

Tất cả phương thức API trả về `*sepay.Response`:

```go
resp, err := client.Order.Retrieve(ctx, "DH0001")
if err != nil {
	log.Fatal(err)
}

fmt.Println(resp.StatusCode) // HTTP status code
fmt.Println(resp.Header)     // HTTP headers
fmt.Println(string(resp.Body)) // Raw response body

// Giải mã JSON vào struct tuỳ chỉnh
var result map[string]any
if err := resp.DecodeJSON(&result); err != nil {
	log.Fatal(err)
}
```

## Giấy phép sử dụng

Thư viện sử dụng giấy phép MIT. Xem chi tiết [LICENSE](LICENSE).

## Hỗ trợ

- Email: <info@sepay.vn>
- Tài liệu: <https://docs.sepay.vn>
