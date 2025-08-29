## Parte 1: Introducción a Docker
En esta primera parte del trabajo práctico se plantean una serie de ejercicios que sirven para introducir las herramientas básicas de Docker que se utilizarán a lo largo de la materia. El entendimiento de las mismas será crucial para el desarrollo de los próximos TPs.

### Ejercicio N°1:
Definir un script de bash `generar-compose.sh` que permita crear una definición de Docker Compose con una cantidad configurable de clientes.  El nombre de los containers deberá seguir el formato propuesto: client1, client2, client3, etc. 

El script deberá ubicarse en la raíz del proyecto y recibirá por parámetro el nombre del archivo de salida y la cantidad de clientes esperados:

`./generar-compose.sh docker-compose-dev.yaml 5`

Considerar que en el contenido del script pueden invocar un subscript de Go o Python:

```
#!/bin/bash
echo "Nombre del archivo de salida: $1"
echo "Cantidad de clientes: $2"
python3 mi-generador.py $1 $2
```

En el archivo de Docker Compose de salida se pueden definir volúmenes, variables de entorno y redes con libertad, pero recordar actualizar este script cuando se modifiquen tales definiciones en los sucesivos ejercicios.

## Solucion

El primer paso de mi solución fue conservar el .yaml original, renombrándolo como `docker-compose-dev-original.yaml`. Lo mismo con el README.md (renombrado a `README-original.md`).
Mi decisión para este ejercicio fue la de invocar un subscript de Python, tal como es sugerido en el enunciado y por una cuestión de practicidad. Por eso, primero creé el script de bash tal como se lo muestra en el enunciado `generar-compose.sh`. Este script recibe los parámetros (nombre del archivo de salida y cantidad de clientes) y se encarga de ejecutar el subscript de Python con dichos parámetros. 

Nuestro `mi-generador.py` utiliza las librerías sys y yaml (sys para el chequeo de argumentos, y yaml para parsear facilmente el archivo de salida). El script es bastante simple, dado que la primera parte del .yaml no cambia. Toda esa porción del archivo es guardada en el diccionario `compose`. Y finalmente, para los clientes que son variables por parámetro, utilizamos un loop for para ir cargando a cada uno en el diccionario. Respetamos el formato propuesto (client1, client2, client3, etc), tal como lo pide el enunciado. Y finalmente, luego de tener listo el diccionario con nuestra definición de Docker Compose, lo guardamos facilmente gracias a la librería `yaml` utilizando el nombre de salida que fue recibido por parámetro.