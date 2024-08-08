![logo](https://i3.wp.com/raw.githubusercontent.com/Quamagi/htmlplus/main/logo.png)
```markdown

# HTML+

HTML+ es un servidor web escrito en Go que permite la ejecución de comandos embebidos en plantillas HTML. Este proyecto está diseñado para ser una herramienta flexible y extensible para desarrollar aplicaciones web dinámicas utilizando el lenguaje de programación Go.

## Características

- **Comandos Embebidos**: Usa etiquetas especiales en tus plantillas HTML para ejecutar comandos Go.
- **Registro de Comandos**: Los comandos se registran automáticamente y están disponibles para su uso en plantillas.
- **Soporte para CORS**: Middleware incluido para manejar CORS.
- **Recolección de Basura**: Soporte para la limpieza automática de variables almacenadas en memoria.

## Estructura del Proyecto

```
HTML+/
│   go.mod
│   main.go
│
├───cmd
│   │   variables.go
│   │
│   ├───towrite
│   │   │   command.go
│   │   │
│   │   └───docs
│   │           examples.txt
│   │
└───templates
        example.txt
        extra.gweb
        index.gweb
        index.html
```

## Instalación

1. Clona el repositorio:
    ```sh
    git clone https://github.com/tu-usuario/HTML+.git
    cd HTML+
    ```

2. Inicializa el módulo Go:
    ```sh
    go mod tidy
    ```

## Uso

1. Ejecuta el servidor:
    ```sh
    go run main.go
    ```

2. Abre tu navegador y navega a `http://127.0.0.1:8080/`.

3. Crea y edita archivos `.gweb` en la carpeta `templates` para incluir comandos embebidos. Por ejemplo, `index.gweb`:
    ```html
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Go Web Server</title>
        <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.2.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-rbsA2VBKQhggwzxH7pPCaAqO46MgnOM80zW1RWuH61DGLwZJEdK2Kadq2F9CUG65" crossorigin="anonymous">
    </head>
    <body>
        <h1>Hello, World!</h1>
        <div><?go towrite "Hello, World!" ?></div>
        <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.2.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-kenU1KFdBIe4zVF0s0G1M5b4hcpxyD9F7jL+jjXkk+Q2h455rYXK/7HAuoJl+0I4" crossorigin="anonymous"></script>
    </body>
    </html>
    ```

## Añadir Nuevos Comandos

1. Crea una nueva carpeta para el comando en la carpeta `cmd`, por ejemplo `cmd/nuevo_comando`.
2. Añade un archivo `command.go` en esta carpeta y define el comando siguiendo la estructura de otros comandos.
3. Registra el nuevo comando utilizando `cmd.Register` en la función `init`.

Ejemplo de `cmd/nuevo_comando/command.go`:

```go
package nuevo_comando

import (
    "fmt"
    "io"
    "HTML+/cmd"
)

type NuevoComando struct{}

func (c *NuevoComando) Execute(args []string, w io.Writer) error {
    if len(args) < 1 {
        return fmt.Errorf("usage: nuevo_comando <message>")
    }
    message := args[0]
    fmt.Fprintf(w, "<p>%s</p>", message)
    return nil
}

func init() {
    cmd.Register("nuevo_comando", &NuevoComando{})
}
```

## Contribuciones

Las contribuciones son bienvenidas. Por favor, abre un issue o envía un pull request.

## Licencia

Este proyecto está bajo la licencia MIT. Consulta el archivo `LICENSE` para obtener más información.
```

### Contacto

Para cualquier consulta o comentario, puedes contactar al desarrollador en `tu-email@example.com`.

---

Este `README.md` proporciona una descripción completa del proyecto, incluyendo su estructura, instalación, uso y cómo añadir nuevos comandos. Ajusta cualquier detalle específico según tus necesidades.
