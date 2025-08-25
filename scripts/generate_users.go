package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"strings"
)

func main() {
	// 创建文件并准备写入
	file, err := os.Create("scripts/generate_users.sql")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	// 使用 strings.Builder 来构建 SQL 插入语句
	var sb strings.Builder
	sb.WriteString("INSERT INTO `users` (`id`, `username`, `password`, `nickname`, `phone`, `email`, `avatar`, `client_ip`, `client_port`, `login_time`, `heartbeat_time`, `logout_time`, `status`, `online_status`, `device_info`) VALUES\n")

	// 从 ID 1001 开始生成用户
	for i := 1001; i <= 2000; i++ {
		username := fmt.Sprintf("user%d", i)
		password := "123456"
		// 加密密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Error generating password: %v", err)
		}

		// 构造 SQL 插入语句
		sql := fmt.Sprintf(
			"(%d, '%s', '%s', 'Nickname%d', '123-456-789%d', 'user%d@example.com', 'https://example.com/avatar%d.jpg', '192.168.0.%d', '555%d', 1635630000, 1645678900, 1656789000, 1, 0, 'Mozilla/5.0')",
			i, username, string(hashedPassword), i, i, i, i, i, i,
		)

		// 将每个插入语句添加到构建器中
		sb.WriteString(sql)
		if i < 2000 {
			sb.WriteString(",\n")
		}
	}

	// 将生成的 SQL 插入语句写入文件
	_, err = file.WriteString(sb.String())
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	fmt.Println("SQL 插入语句已成功写入 ./generate_users.sql 文件")
}
