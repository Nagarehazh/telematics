
# 🚀 Welcome to Telematics

To get started with Kafka, you can use the following Docker command:
```sh
docker run --name kafka -p 9092:9092 -e ALLOW_PLAINTEXT_LISTENER=yes -e KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE=true bitnami/kafka:latest
```

## 🛠️ Installing protobuf compiler (protoc compiler)

### For Linux Users (or WSL2):
```sh
sudo apt install -y protobuf-compiler
```

### For Mac Users:
You can use Homebrew:
```sh
brew install protobuf
```

## 🔧 Installing gRPC and Protobuf Plugins for Golang

### 1. Protobuf
```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
```

### 2. gRPC
```sh
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

### 3. Setting the PATH
Ensure that the `/go/bin` directory is in your PATH:
```sh
PATH="${PATH}:${HOME}/go/bin"
```

### 4. Install Package Dependencies

#### 4.1 Protobuf Package
```sh
go get google.golang.org/protobuf
```

#### 4.2 gRPC Package
```sh
go get google.golang.org/grpc/
```

## 📈 Installing Prometheus

### Using Docker:
```sh
docker run -p 9090:9090 -v ./.config/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
```

### Prometheus Golang Client:
```sh
go get github.com/prometheus/client_golang/prometheus
```

### Installing Prometheus Natively:

#### 1. Clone the Repository:
```sh
git clone https://github.com/prometheus/prometheus.git
```

#### 2. Install:
```sh
cd prometheus
make build
```

#### 3. Run the Prometheus Daemon:
```sh
./prometheus --config.file=<your_config_file>.yml
```

#### 4. Running from Inside the Project Directory:
```sh
../prometheus/prometheus --config.file=.config/prometheus.yml
```