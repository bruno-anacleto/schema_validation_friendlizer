package main

import (
	"fmt"
	"regexp"
	"strings"
)

// 1. O DICIONÁRIO "DE-PARA"
// Mapeia a tag XML técnica (chave) para seu nome amigável (valor).
var tagMap = map[string]string{
	"AAMM":                        "AAMM da emissão da NF de produtor",
	"Attr:Id":                     "PL_005d - 11/08/09 - validação do Id",
	"Attr:dia":                    "Número do dia",
	"Attr:nItem":                  "Número do item do NF",
	"Attr:versao":                 "Versão do leiaute (v4.00)",
	"CEP":                         "CEP",
	"CFOP":                        "Código Fiscal de Operações e Prestações",
	"CNPJ":                        "CNPJ",
	"CNPJProd":                    "CNPJ do produtor da mercadoria, quando diferente do emitente. Somente para os casos de exportação direta ou indireta.",
	"CNPJReceb":                   "CNPJ do beneficiário do pagamento",
	"COFINS":                      "Dados do COFINS",
	"COFINSAliq":                  "Código de Situação Tributária do COFINS.  01 – Operação Tributável - Base de Cálculo = Valor da Operação Alíquota Normal (Cumulativo/Não Cumulativo); 02 - Operação Tributável - Base de Calculo = Valor da Operação (Alíquota Diferenciada);",
	"COFINSNT":                    "Código de Situação Tributária do COFINS: 04 - Operação Tributável - Tributação Monofásica - (Alíquota Zero); 06 - Operação Tributável - Alíquota Zero; 07 - Operação Isenta da contribuição; 08 - Operação Sem Incidência da contribuição; 09 - Operação com suspensão da contribuição;",
	"COFINSOutr":                  "Código de Situação Tributária do COFINS: 49 - Outras Operações de Saída 50 - Operação com Direito a Crédito - Vinculada Exclusivamente a Receita Tributada no Mercado Interno 51 - Operação com Direito a Crédito – Vinculada Exclusivamente a Receita Não Tributada no Mercado Interno 52 - Operação com Direito a Crédito - Vinculada Exclusivamente a Receita de Exportação 53 - Operação com Direito a Crédito - Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno 54 - Operação com Direito a Crédito - Vinculada a Receitas Tributadas no Mercado Interno e de Exportação 55 - Operação com Direito a Crédito - Vinculada a Receitas Não-Tributadas no Mercado Interno e de Exportação 56 - Operação com Direito a Crédito - Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno, e de Exportação 60 - Crédito Presumido - Operação de Aquisição Vinculada Exclusivamente a Receita Tributada no Mercado Interno 61 - Crédito Presumido - Operação de Aquisição Vinculada Exclusivamente a Receita Não-Tributada no Mercado Interno 62 - Crédito Presumido - Operação de Aquisição Vinculada Exclusivamente a Receita de Exportação 63 - Crédito Presumido - Operação de Aquisição Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno 64 - Crédito Presumido - Operação de Aquisição Vinculada a Receitas Tributadas no Mercado Interno e de Exportação 65 - Crédito Presumido - Operação de Aquisição Vinculada a Receitas Não-Tributadas no Mercado Interno e de Exportação 66 - Crédito Presumido - Operação de Aquisição Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno, e de Exportação 67 - Crédito Presumido - Outras Operações 70 - Operação de Aquisição sem Direito a Crédito 71 - Operação de Aquisição com Isenção 72 - Operação de Aquisição com Suspensão 73 - Operação de Aquisição a Alíquota Zero 74 - Operação de Aquisição sem Incidência da Contribuição 75 - Operação de Aquisição por Substituição Tributária 98 - Outras Operações de Entrada 99 - Outras Operações.",
	"COFINSQtde":                  "Código de Situação Tributária do COFINS. 03 - Operação Tributável - Base de Calculo = Quantidade Vendida x Alíquota por Unidade de Produto;",
	"COFINSST":                    "Dados do COFINS da Substituição Tributaria;",
	"CPF":                         "CPF Autorizado",
	"CPFRespTec":                  "CPF do Responsável Técnico pelo receituário",
	"CRT":                         "Código de Regime Tributário.  Este campo será obrigatoriamente preenchido com: 1 – Simples Nacional; 2 – Simples Nacional – excesso de sublimite de receita bruta; 3 – Regime Normal. 4 - Simples Nacional - Microempreendedor individual - MEI",
	"CST":                         "Código de Situação Tributária do COFINS: 49 - Outras Operações de Saída 50 - Operação com Direito a Crédito - Vinculada Exclusivamente a Receita Tributada no Mercado Interno 51 - Operação com Direito a Crédito – Vinculada Exclusivamente a Receita Não Tributada no Mercado Interno 52 - Operação com Direito a Crédito - Vinculada Exclusivamente a Receita de Exportação 53 - Operação com Direito a Crédito - Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno 54 - Operação com Direito a Crédito - Vinculada a Receitas Tributadas no Mercado Interno e de Exportação 55 - Operação com Direito a Crédito - Vinculada a Receitas Não-Tributadas no Mercado Interno e de Exportação 56 - Operação com Direito a Crédito - Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno, e de Exportação 60 - Crédito Presumido - Operação de Aquisição Vinculada Exclusivamente a Receita Tributada no Mercado Interno 61 - Crédito Presumido - Operação de Aquisição Vinculada Exclusivamente a Receita Não-Tributada no Mercado Interno 62 - Crédito Presumido - Operação de Aquisição Vinculada Exclusivamente a Receita de Exportação 63 - Crédito Presumido - Operação de Aquisição Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno 64 - Crédito Presumido - Operação de Aquisição Vinculada a Receitas Tributadas no Mercado Interno e de Exportação 65 - Crédito Presumido - Operação de Aquisição Vinculada a Receitas Não-Tributadas no Mercado Interno e de Exportação 66 - Crédito Presumido - Operação de Aquisição Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno, e de Exportação 67 - Crédito Presumido - Outras Operações 70 - Operação de Aquisição sem Direito a Crédito 71 - Operação de Aquisição com Isenção 72 - Operação de Aquisição com Suspensão 73 - Operação de Aquisição a Alíquota Zero 74 - Operação de Aquisição sem Incidência da Contribuição 75 - Operação de Aquisição por Substituição Tributária 98 - Outras Operações de Entrada 99 - Outras Operações.",
	"CSTIS":                       "Código Situação Tributária do Imposto Seletivo",
	"CSTReg":                      "Informar qual seria o CST caso não cumprida a condição resolutória/suspensiva",
	"DFeReferenciado":             "Referenciamento de item de outros DFe",
	"DI":                          "Declaração de Importação (NT 2011/004)",
	"EXTIPI":                      "Código EX TIPI (3 posições)",
	"IBSCBS":                      "Grupo de informações dos tributos IBS, CBS e Imposto Seletivo",
	"IBSCBSTot":                   "Valores totais da NF com IBS / CBS",
	"ICMSTot":                     "Totais referentes ao ICMS",
	"ICMSUFDest":                  "Grupo a ser informado nas vendas interestarduais para consumidor final, não contribuinte de ICMS",
	"IE":                          "Inscrição Estadual (v2.0)",
	"IEST":                        "Inscricao Estadual do Substituto Tributário",
	"IM":                          "Inscrição Municipal do tomador do serviço",
	"IPI":                         "Informação de IPI devolvido",
	"IS":                          "Grupo de informações do Imposto Seletivo",
	"ISSQNtot":                    "Totais referentes ao ISSQN",
	"ISTot":                       "Valores totais da NF com Imposto Seletivo",
	"ISUF":                        "Inscrição na SUFRAMA (Obrigatório nas operações com as áreas com benefícios de incentivos fiscais sob controle da SUFRAMA) PL_005d - 11/08/09 - alterado para aceitar 8 ou 9 dígitos",
	"NCM":                         "Código NCM (8 posições), será permitida a informação do gênero (posição do capítulo do NCM) quando a operação não for de comércio exterior (importação/exportação) ou o produto não seja tributado pelo IPI. Em caso de item de serviço ou item que não tenham produto (Ex. transferência de crédito, crédito do ativo imobilizado, etc.), informar o código 00 (zeros) (v2.0)",
	"NFe":                         "Nota Fiscal Eletrônica",
	"NFref":                       "Grupo de infromações da NF referenciada",
	"NVE":                         "Nomenclatura de Valor aduaneio e Estatístico",
	"PIS":                         "Dados do PIS",
	"PISAliq":                     "Código de Situação Tributária do PIS.  01 – Operação Tributável - Base de Cálculo = Valor da Operação Alíquota Normal (Cumulativo/Não Cumulativo); 02 - Operação Tributável - Base de Calculo = Valor da Operação (Alíquota Diferenciada);",
	"PISNT":                       "Código de Situação Tributária do PIS. 04 - Operação Tributável - Tributação Monofásica - (Alíquota Zero); 06 - Operação Tributável - Alíquota Zero; 07 - Operação Isenta da contribuição; 08 - Operação Sem Incidência da contribuição; 09 - Operação com suspensão da contribuição;",
	"PISOutr":                     "Código de Situação Tributária do PIS. 99 - Outras Operações.",
	"PISQtde":                     "Código de Situação Tributária do PIS. 03 - Operação Tributável - Base de Calculo = Quantidade Vendida x Alíquota por Unidade de Produto;",
	"PISST":                       "Dados do PIS Substituição Tributária",
	"RNTC":                        "Registro Nacional de Transportador de Carga (ANTT)",
	"TAjusteCompet":               "Tipo Ajuste de Competência",
	"TCIBS":                       "Tipo CBS IBS Completo",
	"TCompraGov":                  "Cada DFe que utilizar deverá utilizar esses tipo no grupo ide",
	"TCompraGovReduzido":          "Cada DFe que utilizar deverá utilizar esses tipo no grupo ide",
	"TConsReciNFe":                "Tipo Pedido de Consulta do Recido do Lote de Notas Fiscais Eletrônicas",
	"TCredPres":                   "Tipo Crédito Presumido",
	"TCredPresIBSZFM":             "Tipo Informações do crédito presumido de IBS para fornecimentos a partir da ZFM",
	"TCredPresOper":               "Tipo Crédito Presumido da Operação",
	"TDevTrib":                    "Tipo Devolução Tributo",
	"TDif":                        "Tipo Diferimento",
	"TEnderEmi":                   "Tipo Dados do Endereço do Emitente  // 24/10/08 - desmembrado / tamanho mínimo",
	"TEndereco":                   "Tipo Dados do Endereço  // 24/10/08 - tamanho mínimo",
	"TEnviNFe":                    "Tipo Pedido de Concessão de Autorização da Nota Fiscal Eletrônica",
	"TEstornoCred":                "Tipo Estorno de Crédito",
	"TIBSCBSMonoTot":              "Grupo de informações de totais da CBS/IBS com monofasia",
	"TIBSCBSTot":                  "Grupo de informações de totais da CBS/IBS",
	"TIS":                         "Grupo de informações do Imposto Seletivo",
	"TISTot":                      "Grupo de informações de totais do Imposto Seletivo",
	"TInfRespTec":                 "Grupo de informações do responsável técnico pelo sistema de emissão de DF-e",
	"TIpi":                        "Tipo: Dados do IPI",
	"TLocal":                      "Tipo Dados do Local de Retirada ou Entrega // 24/10/08 - tamanho mínimo // v2.0",
	"TMonofasia":                  "Tipo Monofasia",
	"TNFe":                        "Tipo Nota Fiscal Eletrônica",
	"TNfeProc":                    "Tipo da NF-e processada",
	"TProtNFe":                    "Tipo Protocolo de status resultado do processamento da NF-e",
	"TRed":                        "Tipo Redução Base de Cálculo",
	"TRetConsReciNFe":             "Tipo Retorno do Pedido de  Consulta do Recido do Lote de Notas Fiscais Eletrônicas",
	"TRetEnviNFe":                 "Tipo Retorno do Pedido de Autorização da Nota Fiscal Eletrônica",
	"TTransfCred":                 "Tipo Transferência de Crédito",
	"TTribBPe":                    "Grupo de informações da Tributação do BPe",
	"TTribCTe":                    "Grupo de informações da Tributação do CTe",
	"TTribCompraGov":              "Tipo Tributação Compra Governamental",
	"TTribNF3e":                   "Grupo de informações da Tributação da NF3e",
	"TTribNFAg":                   "Grupo de informações da Tributação da NFAg",
	"TTribNFCe":                   "Grupo de informações da Tributação da NFCe",
	"TTribNFCom":                  "Grupo de informações da Tributação da NFCom",
	"TTribNFGas":                  "Grupo de informações da Tributação da NFGas",
	"TTribNFe":                    "Grupo de informações da Tributação da NFe",
	"TTribRegular":                "Tipo Tributação Regular",
	"TVeiculo":                    "Tipo Dados do Veículo",
	"Type:TAmb":                   "Tipo Ambiente",
	"Type:TCListServ":             "Tipo Código da Lista de Serviços LC 116/2003",
	"Type:TCOrgaoIBGE":            "Tipo Código de orgão (UF da tabela do IBGE + 90 RFB)",
	"Type:TCST":                   "Código Situação Tributária do IBS/CBS",
	"Type:TChNFe":                 "Tipo Chave da Nota Fiscal Eletrônica",
	"Type:TCnpj":                  "Tipo Número do CNPJ",
	"Type:TCnpjOpc":               "Tipo Número do CNPJ Opcional",
	"Type:TCnpjVar":               "Tipo Número do CNPJ tmanho varíavel (3-14)",
	"Type:TCodMunIBGE":            "Tipo Código do Município da tabela do IBGE",
	"Type:TCodUfIBGE":             "Tipo Código da UF da tabela do IBGE",
	"Type:TCompetApur":            "Ano e mês referência do período de apuração (AAAA-MM)",
	"Type:TCpf":                   "Tipo Número do CPF",
	"Type:TCpfVar":                "Tipo Número do CPF de tamanho variável (3-11)",
	"Type:TData":                  "Tipo data AAAA-MM-DD",
	"Type:TDateTimeUTC":           "Data e Hora, formato UTC (AAAA-MM-DDThh:mm:ssTZD, onde TZD = +hh:mm ou -hh:mm)",
	"Type:TDec1104RTC":            "Tipo Decimal com 15 dígitos, sendo 11 de corpo e 4 decimais",
	"Type:TDec1302RTC":            "Tipo Decimal com 15 dígitos, sendo 13 de corpo e 2 decimais",
	"Type:TDec_0104v":             "Tipo Decimal com até 1 dígitos inteiros, podendo ter de 1 até 4 decimais",
	"Type:TDec_0204v":             "Tipo Decimal com até 2 dígitos inteiros, podendo ter de 1 até 4 decimais",
	"Type:TDec_0302Max100":        "Tipo Decimal com 3 inteiros (no máximo 100), com 2 decimais",
	"Type:TDec_0302_04RTC":        "Tipo Decimal com até 3 dígitos inteiros, podendo ter de 2 até 4 decimais",
	"Type:TDec_0302a04":           "Tipo Decimal com até 3 dígitos inteiros, podendo ter de 2 até 4 decimais",
	"Type:TDec_0302a04Max100":     "Tipo Decimal com 3 inteiros (no máximo 100), com até 4 decimais",
	"Type:TDec_0302a04Opc":        "Tipo Decimal com até 3 dígitos inteiros e 2 até 4 decimais. Utilizados em TAGs opcionais, não aceita valor zero.",
	"Type:TDec_0304Max100":        "Tipo Decimal com 3 inteiros (no máximo 100), com 4 decimais",
	"Type:TDec_03v00a04Max100Opc": "Tipo Decimal com 3 inteiros (no máximo 100), com 4 decimais, não aceita valor zero",
	"Type:TDec_0803v":             "Tipo Decimal com 8 inteiros, podendo ter de 1 até 3 decimais",
	"Type:TDec_1104":              "Tipo Decimal com 11 inteiros, podendo ter 4 decimais",
	"Type:TDec_1104OpRTC":         "Tipo Decimal com 11 inteiros, podendo ter 4 decimais (utilizado em tags opcionais)",
	"Type:TDec_1104Opc":           "Tipo Decimal com 11 inteiros, podendo ter 4 decimais (utilizado em tags opcionais)",
	"Type:TDec_1104v":             "Tipo Decimal com 11 inteiros, podendo ter de 1 até 4 decimais",
	"Type:TDec_1110v":             "Tipo Decimal com 11 inteiros, podendo ter de 1 até 10 decimais",
	"Type:TDec_1203":              "Tipo Decimal com 12 inteiros, podendo ter  3 decimais",
	"Type:TDec_1204":              "Tipo Decimal com 12 inteiros e 4 decimais",
	"Type:TDec_1204Opc":           "Tipo Decimal com 12 inteiros com 1 até 4 decimais",
	"Type:TDec_1204temperatura":   "Tipo Decimal com 12 inteiros, 1 a 4 decimais",
	"Type:TDec_1204v":             "Tipo Decimal com 12 inteiros de 1 até 4 decimais",
	"Type:TDec_1302":              "Tipo Decimal com 15 dígitos, sendo 13 de corpo e 2 decimais",
	"Type:TDec_1302Opc":           "Tipo Decimal com 15 dígitos, sendo 13 de corpo e 2 decimais, utilizado em tags opcionais",
	"Type:TEnteGov":               "Tipo de Ente Governamental",
	"Type:TFinNFe":                "Tipo Finalidade da NF-e (1=Normal; 2=Complementar; 3=Ajuste; 4=Devolução/Retorno)",
	"Type:TGuid":                  "Identificador único (Globally Unique Identifier)",
	"Type:TIdLote":                "Tipo Identificação de Lote",
	"Type:TIe":                    "Tipo Inscrição Estadual do Emitente // alterado EM 24/10/08 para aceitar ISENTO",
	"Type:TIeDest":                "Tipo Inscrição Estadual do Destinatário // alterado para aceitar vazio ou ISENTO - maio/2010 v2.0",
	"Type:TIeDestNaoIsento":       "Tipo Inscrição Estadual do Destinatário // alterado para aceitar vazio ou ISENTO - maio/2010 v2.0",
	"Type:TIeST":                  "Tipo Inscrição Estadual do ST // acrescentado EM 24/10/08",
	"Type:TIndDoacao":             "Tipo Indicador de Doação",
	"Type:TJust":                  "Tipo Justificativa",
	"Type:TMed":                   "Tipo temp médio em segundos",
	"Type:TMod":                   "Tipo Modelo Documento Fiscal",
	"Type:TMotivo":                "Tipo Motivo",
	"Type:TNF":                    "Tipo Número do Documento Fiscal",
	"Type:TOperCompraGov":         "Tipo da Operação com Ente Governamental",
	"Type:TProcEmi":               "Tipo processo de emissão da NF-e",
	"Type:TProt":                  "Tipo Número do Protocolo de Status",
	"Type:TRec":                   "Tipo Número do Recibo do envio de lote de NF-e",
	"Type:TSerie":                 "Tipo Série do Documento Fiscal",
	"Type:TServ":                  "Tipo Serviço solicitado",
	"Type:TStat":                  "Tipo Código da Mensagem enviada",
	"Type:TString":                "Tipo string genérico",
	"Type:TStringRTC":             "Tipo string genérico",
	"Type:TTime":                  "Tipo hora HH:MM:SS // tipo acrescentado na v2.0",
	"Type:TTpCredPresIBSZFM":      "Tipo de classificação do Crédito Presumido IBS ZFM",
	"Type:TTpNFCredito":           "Tipo de Nota de Crédito:  01=Multa e juros;  02=Apropriação de crédito presumido de IBS sobre o saldo devedor na ZFM (art. 450, § 1º, LC 214/25); 03=Retorno por recusa na entrega ou por não localização do destinatário na tentativa de entrega; 04=Redução de valores; 05=Transferência de crédito na sucessão;",
	"Type:TTpNFDebito":            "Tipo de Nota de Débito:  01=Transferência de créditos para Cooperativas;  02=Anulação de Crédito por Saídas Imunes/Isentas;  03=Débitos de notas fiscais não processadas na apuração;  04=Multa e juros;  05=Transferência de crédito na sucessão;  06=Pagamento antecipado;  07=Perda em estoque;  08=Desenquadramento do SN;",
	"Type:TUf":                    "Tipo Sigla da UF",
	"Type:TUfEmi":                 "Tipo Sigla da UF de emissor // acrescentado em 24/10/08",
	"Type:TVerAplic":              "Tipo Versão do Aplicativo",
	"Type:TVerNFe":                "Tipo Versão da NF-e - 4.00",
	"Type:Tano":                   "Tipo ano",
	"Type:TcClassTrib":            "Código de Classificação Tributária do IBS e da CBS",
	"Type:TcCredPres":             "Código de Classificação do Crédito Presumido do IBS e da CBS, conforme tabela cCredPres",
	"Type:Torig":                  "Tipo Origem da mercadoria CST ICMS  origem da mercadoria: 0-Nacional exceto as indicadas nos códigos 3, 4, 5 e 8;1-Estrangeira - Importação direta; 2-Estrangeira - Adquirida no mercado interno; 3-Nacional, conteudo superior 40% e inferior ou igual a 70%; 4-Nacional, processos produtivos básicos; 5-Nacional, conteudo inferior 40%; 6-Estrangeira - Importação direta, com similar nacional, lista CAMEX; 7-Estrangeira - mercado interno, sem simular,lista CAMEX;8-Nacional, Conteúdo de Importação superior a 70%.",
	"UF":                          "Sigla da UF",
	"UFDesemb":                    "UF onde ocorreu o desembaraço aduaneiro",
	"UFSaidaPais":                 "Sigla da UF de Embarque ou de transposição de fronteira",
	"UFTerceiro":                  "Sigla da UF do adquirente ou do encomendante",
	"adRemCBS":                    "Alíquota ad rem da CBS",
	"adRemCBSRet":                 "Alíquota ad rem da CBS retida anteriormente",
	"adRemCBSReten":               "Alíquota ad rem da CBS sujeita a retenção",
	"adRemIBS":                    "Alíquota ad rem do IBS",
	"adRemIBSRet":                 "Alíquota ad rem do IBS retido anteriormente",
	"adRemIBSReten":               "Alíquota ad rem do IBS sujeito a retenção",
	"adi":                         "Adições (NT 2011/004)",
	"agropecuario":                "Produtos Agropecurários Animais, Vegetais e Florestais",
	"autXML":                      "Pessoas autorizadas para o download do XML da NF-e",
	"avulsa":                      "Emissão de avulsa, informar os dados do Fisco emitente",
	"cAut":                        "Número de autorização da operação com cartões, PIX, boletos e outros pagamentos eletrônicos",
	"cBarra":                      "Codigo de barras diferente do padrão GTIN",
	"cBarraTrib":                  "Código de barras da unidade tributável diferente do padrão GTIN",
	"cClassTribReg":               "Informar qual seria o cClassTrib caso não cumprida a condição resolutória/suspensiva",
	"cCredPres":                   "Código de Classificação do Crédito Presumido do IBS e da CBS",
	"cCredPresumido":              "Código de Benefício Fiscal de Crédito Presumido na UF aplicado ao item",
	"cDV":                         "Digito Verificador da Chave de Acesso da NF-e",
	"cEAN":                        "GTIN (Global Trade Item Number) do produto, antigo código EAN ou código de barras",
	"cEANTrib":                    "GTIN (Global Trade Item Number) da unidade tributável, antigo código EAN ou código de barras",
	"cEnq":                        "Código de Enquadramento Legal do IPI (tabela a ser criada pela RFB)",
	"cExportador":                 "Código do exportador (usado nos sistemas internos de informação do emitente da NF-e)",
	"cFabricante":                 "Código do fabricante estrangeiro (usado nos sistemas internos de informação do emitente da NF-e)",
	"cMun":                        "Código do município (utilizar a tabela do IBGE)",
	"cMunFG":                      "Código do Município de Ocorrência do Fato Gerador (utilizar a tabela do IBGE)",
	"cMunFGIBS":                   "Informar o município de ocorrência do fato gerador do fato gerador do IBS / CBS. Campo preenchido somente quando “indPres = 5 (Operação presencial, fora do estabelecimento) ”, e não tiver endereço do destinatário (Grupo: E05) ou local de entrega (Grupo: G01).",
	"cNF":                         "Código numérico que compõe a Chave de Acesso. Número aleatório gerado pelo emitente para cada NF-e.",
	"cOperNFF":                    "Código da operação selecionada na NFF e relacionada ao item",
	"cPais":                       "Código de Pais",
	"cProd":                       "Código do produto ou serviço. Preencher com CFOP caso se trate de itens não relacionados com mercadorias/produto e que o contribuinte não possua codificação própria Formato ”CFOP9999”.",
	"cProdFisco":                  "Código Fiscal do Produto",
	"cRegTrib":                    "Código do regime especial de tributação",
	"cSelo":                       "Código do selo de controle do IPI",
	"cStat":                       "Código do status da mensagem enviada.",
	"cUF":                         "código da UF de atendimento",
	"cana":                        "Informações de registro aquisições de cana",
	"card":                        "Grupo de Cartões, PIX, Boletos e outros Pagamentos Eletrônicos",
	"chNFe":                       "Chaves de acesso da NF-e, compostas por: UF do emitente, AAMM da emissão da NFe, CNPJ do emitente, modelo, série e número da NF-e e código numérico+DV.",
	"chaveAcesso":                 "Chave de Acesso do DFe referenciado",
	"cobr":                        "Dados da cobrança da NF-e",
	"competApur":                  "Ano e mês referência do período de apuração (AAAA-MM)",
	"compra":                      "Informações de compras  (Nota de Empenho, Pedido e Contrato)",
	"dCompet":                     "Data da prestação do serviço  (AAAA-MM-DD)",
	"dDI":                         "Data de registro da DI/DSI/DA (AAAA-MM-DD)",
	"dDesemb":                     "Data do desembaraço aduaneiro (AAAA-MM-DD)",
	"dEmi":                        "Data de emissão do DAR (AAAA-MM-DD)",
	"dFab":                        "Data de fabricação/produção. Formato \"AAAA-MM-DD\".",
	"dPag":                        "Data do Pagamento",
	"dPrevEntrega":                "Data da previsão de entrega ou disponibilização do bem (AAAA-MM-DD)",
	"dVal":                        "Data de validade. Informar o último dia do mês caso a validade não especifique o dia. Formato \"AAAA-MM-DD\".",
	"dVenc":                       "Data de vencimento da duplicata (AAAA-MM-DD)",
	"deduc":                       "Deduções - Taxas e Contribuições",
	"defensivo":                   "Defensivo Agrícola / Agrotóxico",
	"dest":                        "Identificação do Destinatário",
	"det":                         "Dados dos detalhes da NF-e",
	"detExport":                   "Detalhe da exportação",
	"detPag":                      "Grupo de detalhamento da forma de pagamento.",
	"dhEmi":                       "Data e Hora de emissão do Documento Fiscal (AAAA-MM-DDThh:mm:ssTZD) ex.: 2012-09-01T13:00:00-03:00",
	"dhRecbto":                    "Data e hora de processamento, no formato AAAA-MM-DDTHH:MM:SSTZD. Em caso de Rejeição, com data e hora do recebimento do Lote de NF-e enviado.",
	"dhSaiEnt":                    "Data e Hora da saída ou de entrada da mercadoria / produto (AAAA-MM-DDTHH:mm:ssTZD)",
	"digVal":                      "Digest Value da NF-e processada. Utilizado para conferir a integridade da NF-e original.",
	"dup":                         "Dados das duplicatas NT 2011/004",
	"email":                       "Informar o e-mail da pessoa a ser contatada na empresa desenvolvedora do sistema.",
	"emit":                        "Identificação do emitente",
	"enderDest":                   "Dados do endereço",
	"enderEmit":                   "Endereço do emitente",
	"entrega":                     "Identificação do Local de Entrega (informar apenas quando for diferente do endereço do destinatário)",
	"esp":                         "Espécie dos volumes transportados",
	"exportInd":                   "Exportação indireta",
	"exporta":                     "Informações de exportação",
	"fat":                         "Dados da fatura",
	"finNFe":                      "Finalidade da emissão da NF-e: 1 - NFe normal 2 - NFe complementar 3 - NFe de ajuste 4 - Devolução/Retorno 5 - Nota de crédito 6 - Nota de débito",
	"fone":                        "Informar o telefone da pessoa a ser contatada na empresa desenvolvedora do sistema. Preencher com o Código DDD + número do telefone.",
	"forDia":                      "Fornecimentos diários",
	"gCBS":                        "Grupo de Tributação da CBS",
	"gCBSCredPres":                "Grupo de Informações do Crédito Presumido referente a CBS, quando aproveitado pelo emitente do documento.",
	"gCompraGov":                  "Grupo de Compras Governamentais",
	"gCred":                       "Grupo de informações sobre o CréditoPresumido",
	"gDevTrib":                    "Grupo de Informações da devolução de tributos",
	"gDif":                        "Grupo de campos do Diferimento",
	"gEstornoCred":                "Totalização do estorno de crédito",
	"gIBS":                        "Totalização do IBS",
	"gIBSCredPres":                "Grupo de Informações do Crédito Presumido referente ao IBS, quando aproveitado pelo emitente do documento.",
	"gIBSMun":                     "Totalização do IBS de competência Municipal",
	"gIBSUF":                      "Totalização do IBS de competência da UF",
	"gMono":                       "Só deverá ser utilizado para DFe modelos 55 e 65",
	"gMonoDif":                    "Grupo de informações do diferimento da Tributação Monofásica",
	"gMonoPadrao":                 "Grupo de informações da Tributação Monofásica padrão",
	"gMonoRet":                    "Grupo de informações da Tributação Monofásica retida anteriormente",
	"gMonoReten":                  "Grupo de informações da Tributação Monofásica sujeita a retenção",
	"gPagAntecipado":              "Informado para abater as parcelas de antecipação de pagamento, conforme Art. 10. § 4º",
	"gRed":                        "Grupo de campos da redução de aliquota",
	"gTribCompraGov":              "Grupo de informações da composição do valor do IBS e da CBS em compras governamental",
	"gTribRegular":                "Grupo de informações da Tributação Regular. Informar como seria a tributação caso não cumprida a condição resolutória/suspensiva. Exemplo 1: Art. 442, §4. Operações com ZFM e ALC. Exemplo 2: Operações com suspensão do tributo.",
	"guiaTransito":                "Guias De Trânsito de produtos agropecurários animais, vegetais e de origem florestal.",
	"idCadIntTran":                "Identificador cadastrado no intermediador",
	"idDest":                      "Identificador de Local de destino da operação (1-Interna;2-Interestadual;3-Exterior)",
	"idTermPag":                   "Identificador do terminal de pagamento",
	"ide":                         "identificação da NF-e",
	"imposto":                     "Tributos incidentes nos produtos ou serviços da NF-e",
	"indBemMovelUsado":            "Indicador de fornecimento de bem móvel usado: 1-Bem Móvel Usado",
	"indDoacao":                   "Indica se a operação é de doação",
	"indFinal":                    "Indica operação com consumidor final (0-Não;1-Consumidor Final)",
	"indIEDest":                   "Indicador da IE do destinatário: 1 – Contribuinte ICMSpagamento à vista; 2 – Contribuinte isento de inscrição; 9 – Não Contribuinte",
	"indIntermed":                 "Indicador de intermediador/marketplace  0=Operação sem intermediador (em site ou plataforma própria)  1=Operação em site ou plataforma de terceiros (intermediadores/marketplace)",
	"indPag":                      "Indicador da Forma de Pagamento:0-Pagamento à Vista;1-Pagamento à Prazo;",
	"indPres":                     "Indicador de presença do comprador no estabelecimento comercial no momento da oepração (0-Não se aplica (ex.: Nota Fiscal complementar ou de ajuste;1-Operação presencial;2-Não presencial, internet;3-Não presencial, teleatendimento;4-NFC-e entrega em domicílio;5-Operação presencial, fora do estabelecimento;9-Não presencial, outros)",
	"indProc":                     "Origem do processo, informar com: 0 - SEFAZ; 1 - Justiça Federal; 2 - Justiça Estadual; 3 - Secex/RFB; 4 - CONFAZ; 9 - Outros.",
	"indSinc":                     "Indicador de processamento síncrono. 0=NÃO; 1=SIM=Síncrono",
	"indSomaCOFINSST":             "Indica se o valor da COFINS ST compõe o valor total da NFe",
	"indSomaPISST":                "Indica se o valor do PISST compõe o valor total da NF-e",
	"indTot":                      "Este campo deverá ser preenchido com:  0 – o valor do item (vProd) não compõe o valor total da NF-e (vProd)  1  – o valor do item (vProd) compõe o valor total da NF-e (vProd)",
	"infAdFisco":                  "Informações adicionais de interesse do Fisco (v2.0)",
	"infAdProd":                   "Informações adicionais do produto (norma referenciada, informações complementares, etc)",
	"infAdic":                     "Informações adicionais da NF-e",
	"infCpl":                      "Informações complementares de interesse do Contribuinte",
	"infIntermed":                 "Grupo de Informações do Intermediador da Transação",
	"infNFe":                      "Informações da Nota Fiscal eletrônica",
	"infNFeSupl":                  "Informações suplementares Nota Fiscal",
	"infProdEmb":                  "Informações mais detalhadas do produto (usada na NFF)",
	"infProdNFF":                  "Informações mais detalhadas do produto (usada na NFF)",
	"infProt":                     "Dados do protocolo de status",
	"infRespTec":                  "Informações do Responsável Técnico pela emissão do DF-e",
	"infSolicNFF":                 "Grupo para informações da solicitação da NFF",
	"marca":                       "Marca dos volumes transportados",
	"matr":                        "Matrícula do agente",
	"mod":                         "Código do modelo do Documento Fiscal  Preencher com \"2B\", quando se tratar de Cupom Fiscal emitido por máquina registradora (não ECF), com \"2C\", quando se tratar de Cupom Fiscal PDV, ou \"2D\", quando se tratar de Cupom Fiscal (emitido por ECF)",
	"modFrete":                    "Modalidade do frete 0- Contratação do Frete por conta do Remetente (CIF); 1- Contratação do Frete por conta do destinatário/remetente (FOB); 2- Contratação do Frete por conta de terceiros; 3- Transporte próprio por conta do remetente; 4- Transporte próprio por conta do destinatário; 9- Sem Ocorrência de transporte.",
	"nAdicao":                     "Número da Adição",
	"nCOO":                        "Informar o Número do Contador de Ordem de Operação - COO vinculado à NF-e",
	"nDAR":                        "Número do Documento de Arrecadação de Receita",
	"nDI":                         "Número do Documento de Importação (DI, DSI, DIRE, DUImp) (NT2011/004)",
	"nDraw":                       "Número do ato concessório de Drawback",
	"nDup":                        "Número da duplicata",
	"nECF":                        "Informar o número de ordem seqüencial do ECF que emitiu o Cupom Fiscal vinculado à NF-e",
	"nFCI":                        "Número de controle da FCI - Ficha de Conteúdo de Importação.",
	"nFat":                        "Número da fatura",
	"nGuia":                       "Número da Guia",
	"nItem":                       "Número do item do documento referenciado. Corresponde ao atributo nItem do elemento det do documento original.",
	"nItemPed":                    "Número do Item do Pedido de Compra - Identificação do número do item do pedido de Compra",
	"nLacre":                      "Número dos Lacres",
	"nLote":                       "Número do lote do produto.",
	"nNF":                         "Número do Documento Fiscal - 1 – 999999999",
	"nProc":                       "Indentificador do processo ou ato concessório",
	"nProt":                       "Número do Protocolo de Status da NF-e. 1 posição (1 – Secretaria de Fazenda Estadual 2 – Receita Federal); 2 - códiga da UF - 2 posições ano; 10 seqüencial no ano.",
	"nRE":                         "Registro de exportação",
	"nRec":                        "Número do Recibo Consultado",
	"nReceituario":                "Número do Receituário ou Receita do Defensivo / Agrotóxico",
	"nSeqAdic":                    "Número seqüencial do item",
	"nVol":                        "Numeração dos volumes transportados",
	"natOp":                       "Descrição da Natureza da Operação",
	"nro":                         "Número",
	"obsCont":                     "Campo de uso livre do contribuinte informar o nome do campo no atributo xCampo e o conteúdo do campo no xTexto",
	"obsFisco":                    "Campo de uso exclusivo do Fisco informar o nome do campo no atributo xCampo e o conteúdo do campo no xTexto",
	"obsItem":                     "Grupo de observações de uso livre (para o item da NF-e)",
	"pAliqEfet":                   "Aliquota Efetiva que será aplicada a Base de Calculo (em percentual)",
	"pAliqEfetRegCBS":             "Informar como seria a Alíquota caso não cumprida a condição resolutória/suspensiva",
	"pAliqEfetRegIBSMun":          "Informar como seria a Alíquota caso não cumprida a condição resolutória/suspensiva",
	"pAliqEfetRegIBSUF":           "Informar como seria a Alíquota caso não cumprida a condição resolutória/suspensiva",
	"pCBS":                        "Aliquota da CBS (em percentual)",
	"pCOFINS":                     "Alíquota do COFINS (em percentual)",
	"pCredPres":                   "Percentual do Crédito Presumido",
	"pCredPresumido":              "Percentual do Crédito Presumido",
	"pDevol":                      "Percentual de mercadoria devolvida",
	"pDif":                        "Percentual do diferimento",
	"pDifCBS":                     "Percentual do diferimento do imposto monofásico",
	"pDifIBS":                     "Percentual do diferimento do imposto monofásico",
	"pFCPUFDest":                  "Percentual adicional inserido na alíquota interna da UF de destino, relativo ao Fundo de Combate à Pobreza (FCP) naquela UF.",
	"pICMSInter":                  "Alíquota interestadual das UF envolvidas: - 4% alíquota interestadual para produtos importados; - 7% para os Estados de origem do Sul e Sudeste (exceto ES), destinado para os Estados do Norte e Nordeste  ou ES; - 12% para os demais casos.",
	"pICMSInterPart":              "Percentual de partilha para a UF do destinatário: - 40% em 2016; - 60% em 2017; - 80% em 2018; - 100% a partir de 2019.",
	"pICMSRet":                    "Alíquota da Retenção",
	"pICMSUFDest":                 "Alíquota adotada nas operações internas na UF do destinatário para o produto / mercadoria.",
	"pPIS":                        "Alíquota do PIS (em percentual)",
	"pRedAliq":                    "Percentual de redução de aliquota do cClassTrib",
	"pRedutor":                    "Percentual de redução de aliquota em compra governamental",
	"pag":                         "Dados de Pagamento. Obrigatório apenas para (NFC-e) NT 2012/004",
	"pesoB":                       "Peso bruto (em kg)",
	"pesoL":                       "Peso líquido (em kg)",
	"placa":                       "Placa do veículo (NT2011/004)",
	"procEmi":                     "Processo de emissão utilizado com a seguinte codificação: 0 - emissão de NF-e com aplicativo do contribuinte; 1 - emissão de NF-e avulsa pelo Fisco; 2 - emissão de NF-e avulsa, pelo contribuinte com seu certificado digital, através do site do Fisco; 3- emissão de NF-e pelo contribuinte com aplicativo fornecido pelo Fisco.",
	"procRef":                     "Grupo de informações do  processo referenciado",
	"prod":                        "Dados dos produtos e serviços da NF-e",
	"protNFe":                     "Protocolo de status resultado do processamento da NF-e",
	"qBCMono":                     "Valor total da quantidade tributada do ICMS monofásico próprio",
	"qBCMonoRet":                  "Valor total da quantidade tributada do ICMS monofásico retido anteriormente",
	"qBCMonoReten":                "Valor total da quantidade tributada do ICMS monofásico sujeito a retenção",
	"qBCProd":                     "Quantidade Vendida (NT2011/004)",
	"qCom":                        "Quantidade Comercial  do produto, alterado para aceitar de 0 a 4 casas decimais e 11 inteiros.",
	"qExport":                     "Quantidade do item efetivamente exportado",
	"qLote":                       "Quantidade de produto no lote.",
	"qSelo":                       "Quantidade de selo de controle do IPI",
	"qTotAnt":                     "Total Anterior",
	"qTotGer":                     "Total Geral",
	"qTotMes":                     "Total do mês",
	"qTrib":                       "Quantidade Tributável - alterado para aceitar de 0 a 4 casas decimais e 11 inteiros",
	"qVol":                        "Quantidade de volumes transportados",
	"qVolEmb":                     "Volume do produto na embalagem",
	"qrCode":                      "Texto com o QR-Code impresso no DANFE NFC-e",
	"qtde":                        "Quantidade em quilogramas - peso líquido",
	"ref":                         "Mês e Ano de Referência, formato: MM/AAAA",
	"refCTe":                      "Utilizar esta TAG para referenciar um CT-e emitido anteriormente, vinculada a NF-e atual",
	"refECF":                      "Grupo do Cupom Fiscal vinculado à NF-e",
	"refNF":                       "Dados da NF modelo 1/1A referenciada ou NF modelo 2 referenciada",
	"refNFP":                      "Grupo com as informações NF de produtor referenciada",
	"refNFe":                      "Chave de acesso da NF-e de antecipação de pagamento",
	"refNFeSig":                   "Referencia uma NF-e (modelo 55) emitida anteriormente pela sua Chave de Acesso com código numérico zerado, permitindo manter o sigilo da NF-e referenciada.",
	"repEmi":                      "Repartição Fiscal emitente",
	"retTransp":                   "Dados da retenção  ICMS do Transporte",
	"retTrib":                     "Retenção de Tributos Federais",
	"retirada":                    "Identificação do Local de Retirada (informar apenas quando for diferente do endereço do remetente)",
	"safra":                       "Identificação da safra",
	"serie":                       "Série do Documento Fiscal, informar zero se inexistentesérie",
	"serieGuia":                   "Série da Guia",
	"tBand":                       "Bandeira da operadora de cartão",
	"tPag":                        "Forma de Pagamento:",
	"total":                       "Dados dos totais da NF-e",
	"tpAmb":                       "Identificação do Ambiente: 1 - Produção 2 - Homologação",
	"tpAto":                       "Tipo do ato concessório Para origem do Processo na SEFAZ (indProc=0), informar o tipo de ato concessório: 08 - Termo de Acordo; 10 - Regime Especial; 12 - Autorização específica; 14 - Ajuste SINIEF; 15 - Convênio ICMS.",
	"tpCredPresIBSZFM":            "Classificação para subapuração do IBS na ZFM",
	"tpEmis":                      "Forma de emissão da NF-e 1 - Normal; 2 - Contingência FS 3 - Regime Especial NFF (NT 2021.002) 4 - Contingência DPEC 5 - Contingência FSDA 6 - Contingência SVC - AN 7 - Contingência SVC - RS 9 - Contingência off-line NFC-e",
	"tpEnteGov":                   "Para administração pública direta e suas autarquias e fundações: 1=União 2=Estados 3=Distrito Federal 4=Municípios",
	"tpGuia":                      "Tipo da Guia: 1 - GTA; 2 - TTA; 3 - DTA; 4 - ATV; 5 - PTV; 6 - GTV; 7 - Guia Florestal (DOF, SisFlora - PA e MT, SIAM - MG)",
	"tpImp":                       "Formato de impressão do DANFE (0-sem DANFE;1-DANFe Retrato; 2-DANFe Paisagem;3-DANFe Simplificado; 4-DANFe NFC-e;5-DANFe NFC-e em mensagem eletrônica)",
	"tpIntegra":                   "Tipo de Integração do processo de pagamento com o sistema de automação da empresa: 1 - Pagamento integrado com o sistema de automação da empresa (Ex.: equipamento TEF, Comércio Eletrônico, POS Integrado); 2 - Pagamento não integrado com o sistema de automação da empresa (Ex.: equipamento POS Simples).",
	"tpIntermedio":                "Forma de Importação quanto a intermediação  1-por conta propria;2-por conta e ordem;3-encomenda",
	"tpNF":                        "Tipo do Documento Fiscal (0 - entrada; 1 - saída)",
	"tpNFCredito":                 "Tipo de Nota de Crédito",
	"tpNFDebito":                  "Tipo de Nota de Débito",
	"tpOperGov":                   "Tipo da operação com ente governamental: 1 - Fornecimento 2 - Recebimento do Pagamento",
	"tpViaTransp":                 "Via de transporte internacional informada na DI ou na Declaração Única de Importação (DUImp): 1-Maritima;2-Fluvial;3-Lacustre;4-Aerea;5-Postal;6-Ferroviaria;7-Rodoviaria;8-Conduto;9-Meios Proprios;10-Entrada/Saida Ficta; 11-Courier;12-Em maos;13-Por reboque.",
	"transp":                      "Dados dos transportes da NF-e",
	"transporta":                  "Dados do transportador",
	"uCom":                        "Unidade comercial",
	"uEmb":                        "Unidade de Medida da Embalagem",
	"uTrib":                       "Unidade Tributável",
	"urlChave":                    "Informar a URL da \"Consulta por chave de acesso da NFC-e\". A mesma URL que deve estar informada no DANFE NFC-e para consulta por chave de acesso.",
	"vAFRMM":                      "Valor Adicional ao frete para renovação de marinha mercante",
	"vAliqProd":                   "Alíquota do COFINS (em reais) (NT2011/004)",
	"vBC":                         "Base de Cálculo do ISS",
	"vBCCredPres":                 "Valor da Base de Cálculo do Crédito Presumido da Operação",
	"vBCFCPUFDest":                "Valor da Base de Cálculo do FCP na UF do destinatário.",
	"vBCIBSCBS":                   "Total Base de Calculo",
	"vBCIRRF":                     "Base de Cálculo do IRRF",
	"vBCRet":                      "BC da Retenção do ICMS",
	"vBCRetPrev":                  "Base de Cálculo da Retenção da Previdêncica Social",
	"vBCST":                       "BC do ICMS ST",
	"vBCUFDest":                   "Valor da Base de Cálculo do ICMS na UF do destinatário.",
	"vCBS":                        "Valor da CBS",
	"vCBSEstCred":                 "Valor da CBS a ser estornada",
	"vCBSMono":                    "Valor da CBS monofásica",
	"vCBSMonoDif":                 "Valor da CBS monofásica diferida",
	"vCBSMonoRet":                 "Valor da CBS retida anteriormente",
	"vCBSMonoReten":               "Valor da CBS monofásica sujeita a retenção",
	"vCOFINS":                     "Valor do COFINS sobre serviços",
	"vCredPres":                   "Total do Crédito Presumido",
	"vCredPresCondSus":            "Total do Crédito Presumido Condição Suspensiva",
	"vCredPresIBSZFM":             "Valor do crédito presumido calculado sobre o saldo devedor apurado",
	"vCredPresumido":              "Valor do Crédito Presumido",
	"vDAR":                        "Valor Total constante no DAR",
	"vDed":                        "valor da dedução",
	"vDeducao":                    "Valor dedução para redução da base de cálculo",
	"vDesc":                       "Valor do desconto da fatura",
	"vDescCond":                   "Valor desconto condicionado",
	"vDescDI":                     "Valor do desconto do item",
	"vDescIncond":                 "Valor desconto incondicionado",
	"vDevTrib":                    "Valor do tributo devolvido. No fornecimento de energia elétrica, água, esgoto e gás natural e em outras hipóteses definidas no regulamento",
	"vDif":                        "Valor do diferimento",
	"vDup":                        "Valor da duplicata",
	"vFCP":                        "Valor Total do FCP (Fundo de Combate à Pobreza).",
	"vFCPST":                      "Valor Total do FCP (Fundo de Combate à Pobreza) retido por substituição tributária.",
	"vFCPSTRet":                   "Valor Total do FCP (Fundo de Combate à Pobreza) retido anteriormente por substituição tributária.",
	"vFCPUFDest":                  "Valor total do ICMS relativo ao Fundo de Combate à Pobreza (FCP) para a UF de destino.",
	"vFor":                        "Valor  dos fornecimentos",
	"vFrete":                      "Valor Total do Frete",
	"vIBS":                        "Valor do IBS",
	"vIBSEstCred":                 "Valor do IBS a ser estornado",
	"vIBSMono":                    "Valor do IBS monofásico",
	"vIBSMonoDif":                 "Valor do IBS monofásico diferido",
	"vIBSMonoRet":                 "Valor do IBS retido anteriormente",
	"vIBSMonoReten":               "Valor do IBS monofásico sujeito a retenção",
	"vIBSMun":                     "Valor total do IBS Municipal",
	"vIBSUF":                      "Valor total do IBS Estadual",
	"vICMS":                       "Valor Total do ICMS",
	"vICMSDeson":                  "Valor Total do ICMS desonerado",
	"vICMSMono":                   "Valor total do ICMS monofásico próprio",
	"vICMSMonoRet":                "Valor do ICMS monofásico retido anteriormente",
	"vICMSMonoReten":              "Valor total do ICMS monofásico sujeito a retenção",
	"vICMSRet":                    "Valor do ICMS Retido",
	"vICMSUFDest":                 "Valor total do ICMS de partilha para a UF do destinatário",
	"vICMSUFRemet":                "Valor total do ICMS de partilha para a UF do remetente",
	"vII":                         "Valor Total do II",
	"vIPI":                        "Valor Total do IPI",
	"vIPIDevol":                   "Valor Total do IPI devolvido. Deve ser informado quando preenchido o Grupo Tributos Devolvidos na emissão de nota finNFe=4 (devolução) nas operações com não contribuintes do IPI. Corresponde ao total da soma dos campos id: UA04.",
	"vIRRF":                       "Valor Retido de IRRF",
	"vIS":                         "Valor Total do Imposto Seletivo",
	"vISS":                        "Valor Total do ISS",
	"vISSRet":                     "Valor Total Retenção ISS",
	"vItem":                       "Valor total do Item, correspondente à sua participação no total da nota. A soma dos itens deverá corresponder ao total da nota.",
	"vLiq":                        "Valor líquido da fatura",
	"vLiqFor":                     "Valor Líquido dos fornecimentos",
	"vNF":                         "Valor Total da NF-e",
	"vNFTot":                      "Valor Total da NF considerando os impostos por fora IBS, CBS e IS",
	"vOrig":                       "Valor original da fatura",
	"vOutro":                      "Valor outras retenções",
	"vPIS":                        "Valor do PIS sobre serviços",
	"vPag":                        "Valor do Pagamento. Esta tag poderá ser omitida quando a tag tPag=90 (Sem Pagamento), caso contrário deverá ser preenchida.",
	"vProd":                       "Valor Total dos produtos e serviços",
	"vRetCOFINS":                  "Valor Retido de COFINS",
	"vRetCSLL":                    "Valor Retido de CSLL",
	"vRetPIS":                     "Valor Retido de PIS",
	"vRetPrev":                    "Valor da Retenção da Previdêncica Social",
	"vST":                         "Valor Total do ICMS ST",
	"vSeg":                        "Valor Total do Seguro",
	"vServ":                       "Valor do Serviço",
	"vTotCBSMonoItem":             "Total da CBS monofásica do item",
	"vTotDed":                     "Valor Total das Deduções",
	"vTotIBSMonoItem":             "Total de IBS monofásico do item",
	"vTotTrib":                    "Valor estimado total de impostos federais, estaduais e municipais",
	"vTribCBS":                    "Valor que seria devido a CBS, sem aplicação do Art. 473. da LC 214/2025",
	"vTribIBSMun":                 "Valor que seria devido ao município, sem aplicação do Art. 473. da LC 214/2025",
	"vTribIBSUF":                  "Valor que seria devido a UF, sem aplicação do Art. 473. da LC 214/2025",
	"vTribRegCBS":                 "Informar como seria o valor do Tributo caso não cumprida a condição resolutória/suspensiva",
	"vTribRegIBSMun":              "Informar como seria o valor do Tributo caso não cumprida a condição resolutória/suspensiva",
	"vTribRegIBSUF":               "Informar como seria o valor do Tributo caso não cumprida a condição resolutória/suspensiva",
	"vTroco":                      "Valor do Troco.",
	"vUnCom":                      "Valor unitário de comercialização  - alterado para aceitar 0 a 10 casas decimais e 11 inteiros",
	"vUnTrib":                     "Valor unitário de tributação - alterado para aceitar 0 a 10 casas decimais e 11 inteiros",
	"verAplic":                    "Versão do Aplicativo que processou a NF-e",
	"verProc":                     "versão do aplicativo utilizado no processo de emissão",
	"vol":                         "Dados dos volumes",
	"xAgente":                     "Nome do agente",
	"xBairro":                     "Bairro",
	"xCont":                       "Informação do contrato",
	"xContato":                    "Informar o nome da pessoa a ser contatada na empresa desenvolvedora do sistema utilizado na emissão do documento fiscal eletrônico.",
	"xCpl":                        "Complemento",
	"xDed":                        "Descrição da Dedução",
	"xEmb":                        "Embalagem do produto",
	"xEnder":                      "Endereço completo",
	"xFant":                       "Nome fantasia",
	"xLgr":                        "Logradouro",
	"xLocDesemb":                  "Local do desembaraço aduaneiro",
	"xLocDespacho":                "Descrição do local de despacho",
	"xLocExporta":                 "Local de Embarque ou de transposição de fronteira",
	"xMotivo":                     "Descrição literal do status do serviço solicitado.",
	"xMun":                        "Nome do município",
	"xNEmp":                       "Informação da Nota de Empenho de compras públicas (NT2011/004)",
	"xNome":                       "Razão Social ou Nome do Expedidor/Recebedor",
	"xOrgao":                      "Órgão emitente",
	"xPag":                        "Descrição do Meio de Pagamento",
	"xPais":                       "Nome do país",
	"xPed":                        "Informação do pedido",
	"xProd":                       "Descrição do produto ou serviço",
	"xSolic":                      "Solicitação do pedido de emissão da NFF",
}

// getFriendlyName é um helper para buscar no map com segurança.
// Se a tag não for encontrada, retorna a própria tag.
func getFriendlyName(tag string, tagMap map[string]string) string {
	if name, ok := tagMap[tag]; ok {
		return name
	}
	return tag // Retorna a tag crua se não houver tradução
}

// 2. O ANALISADOR DE PADRÕES (REGEX)

// ErrorPattern define uma estrutura para cada tipo de erro que queremos capturar.
// Cada padrão terá sua própria Regex e uma função Tradutora.
type ErrorPattern struct {
	Regex      *regexp.Regexp
	Translator func(matches []string, tagMap map[string]string) string
}

// newPatterns inicializa e compila todas as nossas Regex.
// É aqui que a "mágica" de identificação acontece.
func newPatterns() []ErrorPattern {
	return []ErrorPattern{
		{
			// Padrão: Elemento Obrigatório Faltando
			// Ex: "The element 'det' ... has incomplete content. List of possible elements expected: 'cProd'..."
			Regex: regexp.MustCompile(`The element '([^']*)'.* incomplete content.* expected: '([^']*)'`),
			Translator: func(matches []string, tagMap map[string]string) string {
				// matches[0] = String completa
				// matches[1] = Tag Pai (ex: 'det')
				// matches[2] = Tag Faltante (ex: 'cProd')
				parentTag := getFriendlyName(matches[1], tagMap)
				missingTag := getFriendlyName(matches[2], tagMap)
				return fmt.Sprintf("Erro em '%s': O campo obrigatório '%s' não foi preenchido.", parentTag, missingTag)
			},
		},
		{
			// Padrão: Ordem Incorreta dos Elementos
			// Ex: "The element 'enderEmit' ... is not expected. Expected element: 'IE'..."
			Regex: regexp.MustCompile(`The element '([^']*)'.* is not expected.* Expected element: '([^']*)'`),
			Translator: func(matches []string, tagMap map[string]string) string {
				// matches[1] = Tag Inesperada (ex: 'enderEmit')
				// matches[2] = Tag Esperada (ex: 'IE')
				unexpectedTag := getFriendlyName(matches[1], tagMap)
				expectedTag := getFriendlyName(matches[2], tagMap)
				return fmt.Sprintf("Erro de ordem: O campo '%s' foi informado antes do campo '%s'.", unexpectedTag, expectedTag)
			},
		},
		{
			// Padrão: Valor Inválido para Lista (Enumeração)
			// Ex: "The value '3' of element 'tpNF' is not valid..."
			Regex: regexp.MustCompile(`The value '([^']*)' of element '([^']*)' is not valid`),
			Translator: func(matches []string, tagMap map[string]string) string {
				// matches[1] = Valor Inválido (ex: '3')
				// matches[2] = Tag (ex: 'tpNF')
				badValue := matches[1]
				tag := getFriendlyName(matches[2], tagMap)
				return fmt.Sprintf("Erro de preenchimento: O valor '%s' não é uma opção válida para o campo '%s'.", badValue, tag)
			},
		},
		{
			// Padrão: Tipo de Dado Inválido (Formato)
			// (Assumindo um formato de erro que inclui a tag)
			// Ex: "Element 'vProd': The value '120,50' is invalid according to its data type 'TDec_1302'..."
			Regex: regexp.MustCompile(`Element '([^']*)': The value '([^']*)' is invalid`),
			Translator: func(matches []string, tagMap map[string]string) string {
				// matches[1] = Tag (ex: 'vProd')
				// matches[2] = Valor Inválido (ex: '120,50')
				tag := getFriendlyName(matches[1], tagMap)
				badValue := matches[2]

				// Bônus: Dica específica para o erro mais comum de todos
				advice := ""
				if strings.Contains(badValue, ",") {
					advice = " (Lembre-se: use ponto '.' para decimais, não vírgula)."
				}

				return fmt.Sprintf("Erro no campo '%s': O valor '%s' possui um formato inválido.%s", tag, badValue, advice)
			},
		},
		{
			// Padrão: Erro de Tamanho Mínimo (minLength)
			// Ex: "Element 'xNome': [minLength error] The value has a length of '1'; this is less than the required '2'."
			Regex: regexp.MustCompile(`Element '([^']*)'.*minLength.*length of '([^']*)'.*required '([^']*)'`),
			Translator: func(matches []string, tagMap map[string]string) string {
				// matches[1] = Tag (ex: 'xNome')
				// matches[2] = Tamanho Informado (ex: '1')
				// matches[3] = Tamanho Mínimo (ex: '2')
				tag := getFriendlyName(matches[1], tagMap)
				return fmt.Sprintf("Erro de tamanho: O campo '%s' deve ter no mínimo %s caracteres, mas foi informado %s.", tag, matches[3], matches[2])
			},
		},
		{
			// Padrão: Erro de Tamanho Máximo (maxLength)
			// Ex: "Element 'xProd': [maxLength error] The value has a length of '150'; this is greater than the allowed '120'."
			Regex: regexp.MustCompile(`Element '([^']*)'.*maxLength.*length of '([^']*)'.*allowed '([^']*)'`),
			Translator: func(matches []string, tagMap map[string]string) string {
				// matches[1] = Tag (ex: 'xProd')
				// matches[2] = Tamanho Informado (ex: '150')
				// matches[3] = Tamanho Máximo (ex: '120')
				tag := getFriendlyName(matches[1], tagMap)
				return fmt.Sprintf("Erro de tamanho: O campo '%s' deve ter no máximo %s caracteres, mas foi informado %s.", tag, matches[3], matches[2])
			},
		},
		{
			// Padrão: Erro em Atributo (não elemento)
			// Ex: "Value 'ABC' for attribute 'versao' on element 'infNFe' is not valid..."
			Regex: regexp.MustCompile(`Value '([^']*)' for attribute '([^']*)' on element '([^']*)' is not valid`),
			Translator: func(matches []string, tagMap map[string]string) string {
				// matches[1] = Valor Inválido (ex: 'ABC')
				// matches[2] = Atributo (ex: 'versao')
				// matches[3] = Tag Pai (ex: 'infNFe')

				// Para atributos, o ideal é mapear "Attr:versao"
				attrKey := "Attr:" + matches[2]
				tag := getFriendlyName(attrKey, tagMap)

				return fmt.Sprintf("Erro de atributo: O valor '%s' é inválido para o atributo '%s' (do grupo %s).", matches[1], tag, getFriendlyName(matches[3], tagMap))
			},
		},
		// Você pode adicionar quantos novos padrões quiser aqui...
	}
}

// 3. A LÓGICA PRINCIPAL

func main() {

	// Nossa lista de erros "crus" para testar
	rawErrors := []string{
		"Element 'vProd': The value '120,50' is invalid according to its data type 'TDec_1302'",
		"The element 'det' in namespace 'http://www.portalfiscal.inf.br/nfe' has incomplete content. List of possible elements expected: 'cProd' in namespace 'http://www.portalfiscal.inf.br/nfe'.",
		"The value '9' of element 'modFrete' is not valid. The element 'modFrete' with value '9' failed to match the pattern enumeration constraint. [Expected: 0, 1, 2, 3, 4].",
		"The element 'enderEmit' in namespace 'http://www.portalfiscal.inf.br/nfe' is not expected. Expected element: 'IE' in namespace 'http://www.portalfiscal.inf.br/nfe'.",
		"Some unknown error that won't be matched.",
		// Padrão: Valor Inválido para Lista (Enumeração)
		"The value '9' of element 'modFrete' is not valid. The element 'modFrete' with value '9' failed to match the pattern enumeration constraint. [Expected: 0, 1, 2, 3, 4].",
		"The value '3' of element 'tpNF' is not valid. The value '3' does not match one of the expected values: [0, 1].",
		"The value '7' of element 'indPres' is not valid according to its data type.",

		// Padrão: Tipo de Dado Inválido (Formato)
		"Element 'vProd': The value '120,50' is invalid according to its data type 'TDec_1302'",
		"Element 'dhEmi': The value '20/10/2025 10:30:00' is invalid according to its data type 'TDateTimeUTC'",
		"Element 'CEP': The value '90000-000' is invalid. Pattern constraint '[0-9]{8}' is not satisfied.",
		"Element 'vBC': The value 'ABC' is invalid according to its data type 'TDec_1302'.",

		// Padrão: Elemento Obrigatório Faltando
		"The element 'det' in namespace 'http://...' has incomplete content. List of possible elements expected: 'cProd'.",
		"The element 'emit' in namespace 'http://...' has incomplete content. List of possible elements expected: 'xNome'.",
		"The element 'enderEmit' in namespace 'http://...' has incomplete content. List of possible elements expected: 'xLgr'.",

		// Padrão: Ordem Incorreta dos Elementos
		"The element 'enderEmit' in namespace 'http://...' is not expected. Expected element: 'IE'.",
		"The element 'prod' in namespace 'http://...' is not expected. Expected element: 'imposto'.",

		// TO-DO: (Para testar o fallback e mostrar a necessidade de novas Regex)

		// Erro: Atributo Inválido (sua Regex atual só pega 'element')
		"Value 'ABC' for attribute 'versao' on element 'infNFe' is not valid with respect to its type, 'TVerNFe'.",

		// Erro: Tamanho Mínimo/Máximo (padrão de erro diferente)
		"Element 'xNome': [minLength error] The value has a length of '1'; this is less than the required '2'.",
		"Element 'xProd': [maxLength error] The value has a length of '150'; this is greater than the allowed '120'.",

		// Erro: "Token" Inválido (ex: caractere especial não permitido)
		"Element 'xMun': The value 'Porto Alegre ©' is invalid according to its data type 'TString'. Invalid token.",
	}

	// Carrega nossos padrões
	patterns := newPatterns()

	fmt.Println("--- Processando Erros de Schema XML ---")

	// Itera sobre cada erro cru
	for _, rawError := range rawErrors {
		fmt.Printf("\n[ERRO CRU]\n%s\n", rawError)

		foundMatch := false
		// Tenta aplicar cada um dos nossos padrões
		for _, p := range patterns {
			matches := p.Regex.FindStringSubmatch(rawError)

			// Se 'matches' não for nulo, a Regex funcionou!
			if matches != nil {
				// Executa a função Tradutora específica desse padrão
				friendlyError := p.Translator(matches, tagMap)
				fmt.Printf("\n[TRADUÇÃO 💡]\n%s\n", friendlyError)
				foundMatch = true
				break // Para de procurar, já encontramos o padrão correto
			}
		}

		// Se nenhum padrão funcionar, damos um fallback
		if !foundMatch {
			fmt.Printf("\n[TRADUÇÃO ⚠️]\nErro de schema não mapeado: %s\n", rawError)
		}
		fmt.Println(strings.Repeat("-", 40))
	}
}
