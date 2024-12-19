package templates

func GetGoConfigTemplate() string {
	return `# 服务配置
server:
  port: 8080
  host: localhost
  mode: development

# 数据库配置
database:
  driver: mysql
  host: localhost
  port: 3306
  username: root
  password: password
  dbname: myapp

# 日志配置
log:
  level: info
  path: ./logs
  filename: app.log
  max_size: 100
  max_backups: 3
  max_age: 28
  compress: true`
}
