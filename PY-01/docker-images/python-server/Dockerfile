FROM python:3.10-slim

# Se crea y establece el directorio '/app' como directorio de trabajo.
WORKDIR /app

# Se copian los archivos 'main.py' y 'requirements.txt' al directorio '/app'.
COPY . /app

# Se instalan los paquetes necesarios especificados en 'requirements.txt'.
RUN pip install --no-cache-dir -r requirements.txt

# Se indica que el puerto 80 estará disponible para el mundo fuera de este contenedor.
EXPOSE 80

ENV NAME PythonServer

# Se usa 'uvicorn' para levantar la aplicación 'FastAPI' y hacerla disponible para recibir solicitudes HTTP.
# Con 'main:app' se indica que el módulo 'main' contiene la instancia de FastAPI.
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "80"]
