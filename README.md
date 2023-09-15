## Algorítmo para remoção de sufixos 

Baseado [neste artigo](https://lume.ufrgs.br/bitstream/handle/10183/23576/000597277.pdf) de Alexandre R. Coelho.

```mermaid
graph TD
    A(Início)
    A --> X{Termina<br>em 's'?}
    X --> |Sim| B(Redução do Plural)
    X --> |Não| C(Redução Adverbial)
    B --> C
    C --> Y{Termina<br>em 'a'?}
    Y --> |Sim| D(Redução do Feminino)
    Y --> |Não| E(Redução do Aumentativo)
    D --> E
    E --> F(Redução Nominal)
    F --> W{Sufixo<br>removido?}
    W --> |Sim| L(Remoção de Acentuação)
    W --> |Não| G(Redução Verbal)
    G --> Z{Sufixo<br>removido?}
    Z --> |Sim| H(Remoção de Vogais Temáticas)
    Z --> |Não| L
    H --> L
    L --> M(Fim)
```

## Uso

### `Palavra(term string) string`
Para converter uma palavra: 
`Palavra("normais") => "norm"`

### `Frase(doc string) string`
Para conveter uma frase:
`Frase("Ajuste a Valor Justo - Investimentos") => "ajust val just invest"`
