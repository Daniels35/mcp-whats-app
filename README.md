# WhatsApp MCP Server

Servidor Model Context Protocol (MCP) optimizado para WhatsApp que permite integrar WhatsApp con Claude Desktop de forma eficiente.

## 📋 Índice
- [Características](#características)
- [Arquitectura](#arquitectura)
- [Requisitos](#requisitos)
- [Instalación y Configuración](#instalación-y-configuración)
- [Guía de Uso](#guía-de-uso)
- [Herramientas Disponibles](#herramientas-disponibles)
- [Rendimiento y Optimizaciones](#rendimiento-y-optimizaciones)
- [Almacenamiento de Datos](#almacenamiento-de-datos)
- [Seguridad](#seguridad)
- [Solución de Problemas](#solución-de-problemas)
- [Licencia](#licencia)

## 🚀 Características

- **Búsqueda y lectura de mensajes** personales de WhatsApp (imágenes, videos, documentos, audio)
- **Búsqueda de contactos** y envío de mensajes a individuos o grupos
- **Envío de archivos multimedia** (imágenes, videos, documentos, mensajes de audio)
- **Conexión directa** a tu cuenta personal de WhatsApp vía API web multidevice
- **Almacenamiento local** optimizado en SQLite
- **Carga lazy** - no carga historial completo al iniciar
- **Rendimiento optimizado** para consultas grandes

## 🏗️ Arquitectura

### Componentes principales:

1. **Go WhatsApp Bridge** (`whatsapp-bridge/`): 
   - Conecta a API web de WhatsApp
   - Maneja autenticación QR
   - Almacena mensajes en SQLite
   - API REST para comunicación

2. **Python MCP Server** (`whatsapp-mcp-server/`): 
   - Implementa protocolo MCP
   - Herramientas optimizadas para Claude
   - Consultas eficientes a la BD

## 📋 Requisitos

- Go 1.24.1+
- Python 3.11+
- Claude Desktop app
- UV (gestor Python): `curl -LsSf https://astral.sh/uv/install.sh | sh`
- FFmpeg (opcional, para audio)

## 🔧 Instalación y Configuración

### Instalación Automática (Recomendada)

**Ejecutar en orden:**

1. **`instalar-dependencias.bat`** *(solo primera vez)*
   - Instala dependencias Go y Python
   - Crea `.env` desde `.env.example`
   - Verifica configuración

2. **`iniciar-whatsapp-bridge.bat`** *(cada uso)*
   - Inicia servidor Go
   - Primera vez: escanea código QR
   - **MANTENER ABIERTO** durante uso

3. **`verificar-configuracion.bat`** *(troubleshooting)*
   - Verifica que todo funcione
   - Útil para diagnósticos

### Configuración de Claude Desktop

**Archivo de configuración JSON:**

```json
{
  "mcpServers": {
    "whatsapp": {
      "command": "{{PATH_TO_UV}}",
      "args": [
        "--directory",
        "{{PATH_TO_SRC}}/mi-whatsapp-mcp/whatsapp-mcp-server",
        "run",
        "main.py"
      ]
    }
  }
}
```

**Ubicaciones del archivo:**
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Linux**: `~/.config/Claude/claude_desktop_config.json`

### Instalación Manual

```bash
# 1. Clonar repositorio
git clone https://github.com/tuusuario/mi-whatsapp-mcp.git
cd mi-whatsapp-mcp

# 2. Configurar entorno
cp .env.example .env

# 3. Instalar dependencias Go
cd whatsapp-bridge
go mod tidy

# 4. Instalar dependencias Python
cd ../whatsapp-mcp-server
uv sync
```

**Para Windows con CGO:**
```bash
cd whatsapp-bridge
go env -w CGO_ENABLED=1
go run main.go
```

## 📝 Guía de Uso

### Primera configuración:
1. ✅ **Ejecutar**: `instalar-dependencias.bat`
2. ✅ **Configurar**: Claude Desktop config
3. ✅ **Ejecutar**: `iniciar-whatsapp-bridge.bat`
4. ✅ **Escanear**: Código QR con WhatsApp móvil
5. ✅ **Reiniciar**: Claude Desktop

### Uso diario:
1. **Ejecutar**: `iniciar-whatsapp-bridge.bat` (mantener abierto)
2. **Abrir**: Claude Desktop
3. **Usar**: Herramientas WhatsApp en Claude

## 🛠️ Herramientas Disponibles

| Herramienta | Descripción | Optimización |
|-------------|-------------|--------------|
| `search_contacts` | Buscar contactos por nombre/número | Límite 25 resultados, filtros inteligentes |
| `list_messages` | Recuperar mensajes con filtros | Requiere filtros o `force_load=True` |
| `get_last_interaction` | Último mensaje de un contacto | Consulta directa optimizada |
| `get_message_context` | Contexto alrededor de mensaje | Limitado a 5 mensajes para rendimiento |
| `send_message` | Enviar mensaje texto | Validación de entrada |
| `send_file` | Enviar archivos multimedia | Verificación de rutas |
| `send_audio_message` | Enviar mensaje de voz | Conversión automática a Opus |
| `download_media` | Descargar multimedia | Rutas locales seguras |

## ⚡ Rendimiento y Optimizaciones

### Carga Lazy Implementada
- **No carga historial completo** al iniciar
- **Conexión DB solo cuando necesario**
- **Verificación de existencia** de BD antes de conectar

### Optimizaciones de Consultas
```python
# ❌ Antes: Podía cargar todo el historial
list_messages()

# ✅ Ahora: Requiere filtros específicos
list_messages(chat_jid="123456@s.whatsapp.net", limit=20)
list_messages(query="proyecto", after="2024-12-01")
list_messages(force_load=True, limit=10)  # Solo si necesitas forzar
```

### Límites de Rendimiento
- **Mensajes**: Máximo 50 por consulta
- **Contactos**: Máximo 25 resultados
- **Contexto**: Máximo 5 mensajes con contexto
- **Búsqueda**: Mínimo 2 caracteres

### Mejoras de Búsqueda
- **Indexación optimizada** en SQLite
- **Patrones de búsqueda inteligentes**
- **Ordenamiento por relevancia**
- **Filtros NULL eliminados**

## 💾 Almacenamiento de Datos

### Estructura de Base de Datos
```
whatsapp-bridge/store/
├── messages.db     # Mensajes y chats
└── whatsapp.db     # Sesión WhatsApp
```

### Tablas Principales
- **chats**: Información de chats (JID, nombre, último mensaje)
- **messages**: Mensajes completos con multimedia
- **Índices optimizados** para búsquedas rápidas

### Características
- **Almacenamiento local** (no cloud)
- **SQLite con WAL mode** para concurrencia
- **Claves foráneas** para integridad
- **Limpieza automática** de datos antiguos

## 🔒 Seguridad

### Variables de Entorno (.env)
```bash
# API Configuration
WHATSAPP_API_HOST=localhost
WHATSAPP_API_PORT=8080
WHATSAPP_API_BASE_URL=http://localhost:8080/api

# Database Configuration  
MESSAGES_DB_NAME=messages.db
WHATSAPP_DB_NAME=whatsapp.db

# Server Configuration
REST_SERVER_PORT=8080
DEBUG=false
LOG_LEVEL=INFO
```

### Medidas Implementadas
✅ **Variables de entorno** para configuraciones sensibles  
✅ **Archivo .env protegido** en .gitignore  
✅ **Rutas relativas** en lugar de hardcodeadas  
✅ **Configuración por defecto** segura  
✅ **Separación de secretos** del código  
✅ **Validación de entrada** en todas las funciones  
✅ **Manejo seguro de rutas** de archivos  

### Archivos Protegidos (.gitignore)
- `.env` - Variables de entorno reales
- `whatsapp-bridge/store/` - Datos y mensajes
- `*.db` - Bases de datos
- `*.key`, `*.pem` - Claves y certificados

### ⚠️ NUNCA hagas esto:
❌ Commitear archivos `.env` reales  
❌ Hardcodear passwords o tokens  
❌ Usar puertos por defecto en producción  
❌ Exponer API sin autenticación  
❌ Commitear bases de datos con datos  

## ❓ Solución de Problemas

### Problemas Comunes

**🔸 Código QR no se muestra**
```bash
# Reiniciar bridge
iniciar-whatsapp-bridge.bat
# Verificar terminal soporta QR
```

**🔸 "No se cargan mensajes"**
```python
# En Claude, usar filtros específicos:
# ❌ list_messages() 
# ✅ list_messages(chat_jid="contact@s.whatsapp.net")
# ✅ list_messages(query="palabra", limit=10)
```

**🔸 Rendimiento lento**
- Usar filtros específicos en consultas
- Evitar `include_context=True` sin filtros
- Limitar resultados con `limit` y `max_results`

**🔸 WhatsApp ya conectado**
- Bridge se reconecta automáticamente
- No necesita nuevo QR si sesión activa

**🔸 Límite de dispositivos**
- WhatsApp limita dispositivos vinculados
- Eliminar dispositivo desde WhatsApp móvil

**🔸 Problemas de conexión**
1. Ejecutar `verificar-configuracion.bat`
2. Verificar que bridge esté ejecutándose  
3. Reiniciar Claude Desktop
4. Verificar config en `%APPDATA%\Claude\claude_desktop_config.json`

### Comandos de Diagnóstico

```bash
# Verificar configuración
verificar-configuracion.bat

# Verificar procesos
netstat -an | findstr :8080

# Verificar logs
cd whatsapp-bridge
go run main.go

# Verificar base de datos
cd whatsapp-bridge/store
sqlite3 messages.db ".tables"
```

### Logs y Debug
```bash
# Activar modo debug en .env
DEBUG=true
LOG_LEVEL=DEBUG

# Ver logs en tiempo real
tail -f whatsapp-bridge/logs/app.log
```

## 📊 Métricas de Rendimiento

### Tiempos de Respuesta Esperados
- **Inicio MCP server**: < 5 segundos
- **Búsqueda contactos**: < 1 segundo
- **Lista mensajes (filtrados)**: < 2 segundos
- **Envío mensaje**: < 3 segundos

### Uso de Memoria
- **MCP Server**: ~50MB RAM
- **WhatsApp Bridge**: ~100MB RAM
- **Base de datos**: Variable según historial

## 🔔 Webhook opcional para n8n

El bridge puede notificar a n8n cuando recibe un mensaje, evitando que n8n tenga que consultar `list_messages` continuamente. La función queda desactivada si `N8N_WEBHOOK_URL` no está definida.

Configura estas variables en el entorno del bridge:

```env
N8N_WEBHOOK_URL=https://n8n.example.com/webhook/whatsapp-incoming
N8N_WEBHOOK_TOKEN=replace-with-a-long-secret
N8N_WEBHOOK_QUEUE_SIZE=100
N8N_WEBHOOK_TIMEOUT_SECONDS=10
N8N_WEBHOOK_RETRIES=3
```

Solo los mensajes entrantes se envían por defecto. Cada solicitud incluye el encabezado `X-Idempotency-Key` con el ID de WhatsApp y un cuerpo JSON con `id`, `chat_jid`, `sender`, `timestamp`, `content`, `is_from_me`, `media_type` y `filename`. El bridge usa una cola limitada y reintentos para que un n8n temporalmente caído no bloquee la conexión con WhatsApp.

## 🔌 Webhook opcional para cambios de estado de WhatsApp

El bridge también puede notificar a n8n cuando la sesión autenticada de WhatsApp pierde la conexión, se recupera o queda cerrada. Este webhook es independiente del webhook de mensajes y no debe conectarse a los flujos del asistente, cotizaciones o CRM.

La notificación se genera desde los eventos de `whatsmeow` (`Connected`, `Disconnected`, `LoggedOut` y `StreamReplaced`). Los eventos se deduplican por transición de estado: una reconexión automática no produce una alerta nueva hasta que WhatsApp vuelva a emitir `Connected`. El pairing inicial y los fallos anteriores a una autenticación se mantienen silenciosos para evitar falsos avisos.

Variables opcionales:

```env
N8N_STATUS_WEBHOOK_URL=https://n8n.example.com/webhook/whatsapp-status
N8N_STATUS_WEBHOOK_TOKEN=replace-with-a-long-secret
N8N_STATUS_WEBHOOK_QUEUE_SIZE=100
N8N_STATUS_WEBHOOK_TIMEOUT_SECONDS=10
N8N_STATUS_WEBHOOK_RETRIES=3
```

Si no se definen las variables específicas de tamaño, timeout o reintentos, se reutilizan los valores de `N8N_WEBHOOK_QUEUE_SIZE`, `N8N_WEBHOOK_TIMEOUT_SECONDS` y `N8N_WEBHOOK_RETRIES`. El token de estado usa `N8N_STATUS_WEBHOOK_TOKEN` y, si no existe, puede reutilizar `N8N_WEBHOOK_TOKEN`.

El bridge envía el token únicamente en la cabecera `Authorization: Bearer <token>`. El cuerpo no contiene mensajes, números de terceros ni secretos:

```json
{
  "event_id": "whatsapp-status-...",
  "event": "disconnected",
  "timestamp": "2026-10-18T12:00:00Z",
  "reason": "connection_lost",
  "service": "whatsapp-mcp"
}
```

El envío se realiza mediante una cola asíncrona y limitada. Un n8n caído no bloquea la recepción ni la reconexión de WhatsApp; el bridge registra el error, aplica reintentos limitados y descarta eventos si la cola está llena. Este webhook no puede avisar si el contenedor o el VPS mueren por completo. Para ese caso se necesita posteriormente un monitor externo o una alerta del supervisor de Docker.

## 📄 Licencia

MIT License - Copyright (c) 2025

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
