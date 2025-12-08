# Makefile for deployment convenience

.PHONY: help deploy stop start restart logs health status clean

help:
	@echo "Available commands:"
	@echo "  make deploy    - Deploy the application"
	@echo "  make stop      - Stop all services"
	@echo "  make start     - Start all services"
	@echo "  make restart   - Restart all services"
	@echo "  make logs      - View logs (use LOGS=service for specific service)"
	@echo "  make health    - Run health check"
	@echo "  make status    - Show service status"
	@echo "  make clean     - Clean up containers and volumes"

deploy:
	@./deploy.sh

stop:
	@docker compose -f docker-compose.prod.yml down

start:
	@docker compose -f docker-compose.prod.yml up -d

restart:
	@docker compose -f docker-compose.prod.yml restart

logs:
	@if [ -z "$(LOGS)" ]; then \
		docker compose -f docker-compose.prod.yml logs -f; \
	else \
		docker compose -f docker-compose.prod.yml logs -f $(LOGS); \
	fi

health:
	@./health-check.sh

status:
	@docker compose -f docker-compose.prod.yml ps

clean:
	@docker compose -f docker-compose.prod.yml down -v
	@docker system prune -f

