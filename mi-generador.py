import sys
import yaml

def main():
    if len(sys.argv) != 3:
        print("Uso: python3 mi-generador.py <archivo_salida> <cantidad_clientes>")
        sys.exit(1)

    archivo_salida = sys.argv[1]
    try:
        cantidad_clientes = int(sys.argv[2])
    except ValueError:
        print("La cantidad de clientes debe ser un número entero.")
        sys.exit(1)

    compose = {
        "name": "tp0",
        "services": {
            "server": {
                "container_name": "server",
                "image": "server:latest",
                "entrypoint": "python3 /main.py",
                "environment": [
                    "PYTHONUNBUFFERED=1"
                ],
                "networks": ["testing_net"]
            }
        },
        "networks": {
            "testing_net": {
                "ipam": {
                    "driver": "default",
                    "config": [
                        {"subnet": "172.25.125.0/24"}
                    ]
                }
            }
        }
    }

    # Agregar clientes
    for i in range(1, cantidad_clientes + 1):
        compose["services"][f"client{i}"] = {
            "container_name": f"client{i}",
            "image": "client:latest",
            "entrypoint": "/client",
            "environment": [
                f"CLI_ID={i}"
            ],
            "networks": ["testing_net"],
            "depends_on": ["server"]
        }

    # Guardar YAML
    with open(archivo_salida, "w", encoding="utf-8") as f:
        yaml.dump(compose, f, sort_keys=False)

if __name__ == "__main__":
    main()