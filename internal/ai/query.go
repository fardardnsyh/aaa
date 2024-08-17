package ai

import (
	"context"
	"fmt"
	"os"

	chroma_go "github.com/amikos-tech/chroma-go"
	"github.com/labstack/echo/v4"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/chroma"
)

func QueryToVectorDB(querry string, ctx echo.Context, filename string) (map[string]any, error) {

	doclength := 10
	fmt.Println("doclength is", doclength)

	llm, err := ollama.New(ollama.WithModel("mistral"))
	if err != nil {
		return nil, err
	}

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

	vecOpt := []vectorstores.Option{vectorstores.WithScoreThreshold(0)}
	fmt.Println("Before similartiy search")
	fmt.Println(store)

	resultDocs, resultErr := store.SimilaritySearch(ctx.Request().Context(), querry, doclength, vecOpt...)

	if resultErr != nil {
		return nil, err
	}
	for _, docR := range resultDocs {
		fmt.Println("CONTENT")
		fmt.Println(docR.PageContent)
		fmt.Println(docR.Score)
	}

	stuffQAChain := chains.LoadStuffQA(llm)

	answer, err := chains.Call(context.Background(), stuffQAChain, map[string]any{
		"input_documents": resultDocs,
		"question":        querry,
	})

	if err != nil {
		return nil, err
	}
	fmt.Println("ANSWERRRRRR!")
	fmt.Println(answer)

	return answer, nil

}
