#include <linux/init.h>
#include <linux/kernel.h>
#include <linux/mm.h>
#include <linux/module.h>
#include <linux/sched.h>
#include <linux/sched/signal.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Luis César Lizama Quiñónez");
MODULE_DESCRIPTION("Tarea 1-Sistemas Operativos 1");
MODULE_VERSION("1.0");

/* Se usa '__init' para liberar memoria después de la inicialización.
 * Se usa '__exit' para liberar memoria antes de la eliminación. */
static int __init main_init(void) {
    /* Los tipos de datos 'sysinfo', 'task_struct' y 'list_head' son estructuras definidas en el kernel de Linux.
     * El tipo de dato 'sysinfo' se utiliza para obtener información del sistema.
     * El tipo de dato 'task_struct' para obtener información de los procesos.
     * El tipo de dato 'list_head' para recorrer la lista de procesos. */
    struct sysinfo si;

    /* Se usan punteros para no copiar toda la estructura en memoria, ya que son muy grandes y complejas. */
    struct task_struct *task;
    struct task_struct *child;
    struct list_head *list;

    si_meminfo(&si); // Se obtiene la información del sistema y se almacena en la dirección de memoria de 'si'.
    printk(KERN_INFO "Total RAM: %lu MB\n", si.totalram * 4 / 1024 );
    printk(KERN_INFO "Free RAM: %lu MB\n", si.freeram * 4 / 1024 );

    /* El bucle 'for_each_process' es un bucle especial que recorre la lista de procesos del sistema.
     * A este bucle se le pasa un puntero a una estructura 'task_struct' que se utiliza para almacenar la información de cada proceso. */
    for_each_process(task)
    {
        printk(KERN_INFO "Padre: %s [%d]\n", task->comm, task->pid); // De la estructura 'task_struct' se obtiene el nombre y el PID del proceso.

        /* El bucle 'list_for_each' recorre la lista de procesos hijos de un proceso.
         * A este bucle se le pasa un puntero a una estructura 'list_head' que se utiliza para almacenar la información de cada proceso hijo.
         * Se apunta a la dirección de memoria del proceso padre donde comienza la lista de procesos hijos. */
        list_for_each(list, &task->children)
        {
            /* El macro 'list_entry' se usa para obtener un puntero a la estructura 'task_struct' a partir de un puntero a la estructura 'list_head'. */
            child = list_entry(list,struct task_struct, sibling);
            printk(KERN_INFO " Hijo: %s [%d]\n", child->comm, child->pid );
        }
    }

    return 0;
}

static void __exit main_exit(void) {
    printk(KERN_INFO "Módulo removido.\n" );
}

module_init(main_init);
module_exit(main_exit);
