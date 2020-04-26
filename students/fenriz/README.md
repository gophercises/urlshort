# Exercise #2: URL Shortener

[![exercise status: released](https://img.shields.io/badge/exercise%20status-released-green.svg?style=for-the-badge)](https://gophercises.com/exercises/urlshort)



## Exercise details

El objetivo de este ejercicio es crear un [http.Handler] (https://golang.org/pkg/net/http/#Handler) que examinará la ruta de cualquier solicitud web entrante y determinará si debe redirigir el usuario a una nueva página, como lo haría el acortador de URL.

Por ejemplo, si tenemos una configuración de redireccionamiento para `/ dogs` a` https: // www.somesite.com / a-story-about-dogs`, buscaríamos cualquier solicitud web entrante con la ruta `/ dogs` y redirigirlos.

Para completar estos ejercicios, deberá implementar los métodos apagados en [handler.go] (https://github.com/gophercises/urlshort/blob/master/handler.go). Hay una buena cantidad de comentarios que explican qué debe hacer cada método, y también hay una fuente [main / main.go] (https://github.com/gophercises/urlshort/blob/master/main/main.go) archivo que usa el paquete para ayudarlo a probar su código y tener una idea de lo que su programa debería estar haciendo.

Sugiero comentar primero todo el código en main.go relacionado con la función `YAMLHandler` y centrarme en implementar primero la función` MapHandler`.

Una vez que tenga eso funcionando, concéntrese en analizar el YAML utilizando el paquete [gopkg.in/yaml.v2font>(https://godoc.org/gopkg.in/yaml.v2). * Nota: Tendrá que 'obtener' este paquete si aún no lo tiene. *

Después de obtener el análisis YAML, intente convertir los datos en un mapa y luego use MapHandler para finalizar la implementación de YAMLHandler. Por ejemplo, puede terminar con un código como este:

```go
func YAMLHandler(yaml []byte, fallback http.Handler) (http.HandlerFunc, error) {
  parsedYaml, err := parseYAML(yaml)
  if err != nil {
    return nil, err
  }
  pathMap := buildMap(parsedYaml)
  return MapHandler(pathMap, fallback), nil
}
```
Pero para que esto funcione, deberá crear funciones como `parseYAML` y` buildMap` por su cuenta. Esto debería darle una amplia experiencia trabajando con datos YAML.


## Bonus

As a bonus exercises you can also...

1. Actualice el archivo fuente [main / main.go] (https://github.com/gophercises/urlshort/blob/master/main/main.go) para aceptar un archivo YAML como bandera y luego cargue el YAML desde un archivo en lugar de una cadena.
2. Cree un JSONHandler que tenga el mismo propósito, pero que lea datos de JSON.
3. Cree un controlador que no lea desde un mapa, sino que lea desde una base de datos. Ya sea que use BoltDB, SQL u otra cosa, depende completamente de usted.
