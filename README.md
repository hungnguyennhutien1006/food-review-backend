# Food Review Backend API

Backend service cho hệ thống đánh giá món ăn và nhà hàng.

## Công nghệ sử dụng
- Go 1.21+
- Gorilla Mux Router
- In-memory storage

## Cài đặt
```bash
# Clone repository
git clone <your-repo-url>
cd food-review-backend

# Cài đặt dependencies
go mod download

# Chạy server
go run .
```

## API Endpoints

### Users
- `GET /users` - Lấy tất cả users
- `GET /users/{id}` - Lấy user theo ID
- `POST /users` - Tạo user mới

### Categories
- `GET /categories` - Lấy tất cả categories
- `POST /categories` - Tạo category mới

### Restaurants
- `GET /restaurants` - Lấy tất cả restaurants
- `POST /restaurants` - Tạo restaurant mới

### Foods
- `GET /foods` - Lấy tất cả foods
- `GET /restaurants/{restaurant_id}/foods` - Lấy foods theo restaurant
- `POST /foods` - Tạo food mới

### Comments
- `GET /foods/{food_id}/comments` - Lấy comments theo food
- `GET /users/{user_id}/comments` - Lấy comments theo user
- `POST /comments` - Tạo comment mới

## Cấu trúc project
food-review-backend/
├── main.go          # Entry point, định nghĩa routes
├── storage.go       # Quản lý in-memory storage
├── user.go          # Module User
├── category.go      # Module Category
├── restaurant.go    # Module Restaurant
├── food.go          # Module Food
├── comment.go       # Module Comment
├── go.mod           # Go module dependencies
└── README.md        # Documentation

## Chạy server
```bash
go run .
```

Server sẽ chạy tại `http://localhost:8080`

## Test API

Import Postman collection từ file `food-review-backend.json` để test API.
