# net-chat-server

Servidor de chat, soporta multiples rooms de chat, los clientes deben indicar el nombre de usuario 
y el room al que se desean conectar.

## Instrucciones para la ejecución

1. Descargar el repositorio
2. Ejecutar el comando `go mod tidy` (Se debe tener instalado go) para instalar las dependencia
3. Ejecutar el comando `go run .` (Ejecutar en modo desarrollo de lo contrario se deberia crear el build)

## Instrucciones de administrador del server

1. Enviar un msj a todos los users del server `/broadcast mensaje a enviar`
2. Listar los usuarios conectados al server `/listUsers`
3. Enviar un msj a un usuario especifico `/msg userID message`


