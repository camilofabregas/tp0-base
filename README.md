## Parte 1: Introducción a Docker
En esta primera parte del trabajo práctico se plantean una serie de ejercicios que sirven para introducir las herramientas básicas de Docker que se utilizarán a lo largo de la materia. El entendimiento de las mismas será crucial para el desarrollo de los próximos TPs.

### Ejercicio N°4:
Modificar servidor y cliente para que ambos sistemas terminen de forma _graceful_ al recibir la signal SIGTERM. Terminar la aplicación de forma _graceful_ implica que todos los _file descriptors_ (entre los que se encuentran archivos, sockets, threads y procesos) deben cerrarse correctamente antes que el thread de la aplicación principal muera. Loguear mensajes en el cierre de cada recurso (hint: Verificar que hace el flag `-t` utilizado en el comando `docker compose down`).

## Solucion

El objetivo de este ejercicio es iniciar una secuencia de cierre controlada (o graceful shutdown) del proceso, en lugar de dejar que lo haga el SO. Para eso, lo que hacemos es interceptar la señal SIGTERM (enviada por Docker) tanto en el cliente como en el servidor. 

En el cliente, creamos un canal por donde el SO enviará las señales, y asociamos ese canal con la señal que nos interesa (en este caso SIGTERM)
```
sigc := make(chan os.Signal, 1)
signal.Notify(sigc, syscall.SIGTERM)
```

Paralelamente y por medio de un Goroutine (hilo ligero), escuchamos el channel hasta recibir una señal SIGTERM. En ese momento, cerramos la conexion con el cliente y forzmos la salida del programa.
```
if c.conn != nil {
        c.conn.Close()
}
os.Exit(0)
```

En el servidor, nuestra estrategia es utilizar un booleano `shutdown` para indicar cuando queremos realizar el graceful shutdown. Mientras sea False, continua la ejecucion: se aceptan nuevas conexiones de clientes y se handlean sus pedidos. Por eso, ahora en vez de `while True`, tenemos `while not self._shutdown`.

La ejecucion sera interrumpida cuando recibamos la señal, y esa señal la capturamos por medio de:

```
signal.signal(signal.SIGTERM, self.__handle_sigterm)
```
donde basicamente le decimos al SO que, cuando se reciba SIGTERM, llamemos a la funcion `handle_sigterm`. Allí finalmente seteamos nuestro booleano shutdown en True para no volver a iterar el while, y cerramos el socket del servidor.

```
def __handle_sigterm(self, signum, frame):
    self._shutdown = True
    self._server_socket.close()
```