import json
import os

from fastapi import FastAPI, Body
from pydantic import BaseModel

app = FastAPI()


# Se usa 'BaseModel' para definir la estructura de los datos que se reciben en el endpoint '/log'
class Process(BaseModel):
    pid: int
    name: str
    cmd_line: str
    vsz: int
    rss: int
    memory_usage: float
    cpu_usage: float


def append_to_json(new_data: dict):
    path = "./logs.json"

    if not os.path.exists(path):
        with open(path, "w") as file:
            json.dump([new_data], file, indent=4)
    else:
        with open(path, "r+") as file:
            # Se cargan los datos existentes en el archivo JSON, se añade el nuevo dato y se vuelven a escribir.
            existing_data = json.load(file)
            existing_data.append(new_data)
            file.seek(0)
            json.dump(existing_data, file, indent=4)


@app.post("/log")
def log_process(timestamp: str = Body(...), process: Process = Body(...)):
    """Se coloca 'Body(...)' en los parámetros, ya que se enviaron dos datos en el cuerpo de la petición.
    Es necesario descomponerlos en dos parámetros para poder trabajar con ellos."""
    log_entry = {"timestamp": timestamp, "process": process.model_dump()}

    append_to_json(log_entry)

    return {"estado": "éxito"}
