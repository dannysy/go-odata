### **OData Gen**

Этот пакет позволяет генерировать сущности по $metadata описанию OData сервиса.

#### **Установка**
```
go install github.com/dannysy/go-odata/go-odata-gen
```

#### **Использование**
```
go-odata-gen -path=$PATH_TO_EDMX_FILE_OR_URL$ -out=$OUTPUT_FOLDER$
```