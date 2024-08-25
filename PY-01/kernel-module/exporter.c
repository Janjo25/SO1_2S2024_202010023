#include <stdio.h>
#include <stdlib.h>

int main() {
    FILE *proc_file, *output_file; // Se crean dos punteros a 'FILE' para manejar los archivos.
    char buffer[1024]; // Se crea un buffer de 1024 bytes para almacenar el contenido de '/proc/sysinfo_202010023'.

    proc_file = fopen("/proc/sysinfo_202010023", "r"); // Se abre el archivo en modo de solo lectura.

    if (proc_file == NULL) {
        perror("Ocurrió un error al intentar abrir '/proc/sysinfo_202010023'");

        return 1;
    }

    output_file = fopen("sysinfo.txt", "w"); // Se abre el archivo en modo de escritura.

    if (output_file == NULL) {
        perror("Ocurrió un error al intentar abrir 'sysinfo.json'");
        fclose(proc_file);

        return 1;
    }

    /* Se procede a copiar el contenido de '/proc/sysinfo_202010023' al archivo 'sysinfo.json'. */
    while (fgets(buffer, sizeof(buffer), proc_file) != NULL) {
        fputs(buffer, output_file);
    }

    /* Se cierran los archivos. */
    fclose(proc_file);
    fclose(output_file);

    printf("Se ha guardado la información del sistema en 'sysinfo.json'.\n");

    return 0;
}
