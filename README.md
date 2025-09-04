## Parte 3: Repaso de Concurrencia
En este ejercicio es importante considerar los mecanismos de sincronización a utilizar para el correcto funcionamiento de la persistencia.

### Ejercicio N°8:

Modificar el servidor para que permita aceptar conexiones y procesar mensajes en paralelo. En caso de que el alumno implemente el servidor en Python utilizando _multithreading_,  deberán tenerse en cuenta las [limitaciones propias del lenguaje](https://wiki.python.org/moin/GlobalInterpreterLock).

Básicamente el servidor utiliza `procesos` para atender a cada cliente. Cada vez que se abre una nueva conexión con un cliente, 
se crea un nuevo proceso que se va a encargar de atenderla. Así, cada cliente se atiende en un proceso distinto de manera tal 
que es un procesamiento en **paralelo**.

Inicialmente comencé haciendo una implementación que usaba una `ThreadPool`, pero luego cambié a utilizar procesos ya que no 
sufren del problema del `GIL`. De cualquier manera, como el servidor realiza un procesamiento muy intensivo de IO pero no tanto 
de CPU, no sería un gran problema si se utilizan threads en lugar de procesos.

Para sincronizar el acceso a las funciones que no son thread-safe para la persistencia de la información de las apuestas, 
se utilizan `locks` para que sólo se pueda acceder de a un proceso. Así, evitamos cualquier tipo de problema.

Para mantener información compartida entre procesos, en nuestro caso el diccionario con clientes que ya finalizaron, se utiliza 
un `Manager` de la biblioteca de multiprocessing. Este `Manager` nos permite compartir el diccionario de forma segura entre procesos, 
mediante el mecanismo de serialización de `pickle` que ofrece Python. Así, podemos tratar al diccionario como memoria compartida, 
casi sin enterarnos que se está compartiendo entre procesos.
