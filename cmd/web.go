package cmd

import (
	"encoding/json"
	"net/http"

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
		http.ListenAndServe(":9876", server.Server)
	},
}

func (s *vocabularyServer) getStats(ctx *web.Context){
	stats := s.Vocabulary.GetVocabularyStats([]string{})
	jsonEncoder := json.NewEncoder(ctx.ResponseWriter)
	jsonEncoder.Encode(stats)
}