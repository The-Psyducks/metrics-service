# Build stage
FROM python:3.11-alpine AS builder

WORKDIR /home/app

# Copiar dependencias
COPY /server/requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

# Copiar el código fuente
COPY /server ./

# Test stage (opcional)
# FROM builder AS python-test-stage
# RUN pip install pytest
# CMD ["pytest", "--cov=.", "tests/"]

# Run stage
FROM python:3.11-alpine

WORKDIR /home/app

# Copiar la instalación y código desde el builder
COPY --from=builder /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages
COPY --from=builder /usr/local/bin /usr/local/bin
COPY --from=builder /home/app ./

# Comando de entrada para ejecutar el servidor
ENTRYPOINT ["python", "main.py"]
