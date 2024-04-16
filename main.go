package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

// Função para converter a data do formato MM/DD/AAAA para DD/MM/AAAA
func converterData(data string) string {
	t, err := time.Parse("1/2/2006", data)
	if err != nil {
		fmt.Println("Erro ao converter a data:", err)
		return data
	}
	return t.Format("02/01/2006")
}

func main() {
	// Verificar se um nome de arquivo foi fornecido como argumento da linha de comando
	if len(os.Args) < 2 {
		fmt.Println("Uso: ./converte_data arquivo.csv")
		return
	}

	// Obter o nome do arquivo CSV a ser processado a partir do argumento da linha de comando
	inputFile := os.Args[1]

	// Abrir arquivo CSV de entrada
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()

	// Criar leitor CSV
	reader := csv.NewReader(file)
	reader.Comma = ',' // Definir o delimitador como vírgula

	// Ler o cabeçalho do arquivo CSV
	header, err := reader.Read()
	if err != nil {
		fmt.Println("Erro ao ler o cabeçalho do arquivo:", err)
		return
	}

	// Verificar se o cabeçalho está vazio
	if len(header) == 0 {
		fmt.Println("Cabeçalho do arquivo vazio.")
		return
	}

	// Separar o nome do arquivo e sua extensão para gerar o nome do arquivo de saída
	parts := strings.Split(inputFile, ".")
	outputFile := parts[0] + "_modificado." + parts[1]

	// Abrir arquivo CSV de saída
	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Erro ao criar o arquivo de saída:", err)
		return
	}
	defer outFile.Close()

	// Criar escritor CSV
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	// Escrever o cabeçalho no arquivo de saída
	if err := writer.Write(header); err != nil {
		fmt.Println("Erro ao escrever o cabeçalho no arquivo de saída:", err)
		return
	}

	// Ler e escrever linhas do arquivo CSV
	for i := 1; ; i++ {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Println("Erro ao ler o arquivo:", err)
			return
		}

		// Verificar se a linha possui o número esperado de campos
		if len(record) != len(header) {
			fmt.Printf("Linha %d ignorada devido ao número incorreto de campos: %v\n", i, record)
			continue
		}

		// Se a linha não começar com "Issued To", é uma linha de dados
		if record[0] != "Issued To" {
			// Converter a data da coluna 'Expiration Date' para o novo formato
			record[2] = converterData(record[2])
		}

		// Escrever linha modificada no arquivo de saída
		if err := writer.Write(record); err != nil {
			fmt.Println("Erro ao escrever no arquivo de saída:", err)
			return
		}
	}

	fmt.Println("Arquivo CSV modificado foi salvo com sucesso.")
}
