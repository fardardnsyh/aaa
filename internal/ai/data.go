package ai

import (
	"bytes"
	"fmt"
	"os"

	chroma_go "github.com/amikos-tech/chroma-go"
	"github.com/dslipak/pdf"
	"github.com/labstack/echo/v4"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores/chroma"
)

var doclength int

func InsertToVectorDb(ctx echo.Context, filename string) (*chroma.Store, error) {

	content, err := GetTextFromPdf("./uploads/file.pdf")
	if err != nil {
		return nil, err
	}

	// fmt.Println(content)
	docs := GetTextChunks(content)
	llm, err := ollama.New(ollama.WithModel("mistral"))
	if err != nil {
		return nil, err
	}

	doclength = len(docs)

	ollamaEmbedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return nil, err
	}

	namespace := filename

	chromaUrl := os.Getenv("CHROMA_URL")

	store, err := chroma.New(
		chroma.WithChromaURL(chromaUrl),
		chroma.WithEmbedder(ollamaEmbedder),
		chroma.WithNameSpace(namespace),
		chroma.WithDistanceFunction(chroma_go.COSINE),
	)
	if err != nil {
		return nil, err
	}

	errAdd := store.AddDocuments(echo.Context.Request(ctx).Context(), docs)

	if errAdd != nil {
		return nil, err
	}
	fmt.Println("Documents embedded to store!")

	return &store, nil
}

func GetTextChunks(text string) []schema.Document {

	splitter := textsplitter.NewRecursiveCharacter()
	splitter.ChunkSize = 500
	chunks, err := splitter.SplitText(text)
	if err != nil {
		fmt.Println(err)
	}

	docs := []schema.Document{}

	for _, value := range chunks {
		doc := schema.Document{
			PageContent: value,
		}
		docs = append(docs, doc)
	}

	return docs

}

func GetTextFromPdf(path string) (string, error) {
	r, err := pdf.Open(path)
	// remember close file
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}
