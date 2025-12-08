# ====================================================================
# STAGE 1: Build Stage
# ====================================================================
# 1. ЗМІНА БАЗОВОГО ОБРАЗУ: Використовуємо образ, який містить Bun
FROM oven/bun:alpine AS builder

WORKDIR /app

# 2. КОПІЮВАННЯ ЗАЛЕЖНОСТЕЙ: Копіюємо package.json та bun.lockb
# Примітка: Припускаємо, що Dockerfile знаходиться у корені, а файли у frontend/
COPY frontend/package.json ./
COPY frontend/bun.lock ./

# 3. Встановлюємо залежності за допомогою Bun
RUN bun install --frozen-lockfile

# 4. КОПІЮВАННЯ КОДУ
COPY frontend/ .

# 5. Встановлення змінної середовища
ARG VITE_API_BASE_URL
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL

# 6. Збірка програми
RUN bun run build
# ====================================================================
# STAGE 2: Production stage (Залишається незмінним)
# ====================================================================
FROM nginx:alpine

# 1. Встановлення wget для Health Check (ОК)
RUN apk add --no-cache wget

# 2. Копіюємо вбудовані файли
COPY --from=builder /app/dist /usr/share/nginx/html

# 3. Копіюємо конфігурацію Nginx
# Тут є потенційна проблема: якщо вихідний код був скопійований як COPY frontend/ . в builder,
# то файл frontend/nginx.conf може бути відсутній.
# ВИПРАВЛЕННЯ: якщо nginx.conf знаходиться в локальній папці frontend, то:
COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf

# 4. Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost/health || exit 1

# 5. Start nginx
CMD ["nginx", "-g", "daemon off;"]
