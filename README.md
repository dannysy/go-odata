### **OData Query Builder**

Этот пакет предоставляет простой и удобный способ построения запросов OData. Он позволяет легко создавать, изменять и выполнять запросы к источникам данных, поддерживающим протокол OData.

#### **Установка**
```
go get github.com/dannysy/go-odata
```

#### **Использование**
Для начала работы с пакетом необходимо создать объект построителя запросов:
```go
import (
    "github.com/dannysy/go-odata"
)

func main() {
    // Если необходимо получить сущность по идентификатору
    // Для такой операции возможно использование опций select и expand
    entity := NewEntityBuilder("1","Products")
    // Если необходимо получить множество сущностей
    // Для такой операции возможно использование всех опций
    entity = NewListBuilder("Products")
    // Если необходимо получить количество сущностей
    // Для такой операции возможно использование опции filter
    entity = NewCountBuilder("Products")
}
```
Затем можно добавить параметры запроса, используя методы объекта построителя:
```go
entity.
    With(
        NewSelect("id", "name"),
        NewExpandBuilder().With("category", NewSelect("id", "name")).Build(),
        NewFilter(
            NewCombination(
                NewComparison("a", "b", ComparatorEQ),
                NewComparison("c", "d", ComparatorNEQ),
                And,
            ),
    ),
    NewTop(10),
    NewOrderByBuilder().WithColumns("id").WithDirectionalColumn("name", Asc).Build(),
    ).Build()
```
После добавления всех необходимых параметров можно получить query:
```go
entity.CollectToString()
```

#### **[Генерация типов и метаданных](gen/README.md)**
