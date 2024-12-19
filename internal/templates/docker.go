package templates

func GetDockerComposeTemplate() string {
	return `version: '3'
services:
  app:
    image: nginx:latest
    ports:
      - "80:80"
    volumes:
      - ./data:/data
    environment:
      - ENV=production
    networks:
      - app-network

networks:
  app-network:
    driver: bridge`
}
