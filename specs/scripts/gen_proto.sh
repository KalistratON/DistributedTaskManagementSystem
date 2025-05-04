#!/bin/bash

PROTO_DIR="./proto"
OUTPUT_DIR="./go/pkg"

if [ ! -d "$PROTO_DIR" ]; then
  echo "Ошибка: Директория .proto файлов не найдена: $PROTO_DIR"
  exit 1
fi

if [ ! -d "$OUTPUT_DIR" ]; then
  echo "Создание выходной директории: $OUTPUT_DIR"
  mkdir -p "$OUTPUT_DIR"
fi

PROTO_FILES=$(find "$PROTO_DIR" -type f -name "*.proto")
if [ -z "$PROTO_FILES" ]; then
  echo "Ошибка: В директории $PROTO_DIR не найдено .proto файлов."
  exit 1
fi

for proto_file in $PROTO_FILES; do
  echo "Обработка файла: $proto_file"

  protoc \
    --proto_path="$PROTO_DIR" \
    --go_out="$OUTPUT_DIR" \
    --go-grpc_out="$OUTPUT_DIR" \
    "$proto_file"

  if [ $? -eq 0 ]; then
    echo "Успешно сгенерировано: $proto_file"
  else
    echo "Ошибка при генерации: $proto_file"
  fi
done

echo "Генерация завершена."
