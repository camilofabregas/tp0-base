## Parte 2: Repaso de Comunicaciones

Las secciones de repaso del trabajo práctico plantean un caso de uso denominado **Lotería Nacional**. Para la resolución de las mismas deberá utilizarse como base el código fuente provisto en la primera parte, con las modificaciones agregadas en el ejercicio 4.

### Ejercicio N°7:

Modificar los clientes para que notifiquen al servidor al finalizar con el envío de todas las apuestas y así proceder con el sorteo.
Inmediatamente después de la notificacion, los clientes consultarán la lista de ganadores del sorteo correspondientes a su agencia.
Una vez el cliente obtenga los resultados, deberá imprimir por log: `action: consulta_ganadores | result: success | cant_ganadores: ${CANT}`.

El servidor deberá esperar la notificación de las 5 agencias para considerar que se realizó el sorteo e imprimir por log: `action: sorteo | result: success`.
Luego de este evento, podrá verificar cada apuesta con las funciones `load_bets(...)` y `has_won(...)` y retornar los DNI de los ganadores de la agencia en cuestión. Antes del sorteo no se podrán responder consultas por la lista de ganadores con información parcial.

Las funciones `load_bets(...)` y `has_won(...)` son provistas por la cátedra y no podrán ser modificadas por el alumno.

No es correcto realizar un broadcast de todos los ganadores hacia todas las agencias, se espera que se informen los DNIs ganadores que correspondan a cada una de ellas.

## Solucion

Los clientes, luego de enviar todos sus batch de apuestas, envían un mensaje de `Finish` al servidor, que consiste 
únicamente en enviar un mensaje que contenga el uint32 que utilizamos antes para indicar la cantidad de bytes a leer, 
pero con valor `0`. Así, el servidor no sólo sabe que no tiene que leer más, sino que le indica que el cliente finalizó. 
A continuación, el cliente abre una nueva conexión con el servidor, pero se queda bloqueado esperando a recibir los resultados 
del sorteo que realiza el servidor. Una vez que le llegan, los notifica.

El servidor implementa una estrategia de `long polling` para hacer el sorteo. Cada vez que le llega un mensaje `Finish` de un cliente, 
lo agrega a un diccionario con su address y socket. Cuando le lleguen tantos mensaje `Finish` como clientes haya (el diccionario tiene
una cantidad de elementos igual a la cantidad de clientes), sabe que todos terminaron y realiza el sorteo. Finalmente, notifica uno 
por uno a los clientes que se habían quedando esperando por el resultado del sorteo.