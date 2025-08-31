## Parte 1: Introducción a Docker
En esta primera parte del trabajo práctico se plantean una serie de ejercicios que sirven para introducir las herramientas básicas de Docker que se utilizarán a lo largo de la materia. El entendimiento de las mismas será crucial para el desarrollo de los próximos TPs.

### Ejercicio N°2:

Modificar el cliente y el servidor para lograr que realizar cambios en el archivo de configuración no requiera reconstruír las imágenes de Docker para que los mismos sean efectivos. La configuración a través del archivo correspondiente (config.ini y config.yaml, dependiendo de la aplicación) debe ser inyectada en el container y persistida por fuera de la imagen (hint: docker volumes).

## Solucion

La idea de este ejercicio es que los archivos de configuración del cliente y del servidor se inyecten directamente en sus containers. Para eso, simplemente tenemos que modificar nuestro script de Python para incluir `bind mounts` (uno para el servidor, y uno para el cliente) en el Docker Compose. De esta manera, no requerimos rebuildear las imágenes cada vez que se cambian los archivos de configuración. Los cambios se reflejan inmediatamente en el contenedor.
Para el servidor:
```
"volumes": ["./server/config.ini:/config.ini"]
```
Para el/los cliente/s:
```
"volumes": "volumes": ["./client/config.yaml:/config.yaml"]
```