### Thử làm 1 cái trình compile code Javascript super dummy simple

Minh họa kiến trúc app:

![Javascript compiler](img/grpc-compiler.png?raw=true "Javascript compiler")

App sử dụng tính năng server stream của gRPC, gồm các thành phần:

- gRPC web client chạy ở localhost:8087

- gRPC server chạy ở localhost:9001

Luồng hoạt động:

1. Bật server, sau đó bật client

2. Truy cập localhost:8087. Lúc này trình duyệt sẽ thiết lập ***kết nối socket*** đến gRPC client

3. Nhập code rồi bấm Biên dịch. Lúc này trình duyệt gửi code đến gRPC client

4. gRPC client nhận đc code, sẽ gửi cho gRPC server trên kết nối gRPC.

5. gRPC server nhận được code sẽ ***chạy code trên 1 goroutine riêng***. Code được chạy trong 1 NodeJS container, ***kêt quả được gửi vào 1 channel***

6. gRPC server lấy kết quả từ channel rồi ***stream*** ngược lại về cho gRPC client

7. gRPC client gửi kết quả về cho trình duyệt trên kết nối socket