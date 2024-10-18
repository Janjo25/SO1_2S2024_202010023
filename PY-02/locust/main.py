import random

import locust


class FacultyUser(locust.HttpUser):
    """
    Clase que define el comportamiento de los usuarios que simularán el tráfico:
        HttpUser: Clase que permite simular el comportamiento de un usuario que realiza peticiones HTTP.
        task: Decorador que permite definir una tarea que realizará el usuario.
        between: Función que permite definir un rango de tiempo en segundos de manera aleatoria.
    """
    host = random.choice(["http://agronomy.local", "http://engineering.local"])
    wait_time = locust.between(1, 3)  # Los usuarios esperarán entre 1 y 3 segundos entre tareas

    @locust.task
    def send_traffic(self):
        faculties = ['Ingeniería', 'Agronomía']
        disciplines = [1, 2, 3]  # 1: Natación, 2: Atletismo, 3: Boxeo

        # Crea un payload aleatorio para enviar en la petición POST.
        payload = {
            'name': 'Estudiante-' + str(random.randint(1, 10000)),
            'age': random.randint(18, 25),
            'faculty': random.choice(faculties),
            'discipline': random.choice(disciplines),
        }

        print(f"Enviando tráfico: {payload}")

        # Envia una petición POST al endpoint del Ingress de Kubernetes.
        self.client.post('/', json=payload)
