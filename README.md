## Parte 1: Introducción a Docker
En esta primera parte del trabajo práctico se plantean una serie de ejercicios que sirven para introducir las herramientas básicas de Docker que se utilizarán a lo largo de la materia. El entendimiento de las mismas será crucial para el desarrollo de los próximos TPs.

### Ejercicio N°3:
Crear un script de bash `validar-echo-server.sh` que permita verificar el correcto funcionamiento del servidor utilizando el comando `netcat` para interactuar con el mismo. Dado que el servidor es un echo server, se debe enviar un mensaje al servidor y esperar recibir el mismo mensaje enviado.

En caso de que la validación sea exitosa imprimir: `action: test_echo_server | result: success`, de lo contrario imprimir:`action: test_echo_server | result: fail`.

El script deberá ubicarse en la raíz del proyecto. Netcat no debe ser instalado en la máquina _host_ y no se pueden exponer puertos del servidor para realizar la comunicación (hint: `docker network`). `

## Solucion

La solución propuesta levanta un contenedor temporal de `busybox` en la misma red de Docker que el servidor (tp0_testing_net). Esto es debido a que incluye `nc` y por lo tanto no requerimos de netcat por separado, además de que permite una conexión interna en Docker para no exponer los puertos del servidor. Por medio de un pipe y con `nc`, enviamos un mensaje al servidor (utilizando su direccion y puerto del archivo de configuracion) y esperamos que la respuesta sea la misma. Esto es validado en la parte final del script.

Ejecución:
```
./validar-echo-server.sh
```