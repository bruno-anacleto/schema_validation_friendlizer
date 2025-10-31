# Schema Validation Friendlizer (POC)

Este projeto é uma **Prova de Conceito (POC)** que demonstra como "amigabilizar" erros de validação de schema XSD da NF-e, transformando mensagens de erro técnicas em descrições legíveis para o usuário final.

## O que ele faz?

O objetivo é transformar um erro de parser XML cru, que é técnico e difícil de entender em uma mensagem clara e aceitável para o usuário.

## Como Usar

Siga estes passos para executar a POC e atualizar o dicionário de campos.

### Passo 1: Verificar e Atualizar os XSDs

Os schemas (`.xsd`) mudam a cada nova Nota Técnica (NT) da SEFAZ. O dicionário de campos (`tagMap`) depende diretamente deles.

1.  **Verifique a Versão:** Visite o [Portal Nacional da NF-e](https://www.nfe.fazenda.gov.br/portal/principal.aspx) na seção "Documentos" > "Manuais".
2.  **Baixe o Pacote Mais Recente:** Procure pelo "Pacote de Liberação" (PL) mais recente. Baixe o arquivo `.zip`.
3.  **Atualize os Arquivos:** Descompacte o arquivo e copie **TODOS** os arquivos `.xsd` de dentro da pasta "Schemas" (ou similar) para o diretório `schema_validation_friendlizer/schema_dictionary_generator/`, substituindo os arquivos existentes.

*(Os arquivos XSD incluídos neste repositório são da **v4.00**, mas podem estar desatualizados em relação a pacotes de liberação mais recentes).*

### Passo 2: Executar o Gerador do Dicionário

Este script irá ler os `.xsd` e imprimir o `tagMap` formatado em Go no console.

1.  Navegue até o diretório do gerador:
    ```bash
    cd schema_validation_friendlizer/schema_dictionary_generator
    ```
2.  Execute o `main.go` e redirecione a saída para um novo arquivo:
    ```bash
    go run main.go > output_map.txt
    ```

3.  Agora você terá um arquivo `output_map.txt` contendo o `tagMap` mais atualizado.

### Passo 3: Atualizar o Aplicativo Principal

O `tagMap` gerado precisa ser copiado para o aplicativo principal.

1.  Abra o arquivo `schema_validation_friendlizer/schema_dictionary_generator/output_map.go`.
2.  Copie todo o bloco `var tagMap = map[string]string { ... }`.
3.  Abra o arquivo `schema_validation_friendlizer/main.go` (na raiz do projeto).
4.  **Substitua** o `var tagMap` antigo pelo novo `tagMap` que você acabou de copiar.

### Passo 4: Executar a Demonstração

Agora que o tradutor está com o dicionário atualizado, execute a POC.

1.  Volte para o diretório raiz do projeto:
    ```bash
    cd ..
    ```
2.  Execute o `main.go`:
    ```bash
    go run main.go
    ```

Você verá no console a tradução dos erros de exemplo (definidos na lista `rawErrors` dentro da `func main`) sendo formatados de maneira amigável.

---

## Próximos Passos (Evolução da POC)

Para um uso em produção, o `schema_validation_friendlizer/main.go` não seja um script de console, ele seria refatorado para:
* Expor uma função `TranslateError(rawError string) (FriendlyError, error)`.
* Ter seus padrões de Regex (`ErrorPattern`) expandidos para cobrir mais casos de erro (ex: `invalid token`, `element mismatch`, etc.).