#include <linux/fs.h>
#include <linux/init.h>
#include <linux/kernel.h>
#include <linux/mm.h>
#include <linux/module.h>
#include <linux/proc_fs.h>
#include <linux/sched.h>
#include <linux/sched/signal.h>
#include <linux/seq_file.h>
#include <linux/slab.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Luis César Lizama Quiñónez");
MODULE_DESCRIPTION("SO1_2S2024_PY1_202010023");
MODULE_VERSION("1.0");

#define PROC_NAME "sysinfo_202010023"
#define MAXIMUM_CONTAINERS 1200

pid_t pid_list[MAXIMUM_CONTAINERS]; // El tipo de dato 'pid_t' se utiliza para almacenar el ID de un proceso.
int pid_count = 0;                  // Este contador se utiliza para llevar la cuenta de cuántos procesos se han leído.

/* Se crea la función 'is_container_process' que se encargará de verificar si un proceso es un contenedor o no.
 * Se utiliza 'pid' para obtener el ID de un proceso y se retorna un booleano que indica si el proceso es un contenedor o no. */
bool is_container_process(pid_t pid) {
    int i;

    for (i = 0; i < pid_count; i++) {
        if (pid_list[i] == pid) return true;
    }

    return false;
}

/* Se crea la estructura 'sysinfo_fops' que contiene las operaciones que se pueden realizar con el archivo '/proc/sysinfo_202010023'. */
static const struct file_operations sysinfo_fops = {
        .owner = THIS_MODULE,      // Evita que el módulo se descargue mientras las operaciones están activas.
        .open = sysinfo_open,      // Función para abrir el archivo.
        .read = seq_read,          // Función para leer el archivo.
        .llseek = seq_lseek,       // Función para buscar dentro del archivo.
        .release = single_release, // Función para cerrar el archivo.
};

/* Se crea la función 'sysinfo_open' que se encargará de abrir el archivo '/proc/sysinfo_202010023'.
 * Se utiliza 'inode' para obtener información del archivo y 'file' para abrir el archivo.
 * La función retorna un entero que indica si el archivo se abrió correctamente o no. */
static int sysinfo_open(struct inode *inode, struct file *file) {
    return single_open(file, sysinfo_show, NULL);
}

/* Se crea la función 'sysinfo_show' que se encargará de guardar la información del sistema en el archivo '/proc/sysinfo_202010023'.
 * Se utiliza 'output_file' para escribir en el archivo '/proc/sysinfo_202010023' y 'unused' no se utiliza en esta función.
 * A pesar de que 'unused' no se utiliza, es necesario que esté presente en la función para que pueda ser llamada. */
static int sysinfo_show(struct seq_file *output_file, void *unused) {
    struct sysinfo system;

    si_meminfo(&system);

    seq_printf(output_file, "{\n");
    seq_printf(output_file, "  \"total_ram\": %lu,\n", system.totalram * 4 / 1024);
    seq_printf(output_file, "  \"free_ram\": %lu,\n", system.freeram * 4 / 1024);
    seq_printf(output_file, "  \"used_ram\": %lu,\n", (system.totalram - system.freeram) * 4 / 1024);
    seq_printf(output_file, "  \"processes\": [\n");

    struct task_struct *task;

    load_container_pids();

    for_each_process(task) {
        if (task->flags & PF_KTHREAD) continue; // Esta línea se agrega para omitir algunos procesos innecesarios.

        if (is_container_process(task->pid)) {
            seq_printf(output_file, "    {\n");
            seq_printf(output_file, "      \"pid\": %d,\n", task->pid);
            seq_printf(output_file, "      \"name\": \"%s\",\n", task->comm);
            seq_printf(output_file, "      \"cmd_line\": \"%s\",\n", task->comm);
            seq_printf(output_file, "      \"vsz\": %lu,\n", task->mm ? task->mm->total_vm * 4 : 0);
            seq_printf(output_file, "      \"rss\": %lu,\n", task->mm ? get_mm_rss(task->mm) * 4 : 0);
            seq_printf(output_file, "      \"memory_usage\": %lu\n,", task->mm ? task->mm->total_vm * 4 : 0);
            seq_printf(output_file, "      \"cpu_usage\": %lu\n", task->se.exec_start);
            seq_printf(output_file, "    },\n");
        }
    }

    seq_printf(output_file, "  ]\n");
    seq_printf(output_file, "}\n");

    return 0;
}

/* Se crea la función 'load_container_pids' que se encargará de leer los ID de los procesos y almacenarlos en 'pid_list'. */
void load_container_pids(void) {
    char pid_buffer[16];   // Este buffer es para almacenar el ID de un proceso.
    mm_segment_t old_fs;   // Un 'checkpoint' para restaurar el espacio de memoria del kernel luego de leer el archivo.
    pid_t pid;
    struct file *pid_file; // Este struct se utiliza para guardar el puntero que apunta al archivo que se va a leer.

    char *pid_file_path = "containers_pid.txt";

    old_fs = get_fs(); // Se guarda el espacio de memoria actual del kernel.
    set_fs(KERNEL_DS); // Se cambia el espacio de memoria. Esta pasa del kernel al usuario para poder leer el archivo.

    pid_file = filp_open(pid_file_path, O_RDONLY, 0); // Se abre el archivo en modo de solo lectura.

    if (!IS_ERR(pid_file)) {
        /* 'pid_file' es el puntero que apunta al archivo que se va a leer.
         * 'pid_buffer' es el buffer donde se almacenará el ID de un proceso.
         * 'sizeof(pid_buffer)' es el tamaño del buffer.
         * '&pid_file->f_pos' es la posición actual del archivo. */
        while (kernel_read(pid_file, pid_buffer, sizeof(pid_buffer), &pid_file->f_pos) > 0) {
            pid = simple_strtol(pid_buffer, NULL, 10); // Se convierte el ID de un proceso de cadena a entero.

            if (pid_count < MAX_CONTAINERS) {
                pid_list[pid_count++] = pid;
            }
        }

        filp_close(pid_file, NULL);
    }

    set_fs(old_fs); // Se restaura el espacio de memoria del kernel.
}

static int __init sysinfo_init(void) {
    proc_create(PROC_NAME, 0, NULL, &sysinfo_fops);
    printk(KERN_INFO "módulo 'sysinfo' cargado.\n");

    return 0;
}

static void __exit sysinfo_exit(void) {
    remove_proc_entry(PROC_NAME, NULL);
    printk(KERN_INFO "módulo 'sysinfo' descargado.\n");
}

module_init(sysinfo_init);
module_exit(sysinfo_exit);
