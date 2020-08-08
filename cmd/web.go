package cmd

import (
	"encoding/json"
	"net/http"

	"path"

	voc "github.com/JaCoB1123/vocabulary/internal/vocabulary"
	"github.com/JaCoB1123/web"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(webCmd)
}

type vocabularyServer struct {
	*web.Server
	*voc.Vocabulary
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Serve vocabulary trainer via HTTP",
	Long:  "Used to access the vocabulary via a webbrowser",
	Run: func(cmd *cobra.Command, args []string) {
		server := &vocabularyServer{
			Server: web.NewServer(),
			Vocabulary: voc.MustVocabulary(*wordsfilename, *statsfilename),
		}

		server.Server.Get("/vocabulary/stats", server.getStats)
		server.Server.Get("/vocabulary/learn", server.learn)
		server.Server.Get("/(.*)", server.static)
		http.ListenAndServe(":9876", server.Server)
	},
}

func (s *vocabularyServer) getStats(ctx *web.Context){
	stats := s.Vocabulary.GetVocabularyStats([]string{})
	jsonEncoder := json.NewEncoder(ctx.ResponseWriter)
	jsonEncoder.Encode(stats)
}

func (s *vocabularyServer) learn(ctx *web.Context){
	words := s.Vocabulary.GetSortedWords([]string{})
	jsonEncoder := json.NewEncoder(ctx.ResponseWriter)
	jsonEncoder.Encode(words)
}

func (s *vocabularyServer) static(ctx *web.Context, filepath string){
	http.ServeFile(ctx.ResponseWriter, ctx.Request, path.Join("web", filepath))
}