import collections
import json
import os

from fastapi import FastAPI, Body
from matplotlib import pyplot as plt
from pydantic import BaseModel  # Se usa para definir la estructura de los datos que se reciben en el endpoint.

app = FastAPI()


class Memory(BaseModel):
    total_ram: int
    free_ram: int
    used_ram: int


class Process(BaseModel):
    pid: int
    name: str
    cmd_line: str
    vsz: int
    rss: int
    memory_usage: float
    cpu_usage: float


@app.get("/generate-memory-graph")
def generate_memory_graph():
    path = "./memory_logs.json"

    if not os.path.exists(path):
        return {"estado": "no existen registros"}

    with open(path, "r") as file:
        logs = json.load(file)

    timestamps = [log['timestamp'] for log in logs]
    used_memory = [log['memory']['used_ram'] for log in logs]

    plt.figure()  # Se crea una instancia de la gráfica. Esto inicia una figura en blanco.
    plt.plot(timestamps, used_memory, marker='o')  # Se crea un gráfico de líneas con los datos obtenidos.
    plt.title("Uso de Memoria")
    plt.xlabel("Tiempo")
    plt.ylabel("Memoria Usada (KB)")
    plt.xticks(rotation=45)  # Se rota el texto del eje X para que no se superpongan los datos.
    plt.tight_layout()  # Se ajusta el tamaño de la gráfica para que no se recorten los elementos.
    plt.savefig("./memory_graph.png")

    return {"estado": "gráfico generado", "ruta": "./memory_graph.png"}


@app.get("/generate-processes-graph")
def generate_processes_graph():
    path = "./processes_logs.json"

    if not os.path.exists(path):
        return {"estado": "no existen registros"}

    with open(path, "r") as file:
        logs = json.load(file)

    cmd_lines = [log['process']['cmd_line'] for log in logs]
    cmd_count = collections.Counter(cmd_lines)
    cmd_labels = list(cmd_count.keys())
    cmd_frequencies = list(cmd_count.values())

    plt.figure()  # Se crea una instancia de la gráfica. Esto inicia una figura en blanco.
    plt.bar(cmd_labels, cmd_frequencies)  # Se crea un gráfico de barras con los datos obtenidos.
    plt.title("Frecuencia de Procesos Eliminados")
    plt.xlabel("Tipo de Proceso")
    plt.ylabel("Frecuencia")
    plt.tight_layout()  # Se ajusta el tamaño de la gráfica para que no se recorten los elementos.
    plt.savefig("./processes_graph.png")

    return {"estado": "gráfico generado", "ruta": "./processes_graph.png"}


def append_to_json(flag, new_data: dict):
    path = ""

    if flag == "memory":
        path = "./memory_logs.json"
    elif flag == "process":
        path = "./processes_logs.json"

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


@app.post("/memory-log")
def log_memory(timestamp: str = Body(...), memory: Memory = Body(...)):
    """Se coloca 'Body(...)' en los parámetros, ya que se enviaron dos datos en el cuerpo de la petición.
    Es necesario descomponerlos en dos parámetros para poder trabajar con ellos."""
    log_entry = {"timestamp": timestamp, "memory": memory.model_dump()}

    append_to_json("memory", log_entry)

    return {"estado": "éxito"}


@app.post("/process-log")
def log_process(timestamp: str = Body(...), process: Process = Body(...)):
    """Se coloca 'Body(...)' en los parámetros, ya que se enviaron dos datos en el cuerpo de la petición.
    Es necesario descomponerlos en dos parámetros para poder trabajar con ellos."""
    log_entry = {"timestamp": timestamp, "process": process.model_dump()}

    append_to_json("process", log_entry)

    return {"estado": "éxito"}
