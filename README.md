# iso8583

[![Go Reference](https://pkg.go.dev/badge/github.com/tomasdemarco/iso8583.svg)](https://pkg.go.dev/github.com/tomasdemarco/iso8583)

Librería para la codificación y decodificación de mensajes ISO 8583.

## Uso Básico

### 1. Definir un Packager en JSON

La librería utiliza un "packager" para definir la estructura de los mensajes ISO 8583. Este se define en un archivo JSON.

**Ejemplo: `mi_packager.json`**

```json
{
    "description": "Mi Packager ISO 8583",
    "prefix": {
        "type": "LL",
        "encoding": "ASCII"
    },
    "fields": {
        "000": {
            "description": "Message Type Indicator",
            "type": "NUMERIC",
            "length": 4,
            "pattern": "^(0200|0210|0400|0410)$",
            "encoding": "ASCII"
        },
        "001": {
            "description": "Bitmap",
            "type": "BITMAP",
            "length": 16,
            "pattern": "^[0-9a-fA-F]{16,32}$",
            "encoding": "ASCII"
        },
        "003": {
            "description": "Processing Code",
            "type": "NUMERIC",
            "length": 6,
            "pattern": "^[0-9]{6}$",
            "encoding": "ASCII"
        },
        "004": {
            "description": "Transaction Amount",
            "type": "NUMERIC",
            "length": 12,
            "pattern": "^[0-9]{12}$",
            "encoding": "ASCII"
        }
    }
}
```

### 2. Código Go

```go
package main

import (
    "fmt"
    "log"

    "github.com/tomasdemarco/iso8583/message"
    "github.com/tomasdemarco/iso8583/packager"
)

func main() {
    // Cargar la definición del packager desde el archivo JSON
    pkg, err := packager.LoadFromJsonV2("./packager_json_files", "mi_packager.json")
    if err != nil {
        log.Fatalf("Error cargando el packager: %v", err)
    }

    // Crear un nuevo mensaje
    msg := message.NewMessage(pkg)

    // Añadir campos al mensaje
    msg.SetField("000", "0200")
    msg.SetField("003", "000000")
    msg.SetField("004", "000000010000") // 100.00

    // Empaquetar el mensaje
    packedMsg, err := msg.Pack()
    if err != nil {
        log.Fatalf("Error empaquetando el mensaje: %v", err)
    }
    fmt.Printf("Mensaje empaquetado: %x
", packedMsg)

    // Desempaquetar un mensaje
    unpackedMsg := message.NewMessage(pkg)
    err = unpackedMsg.Unpack(packedMsg)
    if err != nil {
        log.Fatalf("Error desempaquetando el mensaje: %v", err)
    }

    // Leer campos del mensaje desempaquetado
    mti, _ := unpackedMsg.GetField("000")
    processingCode, _ := unpackedMsg.GetField("003")
    amount, _ := unpackedMsg.GetField("004")

    fmt.Printf("MTI: %s
", mti)
    fmt.Printf("Processing Code: %s
", processingCode)
    fmt.Printf("Amount: %s
", amount)
}
```

## Encoding

La propiedad `encoding` en la definición de un campo especifica cómo se codifican y decodifican los datos. La librería soporta las siguientes codificaciones:

*   **`ASCII`**: Codificación estándar de caracteres. Los datos se representan como cadenas de texto.
*   **`BCD`** (Binary-Coded Decimal): Cada dígito decimal se representa con 4 bits. Es más compacto que ASCII para datos numéricos.
*   **`EBCDIC`** (Extended Binary Coded Decimal Interchange Code): Otro esquema de codificación de caracteres, común en mainframes de IBM.
*   **`BINARY`**: Representación binaria cruda de los datos. La cadena de entrada se trata como una representación hexadecimal de los datos binarios.
*   **`HEX`**: Similar a `BINARY`, pero explícitamente para cadenas hexadecimales.

La elección de la codificación correcta es crucial para asegurar la interoperabilidad con otros sistemas ISO 8583.

## Prefijos de Longitud (prefix)

Para campos de longitud variable, la propiedad `prefix` en la definición del campo especifica cómo se codifica y decodifica la longitud de los datos.

*   **`type`**: Define el número de dígitos que componen el prefijo. Los valores comunes son:
    *   `L`: 1 dígito.
    *   `LL`: 2 dígitos.
    *   `LLL`: 3 dígitos.
    *   `LLLL`: 4 dígitos.
    *   `FIXED`: Indica que el campo tiene una longitud fija y no lleva prefijo.
*   **`encoding`**: La codificación del propio prefijo (puede ser diferente a la del campo). Soporta los mismos valores que la codificación de campos: `ASCII`, `BCD`, `EBCDIC`, `BINARY`.
*   **`isInclusive`**: (booleano, opcional) Si es `true`, la longitud indicada por el prefijo incluye la longitud del propio prefijo. Por defecto es `false`.
*   **`hex`**: (booleano, opcional) Si es `true`, la longitud se interpreta como un número hexadecimal. Por defecto es `false`.

## Relleno (padding)

Para campos de longitud fija, la propiedad `padding` en la definición del campo asegura que los datos siempre tengan la longitud correcta.

*   **`type`**: Define el tipo de relleno a aplicar.
    *   `FILL`: Rellena el campo con un caracter específico hasta alcanzar la longitud definida.
    *   `PARITY`: Asegura que la longitud del campo sea par, agregando un caracter si es necesario. Útil para codificaciones como BCD donde los datos se agrupan en pares.
    *   `NONE`: No se aplica ningún relleno.
*   **`position`**: Indica de qué lado se debe agregar el relleno.
    *   `LEFT`: El relleno se agrega a la izquierda.
    *   `RIGHT`: El relleno se agrega a la derecha.
*   **`char`**: El caracter que se usará para el relleno (ej. "0", " ").

## Datos EMV (Tag-Length-Value)

El paquete `emv` proporciona herramientas para trabajar con datos codificados en formato TLV (Tag-Length-Value), que son comunes en campos como el DE-55 (Datos ICC).

### Desempaquetar (Unpack)

La función `emv.Unpack` permite parsear una cadena de datos TLV (en formato hexadecimal) y extraer los valores de tags específicos.

```go
package main

import (
    "fmt"
    "log"

    "github.com/tomasdemarco/iso8583/emv"
)

func main() {
    tlvData := "9F02060000000012349F03060000000000009F26082B84525423221105"

    // Desempaquetar todos los tags
    tags, err := emv.Unpack(tlvData)
    if err != nil {
        log.Fatalf("Error desempaquetando EMV: %v", err)
    }
    fmt.Printf("Tags desempaquetados: %v
", tags)

    // Desempaquetar solo tags específicos
    filteredTags, err := emv.Unpack(tlvData, "9F02", "9F26")
    if err != nil {
        log.Fatalf("Error desempaquetando EMV: %v", err)
    }
    fmt.Printf("Tags filtrados: %v
", filteredTags)
}
```

### Empaquetar (Pack)

La función `emv.Pack` construye una cadena de datos TLV a partir de un mapa de tags y valores. Los tags se ordenan alfabéticamente antes de ser empaquetados.

```go
package main

import (
    "fmt"

    "github.com/tomasdemarco/iso8583/emv"
)

func main() {
    tags := map[string]string{
        "9F02": "000000001234",
        "9F03": "000000000000",
        "9F26": "2B84525423221105",
    }

    tlvData, err := emv.Pack(tags)
    if err != nil {
        log.Fatalf("Error empaquetando EMV: %v", err)
    }

    fmt.Printf("Datos TLV empaquetados: %s
", tlvData)
}
```

## Estructura del JSON del Packager

*   `description`: (string) Una breve descripción del packager.
*   `prefix`: (object) Define el prefijo de longitud del mensaje.
    *   `type`: (string) El tipo de prefijo (ej. "LL", "LLL", "LLLL").
    *   `encoding`: (string) La codificación del prefijo (ej. "ASCII", "BCD").
*   `fields`: (object) Un mapa de las definiciones de los campos. La clave es el número del campo (ej. "000", "003").
    *   `description`: (string) Descripción del campo.
    *   `type`: (string) Tipo de dato del campo (ej. "NUMERIC", "STRING", "BITMAP").
    *   `length`: (int) La longitud máxima del campo.
    *   `pattern`: (string) Una expresión regular para validar el formato del campo.
    *   `encoding`: (string) La codificación del campo (ej. "ASCII", "BCD", "EBCDIC", "BINARY", "HEX").
    *   `prefix`: (object, opcional) Define un prefijo de longitud para el campo (similar al prefijo del mensaje).
    *   `padding`: (object, opcional) Define el relleno del campo.
        *   `type`: (string) "PARITY" o "FILL".
        *   `position`: (string) "LEFT" o "RIGHT".
        *   `char`: (string) El caracter a usar para el relleno.

