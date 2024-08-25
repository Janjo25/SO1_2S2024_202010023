#include <linux/fs.h>
#include <linux/init.h>
#include <linux/jiffies.h>
#include <linux/kernel.h>
#include <linux/ktime.h>
#include <linux/mm.h>
#include <linux/module.h>
#include <linux/proc_fs.h>
#include <linux/sched.h>
#include <linux/sched/signal.h>
#include <linux/seq_file.h>
#include <linux/slab.h>
#include <linux/uaccess.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Luis César Lizama Quiñónez");
MODULE_DESCRIPTION("SO1_2S2024_PY1_202010023");
MODULE_VERSION("1.0");

#define PROC_NAME "sysinfo_202010023"
#define MAXIMUM_CONTAINERS 1200

pid_t pid_list[MAXIMUM_CONTAINERS]; // El tipo de dato 'pid_t' se utiliza para almacenar el ID de un proceso.
int pid_count = 0;                  // Este contador se utiliza para llevar la cuenta de cuántos procesos se han leído.

/* Se crea la función 'load_container_pids' que se encargará de leer los ID de los procesos y almacenarlos en 'pid_list'. */
void load_container_pids(void) {
    char *file_buffer;     // Este buffer es para almacenar temporalmente el contenido completo del archivo.
    char *pid_file_path = "containers_pid.txt";
    int bytes_read;
    loff_t position = 0;   // Esta variable indica la posición en la que se comenzará a leer el archivo.
    struct file *pid_file; // Este struct se utiliza para guardar el puntero que apunta al archivo que se va a leer.

    /* Se reinician 'pid_list' y 'pid_count' para asegurarse de que no haya basura en la memoria. */
    memset(pid_list, 0, sizeof(pid_list)); // Se utiliza 'memset' para llenar un bloque de memoria con un valor específico.
    pid_count = 0;

    /* Se asigna memoria para el buffer que se utilizará para leer el archivo.
     * Se utiliza 'kzalloc' para reservar memoria en el kernel e inicializarla, asegurando que está limpia y sea segura de usar.
     * Se utiliza 'GFP_KERNEL' para indicarle al kernel que puede esperar a que la memoria esté disponible, ya que no es una operación crítica. */
    file_buffer = kzalloc(1024, GFP_KERNEL);

    if (!file_buffer) {
        printk(KERN_ERR "Fallo al asignar memoria para la lectura de PIDs.\n");

        return;
    }

    pid_file = filp_open(pid_file_path, O_RDONLY, 0); // Se abre el archivo en modo de solo lectura.

    if (IS_ERR(pid_file)) {
        printk(KERN_ERR "Ocurrió un error al intentar abrir 'container_pids.txt'.\n");
        kfree(file_buffer);

        return;
    }

    bytes_read = kernel_read(pid_file, file_buffer, 1024 - 1, &position); // Se utiliza 'kernel_read' para leer el contenido del archivo.

    if (bytes_read < 0) {
        printk(KERN_ERR "Ocurrió un error al intentar leer el contenido de 'container_pids.txt'.\n");
        filp_close(pid_file, NULL);
        kfree(file_buffer);

        return;
    }

    file_buffer[bytes_read] = '\0'; // Se agrega un carácter nulo al final del buffer para indicar el final del contenido.
    filp_close(pid_file, NULL);     // Se cierra el archivo.

    /* Se utiliza 'strsep' para separar el contenido del buffer por saltos de línea.
     * Se utiliza 'simple_strtol' para convertir el contenido de una línea a un entero. */
    char *current_pos = file_buffer; // Se usará una variable temporal donde el puntero será modificado mientras se lee el contenido del buffer.
    char *current_line;

    while ((current_line = strsep(&current_pos, "\n")) != NULL) {
        if (strlen(current_line) > 0) {
            pid_t pid = simple_strtol(current_line, NULL, 10); // Convierte el contenido de la línea a un entero.
            printk(KERN_INFO "PIDs leído: %d\n", pid);

            if (pid_count < MAXIMUM_CONTAINERS && pid != 0) {
                pid_list[pid_count++] = pid;
            }
        }
    }

    kfree(file_buffer); // Se libera la memoria que se asignó para el buffer.
}

/* Se crea la función 'is_container_process' que se encargará de verificar si un proceso es un contenedor o no.
 * Se utiliza 'pid' para obtener el ID de un proceso y se retorna un booleano que indica si el proceso es un contenedor o no. */
bool is_container_process(pid_t pid) {
    int i;

    for (i = 0; i < pid_count; i++) {
        if (pid_list[i] == pid) return true;
    }

    return false;
}

/* Se crea la función 'sysinfo_show' que se encargará de guardar la información del sistema en el archivo '/proc/sysinfo_202010023'.
 * Se utiliza 'output_file' para escribir en el archivo '/proc/sysinfo_202010023' y 'unused' no se utiliza en esta función.
 * A pesar de que 'unused' no se utiliza, es necesario que esté presente en la función para que pueda ser llamada. */
static int sysinfo_show(struct seq_file *output_file, void *unused) {
    struct sysinfo system;

    si_meminfo(&system);

    unsigned long total_ram = system.totalram * 4 / 1024; // Está dado en KB.
    unsigned long free_ram = system.freeram * 4 / 1024;   // Está dado en KB.
    unsigned long used_ram = total_ram - free_ram;

    seq_printf(output_file, "{\n");
    seq_printf(output_file, "  \"total_ram\": %lu,\n", total_ram);
    seq_printf(output_file, "  \"free_ram\": %lu,\n", free_ram);
    seq_printf(output_file, "  \"used_ram\": %lu,\n", used_ram);
    seq_printf(output_file, "  \"processes\": [\n");

    struct task_struct *task;

    load_container_pids();
    int written_pid_count = 0;

    for_each_process(task) {
        if (task->flags & PF_KTHREAD) continue; // Esta línea se agrega para omitir algunos procesos innecesarios.

        if (is_container_process(task->pid)) {
            /* Se calcula el uso de memoria de un proceso.
             * Se utiliza 'get_mm_rss' para obtener la cantidad de memoria que está siendo utilizada por un proceso.
             * Se utiliza 'task->mm' para obtener la estructura de memoria de un proceso. */
            unsigned long rss = task->mm ? get_mm_rss(task->mm) * 4 : 0; // Está dado en KB.
            unsigned long memory_usage = (rss * 1000) / total_ram;

            if (memory_usage > 1000) memory_usage = 1000; // Se asegura de que el uso de memoria no sea mayor al 100%.

            /* Se calcula el uso de CPU de un proceso.
             * Se utiliza 'task->utime' para obtener el tiempo de CPU que ha pasado en modo de usuario.
             * Se utiliza 'task->stime' para obtener el tiempo de CPU que ha pasado en modo de kernel.
             * Se utiliza 'get_jiffies_64' para obtener el tiempo total que ha pasado desde que se inició el sistema. */
            unsigned long total_cpu_time = task->utime + task->stime;
            unsigned long total_time = jiffies_to_msecs(get_jiffies_64()); // El tiempo total está dado en milisegundos.
            unsigned long cpu_usage = (total_cpu_time * 1000) / total_time;

            if (cpu_usage > 1000) cpu_usage = 1000; // Se asegura de que el uso de CPU no sea mayor al 100%.

            seq_printf(output_file, "    {\n");
            seq_printf(output_file, "      \"pid\": %d,\n", task->pid);
            seq_printf(output_file, "      \"name\": \"%s\",\n", task->comm);
            seq_printf(output_file, "      \"cmd_line\": \"%s\",\n", task->comm);
            seq_printf(output_file, "      \"vsz\": %lu,\n", task->mm ? task->mm->total_vm * 4 : 0);
            seq_printf(output_file, "      \"rss\": %lu,\n", rss);
            seq_printf(output_file, "      \"memory_usage\": %lu.%lu,\n", memory_usage / 10, memory_usage % 10);
            seq_printf(output_file, "      \"cpu_usage\": %lu.%lu\n", cpu_usage / 10, cpu_usage % 10);

            if (written_pid_count != pid_count - 1) {
                seq_printf(output_file, "    },\n");
                written_pid_count++;
            } else {
                seq_printf(output_file, "    }\n");
            }
        }
    }

    seq_printf(output_file, "  ]\n");
    seq_printf(output_file, "}\n");

    return 0;
}

/* Se crea la función 'sysinfo_open' que se encargará de abrir el archivo '/proc/sysinfo_202010023'.
 * Se utiliza 'inode' para obtener información del archivo y 'file' para abrir el archivo.
 * La función retorna un entero que indica si el archivo se abrió correctamente o no. */
static int sysinfo_open(struct inode *inode, struct file *file) {
    return single_open(file, sysinfo_show, NULL);
}

/* Se crea la estructura 'sysinfo_ops' que contiene las operaciones que se pueden realizar con el archivo '/proc/sysinfo_202010023'. */
static const struct proc_ops sysinfo_ops = {
    .proc_lseek   = seq_lseek,      // Función para buscar dentro del archivo.
    .proc_open    = sysinfo_open,   // Función para abrir el archivo.
    .proc_read    = seq_read,       // Función para leer el archivo.
    .proc_release = single_release, // Función para cerrar el archivo.
};

static int __init sysinfo_init(void) {
    proc_create(PROC_NAME, 0, NULL, &sysinfo_ops);
    printk(KERN_INFO "módulo 'sysinfo' cargado.\n");

    return 0;
}

static void __exit sysinfo_exit(void) {
    remove_proc_entry(PROC_NAME, NULL);
    printk(KERN_INFO "módulo 'sysinfo' descargado.\n");
}

module_init(sysinfo_init);
module_exit(sysinfo_exit);
