package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// --- Structs para Parsear o XSD ---
// Estas structs espelham a estrutura de um arquivo XSD (XML Schema Definition)
// para que a biblioteca `encoding/xml` possa lê-los.

// Schema é o elemento raiz do XSD
type Schema struct {
	XMLName      xml.Name      `xml:"schema"`
	Elements     []Element     `xml:"element"`
	ComplexTypes []ComplexType `xml:"complexType"`
	SimpleTypes  []SimpleType  `xml:"simpleType"`
}

// Element representa <xs:element name="..." ...>
type Element struct {
	Name        string       `xml:"name,attr"`
	Annotation  Annotation   `xml:"annotation"`
	ComplexType *ComplexType `xml:"complexType"`
}

// ComplexType representa <xs:complexType name="..." ...>
type ComplexType struct {
	Name       string      `xml:"name,attr"`
	Annotation Annotation  `xml:"annotation"`
	Sequence   *Sequence   `xml:"sequence"`
	Choice     *Choice     `xml:"choice"`
	Attributes []Attribute `xml:"attribute"`
}

// SimpleType representa <xs:simpleType name="..." ...>
type SimpleType struct {
	Name       string     `xml:"name,attr"`
	Annotation Annotation `xml:"annotation"`
}

// Attribute representa <xs:attribute name="..." ...>
type Attribute struct {
	Name       string     `xml:"name,attr"`
	Annotation Annotation `xml:"annotation"`
}

// Sequence representa <xs:sequence ...>
type Sequence struct {
	Elements []Element `xml:"element"`
}

// Choice representa <xs:choice ...>
type Choice struct {
	Elements []Element `xml:"element"`
}

// Annotation representa <xs:annotation ...>
type Annotation struct {
	Documentation string `xml:"documentation"`
}

// --- Lógica Principal ---

// tagDocs é nosso mapa global que armazena "tag" -> "documentação"
var tagDocs = make(map[string]string)

func main() {
	fmt.Println("// Script iniciado. Lendo arquivos XSD no diretório atual...")

	// Encontra todos os arquivos .xsd no diretório
	files, err := filepath.Glob("*.xsd")
	if err != nil {
		fmt.Println("Erro ao procurar arquivos XSD:", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("Nenhum arquivo .xsd encontrado. Você baixou o Pacote de Liberação (PL) da NF-e e o colocou nesta pasta?")
		return
	}

	// Processa cada arquivo XSD encontrado
	for _, file := range files {
		fmt.Printf("// Processando: %s\n", file)
		parseXSD(file)
	}

	fmt.Println("\n// Processamento concluído. Gerando mapa Go...")
	fmt.Println("// NOTA: 'Attr:' = Atributo XML, 'Type:' = Definição de Tipo")
	fmt.Println("")

	// Imprime o mapa Go formatado
	printGoMap()
}

// parseXSD abre e faz o unmarshal do arquivo XSD
func parseXSD(filename string) {
	xmlFile, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Erro ao abrir %s: %v\n", filename, err)
		return
	}
	defer xmlFile.Close()

	byteValue, _ := io.ReadAll(xmlFile)

	var schema Schema
	err = xml.Unmarshal(byteValue, &schema)
	if err != nil {
		fmt.Printf("Erro ao fazer Unmarshal de %s: %v\n", filename, err)
		return
	}

	// Inicia a varredura recursiva
	for _, el := range schema.Elements {
		processElement(el)
	}
	for _, ct := range schema.ComplexTypes {
		processComplexType(ct)
	}
	for _, st := range schema.SimpleTypes {
		processSimpleType(st)
	}
}

// --- Funções Recursivas de Processamento ---

// processElement varre um elemento e seus sub-elementos
func processElement(el Element) {
	doc := strings.TrimSpace(el.Annotation.Documentation)
	if el.Name != "" && doc != "" {
		tagDocs[el.Name] = doc
	}

	// Se o elemento tiver um complexType aninhado
	if el.ComplexType != nil {
		processComplexType(*el.ComplexType)
	}
}

// processComplexType varre um tipo complexo (grupo de tags/atributos)
func processComplexType(ct ComplexType) {
	doc := strings.TrimSpace(ct.Annotation.Documentation)
	if ct.Name != "" && doc != "" {
		tagDocs[ct.Name] = doc
	}

	// Varre elementos dentro de <sequence>
	if ct.Sequence != nil {
		for _, el := range ct.Sequence.Elements {
			processElement(el)
		}
	}

	// Varre elementos dentro de <choice>
	if ct.Choice != nil {
		for _, el := range ct.Choice.Elements {
			processElement(el)
		}
	}

	// Varre atributos
	for _, at := range ct.Attributes {
		processAttribute(at)
	}
}

// processSimpleType varre um tipo simples (campo com restrições)
func processSimpleType(st SimpleType) {
	doc := strings.TrimSpace(st.Annotation.Documentation)
	if st.Name != "" && doc != "" {
		// Prefixamos para saber que é um TIPO, não uma TAG
		tagDocs["Type:"+st.Name] = doc
	}
}

// processAttribute varre um atributo
func processAttribute(at Attribute) {
	doc := strings.TrimSpace(at.Annotation.Documentation)
	if at.Name != "" && doc != "" {
		// Prefixamos para saber que é um ATRIBUTO, não uma TAG
		tagDocs["Attr:"+at.Name] = doc
	}
}

// printGoMap formata a saída como um mapa Go
func printGoMap() {
	// Pega as chaves e ordena para uma saída limpa
	keys := make([]string, 0, len(tagDocs))
	for k := range tagDocs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Println("var tagMap = map[string]string{")
	for _, k := range keys {
		// Limpa a string de documentação para caber no Go
		doc := tagDocs[k]
		doc = strings.ReplaceAll(doc, "\n", " ")    // Remove quebras de linha
		doc = strings.ReplaceAll(doc, "\"", "\\\"") // Escapa aspas
		doc = strings.ReplaceAll(doc, "\t", "")     // Remove tabs
		doc = strings.TrimSpace(doc)                // Remove espaços extras

		fmt.Printf("\t\"%s\": \"%s\",\n", k, doc)
	}
	fmt.Println("}")
}
