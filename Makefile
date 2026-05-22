.PHONY: dev prod down logs

# Local için tüm servisleri kaldırır
dev:
	docker-compose up -d --build

# Production için sadece API'yi bulut konfigürasyonlarıyla kaldırır
prod:
	docker-compose -f docker-compose.prod.yml up -d --build

# Tüm çalışanları indirir
down:
	docker-compose down -v
	docker-compose -f docker-compose.prod.yml down -v

logs:
	docker-compose logs -f api

swag:
	swag init -g cmd/server/main.go --parseDependency --parseInternal
