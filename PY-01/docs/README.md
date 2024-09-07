# Proyecto 1: Gestor de Contenedores

## Introducción

El propósito de este proyecto es la creación de un **gestor de contenedores** que automatice la creación, eliminación y
monitoreo de contenedores utilizando tecnologías como **Docker**, **Rust**, y módulos del **Kernel de Linux**. A través
de este gestor, se puede observar y gestionar de manera eficiente los recursos de los contenedores en un entorno de
Linux. El proyecto incluye la captura de métricas del sistema mediante un módulo de kernel, un servicio en **Rust** que
coordina la operación de los contenedores, y un servidor en **Python** que administra los registros generados y produce
gráficas para la visualización de datos.

## 1. Modo de Uso

El gestor se puede utilizar para administrar automáticamente contenedores de alto y bajo consumo. A través de un
CronJob, se generan contenedores aleatorios que se gestionan según el uso de memoria y CPU. El servicio en Rust captura
métricas del sistema a través del módulo de kernel, analiza los procesos de los contenedores y selecciona los dos
contenedores que más recursos consumen y los tres que menos recursos utilizan. Estos contenedores son gestionados para
mantener un balance, eliminando los que no cumplen con estas características. Finalmente, el sistema envía registros al
servidor Python, el cual los almacena y genera gráficos con las métricas.

Pasos principales para su uso:

1. **Construir las imágenes de Docker**:
   Antes de poder crear y gestionar los contenedores, es necesario construir las imágenes de Docker que se utilizarán.
   Esto se realiza con el archivo `Dockerfile`:

   ```bash
   docker build -t <nombre-imagen> .
   ```

2. **Compilar el módulo del Kernel**:
   El siguiente paso es compilar el módulo del kernel en C, el cual se encargará de capturar métricas del sistema:

   ```bash
   make
   sudo insmod sysinfo.ko
   ```

3. **Compilar el exportador en C**:
   El exportador en C se encarga de leer el archivo en `/proc` y convertir los datos de métricas en un formato **JSON**
   válido para que el servicio en Rust los pueda procesar:

   ```bash
   gcc exporter.c -o exporter
   ```

4. **Levantar el servidor en Python**:
   Es necesario levantar el servidor en Python, el cual estará esperando peticiones HTTP para recibir y almacenar
   registros, y también generar las gráficas:

   ```bash
   docker run -d -p <puerto_anfitrión>:<puerto_contenedor> -v <directorio-anfitrión>:/app <nombre-imagen>
   ```

5. **Ejecutar el servicio en Rust**:
   Una vez que el servidor está listo, se puede ejecutar el servicio en Rust. Este servicio creará automáticamente el
   CronJob que gestionará los contenedores. El servicio en Rust ordena los contenedores por consumo de recursos y
   elimina aquellos que no cumplen con los criterios establecidos (2 de alto consumo, 3 de bajo consumo). También envía
   los registros al servidor Python mediante peticiones HTTP.

   ```bash
   cargo run
   ```

6. **Finalización del servicio**:
   Cuando se detiene el servicio con `CTRL+C`, se envía una última petición HTTP al servidor para generar las gráficas
   finales. Además, se eliminan el CronJob y los archivos temporales.

## 2. Instalación

Los siguientes pasos detallan cómo instalar y preparar el sistema:

1. **Requisitos del sistema**:
    - Compilador de C para módulos del kernel.
    - Docker instalado.
    - Rust (con Cargo).
    - Sistema operativo Linux.

2. **Instalación de dependencias**:
    - Instalar `base-devel` (GCC, Make y más herramientas de desarrollo):
      ```bash
       sudo pacman -S base-devel
      ```
    - Instalar Docker:
      ```bash
      sudo pacman -S docker
      sudo systemctl start docker
      sudo systemctl enable docker
      ```
    - Instalar las cabeceras del kernel:
      ```bash
      sudo pacman -S linux<versión-kernel>-headers
      ```
    - Instalar Rust y Cargo:
      ```bash
      curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
      ```

## 3. Ejemplos con Explicación

A continuación se muestra un ejemplo de la ejecución del servicio en Rust, incluyendo la captura de métricas del sistema
y la gestión de contenedores.

1. **Ejecución del servicio**:
   Al ejecutar el servicio, se utiliza el siguiente comando para iniciarlo:

   ```bash
   cargo run
   ```

2. **Gestión de contenedores**:
   Cada 10 segundos, el servicio captura información del sistema, y elimina los contenedores que no cumplen con los
   criterios de alto y bajo consumo de recursos.

   Ejemplo de captura de contenedores:

   ```bash
   Se ha guardado la información del sistema en 'sysinfo.json'.
   
   Total de RAM: 11757 KB
   RAM libre: 1198 KB
   RAM usada: 10559 KB
   Procesos de alto consumo:
   ("N2Cnmb5c", 67073)
   ("2JYt69n3", 67498)
   Procesos de bajo consumo:
   ("1Uuwt0JP", 66913)
   ("9e61cMe6", 66992)
   ("MoMrZl8r", 67157)
   
   Eliminando contenedor con nombre 'raBY9dSO'
   Eliminando contenedor con nombre 'V9hAynD8'
   Eliminando contenedor con nombre 'EHgpJPqC'
   Eliminando contenedor con nombre '5U1UOIzH'
   Eliminando contenedor con nombre 'GJXWzHJq'
   ```

3. **Actualización de métricas**:
   A medida que el servicio continúa, las métricas se actualizan cada 10 segundos, y los procesos de alto y bajo consumo
   son identificados y gestionados:

   ```bash
   Esperando 10 segundos...
   Se ha guardado la información del sistema en 'sysinfo.json'.
   
   Total de RAM: 11757 KB
   RAM libre: 4702 KB
   RAM usada: 7055 KB
   Procesos de alto consumo:
   ("Nr82s4cQ", 68204)
   ("prFZA0gr", 68449)
   Procesos de bajo consumo:
   ("MoMrZl8r", 67157)
   ("9e61cMe6", 66992)
   ("N2Cnmb5c", 67073)
   ```

4. **Finalización del servicio**:
   Al finalizar la ejecución, el servicio puede ser detenido con `CTRL+C`, momento en el que se realizan las últimas
   eliminaciones de contenedores y se termina el proceso. Además, se genera una última petición HTTP para la creación de
   gráficas, se eliminan los archivos temporales y se desactiva el CronJob.

   ```bash
   ^C
   Saliendo del programa...
   Process finished with exit code 0
   ```
